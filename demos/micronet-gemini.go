// Jonathan Frech, 2021-02-27

/*
   A Jemini demo.
   See: `gemini://jemini-demo.jfrech.com` and `https://github.com/jfrech/jemini`

   LetsEncrypt is not required as a certificate authority, yet using their
   service may be the most convenient way to get a TLS certificate:

   # echo "assuming a Debian-esque server"
   # echo "replace 'jemini-demo.jfrech.com' with this machine's domain"
   # apt install certbot
   # certbot certonly --standalone --preferred-challenges http -d "jemini-demo.jfrech.com"
   # go get pkg.jfrech.com/jemini
   # cd ~/go/src/pkg.jfrech.com/jemini/standalone-demo/ && go run standalone-demo.go
*/

package main
import "pkg.jfrech.com/jemini"


func main() {
    jemini.RunDemo(
        "jemini-demo.jfrech.com",
        "/etc/letsencrypt/live/jemini-demo.jfrech.com/fullchain.pem",
        "/etc/letsencrypt/live/jemini-demo.jfrech.com/privkey.pem",
    )
}
