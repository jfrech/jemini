package gemini

import (
    "log"
    "time"
)


func DemoHandler(gc GeminiConnection) error {
    t := time.Now()

    gc.Header(StatusSuccess, "text/gemini; charset=utf-8; lang=en")

    gc.Bodyln(`# micronet/gemini: a go implementation of a Gemini server`)
    gc.Bodyln(`micronet/gemini implements the gopher-like but TSL supporting Gemini protocol's server side in go, attempting to preserve its minimalistic aspirations. It is implemented in less than 512 lines of code and can be seemlessly integrated into an existing go web server and TLS certificate environment, allowing to both serve the dazzlingly white web as well as the possibly a tad fusty micro-web.`)
    gc.Bodyln(`=> https://github.com/jfrech/micronet/tree/master/gemini micronet/gemini on GitHub`)
    gc.Bodyln(`=> gemini://gemini.circumlunar.space/ Project Gemini`)
    gc.Bodyln(`=> https://go.dev go.dev`)
    gc.Bodyln("")

    gc.Bodyln(`## Author`)
    gc.Bodyln(`=> gemini://jfrech.com gemini://jfrech.com`)
    gc.Bodyln(`=> https://www.jfrech.com https://www.jfrech.com`)
    gc.Bodyln("")

    gc.Bodyln("## micronet/gemini demo")
    gc.Bodyf("* remote address: %v\n", gc.RemoteAddr())
    gc.Bodyf("* requested url: %v\n", gc.Url())
    gc.Bodyf("* request timestamp: %04d-%02d-%02d, %02d:%02d:%02d.%09d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
    gc.Bodyln("```")
    gc.Bodyln(`package main`)
    gc.Bodyln(`import "pkg.jfrech.com/micronet/gemini"`)
    gc.Bodyln(``)
    gc.Bodyln(``)
    gc.Bodyln(`func main() {`)
    gc.Bodyln(`    gemini.RunDemo(`)
    gc.Bodyln(`        "micronet-gemini.jfrech.com",`)
    gc.Bodyln(`        "/etc/letsencrypt/live/micronet-gemini.jfrech.com/fullchain.pem",`)
    gc.Bodyln(`        "/etc/letsencrypt/live/micronet-gemini.jfrech.com/privkey.pem",`)
    gc.Bodyln(`    )`)
    gc.Bodyln(`}`)
    gc.Bodyln("```")

    return gc.Err
}

func RunDemo(domain, certFullChain, certPrivKey string) {
    log.Printf("launching a demo gemini server on port :1965 and domain %q", domain)

    log.Fatal(Run([]Capsule{
        Capsule{
            Domain: domain,
            CertFullChain: certFullChain,
            CertPrivKey: certPrivKey,
            Handler: DemoHandler,
        },
    }))
}
