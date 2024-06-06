#!python

import os
import click
import json


i18n_files = [
'en_US',
'de_DE',
'es_ES',
'fr_FR',
'it_IT',
'ja_JP',
'ko_KR',
'pt_BR',
'zh_Hans',
'zh_Hant',
]




def reformat(filename: str):
    plugin_dir = os.path.join(os.getcwd(), 'plugin')
    i18n_path = os.path.join(plugin_dir, 'i18n', 'v1Resources', f"{filename}.json")
    new_i18n = os.path.join(plugin_dir, 'i18n', 'v2Resources', f"active.{filename}.json")
    original = get_source_data(i18n_path)
    reformatted = {}
    for key in original:
        translation = original[key].get('translation', '')
        if translation == '':
            translation = key
        reformatted[key] = {"other": translation}
    with open(new_i18n, 'w', encoding="utf8", newline='\n') as f:
       json.dump(reformatted, f, sort_keys=True, separators=(',', ': ',), ensure_ascii=False, indent=2)
    f.close()

def get_source_data(file_name: str) -> dict:
    """Reads from the i18n files and returns a formatted dict"""
    source_i18n = {}
    with open(file_name, encoding="utf8") as f:
        data = json.load(f)
        for i in data:
            source_i18n[i.get('id')] = i
    f.close()
    return source_i18n

@click.command()
def cli():
    cwd = os.getcwd()
    if not cwd.endswith('softlayer-cli'):
        raise Exception(f"Working Directory should be softlayer-cli, is currently {self.cwd}")
    for file in i18n_files:
    	click.echo(f"Working on {file}")
    	reformat(file)


if __name__ == '__main__':
    cli()