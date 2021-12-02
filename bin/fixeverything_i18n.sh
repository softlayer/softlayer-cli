# This script is used to resolve any translation related trings.
# It does the following:
# 1. Finds the mismatched strings.
#    - Anything in the code, but not in en_US will get added to ./old-i18n/add_these.json
#    - Anything in en_US but NOT in the code, will get added to ./old-i18n/remove_these.json
# 2. ./bin/split_i18n.py will be run, adding/removing strings from all i18n files as specified in the add_these/remove_these files
# 3. ./bin/generate_i18n_resources.sh will be run, generating the golang resources required for translation
# 4. A quick commit of the expected changed files will be run
#     - `git add ./plugin/i18n/resources/*.json`
#     - `git add ./plugin/resources/i18n_resources.go`
#     - `git commit --message="Translation fixes from ./bin/fixeverything_i18n.sh`


ROOT_DIR=${PWD##*/}
if [ "$ROOT_DIR" != "softlayer-cli" ]
then
    echo "Please run this command in the base softlayer-cli directory"
    exit 1
fi

echo "Running: ./bin/catch-i18n-mismatch.sh"
./bin/catch-i18n-mismatch.sh
STATUS=$?

if [ $STATUS -eq 0 ]
then
    echo "I18N files are fine, nothing to do it seems."
    exit 0
fi

echo "Running: python ./bin/split_i18n.py"
python ./bin/split_i18n.py

echo "Running: ./bin/generate-i18n-resources.sh"
./bin/generate-i18n-resources.sh


echo "Running: git checkout ./old-i18n/*.json"
git checkout ./old-i18n/*.json


echo "Running: ./bin/catch-i18n-mismatch.sh a second time."
./bin/catch-i18n-mismatch.sh
STATUS=$?

if [ $STATUS -ne 0 ]
then
    echo "I18N files are still broken, please fix manually."
    exit $STATUS
fi

echo "Done"

