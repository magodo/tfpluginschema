#!/bin/bash

die() {
    echo "$@" >&2
    exit 1
}

usage() {
    cat << EOF
Dump the provider schema into a json file. The script is expected to be run in the root directory of the provider repo.

Usage: ./${MYNAME} [options] provider_type

Options:
    -h|--help           Show this message
    -o                  Output file

Arguments:
    provider_type       The fully qualified address to the provider type.
                        E.g. github.com/Azure/terraform-provider-azapi/internal/provider.Provider
EOF
}

while :; do
    case $1 in
        -h|--help)
            usage 
            exit 1
            ;;
        -o)
            shift
            output_file=$1
            shift
            break
            ;;
        --)
            shift
            break
            ;;
        *)
            break
            ;;
    esac
    shift
done

expect_n_arg=1
[[ $# = "$expect_n_arg" ]] || die "wrong arguments (expected: $expect_n_arg, got: $#)"

provider_func=$1

type_fqname=${provider_func##*/} # e.g. provider.Provider
type_name=${provider_func##*.} # e.g. Provider
package_path=${provider_func%.${type_name}} # e.g. github.com/Azure/terraform-provider-azapi/internal/provider

[[ -d tfpluginschema ]] || mkdir tfpluginschema
cat << EOF > ./tfpluginschema/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"${package_path}"
	"github.com/magodo/tfpluginschema"
)

func main() {
    sch, err := tfpluginschema.FromFWProvider(&${type_fqname}{})
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(sch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
EOF

go mod tidy
go mod vendor
if [[ -n $output_file ]]; then
    go run ./tfpluginschema/main.go > $output_file
else
    go run ./tfpluginschema/main.go
fi
