#!python

import click
import json
from pathlib import Path
import os
import re
import subprocess
import requests
import platform
import hashlib
import glob
from rich import print
from rich.markup import escape


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



def buildArchs() -> dict:
    """Returns the list of binaries we should build"""
    buildArchs =  {
        'darwin': ['arm64', 'amd64'],
        'linux': ['amd64', '386', 'arm64'],
        'windows': ['386', 'amd64'],
    }
    return buildArchs

def cgoEnable(theOs: str, theArch: str) -> int:
    """Disable cgo for these archs"""
    cgo_disable = ["amd64", "386", "arm64"]
    if theOs == "linux" and theArch in cgo_disable:
        return 0
    return 1



def runTests() -> None:
    """Runs unit tests"""

    go_ven = ['go', 'mod', 'vendor']
    print("[turquoise2]Running: " + " ".join(go_ven))
    subprocess.run(go_ven, check=True)
    # We can't use the `| grep -v "fixtures" | grep -v "vendor"` stuff because working with pipes in 
    # subprocess is tricky, doing this is easier to me.
    go_mods = ['go', 'list', './...']
    print(f"[turquoise2]Running: " + " ".join(go_mods))
    mods = subprocess.run(go_mods, capture_output=True, check=True, text=True)
    clean_mods = []
    for mod in mods.stdout.split("\n"):
        if re.match(r"fixtrues|vendor", mod) is None:
            clean_mods.append(mod)

    go_vet = ['go', 'vet'] +  clean_mods
    # Not using the 'real' command here because this looks neater.
    print(f'[turquoise2]Running: go vet $(go list ./... | grep -v "fixtures" | grep -v "vendor")')
    subprocess.run(go_vet, check=True)
    go_test = ['go', 'test'] +  clean_mods
    print(f'[turquoise2]Running: go test $(go list ./... | grep -v "fixtures" | grep -v "vendor")')
    subprocess.run(go_test, check=True)
    go_sec = ['gosec', '-exclude-dir=fixture', '-exclude-dir=plugin/resources', '-quiet', './...']
    # Not using the 'real' command because this is more copy/pasteable.
    print('[turquoise2]Running: ' + " ".join(go_sec)) 
    try:
        subprocess.run(go_sec, check=True)
    except FileNotFoundError:
        gosec_instal = "curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $GOPATH/bin"
        print(f"[red]gosec not found. Try running:\n{gosec_instal}")
    



### Section for i18n4go stuff ###
def runI18n4go(path: str) -> None:
    """Runs the i18n4go program, and fixes the missing entries

    :param str path: Base path of the repo
    """
    plugin_dir = os.path.join(path, 'plugin')
    binary = os.path.join(path, 'bin', 'i18n4go')
    # TODO: Support linux too I guess.
    if platform.system() == 'Windows':
        binary = f"{binary}.exe"
    elif platform.system() == "Darwin":
        binary = f"{binary}_mac"
    cmd = [binary, "-c=checkup", "-q=i18n", "-v", f"-d={plugin_dir}"]
    os.chdir(os.path.join(path, 'plugin'))
    print("[turquoise2]Running: "  + " ".join(cmd))
    result = subprocess.run(cmd, capture_output=True, text=True)
    os.chdir(path)
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

    if result.stderr:
        print(f"[red]Error: {result.stderr}")

    if result.returncode > 0:
        print("\t[yellow] ====== ADD ======== ")
        add_results = add_re.findall(result.stdout)
        for line in add_results:
            print(f"\t[yellow]|{escape(line)}|")
            # We use the whole line as the key to make search easier.
            add_json[line] = {"id": line, "translation": line}
        del_results = del_re.findall(result.stdout)
        print("\t[red] ====== REMOVE ========")
        for line in del_results:
            print(f"\t[red]|{escape(line)}|")
            del_json[line]= {"id": line, "translation": line}
        
        print("\t[blue] ====== MISMATCH (REMOVE PART 2) ========")
        mxm_results = mxm_re.findall(result.stdout)
        for line in mxm_results:
            print(f"\t[blue]|{escape(line[0])}| is missing from {line[1]}")
            del_json[line[0]] = {"id": line[0], "translation": line[0]}
        # We only want to ADD things to en_US so that our translators have an easier time figuring out what they need
        # to add to the other languages.
        en_us = os.path.join(plugin_dir, 'i18n', 'resources', 'en_US.all.json')
        add_i18n(en_us, add_json)
        for f in i18n_files:
            i18n_path = os.path.join(plugin_dir, 'i18n', 'resources', f)
            del_i18n(i18n_path, del_json)

    else:
        print(f"\t[green]No Changes Needed!")

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

def genBinData() -> None:
    """Generates the I18N Binary data required for translations"""
    goBindata = './bin/go-bindata'
    if isWindows():
        goBindata = f'{goBindata}.exe'
    elif platform.system() == "Darwin":
        goBindata = f'{goBindata}_mac'
    goBindata = [goBindata, "-pkg=resources", "-o=plugin/resources/i18n_resources.go",  "plugin/i18n/resources"]
    print("[turquoise2]Building I18N: " + " ".join(goBindata))
    result = subprocess.run(goBindata, capture_output=True)
    if result.returncode > 0:
        print(f"\t[red]{result.stderr.decode('ascii')}")
    else:
        print(f"\t[green]OK!")


### END i18n4go stuff ###

class Builder(object):
    def __init__(self):
        self.cwd = os.getcwd()
        self.version = '0.0.1'
        self.debug = True
        self.cnd_url = 'https://s3.us-east.cloud-object-storage.appdomain.cloud/softlayer-cli-binaries/'
        if not self.cwd.endswith('softlayer-cli'):
            raise Exception(f"Working Directory should be softlayer-cli, is currently {self.cwd}")

    def getdir(self) -> str:
        return self.cwd

    def setVersion(self, version: str) -> None:
        if not re.match(r'\d+\.\d+\.\d+', version):
            raise Exception(f"{version} is not valid, needs to be in the format of Major.Minor.Revision")
        self.version = version
        sl_go = Path(os.path.join(self.cwd, 'plugin', 'metadata', 'sl.go'))
        
        data = sl_go.read_text()

        old_v = re.search(r'^\W+PLUGIN_VERSION\W+= \"([0-9]+\.[0-9]+\.[0-9]+)\"', data, re.M)
        if old_v is None:
            raise Exception(f"[red]Can't find old version!")
        print(f"[turquoise2]Old Version: {old_v[1]}")
        updated = re.sub(
            r'PLUGIN_VERSION\W+= \"([0-9]+\.[0-9]+\.[0-9]+)\"',
            f"PLUGIN_VERSION         = \"{self.version}\"",
            data)
        # print(updated)
        sl_go.write_text(updated)   
        print(f"[turquoise2]Updated {sl_go} PLUGIN_VERSION {old_v[1]} -> {version}") 


    def deploy(self):
        """Uploads binaries to IBM COS"""
        apikey = os.getenv("IBMCLOUD_APIKEY")
        # if IBMCLOUD_TRACE is true the upload will print out the binary file data to the screen.
        os.environ["IBMCLOUD_TRACE"] = False
        if not apikey:
            raise Exception("IBMCLOUD_APIKEY needs to be set to the proper API key first.")
        login_cmd = ["ibmcloud", "login", f"--apikey={apikey}"]
        print(f"[yellow]Running: ibmcloud login --apikey $IBMCLOUD_APIKEY")
        subprocess.run(login_cmd)
        files = glob.glob(os.path.join(self.cwd, 'out', f"sl-{self.version}-*"))
        for f in files:
            upload_cmd = ["ibmcloud", "cos", "upload", "--bucket=softlayer-cli-binaries",
                          f"--file={f}", f"--key={os.path.basename(f)}"]
            print(f"[yellow]Running: " + " ".join(upload_cmd))
            subprocess.run(upload_cmd)

    def getChecksums(self) -> dict:
        """Calcs checksums for all files we generated"""
        checksums = {
            'Linux_X86': {'file': f"sl-{self.version}-linux-386", 'checksum': ''},
            'Linux_X64': {'file': f"sl-{self.version}-linux-amd64", 'checksum': ''},
            'Linux_arm64' : {'file': f"sl-{self.version}-linux-arm64", 'checksum': ''},
            'MacOS': {'file': f"sl-{self.version}-darwin-amd64", 'checksum': ''},
            'MacOS_arm64': {'file': f"sl-{self.version}-darwin-arm64", 'checksum': ''},
            'Win_X86': {'file': f"sl-{self.version}-windows-386.exe", 'checksum': ''},
            'Win_X64': {'file': f"sl-{self.version}-windows-amd64.exe", 'checksum': ''}
        }

        for x in checksums.keys():
            with open(os.path.join(self.cwd, 'out', checksums[x]['file']), "rb") as f:
                digest = hashlib.file_digest(f, "sha1")
                print(f"[yellow]{checksums[x]['file']} => {digest.hexdigest()}")
                checksums[x]['checksum'] = digest.hexdigest()
        return checksums

    def runJenkins(self):
        """Starts up a Jenkins build.
        Copied some of the logic from https://github.ibm.com/coligo/cli/blob/main/script/publish.to.repo.sh so this might not be the best solution.

This CURL command works (missing all the OS options but you get the idea). We should be using 'buildWithParameters' I think, but I couldn't get it to work.
curl -X POST https://wcp-cloud-foundry-jenkins.swg-devops.com/job/Publish%20Plugin%20to%20YS1/build \
--user $JENKINS_USER:$JENKINS_TOKEN \
--form json='{"parameter": [{"name":"Plugin_Name", "value":"sl"},'\
'{"name":"Version", "value": "1.4.2"},'\
'{"name":"Description", "value": "Manage Classic infrastructure services"},'\
'{"name":"min_cli_version", "value": "2.18.0"},'\
'{"name":"private_endpoint_supported", "value": true},'\
'{"name":"Checksum_Linux_X86", "value":"d069b943532f69feadd292d6d1c62eece8ca1112"},'\
'{"name":"Url_Linux_X86", "value":"https://s3.us-east.cloud-object-storage.appdomain.cloud/softlayer-cli-binaries/sl-1.4.2-linux-386"}]}' -v
        """
        checksums = self.getChecksums()
        jenkinsUrl = 'https://wcp-cloud-foundry-jenkins.swg-devops.com/job/Publish%20Plugin%20to%20YS1'
        jenkins_token = os.getenv('JENKINS_TOKEN')
        if not jenkins_token:
            raise Exception("JENKINS_TOKEN is not set to an API key")
        auth = f"cgallo@us.ibm.com:{jenkins_token}"
        form_json = {
            "parameter": [
                {"name": "Plugin_Name", "value":"sl"},
                {"name": "Version", "value":self.version},
                {"name": "Description", "value":"Manage Classic infrastructure services"},
                {"name": "min_cli_version", "value":"2.18.0"},
                {"name": "private_endpoint_supported", "value":True}
            ]
        }

        for x in checksums.keys():
            urlData = {'name': f"Url_{x}", "value": self.cnd_url + checksums[x]['file']}
            checkData = {'name': f"Checksum_{x}", "value": checksums[x]['checksum']}
            form_json['parameter'].append(urlData)
            form_json['parameter'].append(checkData)

        print(f"[yellow]Trying to start jenkins job on {jenkinsUrl}")
        print(form_json)
        result = requests.post(f"{jenkinsUrl}/build",  auth=('cgallo@us.ibm.com', jenkins_token), data={'json':json.dumps(form_json)})
        if (reult.status_code == 201 ):
            print(f"[green] Created Job! Check {jenkinsUrl}" )
        else:
            print(f"[red]Error: {result.status_code} {result.reason}")
            raise Excetion("Error in runJenkins()")



    def goBuild(self, theOs: str, theArch: str) -> None:
        """Runs the go build command
        
        :param str cwd: The current working directory
        :param str theOs: OS to build for
        :param str theArch: Architecture to build for
        """
        os.environ["GOOS"] = theOs
        os.environ["GOARCH"] = theArch
        os.environ["CGO_ENABLED"] = str(cgoEnable(theOs, theArch))

        print(f"[green]Building {theOs}-{theArch}")
        binaryName = os.path.join(self.cwd, 'out', f"sl-{self.version}-{theOs}-{theArch}")
        if theOs == "windows":
            binaryName = f"{binaryName}.exe"
        buildCmd = f"go build -ldflags \"-s -w\" -o {binaryName} ."
        print(f"[turquoise2]Running {buildCmd}")
        # This command basically requires shell=True on mac because -ldflags doesn't get parsed properly withoutit.
        subprocess.run(buildCmd, shell=True)

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
    ctx.obj.setVersion(version)
    toBuild = buildArchs()
    for os in toBuild.keys():
        for arch in toBuild[os]:
            ctx.obj.goBuild(os, arch)

    

@cli.command()
@click.argument("version")
@click.pass_context
def deploy(ctx, version):
    """Deploys the SL binaries"""
    click.echo("Deploying...")
    ctx.obj.setVersion(version)
    ctx.obj.deploy()


@cli.command()
@click.argument("version")
@click.pass_context
def release(ctx, version):
    """Builds, then deploys the release"""
    click.echo("Performing a Release ...")

@cli.command()
@click.argument("version")
@click.pass_context
def jenkins(ctx, version):
    """Trigger a Jenkins build with existing files."""
    ctx.obj.setVersion(version)
    ctx.obj.runJenkins()

@cli.command()
@click.pass_context
def test(ctx):
    """Runs the tests"""
    runTests()

@cli.command()
@click.pass_context
def i18n(ctx):
    """Checks and builds the i18n files"""
    runI18n4go(ctx.obj.getdir())
    genBinData()

if __name__ == '__main__':
    cli()
    # try:
    #     cli()
    # except Exception as e:
    #     print(f"[red]{e}")


