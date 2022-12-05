#!/bin/bash

set -e

echo "Generating i18n resource file ..."
./bin/go-bindata -pkg resources -o plugin/resources/i18n_resources.go plugin/i18n/resources
echo "Done."
