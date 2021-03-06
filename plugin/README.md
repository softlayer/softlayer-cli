

## Adding new actions to slplugin

> Terminology:
> `ibmcloud sl <COMMAND> <ACTION>`
> *COMMAND*: is a collection of actions here.
> *ACTION*: What part of the command you are running.

1. Add command metadata func in `bluemix-cli\bluemix\slplugin\metadata\<COMMAND>.go`
- Add command metadata func

```
func BlockVolumeLimitsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_BLOCK_NAME,
		Name:        CMD_BLK_VOLUME_LIMITS_NAME,
		Description: T("Lists the storage limits per datacenter for this account."),
		Usage: T(`${COMMAND_NAME} sl block volume-limits [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl block volume-limits
	This command lists the storage limits per datacenter for this account.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}
```

- Add the command metadata to the `BlockMetaData()` top level func

```

func BlockMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_BLOCK_NAME,
		Description: T("Classic infrastructure Block Storage"),
		Usage:       "${COMMAND_NAME} sl block",
		Subcommands: []cli.Command{
...
			BlockVolumeLimitsMetaData(),
...
		},
	}
}
```

- Add the `BlockMetaData()` to the func `getCLITopCommands()` in `bluemix-cli\bluemix\slplugin\softlayer_plugin.go`
2. Add action mapping in `bluemix-cli\bluemix\slplugin\actions.go`
3. Add actual CLI command code in `bluemix-cli\bluemix\slplugin\commands\<COMMAND>\<ACTION>.go`
    - `runtime error: invalid memory address or nil pointer dereference` means you forgot this step.
4. Add any manager code in `bluemix-cli\bluemix\slplugin\managers\`



## i18n stuff

anything with `T("some string here")` uses the internationalization system. Definitions are located in `bluemix-cli\bluemix\i18n\en_US.all.json` for english.

[i18n4go](https://github.com/maximilien/i18n4go) is used to make sure all strings being transalted have translations. To test run this command

Your working directory should be in `go/src/github.ibm.com/Bluemix/bluemix-cli/`

```

# This command will check if there were any mismatch between `en_US.all.json` and the string with `T("some string here")` in your code. It will output details about the mismatch, fix these mismatch manually.
$ ./bin/catch-i18n-mismatch.sh  
OKTotal time: 372.753966ms



# This command will format en_US.all.json and other language json file.
$ ./bin/format-translation-files 
*** Process zh_Hans.all.json
*** Process ko_KR.all.json
*** Process es_ES.all.json
*** Process en_US.all.json
*** Process ja_JP.all.json
*** Process zh_Hant.all.json
*** Process it_IT.all.json
*** Process fr_FR.all.json
*** Process pt_BR.all.json
*** Process de_DE.all.json

# This command will generate/update i18n_resources.go file

$ ./bin/generate-i18n-resources 
Generating i18n resource file ...
Done.
```


## Unit Tests




### Running Tests

To actually run the tests, do `go test <PACKAGE>`. Use `-coverprofile=coverage.out` to produce a coverage.out file that you can then use to figure out what lines are missing coverage.

```
$> go test -coverprofile=coverage.out github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file
ok      github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file       1.225s
```

Coverage report

For basic information
```
go tool cover -func=coverage.out
```

Detailed HTML output

```
go tool cover -html=coverage.out
```

Specific Tests

```
go test -v -coverprofile=coverage.out github.ibm.com/SoftLayer/softlayer-cli/plugin/managers -ginkgo.focus Issues3190
```
### Fake Managers

CLI calls to manager functions need an entry in `bluemix-cli\bluemix\slplugin\testhelpers\fake_manager.go`


Managers have a fake/test interface that is autogenerate with a program called [couterfieter](https://github.com/maxbrunsfeld/counterfeiter)

```
# From /github.ibm.com/Bluemix/bluemix-cli
cd bluemix/slplugin/managers
counterfeiter.exe -o ../testhelpers/fake_storage_manager.go . StorageManager
```


### `[no tests to run]`
New commands needs a `command_test.go` file in the CLI directory.

If you added `slplugin/commands/new/` then there needs to be a `slplugin/commands/new/new_test.go` file. Copy the content from one of the other command test files and just change the name and package.

# Vendor

Vendor files are now managed by `go mod vendor`, I had to set these environment variables to download github.ibm.com vendor objects. To update the github.com/softlayer/softlayer-go dependancy, update `go.mod` file.

```bash
export GOPROXY=direct
export GOSUMDB=off
go mod vendor
```

If you get this error, turning off GOSUMDB and setting GOPROXY=direct seemed to work.
```
$ go mod vendor
go: github.ibm.com/Bluemix/cf-admin-cli@v0.0.0-20200515160705-accb00409d86: verifying go.mod: github.ibm.com/Bluemix/cf-admin-cli@v0.0.0-20200515160705-accb00409d86/go.mod: reading https://sum.golang.org/lookup/github.ibm.com/!bluemix/cf-admin-cli@v0.0.0-20200515160705-accb00409d86: 410 Gone
        server response:
        not found: github.ibm.com/Bluemix/cf-admin-cli@v0.0.0-20200515160705-accb00409d86: invalid version: git fetch -f origin refs/heads/*:refs/heads/* refs/tags/*:refs/tags/* in /tmp/gopath/pkg/mod/cache/vcs/7c3b4597f53c7708d8d63068430570d5325f6ceef4fb0e2076cc6c593df4c01a: exit status 128:
                fatal: could not read Username for 'https://github.ibm.com': terminal prompts disabled

```



## Detect Secrets
Uses https://github.ibm.com/Whitewater/whitewater-detect-secrets NOT the normal yelp/detect-secrets