package gemini

import (
    "fmt"
    "net"
    "net/url"
    "strings"
    "unicode/utf8"
)


type Connection struct {
    conn   net.Conn
    url    *url.URL

    status int
    Err    error
}

type handler func(Connection)error


func (gc *Connection) Header(status int, meta string) error {
    if gc.status == StatusInternalErroneous || gc.Err != nil {
        return gc.Err
    }

    if gc.status != StatusInternalNone {
        gc.status = StatusInternalErroneous
        return fmt.Errorf("status already set to: %v", gc.status)
    }
    if !ValidStatus(status) {
        return gc.Errorf("invalid status: %v", status)
    }
    if !utf8.Valid([]byte(meta)) {
        return gc.Errorf("invalid meta string encoding: %v", meta)
    }
    if strings.TrimSpace(meta) != meta || strings.ContainsAny(meta, "\r\n") {
        return gc.Errorf("invalid meta string: %v", meta)
    }

    if _, err := fmt.Fprintf(gc.conn, "%02d %s\r\n", status, meta); err != nil {
        return gc.Errorf("failed to write header")
    }

    gc.status = status
    return gc.Err
}

func (gc *Connection) RawBody(body []byte) error {
    if gc.status == StatusInternalErroneous || gc.Err != nil {
        return gc.Err
    }

    if gc.status != StatusSuccess {
        return gc.Errorf("failed to write body: invalid status")
    }

    if _, err := gc.conn.Write(body); err != nil {
        return gc.Errorf("failed to write body: %v", err)
    }

    return gc.Err
}

func (gc *Connection) Body(body string) error {
    if !utf8.Valid([]byte(body)) {
        return gc.Errorf("invalid body encoding")
    }
    for _, c := range body {
        if '\x00' <= c && c != '\t' && c != '\n' && c < ' ' {
            return gc.Errorf("invalid body")
        }
    }

    return gc.RawBody([]byte(body))
}

func (gc *Connection) Bodyf(format string, a ...interface{}) error {
    return gc.Body(fmt.Sprintf(format, a...))
}

func (gc *Connection) Bodyln(a ...interface{}) error {
    return gc.Body(fmt.Sprintln(a...))
}


func (gc *Connection) Url() *url.URL {
    if gc.url == nil {
        buf := make([]byte, BufferSize)
        n, err := gc.conn.Read(buf)
        if err != nil || n >= BufferSize {
            gc.ClientErrorf(StatusBadRequest, "could not read url: too long")
            return nil
        }
        rawurl := string(buf[:n])
        if !strings.HasSuffix(rawurl, "\r\n") {
            gc.ClientErrorf(StatusBadRequest, "could not read url: bad terminating sequence")
            return nil
        }
        rawurl = strings.TrimSuffix(rawurl, "\r\n")

        url, err := url.Parse(rawurl)
        if err != nil {
            gc.ClientErrorf(StatusBadRequest, "could not read url: parse error")
            return nil
        }

        if url.Scheme != "gemini" && url.Scheme != "gemini:1965" {
            gc.ClientErrorf(StatusBadRequest, "invalid scheme: %q", url.Scheme)
            return nil
        }

        gc.url = url
    }

    return gc.url
}

func (gc *Connection) RemoteAddr() net.Addr {
    return gc.conn.RemoteAddr()
}


func (gc *Connection) Errorf(format string, a ...interface{}) error {
    if gc.status == StatusInternalNone {
        gc.Header(StatusTemporaryFailure, fmt.Sprintf("internal server error: " + format, a...))
    }

    gc.status = StatusInternalErroneous
    gc.Err = fmt.Errorf(format, a...)
    return gc.Err
}

func (gc *Connection) ClientErrorf(status int, format string, a ...interface{}) error {
    if gc.status == StatusInternalNone {
        gc.Header(status, fmt.Sprintf(format, a...))
    }

    gc.Err = fmt.Errorf(format, a...)
    return gc.Err
}
