// Jonathan Frech, 2021-02-27, 2021-03-08

/*
   A micronet/gemini demo. See:
   * `gemini://micronet-gemini.jfrech.com`
   * `https://github.com/jfrech/micronet/tree/master/gemini`

   LetsEncrypt is not required as a certificate authority, yet using their
   service may be the most convenient way to get a TLS certificate:
   # echo "assuming a Debian-esque server"
   # echo "replace 'micronet-gemini.jfrech.com' with this machine's domain"
   # apt install certbot
   # certbot certonly --standalone --preferred-challenges http -d "micronet-gemini.jfrech.com"
   # go get pkg.jfrech.com/micronet/gemini
   # cd ~/go/src/pkg.jfrech.com/micronet/gemini/demos && go run micronet-gemini-demo.go
*/

package main
import "pkg.jfrech.com/micronet/gemini"


func main() {
    gemini.RunDemo(
        "micronet-gemini.jfrech.com",
        "/etc/letsencrypt/live/micronet-gemini.jfrech.com/fullchain.pem",
        "/etc/letsencrypt/live/micronet-gemini.jfrech.com/privkey.pem",
    )
}
