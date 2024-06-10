#!python

import click
import json
from pathlib import Path
import os
import sys
import re
import subprocess
import requests
import platform
import hashlib
import glob
from rich import print
from rich.markup import escape


i18n_files = [
'en_US.json',
'de_DE.json',
'es_ES.json',
'fr_FR.json',
'it_IT.json',
'ja_JP.json',
'ko_KR.json',
'pt_BR.json',
'zh_Hans.json',
'zh_Hant.json',
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
        'linux': ['amd64', '386', 'arm64', 'ppc64le', 's390x'],
        'windows': ['386', 'amd64'],
    }
    return buildArchs


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

    ## GO GENERATE
    go_generate = ['go', 'generate', './...']
    # Not using the 'real' command here because this looks neater.

    print(f'[turquoise2]Running: go generate ./...')
    try:
        subprocess.run(go_generate, check=True)
    except subprocess.CalledProcessError as e:
        print(f"[red]>>> Go Generate failed <<<")
        sys.exit(e.returncode)

    ## GO VET
    go_vet = ['go', 'vet'] +  clean_mods
    # Not using the 'real' command here because this looks neater.

    print(f'[turquoise2]Running: go vet $(go list ./... | grep -v "fixtures" | grep -v "vendor")')
    try:
        subprocess.run(go_vet, check=True)
    except subprocess.CalledProcessError as e:
        print(f"[red]>>> Go Vet failed <<<")
        sys.exit(e.returncode)



    ## GO TEST
    go_test = ['go', 'test'] +  clean_mods
    print(f'[turquoise2]Running: go test $(go list ./... | grep -v "fixtures" | grep -v "vendor")')
    try:
        subprocess.run(go_test, check=True)
    except subprocess.CalledProcessError as e:
        print(f"[red]>>> Go Test failed <<<")
        sys.exit(e.returncode)
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
    binary = os.path.join(path, 'bin', 'goi18n2')
    if platform.system() == 'Windows':
        binary = f"{binary}.exe"
    elif platform.system() == "Darwin":
        binary = f"{binary}_mac"
    cmd = [binary, "extract", "-outdir=plugin/i18n/v2Resources/", "-format=json",
           "-sourceLanguage=en_US", plugin_dir]
    # os.chdir(os.path.join(path, 'plugin'))
    print("[turquoise2]Running: "  + " ".join(cmd))
    result = subprocess.run(cmd, capture_output=True, text=True)
    # os.chdir(path)

    if result.stderr:
        print(f"[red]Error: {result.stderr}")
    else:
        print(f"\t[green]Generated en-US translation file")


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

        old_v = re.search(r'^\W+PLUGIN_VERSION\W+= \"([0-9]+\.[0-9]+\.[0-9]+(-[a-z]+)?)\"', data, re.M)
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
        os.environ["IBMCLOUD_TRACE"] = "false"
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
            'Linux_Ppc64le': {'file': f"sl-{self.version}-linux-ppc64le", 'checksum': ''},
            'Linux_s390x': {'file': f"sl-{self.version}-linux-s390x", 'checksum': ''},
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
        # Set this to true if you want to use the Refresh-Plugin-Version-on-YS1 job
        refresh = False
        # This create a new version
        jenkinsUrl = 'https://wcp-cloud-foundry-jenkins.swg-devops.com/job/Publish%20Plugin%20to%20YS1'
        
        if refresh:
            # This updates an existing version
            jenkinsUrl = 'https://wcp-cloud-foundry-jenkins.swg-devops.com/job/Refresh-Plugin-Version-on-YS1'
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
            # This one doesn't follow the normal pattern so it needs a special case.
            if x == "Linux_s390x" and not refresh:
                checkData = {'name': f"Checksum_Linux_S390x", "value": checksums[x]['checksum']}
            else:
                checkData = {'name': f"Checksum_{x}", "value": checksums[x]['checksum']}
            form_json['parameter'].append(urlData)
            form_json['parameter'].append(checkData)

        print(f"[yellow]Trying to start jenkins job on {jenkinsUrl}")
        print(form_json)
        result = requests.post(f"{jenkinsUrl}/build",  auth=('cgallo@us.ibm.com', jenkins_token), data={'json':json.dumps(form_json)})
        if (result.status_code == 201 ):
            print(f"[green] Created Job! Check {jenkinsUrl}" )
        else:
            print(f"[red]Error: {result.status_code} {result.reason}")
            print(f"[yellow] {result.text}")
            print(f"[yellow] {result.url}")
            print(f"[yellow] {result.request}")
            raise Exception("Error in runJenkins()")



    def goBuild(self, theOs: str, theArch: str) -> None:
        """Runs the go build command
        
        :param str cwd: The current working directory
        :param str theOs: OS to build for
        :param str theArch: Architecture to build for
        """
        cgo_enabled = 0
        os.environ["GOOS"] = theOs
        os.environ["GOARCH"] = theArch
        os.environ["CGO_ENABLED"] = str(cgo_enabled)

        print(f"[green]Building {theOs}-{theArch}")
        binaryName = os.path.join(self.cwd, 'out', f"sl-{self.version}-{theOs}-{theArch}")
        if theOs == "windows":
            binaryName = f"{binaryName}.exe"
        buildCmd = f" go build -ldflags \"-s -w\" -o {binaryName} ."
        print(f"[turquoise2]Running GOOS={theOs} GOARCH={theArch} CGO_ENABLED={cgo_enabled} {buildCmd}")
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
    # genBinData()

if __name__ == '__main__':
    cli()
    # try:
    #     cli()
    # except Exception as e:
    #     print(f"[red]{e}")

