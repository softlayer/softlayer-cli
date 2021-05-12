#!/bin/bash

set -e

echo "Generating i18n resource file ..."
$GOPATH/bin/go-bindata -pkg resources -o plugin/resources/i18n_resources.go -nocompress plugin/i18n/resources
echo "Done."
