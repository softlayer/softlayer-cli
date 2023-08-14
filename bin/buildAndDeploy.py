#!python

import click
import os
import subprocess
from rich import print

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

def goBuild(theOs: str, theArch: str) -> None:
    cwd = os.getcwd()
    print(f"[yellow]CWD: {cwd}")
    if not cwd.endswith('softlayer-cli'):
        raise Exception(f"Working Directory should be githubio_source, is currently {cwd}")
    envVars = {
        "GOOS": theOs,
        "GOARCH": theArch,
        "CGO_ENABLED": cgoEnable(theOs, theArch)
    }
    print(f"[green]Building {theOs}-{theArch}")
    binaryName = f"{cwd}/out/sl-{theOs}-{theArch}"
    if theOs == "windows":
        binaryName = f"{binaryName}.exe"
    buildCmd = f"go build -ldflags '-s -w' -o {binaryName} ."
    subprocess.run()


@click.group()
@click.pass_context
def cli(ctx):
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
            goBuild(os, arch)

    

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

if __name__ == '__main__':
    cli()


