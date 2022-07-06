

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


def sync_en_US(file_name, file_to_update):
    """Adds keys that are in en_US to all the other all.json files

    :param file_name: location of en_US.all.json
    :param file_to_update: location of target language file
    """
    print(f"Syncing {file_to_update} ... ", end='')
    # Need to be a dict, so the `id` property can be the key
    en_US = {}
    other_i18n = {}

    keys_updated = 0
    # Open en_US and import as json
    with open(file_name, encoding="utf8") as f1:
        data = json.load(f1)
        for i in data:
            en_US[i.get('id')] = i
    f1.close()

    # open other i18n file as json
    with open(file_to_update, encoding="utf8") as f2:
        data = json.load(f2)
        for i in data:
            other_i18n[i.get('id')] = i
    f2.close()

    # Iterate over the `id` of each translation
    for xlation in en_US.keys():
        # if its missing from the other file, add it
        if xlation not in other_i18n:
            other_i18n[xlation] = en_US[xlation]
            keys_updated = keys_updated + 1

    synced_i18n = []
    # Keeps everything sorted, and coverts back to a normal list
    for i in sorted(other_i18n.keys()):
        synced_i18n.append(other_i18n[i])

    with open(file_to_update, 'w', encoding="utf8", newline='\n') as f3:
       json.dump(synced_i18n, f3, sort_keys=True, separators=(',', ': ',), ensure_ascii=False, indent=2)
    f3.close()
    print(f"Done! Updated {keys_updated} entries")


for i18n in files:
    base_path = './plugin/i18n/resources/'
    # SYNC
    
    sync_en_US(base_path + 'en_US.all.json', base_path + i18n)
