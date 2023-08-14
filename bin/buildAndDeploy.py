#!python

import click
import json
import os
import re
import subprocess
from rich import print


i18n_files = [
'en_US.all.json',
'de_DE.all.json',
'es_ES.all.json',
'fr_FR.all.json',
'it_IT.all.json',
'ja_JP.all.json',
'ko_KR.all.json',
'pt_BR.all.json',
'zh_Hans.all.json',
'zh_Hant.all.json',
]

def isWindows() -> bool:
    """returns if we are running on windows or not"""
    if os.name == 'nt':
        return True
    return False

def genBinData() -> None:
    """Generates the I18N Binary data required for translations"""
    goBindata = './bin/go-bindata'
    if isWindows():
        goBindata = f'{goBindata}.exe'
    goBindata = f'{goBindata} -pkg resources -o plugin/resources/i18n_resources.go plugin/i18n/resources'
    print("[green]Building I18N ...")
    print(f"\t[yellow]{goBindata}")
    result = subprocess.run(goBindata, capture_output=True)
    if result.returncode > 0:
        print(f"[red]{result.stderr.decode('ascii')}")
    else:
        print(f"\t[green]OK")

def buildArchs() -> dict:
    """Returns the list of binaries we should build"""
    buildArchs =  {
        'darwin': ['arm64', 'amd64'],
        'linux': ['amd64', '386', 'arm64', 'ppc64le', 's390x'],
        'windows': ['386', 'amd64'],
    }
    return buildArchs

def cgoEnable(theOs: str, theArch: str) -> int:
    """Disable cgo for these archs"""
    cgo_disable = ["amd64", "386", "arm64"]
    if theOs == "linux" and theArch in cgo_disable:
        return 0
    return 1

def goBuild(cwd: str, theOs: str, theArch: str) -> None:
    """Runs the go build command
    
    :param str cwd: The current working directory
    :param str theOs: OS to build for
    :param str theArch: Architecture to build for
    """
    envVars = {
        "GOOS": theOs,
        "GOARCH": theArch,
        "CGO_ENABLED": cgoEnable(theOs, theArch)
    }
    print(f"[green]Building {theOs}-{theArch}")
    binaryName = os.path.join(cwd, 'out', f"sl-{theOs}-{theArch}")
    if theOs == "windows":
        binaryName = f"{binaryName}.exe"
    buildCmd = f"go build -ldflags \"-s -w\" -o {binaryName} ."
    print(f"[yellow]Running {buildCmd}")
    subprocess.run(buildCmd)


### Section for i18n4go stuff ###
def runI18n4go(path: str) -> None:
    """Runs the i18n4go program, and fixes the missing entries

    :param str path: Base path of the repo
    """
    plugin_dir = os.path.join(path, 'plugin')
    binary = os.path.join(path, 'bin', 'i18n4go')
    # TODO: Support linux too I guess.
    if not isWindows():
        binary = f"{binary}_mac"
    cmd = f"{binary} -c checkup -q i18n -v -d {plugin_dir}"
    print(f"[yellow]Running: {cmd}")
    result = subprocess.run(cmd, capture_output=True, text=True)
    # We have some mismatching lines, lets fix that.
    missmatch = ""
    # These strings need to be added
    # The ? at the end of `(.+?)` makes the search NON-greedy, which is required for this regex to work
    # https://docs.python.org/3/library/re.html#regular-expression-syntax
    add_re = re.compile(r"^\"(.+?)\" exists in the code, but not in en_US$", flags=re.M | re.DOTALL)
    add_json = {}
    # These strings need to be removed
    del_re = re.compile(r"^\"(.+?)\" exists in en_US, but not in the code$", flags=re.M | re.DOTALL)
    del_json = {}
    # These strings need to remain missing until they are translated.
    unt_re = re.compile(r"^\"(.+?)\" exists in en_US, but not in ([a-zA-Z_]{5,7})$", flags=re.M | re.DOTALL)
    # These also need to be removed
    mxm_re = re.compile(r"^\"(.+?)\" exists in ([a-zA-Z_]{5,7}), but not in en_US$", flags=re.M | re.DOTALL)
    # This should be the last line
    fin_re = re.compile(r"Could not checkup, err: Strings don't match$")

    if result.returncode > 0:
        add_results = add_re.findall(result.stdout)
        print("[yellow] ====== ADD ========")
        for line in add_results:
            print(f"[yellow]|{line}|")
            # We use the whole line as the key to make search easier.
            add_json[line] = {"id": line, "translation": line}
        del_results = del_re.findall(result.stdout)
        print("[red] ====== REMOVE ========")
        for line in del_results:
            print(f"[red]|{line}|")
            del_json[line]= {"id": line, "translation": line}
        
        print("[blue] ====== MISMATCH (REMOVE PART 2) ========")
        mxm_results = mxm_re.findall(result.stdout)
        for line in mxm_results:
            print(f"[blue]|{line[0]}| is missing from {line[1]}")
            del_json[line[0]] = {"id": line[0], "translation": line[0]}
        # We only want to ADD things to en_US so that our translators have an easier time figuring out what they need
        # to add to the other languages.
        en_us = os.path.join(plugin_dir, 'i18n', 'resources', 'en_US.all.json')
        add_i18n(en_us, add_json)
        for f in i18n_files:
            i18n_path = os.path.join(plugin_dir, 'i18n', 'resources', f)
            del_i18n(i18n_path, del_json)
    else:
        print(f"[green]I18N files are ok!")

def add_i18n(file_name: str, updates: dict) -> None:
    """Adds the new translations

    :param str file_name: file_name we are going to add translations to
    :param dict new_i18n: A dictionary, the key of each entry is the 'id' of the line to be added
    """
    # The original data
    print(f"Adding to {file_name}...")
    source_i18n = get_source_data(file_name)

    # This ensures we don't accidently update a translation file if there is already a translation for this key
    for key in updates.keys():
        if key not in source_i18n:
            source_i18n[key] = updates[key]

    write_source_data(file_name, source_i18n)


def del_i18n(file_name: str, updates: dict) -> None:
    """Removes unneeded translations

    :param str file_name: file_name we are going to remove translations froms
    :param dict new_i18n: A dictionary, the key of each entry is the 'id' of the line to be added
    """
    # The original data
    print(f"Removing from {file_name}...")
    source_i18n = get_source_data(file_name)

    # Remove any matches we find
    for key in updates.keys():
        if key in source_i18n:
            source_i18n.pop(key, None)

    write_source_data(file_name, source_i18n)


def get_source_data(file_name: str) -> dict:
    """Reads from the i18n files and returns a formatted dict"""
    source_i18n = {}
    with open(file_name, encoding="utf8") as f:
        data = json.load(f)
        for i in data:
            source_i18n[i.get('id')] = i
    f.close()
    return source_i18n

def write_source_data(file_name: str, updates: dict) -> None:
    """writes updates to the file_name, converts to a list and sorts first"""
    updated_i18n = []
    for i in sorted(updates.keys()):
        updated_i18n.append(updates[i])

    # Write out the new updates
    with open(file_name, 'w', encoding="utf8", newline='\n') as f:
       json.dump(updated_i18n, f, sort_keys=True, separators=(',', ': ',), ensure_ascii=False, indent=2)
    f.close()

### END i18n4go stuff ###

class Builder(object):
    def __init__(self):
        self.cwd = os.getcwd()
        if not self.cwd.endswith('softlayer-cli'):
            raise Exception(f"Working Directory should be softlayer-cli, is currently {self.cwd}")

    def getdir(self):
        return self.cwd

@click.group()
@click.pass_context
def cli(ctx):
    ctx.obj = Builder()
    pass

@cli.command()
@click.argument("version")
@click.pass_context
def build(ctx, version):
    """Builds the SL binaries"""
    # genBinData()
    toBuild = buildArchs()
    for os in toBuild.keys():
        for arch in toBuild[os]:
            goBuild(ctx.obj.getdir(), os, arch)

    

@cli.command()
@click.argument("version")
@click.pass_context
def deploy(ctx, version):
    """Deploys the SL binaries"""
    click.echo("Deploying...")


@cli.command()
@click.argument("version")
@click.pass_context
def release(ctx, version):
    """Builds, then deploys the release"""
    click.echo("Performing a Release ...")

@cli.command()
@click.pass_context
def i18n(ctx):
    """Checks and builds the i18n files"""
    click.echo("I18N STUFF ...")
    runI18n4go(ctx.obj.getdir())

if __name__ == '__main__':
    cli()


