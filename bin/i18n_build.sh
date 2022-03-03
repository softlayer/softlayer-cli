#!/bin/bash

# This file will build the EN_US.all.json file from the soruce files.
# ./bin/.i18n_build.sh

DIRECTORIES="plugin/commands/ plugin/metadata/ plugin/version/ plugin/client/ plugin/errors/ plugin/managers/ plugin/testfixtures/ plugin/utils/"

# for i in $DIRECTORIES
# do
#     echo "Building |$i|"
#     i18n4go -c extract-strings -v -o plugin/i18n/tmp_resources/$i -d $i  -r --ignore-regexp  ".*test.*"

# done
./bin/i18n4go -c extract-strings -v -o plugin/i18n/tmp_resources/ -d plugin  -r --ignore-regexp  ".*test.*" --output-match-package
./bin/i18n4go -c merge-strings -v -r -d plugin/i18n/tmp_resources/
# 
# i18n4go -c merge-strings -v -d plugin/i18n/tmp_resources/ 
# mv plugin/i18n/tmp_resources/en.all.json plugin/i18n/resources/en_US.all.json
# rm plugin/i18n/tmp_resources/*.go.en.json