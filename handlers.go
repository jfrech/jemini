package jemini

import (
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)


func FileServerHandler(root string, language string) geminiConnectionHandler {
    if language == "" {
        language = DefaultMimeLanguage
    }

    return func(gc GeminiConnection) error {
        if !filepath.IsAbs(root) {
            return gc.Errorf("misconfigured server-side root path")
        }

        /* syntactic validation */
        path := path.Join(root, path.Clean("/" + gc.Url().Path))
        if path == root {
            path = root + "/"
        }
        if !filepath.IsAbs(path) || !strings.HasPrefix(path, root + "/") {
            return gc.ClientErrorf(StatusBadRequest, "not an absolute path after cleaning")
        }
        for _, component := range strings.Split(strings.TrimPrefix(path, root + "/"), "/") {
            if len(component) > 0 && component[0] == '.' {
                return gc.ClientErrorf(StatusNotFound, "syntactically hidden")
            }
        }
        if path == root + "/" {
            path = filepath.Join(path, "/start.gmi")
        }


        /* directory listing */
        if info, err := os.Stat(path); err == nil && info.Mode().IsDir() {
            gc.Header(StatusSuccess, "text/gemini; charset=utf-8; lang=" + language)
            gc.Body("# Contents of the directory »" + GeminiEscape(strings.TrimPrefix(path, root)))
            if dir, err := ioutil.ReadDir(path); err != nil {
                gc.Body("(Directory listing failed.)\n")
            } else {
                for _, fp_ := range dir {
                    fp := fp_.Name()

                    if info, err := os.Stat(filepath.Join(path, fp)); err == nil && info.Mode().IsRegular() {
                        gc.Body("=> " + gc.Url().Path + fp + " [f] " + fp + "\n")
                    } else if info, err := os.Stat(filepath.Join(path, fp)); err == nil && info.Mode().IsDir() {
                        gc.Body("=> " + gc.Url().Path + fp + " [d] " + fp + "\n")
                    } else {
                        gc.Body("[e]" + fp + "\n")
                    }
                }
            }

            return gc.Err
        }


        /* file serving */
        if info, err := os.Stat(path); err != nil || !info.Mode().IsRegular() {
            return gc.ClientErrorf(StatusNotFound, "not a regular file found on disk")
        }

        if bytes, err := ioutil.ReadFile(path); err != nil {
            return gc.ClientErrorf(StatusNotFound, "server-side disk error")
        } else {
            mtype, isText := ExtensionAndLanguageToMimeTypeAndTextStatus(filepath.Ext(path), language)
            if err := gc.Header(StatusSuccess, mtype); err != nil {
                return err
            }

            if isText {
                return gc.Body(string(bytes))
            }

            return gc.RawBody(bytes)
        }
    }
}


func RedirectHandler(status int, url string) geminiConnectionHandler {
    return func(gc GeminiConnection) error {
        if status != StatusTemporaryRedirect && status != StatusPermanentRedirect {
            return gc.Errorf("misconfigured server-side redirection")
        }

        return gc.Header(status, url)
    }
}
