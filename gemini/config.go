package gemini

import "mime"


const (
    BufferSize = 1024
    ConnectionDeadlineSeconds = 8

    FileServerStartGeminiFilename = "start.gmi"

    /* TODO What is the deal with `language`? Is this even correct MIME usage? */
    DefaultMimeLanguage = "en"
)


func ExtensionAndLanguageToMimeTypeAndTextStatus(ext string, language string) (string, bool) {
    switch ext {
    case ".txt": return "text/plain; charset=utf-8; lang=" + language, true
    case ".md": return "text/markdown; charset=utf-8; lang=" + language, true
    case ".gmi": return "text/gemini; charset=utf-8; lang=" + language, true
    /* One might add more MIME types. */
    }

    return mime.TypeByExtension(ext), false
}
