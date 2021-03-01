#! /bin/sh

cd "$(dirname "$(realpath "$0")")/.."

printf '#! /bin/sh\n'
printf '\n'
printf 'install_target="$HOME/go/src/pkg.jfrech.com/jemini"\n'
printf 'rm -rf "$install_target" && mkdir -p "$install_target"\n'
printf '\n'
printf 'f () {\n'
printf '%s%s\n' "    printf '%s\\n' " '"$1"'
printf '    mkdir -p "$install_target/$(dirname "$1")"\n'
printf '    echo "$2" | base64 -d > "$install_target/$1"\n'
printf '}\n'
printf '\n'

ls go.mod *.go _standalone-demo/*.go | sort | while read gosrc; do
    printf "f '%s' '" "$gosrc"
    base64 -w0 "$gosrc"
    printf "'\n"
done
