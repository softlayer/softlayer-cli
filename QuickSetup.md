Setting up golang env

0: Make sure you have your SSH keys setup on github.ibm.com https://github.ibm.com/settings/keys

1: Know your GOPATH: https://go.dev/doc/gopath_code#GOPATH

```
$> echo $GOPATH
C:\Users\allmi\go
$> cd ~/go/
```

2: Get the softlayer-cli code  https://github.ibm.com/softlayer/softlayer-cli
Read the README as well.

```
$> mkdir -p src/github.ibm.com/softlayer
$> cd src/github.ibm.com/softlayer
$> pwd
/c/Users/allmi/go/src/github.ibm.com/softlayer
$> git clone https://github.ibm.com/softlayer/softlayer-cli
```

At this point you have the code, but need to setup the required libraries

3: run go mod vendor  https://go.dev/ref/mod#go-mod-vendor This command will download everything listed in `go.mod`, which are the libraries this codebase uses.

```
$> go mod vendor
go: downloading github.com/IBM-Cloud/ibm-cloud-cli-sdk v1.0.1
go: downloading github.com/softlayer/softlayer-go v1.1.2
go: downloading github.com/miekg/dns v1.1.50
go: downloading github.com/onsi/ginkgo v1.16.2
go: downloading github.com/onsi/gomega v1.11.0
go: downloading github.com/spf13/cobra v1.5.0
go: downloading github.com/spf13/pflag v1.0.5
go: downloading github.com/Xuanwo/go-locale v1.1.0
go: downloading github.com/nicksnyder/go-i18n v1.10.1
go: downloading golang.org/x/text v0.7.0
go: downloading github.com/stretchr/testify v1.7.0
go: downloading golang.org/x/net v0.7.0
go: downloading golang.org/x/sys v0.5.0
go: downloading golang.org/x/tools v0.1.12
go: downloading github.com/fatih/color v1.10.0
go: downloading github.com/mattn/go-colorable v0.1.8
go: downloading github.com/mattn/go-runewidth v0.0.12
go: downloading golang.org/x/crypto v0.6.0
go: downloading github.com/fatih/structs v1.1.0
go: downloading github.com/softlayer/xmlrpc v0.0.0-20200409220501-5f089df7cb7e
go: downloading github.com/inconshreveable/mousetrap v1.0.0
go: downloading github.com/davecgh/go-spew v1.1.1
go: downloading github.com/pmezard/go-difflib v1.0.0
go: downloading gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
go: downloading gopkg.in/yaml.v2 v2.4.0
go: downloading github.com/nicksnyder/go-i18n/v2 v2.2.0
go: downloading github.com/mattn/go-isatty v0.0.12
go: downloading github.com/rivo/uniseg v0.1.0
go: downloading github.com/nxadm/tail v1.4.8
go: downloading github.com/pelletier/go-toml v1.2.0
go: downloading golang.org/x/term v0.5.0
go: downloading gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
go: downloading github.com/fsnotify/fsnotify v1.4.9
go: downloading golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
```

4: You can now build the softlayer-cli: https://go.dev/ref/mod#build-commands

```
$> go build
$> ./softlayer-cli.exe --help
Manage Classic infrastructure services

Usage:
  sl [command]
```

At this point you can test the softlayer-cli plugin independent of the `ibmcloud` command. If you want to test it all together, you need to ALSO build the `ibmcloud` binary.

5: Get the ibmcloud code https://github.ibm.com/ibmcloud-cli/bluemix-cli

```
$> mkdir -p $GOPATH/src/github.ibm.com/ibmcloud-cli
$> cd $GOPATH/src/github.ibm.com/ibmcloud-cli
$> git@github.ibm.com:ibmcloud-cli/bluemix-cli.git
```

6: Edit the ibmcloud-cli/go.mod file to force it to use the copy of softlayer-cli on your computer
Add `github.ibm.com/SoftLayer/softlayer-cli => ../../SoftLayer/softlayer-cli` to the `replace` section in `github.ibm.com/ibmcloud-cli/bluemix-cli/go.mod`

```
module github.ibm.com/ibmcloud-cli/bluemix-cli

go 1.20

replace (
        github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
        github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20191121165722-d1d5f6476656+incompatible
        github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.8
        github.ibm.com/SoftLayer/softlayer-cli => ../../SoftLayer/softlayer-cli
)
```

7: vendor and build

You'll need to change some git settings first to access private repositories. See https://github.ibm.com/softlayer/softlayer-cli#vendor

```
# I usually add these to ~/.bashrc so they are always set
$> export GOPROXY=direct
$> export GOPRIVATE=github.ibm.com/*
# Make sure you gitconfig has these lines
$> git config --global url."ssh://git@github.ibm.com/".insteadof https://github.ibm.com/
$> go mod vendor
```

Need to run `go mod tidy` when you make changes to go.mod
```
$> go mod tidy
go: downloading gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
go: downloading github.com/onsi/ginkgo/v2 v2.1.6
go: downloading github.com/jarcoal/httpmock v1.0.5
go: downloading github.com/smartystreets/goconvey v1.6.7
go: downloading github.com/BurntSushi/toml v1.1.0
go: downloading golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
go: downloading github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e
go: downloading sigs.k8s.io/yaml v1.2.0
go: downloading github.com/jtolds/gls v4.20.0+incompatible
go: downloading github.com/smartystreets/assertions v1.0.0
go: downloading github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1
go: downloading github.com/kr/pretty v0.3.0
go: downloading github.com/kr/text v0.2.0
go: downloading github.com/rogpeppe/go-internal v1.6.1
$> go mod vendor
go: downloading github.ibm.com/arf/cli-dev-plugin v1.3.4-0.20230220210426-acb097801062
go: downloading github.ibm.com/SoftLayer/softlayer-cli v1.4.1
go: downloading github.ibm.com/Bluemix/resource-catalog-cli v0.0.0-20220906182229-845aab607438
go: downloading github.com/pelletier/go-toml v1.9.5
go: downloading github.com/parnurzeal/gorequest v0.2.16
go: downloading github.com/briandowns/spinner v0.0.0-20190311160019-998b3556fb3f
go: downloading k8s.io/api v0.25.3
go: downloading moul.io/http2curl v1.0.0
go: downloading github.com/miekg/dns v1.1.25
go: downloading k8s.io/apimachinery v0.25.3
go: downloading github.com/gogo/protobuf v1.3.2
go: downloading k8s.io/klog/v2 v2.80.0
go: downloading k8s.io/utils v0.0.0-20220823124924-e9cbc92d1a73
go: downloading github.com/google/gofuzz v1.2.0
go: downloading gopkg.in/inf.v0 v0.9.1
go: downloading sigs.k8s.io/structured-merge-diff/v4 v4.2.3
go: downloading sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2
go: downloading github.com/go-logr/logr v1.2.3
go: downloading github.com/json-iterator/go v1.1.12
go: downloading github.com/modern-go/reflect2 v1.0.2
go: downloading github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd

$> go build
$> ./bluemix-cli
```

8: Login to use it. I like to use an API key. Use https://cloud.ibm.com/iam/apikeys to create a "My IBM Cloud API Keys". It will be 32 characters long, and different from a classic infrastructure api key.


```
$> ibmcloud.exe login --apikey <keyhere>
$> ibmcloud sl vs list (or some other command)
```



### Uploading binaries for testing

Sometimes I'll upload binaries build with local softlayer-cli changes, mostly for the translation team to test with. That process is as follows:


URL: https://s3.us-east.cloud-object-storage.appdomain.cloud/softlayer-cli-binaries/index.html
CLI Docs: https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-cli-plugin-ic-cos-cli

COS Plugin install and setup
```
$> ibmcloud plugin install cloud-object-storage
$> ibmcloud cos config region us-east
$> ibmcloud cos config crn --crn 0859c995-bf8b-46fe-9d8b-0f3405e25359
# be in the bluemix-cli directory, make sure you have the replace directive in go.mod and its ready to be build
$> cd $GOPATH/src/github.ibm.com/ibmcloud-cli/bluemix-cli
$> ./bin/build-all
$> for i in `\ls out/`; do echo UPLOADING $i; ibmcloud.exe cos object-put --bucket=softlayer-cli-binaries --key=$i  --body=./out/$i; done
```
Hopefully that will work.