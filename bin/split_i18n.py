

"""
This program will take a json file (./plugin/i18n/resources/bad.json) and use it to remove those values
from ./plugin/i18n/resources/*.all.json


HELPFUL REGEX:
Takes output from `i18n4go -c checkup -q i18n -v` and helps turn it into json, poorly.
1: /(?s)^"(.+?)" (exists in the code, but not in en_US){1}/{"id": "\1", "translation": "\1"},
2: /},\\n{/},\n{
3: Make sure there are not too many "" lying around.
"""

import json
from pprint import pprint as pp
  
# Opening JSON file

# returns JSON object as 
# a dictionary

files = [
'de_DE.all.json',
'en_US.all.json',
'es_ES.all.json',
'fr_FR.all.json',
'it_IT.all.json',
'ja_JP.all.json',
'ko_KR.all.json',
'pt_BR.all.json',
'zh_Hans.all.json',
'zh_Hant.all.json',
]

def prune_bad_matches(bad_dct, mixed_dct):

    for key in bad_dct.keys():
        if key in mixed_dct:
            # print("|{}| is bad".format(key))
            mixed_dct.pop(key, None)

    clean_i18n = []
    for i in sorted(mixed_dct.keys()):
        clean_i18n.append(mixed_dct[i])

    return clean_i18n

def keep_good_matches(bad_dct, mixed_dct):

    clean_i18n = []
    for key in bad_dct.keys():
        if key not in mixed_dct:
            mixed_dct[key] = bad_dct[key]

    for i in sorted(mixed_dct.keys()):
        clean_i18n.append(mixed_dct[i])

    return clean_i18n


def cleanup_i18n_file(file_name, bad_file='bad.json', bad=True):
    """

    file_name: path to i18n file you want to modify
    bad_file: file with the reference ids you want to remove (or keep)
    base_path: folder with all of these files in it
    bad: True if you want to remove matches, true if you want to remove non-matches
    """

    bad_i18n = {}
    mixed_i18n = {}

    with open(bad_file, encoding="utf8") as f:
        data = json.load(f)
        for i in data:
            # pp(i.get(['id']))
            # print("==============================")
            bad_i18n[i.get('id')] = i
    f.close()


    with open(file_name, encoding="utf8") as f:
        data = json.load(f)
        for i in data:
            mixed_i18n[i.get('id')] = i
    f.close()

    if bad:
        clean_i18n = prune_bad_matches(bad_i18n, mixed_i18n)
    else:
        clean_i18n = keep_good_matches(bad_i18n, mixed_i18n)

    with open(file_name, 'w', encoding="utf8") as f:
       json.dump(clean_i18n, f, sort_keys=True, separators=(',\n', ': '), ensure_ascii=False)
    f.close()



for i18n in files:
    # Removes everything not in en_US.all.json
    # cleanup_i18n_file('./plugin/i18n/resources/' + i18n, bad_file='en_US.all.json', bad=False)

    # cleans up github.ibm.com/bluemix/bluemix-cli
    base_path = '/Users/allmi/go/src/github.ibm.com/Bluemix/bluemix-cli/bluemix/i18n/resources/'
    # Remove everything in this project
    cleanup_i18n_file(base_path + i18n, bad_file='./plugin/i18n/resources/en_US.all.json')
    # Add these back in
    cleanup_i18n_file(base_path + i18n, bad_file='./old-i18n/bad2.json', bad=False)
    # remove these again
    cleanup_i18n_file(base_path + i18n, bad_file='./old-i18n/bad3.json')