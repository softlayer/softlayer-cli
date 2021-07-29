#!/bin/bash


ROOT_DIR=${PWD##*/}
if [ "$ROOT_DIR" != "softlayer-cli" ]
then
    echo "Please run this command in the base softlayer-cli directory"
    exit 1
fi

# Switches to plugin directory or exits if it doesn't exist.
pushd "./plugin" > /dev/null || exit 1

RESULTS=`i18n4go -c checkup -q i18n -v`
# RESULTS="OKTotal time:"
OUTPUT=""
ADD_JSON_OUT=$'[\n'
DEL_JSON_OUT=$'[\n'
if [[ "$RESULTS" =~ "OKTotal time:" ]]
then
    echo $RESULTS
    exit 0
fi

    
while IFS= read -r line; do
    # This is so we can handle translation strings that span multiple lines
    OUTPUT="${OUTPUT}${line}"

    # These strings need to be ADDED
    if [[ $line =~ "exists in the code, but not in en_US" ]]
    then
        
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in the code, but not in en_US/{\"id\": \1, \"translation\": \1},/g"`
        echo ">>> |$OUTPUT| <<<"
        OUTPUT=""
        
        # This printf bit is required for bash to put real newlines between missing translations, but preserve
        # the newline '\n' characters in the actual translation strings themselves.
        ADD_JSON_OUT=`printf "%s\n\t%s" "${ADD_JSON_OUT}" "${JSON}"`
        JSON=""

    # Need to REMOVE these strings
    elif [[ $line =~ "exists in en_US, but not in the code" ]]
    then
        
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in en_US, but not in the code/{\"id\": \1, \"translation\": \1},/g"`
        echo ">>> |$OUTPUT| <<<"
        OUTPUT=""
        
        # This printf bit is required for bash to put real newlines between missing translations, but preserve
        # the newline '\n' characters in the actual translation strings themselves.
        DEL_JSON_OUT=`printf "%s\n\t%s" "${DEL_JSON_OUT}" "${JSON}"`
        JSON=""
    # Should be the last line
    elif [[ $line =~ "Could not checkup, err: Strings don't match" ]]
    then
        # Add these strings
        # Remove the ending ","
        ADD_JSON_OUT=${ADD_JSON_OUT%,}
        ADD_JSON_OUT=`printf "%s\n]" "${ADD_JSON_OUT}"`
        printf "====== ADD THESE =======\n"
        echo "$ADD_JSON_OUT"
        echo "$ADD_JSON_OUT" > ../old-i18n/add_these.json

        # Remove these strings
        # Remove the ending ","
        DEL_JSON_OUT=${DEL_JSON_OUT%,}
        DEL_JSON_OUT=`printf "%s\n]" "${DEL_JSON_OUT}"`
        printf "====== DEL THESE =======\n"
        echo "$DEL_JSON_OUT"
        echo "$DEL_JSON_OUT" > ../old-i18n/remove_these.json
        # python ./bin/split_i18n.py
        exit 3
    # There is a newline in our string, so we need to add it back in
    else
        OUTPUT="${OUTPUT}\n"
        # echo "|${OUTPUT}|  <- OUTPUT "
    fi
    
done <<< "$RESULTS"

popd > /dev/null || exit 1