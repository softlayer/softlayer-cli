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
UNTRANSLATED=$'[\n'
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

    elif [[ $line =~ "exists in en_US, but not in "([a-zA-Z_]{5,7})"" ]]
    then
        # This is "normal", means the text has not been translated.
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in en_US, but not in ([a-zA-Z_]{5,7})/{\"id\": \1, \"translation\": \1, \"file\": \2.all.json},/g"`
        OUTPUT=""
        UNTRANSLATED=`printf "%s\n    %s" "${UNTRANSLATED}" "${JSON}"`
        JSON=""
        # EXIT_CODE=5

    # A translation exists in some files, but not in en_US
    elif [[ $line =~ "exists in "([a-zA-Z_]{5,7})", but not in en_US" ]]
    then
        JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in [a-zA-Z_]{5,7}, but not in en_US/{\"id\": \1, \"translation\": \1},/g"`
        # JSON=`echo "$OUTPUT" | sed -E "s/(.+) exists in ([a-zA-Z_]{5,7}), but not in ([a-zA-Z_]{5,7})/id: \1 is missing from \3.all.json, but exists in \2.all.json/g"`
        printf "\033[0;33m>>>  $JSON <<<\033[0m\n"
        OUTPUT=""
        DEL_JSON_OUT=`printf "%s\n    %s" "${DEL_JSON_OUT}" "${JSON}"`
        JSON=""
        # printf "\033[0;31m>>> cd plugin; ../bin/i18n4go -c checkup -q i18n -v <<<\033[0m\n"
        # printf "\033[0;31m>>> The translation files are out of sync. Run \`python bin/sync_enUS.py\`.<<<\033[0m\n"
        # echo $RESULTS
        EXIT_CODE=4

    # Should be the last line
    elif [[ $line =~ "Could not checkup, err: Strings don't match" ]]
    then
        printf "\033[0;31m>>> Could not checkup, err: Strings don't match.  <<<\033[0m\n"
        printf "\033[0;31m>>> Make sure any strings in  old-i18n/add_these.json are added <<<\033[0m\n"
        printf "\033[0;31m>>> cd plugin; ../bin/i18n4go -c checkup -q i18n -v <<<\033[0m\n"
    # There is a newline in our string, so we need to add it back in
    else
        OUTPUT="${OUTPUT}\n"
        # echo "|${OUTPUT}|  <- OUTPUT "
    fi
    
done <<< "$RESULTS"

# Junk
if [[ $EXIT_CODE -ge 1 ]]
then

    ADD_JSON_OUT=${ADD_JSON_OUT%,}
    ADD_JSON_OUT=`printf "%s\n]" "${ADD_JSON_OUT}"`
    # JSON panics when it hits tab characters.
    ADD_JSON_OUT=`echo "$ADD_JSON_OUT" | sed -E 's/\t/\\\t/g'`
    printf "\033[0;34m ====== ADD THESE ======= \033[0m \n"
    echo "$ADD_JSON_OUT"
    echo "$ADD_JSON_OUT" > ../old-i18n/add_these.json

    DEL_JSON_OUT=${DEL_JSON_OUT%,}
    DEL_JSON_OUT=`printf "%s\n]" "${DEL_JSON_OUT}"`
    # JSON panics when it hits tab characters.
    DEL_JSON_OUT=`echo "$DEL_JSON_OUT" | sed -E 's/\t/\\\t/g'`
    printf "\033[0;34m ====== DEL THESE ======= \033[0m \n"
    echo "$DEL_JSON_OUT"
    echo "$DEL_JSON_OUT" > ../old-i18n/remove_these.json


    # UNCOMMENT this line if you want to see what needs to be translated.
    # It will likely have an entry for every language though.
    # UNTRANSLATED=${UNTRANSLATED%,}
    # UNTRANSLATED=`printf "%s\n]" "${UNTRANSLATED}"`
    # UNTRANSLATED=`echo "$UNTRANSLATED" | sed -E 's/\t/\\\t/g'`
    # printf "\033[0;34m ====== TRANSLATE THESE ======= \033[0m \n"
    # echo "$UNTRANSLATED"
    # echo "$UNTRANSLATED" > ../old-i18n/translate_these.json

    exit $EXIT_CODE
else
    printf "\033[0;32mEverything looks good to me.\033[0m\n"
fi