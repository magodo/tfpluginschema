#!/bin/bash

die() {
    echo "$@" >&2
    exit 1
}

usage() {
    cat << EOF
Dump the provider schema into a json file. The script is expected to be run in the root directory of the provider repo.

Usage: ./${MYNAME} [options] provider_func

Options:
    -h|--help           Show this message
    -o                  Output file

Arguments:
    provider_func       The fully qualified address to the provider function.
                        E.g. github.com/hashicorp/terraform-provider-azurerm/internal/provider.AzureProvider
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

func_call=${provider_func##*/} # e.g. provider.AzureProvider
func_name=${provider_func##*.} # e.g. AzureProvider
package_path=${provider_func%.${func_name}} # e.g. github.com/hashicorp/terraform-provider-azurerm/internal/provider

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
    sch := tfpluginschema.FromSDKv2Provider(${func_call}())
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
