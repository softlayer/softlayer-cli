#!/bin/bash


ROOT_DIR=${PWD##*/}
if [ "$ROOT_DIR" != "softlayer-cli" ]
then
    echo "Please run this command in the base softlayer-cli directory"
    exit 1
fi

# Switches to plugin directory or exits if it doesn't exist.
cd "./plugin" 

RESULTS=`../bin/i18n4go -c checkup -q i18n -v`
# RESULTS="OKTotal time:"
OUTPUT=""
ADD_JSON_OUT=$'[\n'
DEL_JSON_OUT=$'[\n'
if [[ "$RESULTS" =~ ^"OKTotal time: "[0-9\.]*ms$ ]]
then
    echo $RESULTS
    exit 0
fi

EXIT_CODE=0
    
while IFS= read -r line; do
    # This is so we can handle translation strings that span multiple lines
    OUTPUT="${OUTPUT}${line}"

    # These strings need to be ADDED
    if [[ $line =~ "exists in the code, but not in en_US" ]]
    then
        
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in the code, but not in en_US/{\"id\": \1, \"translation\": \1},/g"`
        # echo ">>> |$JSON| <<<"
        OUTPUT=""
        
        # This printf bit is required for bash to put real newlines between missing translations, but preserve
        # the newline '\n' characters in the actual translation strings themselves.
        ADD_JSON_OUT=`printf "%s\n    %s" "${ADD_JSON_OUT}" "${JSON}"`
        JSON=""
        EXIT_CODE=1

    # Need to REMOVE these strings
    elif [[ $line =~ "exists in en_US, but not in the code" ]]
    then
        
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in en_US, but not in the code/{\"id\": \1, \"translation\": \1},/g"`
        # echo ">>> |$JSON| <<<"
        OUTPUT=""
        
        # This printf bit is required for bash to put real newlines between missing translations, but preserve
        # the newline '\n' characters in the actual translation strings themselves.
        DEL_JSON_OUT=`printf "%s\n    %s" "${DEL_JSON_OUT}" "${JSON}"`
        JSON=""
        EXIT_CODE=2

    # A translation exists in en_US but missing from some other file
    elif [[ $line =~ "exists in en_US, but not in "[a-zA-Z_]{5,7} ]]
    then
        # JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in en_US, but not in [a-zA-Z_]{5,7}/{\"id\": \1, \"translation\": \1},/g"`
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in en_US, but not in ([a-zA-Z_]{5,7})/id: \1 is missing from \2.all.json/g"`
        printf "\t\033[0;33m>>> $JSON <<<\033[0m\n"
        OUTPUT=""
        # ADD_JSON_OUT=`printf "%s\n    %s" "${ADD_JSON_OUT}" "${JSON}"`
        JSON=""
        # printf "\033[0;31m>>> cd plugin; ../bin/i18n4go -c checkup -q i18n -v <<<\033[0m\n"
        # printf "\033[0;31m>>> The translation files are out of sync. Run \`python bin/sync_enUS.py\`.<<<\033[0m\n"
        # echo $RESULTS
        # exit 4

    # Should be the last line
    elif [[ $line =~ "Could not checkup, err: Strings don't match" ]]
    then
        # Add these strings
        # Remove the ending ","
        ADD_JSON_OUT=${ADD_JSON_OUT%,}
        ADD_JSON_OUT=`printf "%s\n]" "${ADD_JSON_OUT}"`
        # JSON panics when it hits tab characters.
        ADD_JSON_OUT=`echo "$ADD_JSON_OUT" | sed -E 's/\t/\\\t/g'`
        printf "\033[0;34m ====== ADD THESE ======= \033[0m \n"
        echo "$ADD_JSON_OUT"
        echo "$ADD_JSON_OUT" > ../old-i18n/add_these.json

        # Remove these strings
        # Remove the ending ","
        DEL_JSON_OUT=${DEL_JSON_OUT%,}
        DEL_JSON_OUT=`printf "%s\n]" "${DEL_JSON_OUT}"`
        # JSON panics when it hits tab characters.
        DEL_JSON_OUT=`echo "$DEL_JSON_OUT" | sed -E 's/\t/\\\t/g'`
        printf "\033[0;34m ====== DEL THESE ======= \033[0m \n"
        echo "$DEL_JSON_OUT"
        echo "$DEL_JSON_OUT" > ../old-i18n/remove_these.json
        # python ./bin/split_i18n.py
        exit $EXIT_CODE
    # There is a newline in our string, so we need to add it back in
    else
        OUTPUT="${OUTPUT}\n"
        # echo "|${OUTPUT}|  <- OUTPUT "
    fi
    
done <<< "$RESULTS"
