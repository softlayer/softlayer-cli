#!/bin/bash

ROOT_DIR=$(cd $(dirname $(dirname $0)) && pwd)

pushd "${ROOT_DIR}/plugin" > /dev/null || exit 1
(i18n4go -c checkup -q i18n -v | sed -E 's/(.+) exists in the code, but not in en_US/{"id": \1, "translation": \1},/g')
popd > /dev/null || exit1