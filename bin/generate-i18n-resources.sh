#!/bin/bash

set -e

echo "Generating i18n resource file ..."
if [[ "$OSTYPE" == "msys"* ]]; then
    # For SO windows
    ./bin/go-bindata.exe -pkg resources -o plugin/resources/i18n_resources.go plugin/i18n/resources
else
    ./bin/go-bindata -pkg resources -o plugin/resources/i18n_resources.go plugin/i18n/resources
fi
echo "Done."
