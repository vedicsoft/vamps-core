#!/usr/bin/env bash

BINARY_NAME="vamps-core"

## install realize if not already installed
if ! [ -x "$(command -v realize)" ]; then
      echo "[INFO] installing realize for the first time"
      go get github.com/oxequa/realize
fi

function error_exit() {
	echo "[ERROR] $1" 1>&2
	exit 1
}

# checking code quality
function checkCC(){
    echo "[INFO] Vetting go files"
    if go vet $(go list ./... | grep -v /vendor/); then
        echo "[INFO] Go vetting is successful"
    else
        error_exit "Go vetting denied the changes"
    fi

    if ! [ -x "$(command -v goimports)" ]; then
      error_exit "goimports command not installed"
    fi
    # checking code formatting
    echo "[INFO] checking go import formatting"
    if [[ $(goimports -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./resources/*")) ]]; then
        error_exit "go code is not properly formatted. Please run build.sh before committing"
    else
        echo "[INFO] code quality check passed"
    fi
}

# format golang code
function fmt(){
    echo "[INFO] Vetting go files"
    if go vet $(go list ./... | grep -v /vendor/); then
        echo "[INFO] Go vetting is successful"
    else
        error_exit "Go vetting denied the changes"
    fi

    echo "[INFO] formatting go code"
    goimports -w $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./resources/*")
}

# build the binary of the server
function build_(){
    echo "[INFO] Building the binary with go configs:"
    go env
    echo "[INFO] Installing dependencies"
    if glide install; then
        echo "[INFO] Installed dependencies successfully"
    else
        error_exit "glide dependency installation failed"
    fi

    if go build .; then
        echo "[INFO] Binary generated successfully"
        echo "[INFO] Binary size: " $(du -hsk ${BINARY_NAME})
        if [[ $(du -hsk ${BINARY_NAME} | cut -f1) -lt 14000 ]]; then
         error_exit "generated binary size is too small"
        fi
        echo "[INFO] Binary Info: " $(file ${BINARY_NAME})
    else
        error_exit "Binary generation failed"
    fi
}

case "$1" in
    "")
        fmt
        runTests
        genDocs
    ;;
    fmt)
        fmt
    ;;
    cc)
        checkCC
    ;;
    ci)
        runTests
        genDocs
        build_
    ;;
    build)
        build_
    ;;
    docs)
        genDocs
    ;;
    run)
        realize start
    ;;
    *)
      echo $"usage: $0 {docs|ci|cc|run}"
            exit 1
    ;;
esac
