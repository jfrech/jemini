package gemini

import (
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "sort"
    "strings"
)


func FileServerHandler(root string, language string) handler {
    if language == "" {
        language = DefaultMimeLanguage
    }

    return func(gc Connection) error {
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
            gc.Bodyln("# Contents of the directory »" + Escape(strings.TrimPrefix(path, root)) + "«")
            if dir, err := ioutil.ReadDir(path); err != nil {
                gc.Bodyln("(Directory listing failed.)")
            } else {
                var files, directories, erroneous []string
                for _, fp_ := range dir {
                    fp := fp_.Name()

                    if info, err := os.Stat(filepath.Join(path, fp)); err == nil && info.Mode().IsRegular() {
                        files = append(files, fp)
                    } else if info, err := os.Stat(filepath.Join(path, fp)); err == nil && info.Mode().IsDir() {
                        directories = append(directories, fp)
                    } else {
                        erroneous = append(erroneous, fp)
                    }
                }

                total := 0
                write := func(fps []string, name string) {
                    sort.Strings(fps)
                    for _, fp := range fps {
                        gc.Bodyln("=> " + gc.Url().Path + "/" + fp + " [" + name + "] " + fp)
                        total++
                    }
                }
                write(files, "f")
                write(directories, "d")
                write(erroneous, "e")

                gc.Bodyf("(total: %d)\n", total)
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


func RedirectHandler(status int, url string) handler {
    return func(gc Connection) error {
        if status != StatusTemporaryRedirect && status != StatusPermanentRedirect {
            return gc.Errorf("misconfigured server-side redirection")
        }

        return gc.Header(status, url)
    }
}
