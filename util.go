package jemini

import (
    "fmt"
    "strings"
)


func GeminiEscape(raw string) string {
    var b strings.Builder
    for _, c := range []byte(raw) {
        if ' ' <= c && c <= '~' {
            b.WriteByte(c)
        } else {
            fmt.Fprintf(&b, "\\x%02X", c)
        }
    }

    return b.String()
}
