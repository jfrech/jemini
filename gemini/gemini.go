// Jonathan Frech, 2021-02-26, 2021-02-27, 2021-03-08

package gemini

import (
    "crypto/tls"
    "fmt"
    "log"
    "net"
    "time"
)


type Capsule struct {
    Domain string
    CertFullChain, CertPrivKey string
    Handler handler
}


func ListenAndServe(config *tls.Config, lowLevelHandler func(net.Conn)error) error {
    listener, err := tls.Listen("tcp", ":1965", config)
    if err != nil {
        return err
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Print(err)
            continue
        }

        go lowLevelHandler(conn)
    }
}

func Run(grealms []Capsule) error {
    certificates := make(map[string]*tls.Certificate)
    capsules := make(map[string]Capsule)

    for _, grealm := range grealms {
        if _, ok := certificates[grealm.Domain]; ok {
            return fmt.Errorf("duplicate gemini realm domain: %v", grealm.Domain)
        }
        cert, err := tls.LoadX509KeyPair(grealm.CertFullChain, grealm.CertPrivKey)
        if err != nil {
            return err
        }

        certificates[grealm.Domain] = &cert
        capsules[grealm.Domain] = grealm
    }

    getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
        if cert, ok := certificates[chi.ServerName]; ok {
            return cert, nil
        }

        return nil, fmt.Errorf("no certificate known for: %v", chi.ServerName)
    }

    config := &tls.Config{
        GetCertificate: getCertificate,

        /* see `https://golang.org/pkg/crypto/tls/#pkg-constants` */
        CipherSuites: []uint16{
            tls.TLS_AES_128_GCM_SHA256,
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256},
    }

    lowLevelHandler := func(conn net.Conn) error {
        defer conn.Close()
        conn.SetDeadline(time.Now().Add(ConnectionDeadlineSeconds * time.Second))
        gc := Connection{conn: conn,}

        if url := gc.Url(); url == nil {
            return gc.ClientErrorf(StatusPermanentFailure, "unknown host")
        } else if grealm, ok := capsules[url.Host]; !ok {
            return gc.ClientErrorf(StatusPermanentFailure, "unknown host: %q", url.Host)
        } else {
            return grealm.Handler(gc)
        }
    }

    return ListenAndServe(config, lowLevelHandler)
}
