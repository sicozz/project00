#!/bin/bash

usage() {
    echo "Usage: $0 <version>"
    exit 1
}

if [ -z "$1" ]; then
    usage
fi

# Params
version=$1
makefile='Makefile'
dockerenv='docker/.env'

replace_third_word() {
    local file=$1
    local version=$2
    awk -v version="$version" '
    NR == 2 {
        $3 = version
    }
    {
        print
    }
    ' "$file" > tmpfile && mv tmpfile "$file"
    echo "[+] $file"
}

replace_value_after_equal() {
    local file=$1
    local version=$2
    awk -v version="\"$version\"" '
    NR == 2 {
        sub(/=.*/, "=" version)
    }
    {
        print
    }
    ' "$file" > tmpfile && mv tmpfile "$file"
    echo "[+] $file"
}

replace_third_word "$makefile" "$version"
replace_value_after_equal "$dockerenv" "$version"
