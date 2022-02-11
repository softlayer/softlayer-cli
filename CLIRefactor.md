# The CLI Refactor

I'd like to refactor how the softlayer-cli project is structured to be more like the softlayer-python project, easier to manage and work with. Mainly moving code from single, high traffic files, to files closer to the CLI code itself.

First I will explain the project as a whole, then go over the old pattern, then how the new pattern should look and what steps need to be done to get there. Most of this will just be copy/paste code from one area to another, but that is still quite a lot of work to do.


## Project structure

+ `softlayer-cli/` The main project folder
    * `plugin/` Where all the code goes
        - `client/` Deals with creating the softlayer-go client instance
        - `commands/` All the CLI command code goes here
            + `top level command` A folder for each groups of commands.
                * `action.go` code thet gets run when you do `ibmcloud sl top_level_command action`
                * `action_test.go` tests the action
                * `top_level_command.go` defines the class that all actions will go into
                * `top_level_command_test.go` tests the actions
        - `errors/`  Some error handling
        - `i18n/` translation files, try not to edit these manually, use the bin/fixeverything_i18n.sh script
            + `i18n.go` sets up the i18n resources and translations
            + `resources/*.all.json` contains a json definition of all strings and their translations in a variety of languages
        - `managers/` Used by the CLI to get data from the SL API. Make all actual API calls here (not in the commands)
            + `manager.go` Groups API calls by what area of the API they interact with. 
            + `manager_test.go` Tests the specific manager
        - `metadata/` Contains the CLI definitions (flags, descriptions, etc). I want to move away from these files
        - `resources/` The compiled i18n file, do not edit this manually (run bin/generate-i18n-resources.sh if needed)
        - `testfixtures/` Fake API data for tests. Format is `testfixtures/SOFTLAYER_SERVICE/method.json`
        - `testhelpers/` Fake manager definitions. Uses `couterfieter` to generate them from real managers. Use these in testing the CLI commands
        - `utils/` some utils for dealing with softlayer-go client and other nice to have things
        - `version/` Not really used
        - `actions.go` Responsible for loading all command definitions, needs to be refactored badly.
        - `plugin.go` Sets up the plugin, loads metadata and command definitions.



## The Old Pattern

To illustrate all the files involved in building a CLI command, lets assume you want to make a new command: `ibmcloud sl newGroup newAction`


1. Create a new folder `softlayer-cli/plugin/commands/newGroup`
2. Create a new file `softlayer-cli/plugin/commands/newGroup/newGroup_test.go`. This will ensure ginkogo runs the CLI command tests

```go
package newGroup_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestUser(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "NewGroup Suite")
}

```

3. Create a new file `softlayer-cli/plugin/commands/newGroup/newAction.go`

```go
package newGroup

import (
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
    "github.com/urfave/cli"
    . "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

// Define any managers your command needs here, and pass them into the New<Action> function
type NewActionCommand struct {
    UI             terminal.UI
    NetworkManager managers.NetworkManager
}

// The format here is New<action> which I realize might be hard to read when the action is called NewAction, but what can you do.
func NewNewActionCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *NewActionCommand) {
    return &ListCommand{
        UI:             ui,
        NetworkManager: networkManager,
    }
}

func (cmd *ListCommand) Run(c *cli.Context) error {
    // Do the actual things you need to do here
    return nil
}
```


4. Create a new file `softlayer-cli/plugin/commands/newGroup/newAction_test.go`

```go
package newGroup_test

import (
    "errors"
    "strings"

    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/softlayer/softlayer-go/datatypes"
    "github.com/softlayer/softlayer-go/sl"
    "github.com/urfave/cli"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/newGroup"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("NewGroup NewAction Tests", func() {
    var (
        fakeUI             *terminal.FakeUI
        fakeNetworkManager *testhelpers.FakeNetworkManager
        cmd                *newGroup.NewActionCommand
        cliCommand         cli.Command
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeNetworkManager = new(testhelpers.FakeNetworkManager)
        cmd = newGroup.NewActionCommand(fakeUI, fakeNetworkManager)
        cliCommand = cli.Command{
            Name:        metadata.NewGroupNewActionMetaData().Name,
            Description: metadata.NewGroupNewActionMetaData().Description,
            Usage:       metadata.NewGroupNewActionMetaData().Usage,
            Flags:       metadata.NewGroupNewActionMetaData().Flags,
            Action:      cmd.Run,
        }
    })

    Describe("NewGroup NewAction", func() {
        Context("Run the command normally", func() {
            It("return no error", func() {
                err := testhelpers.RunCommand(cliCommand, "")
                Expect(err).NotTo(HaveOccurred())
                Expect(fakeUI.Outputs()).To(ContainSubstring("something we expected"))
            })
        })
    })
})
```

5. Setup the metadata `softlayer-cli/plugin/metadata/newGroup.go`

```go
package metadata

import (
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
    "github.com/urfave/cli"
    . "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
    NS_NEWGROUP_NAME = "newGroup"
    NS_NEWGROUP_DESC = T("Classic infrastructure newGroup")

    CMD_NEWGROUP_NAME  = "newgroup"
    CMD_NEWGROUP_DESC  = "Classic infrastructure newGroup"
    CMD_NEWGROUP_USAGE = "${COMMAND_NAME} sl newGroup"
    //sl-newgroup
    CMD_NEWGROUP_NEWACTION_NAME    = "newAction"
    CMD_NEWGROUP_NEWACTION_DESC    = T("Create a new newAction")
    CMD_NEWGROUP_NEWACTION_USAGE    = "${COMMAND_NAME} sl newGroup newAction  [OPTIONS]"
    CMD_NEWGROUP_NEWACTION_OPT1      = "type"
    CMD_NEWGROUP_NEWACTION_OPT1_DESC = T("newGroup type  [required]. Options are: vlan,vs,hardware")
)

var NS_NEWGROUP = plugin.Namespace{
    ParentName:  NS_SL_NAME,
    Name:        NS_NEWGROUP_NAME,
    Description: NS_NEWGROUP_DESC,
}

var CMD_FW = cli.Command{
    Category:    NS_SL_NAME,
    Name:        CMD_NEWGROUP_NAME,
    Description: CMD_NEWGROUP_DESC,
    Usage:       CMD_NEWGROUP_USAGE,
    Subcommands: []cli.Command{
        CMD_NEWGROUP_NEWACTION,
    },
}

var CMD_NEWGROUP_NEWACTION = cli.Command{
    Category:    CMD_NEWGROUP_NAME,
    Name:        CMD_NEWGROUP_NEWACTION_NAME,
    Description: CMD_NEWGROUP_NEWACTION_DESC,
    Usage:       CMD_NEWGROUP_NEWACTION_USAGE,
    Flags: []cli.Flag{
        cli.StringFlag{
            Name:  CMD_NEWGROUP_NEWACTION_OPT1,
            Usage: CMD_NEWGROUP_NEWACTION_OPT1_DESC,
        },
        ForceFlag(),
    },
}
```

6. Edit `softlayer-cli/plugin/actions.go`
    + add an import to your new `newGroup` command `"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/newGroup"`
    + (optional) define a new manager if you need it
    + Add each command to the `CommandActionBindings` map, something like this. The Constants will be defined in the metadata files
```go
    NS_NEWGROUP_NAME + "-" + CMD_NEWGROUP_NEWACTION_NAME: func(c *cli.Context) error {
        return newGroup.NewNewActionCommand(ui, dnsManager).Run(c)
},
        ```

7. After adding your command and tests, build the i18n files `sh bin/fixeverything_i18n.sh`
8. Commit your code, push to a branch and submit a pull request

The problem with this pattern is that it uses contants VERY heavily, all over the place, and that makes the code a lot harder to read, and hard to work with.
It also tends to having a few large, high traffic files, like `actions.go` and `metadata/newGroup.go` which generate merge conflicts.
When working on a command, you have to keep scrolling through the metadata to remember what actions/flag/etc your command has.slcli 

## The New Pattern

The purpose of this refactor is to reduce the number of global constants that define actions, as they are only really used in 1-2 places and make development more challenging. To spread out the CLI command definitions into command specific files so that there are less high traffic files that every command needs to modify.

To demonstrate this pattern, lets assume we are making a new command group and action like before, `ibmcloud sl newGroup newAction`.


1. Create a new folder for the new group of commands `softlayer-cli/plugin/commands/newGroup`
2. Create a new file `softlayer-cli/plugin/commands/newGroup/newGroup_test.go`. This will ensure ginkogo runs the CLI command tests. 

```go
package newGroup_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestUser(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "NewGroup Suite")
}

```

3. Create a new file `softlayer-cli/plugin/commands/newGroup/newGroup.go`. This is where you will define the CommandActionBindings, Namespace and Metadata for the newGroup of actions. The metadata for each action will be defined in that actions command file.
Unlike the old pattern, here we should just pass the `session` variable to each command, and the command itself can create managers if needed.

```go
package newGroup

import (
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
    "github.com/softlayer/softlayer-go/session"
    "github.com/urfave/cli"

    . "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
    // Each new command in this group needs to be added here
    CommandActionBindings := map[string]func(c *cli.Context) error{
        "newGroup-newAction": func(c *cli.Context) error {
            return NewNewActionCommand(ui, session).Run(c)
        },
    }
    return CommandActionBindings
}

func NewGroupNamespace() plugin.Namespace {
    return plugin.Namespace{
        ParentName:  "sl",
        Name:        "newGroup",
        Description: T("Classic infrastructure newGroup commands"),
    }
}

func NewGroupMetaData() cli.Command {
    return cli.Command{
        Category:       "sl",
        Name:           "newGroup",
        Description:    T("Classic infrastructure newGroup commands"),
        Usage:          "${COMMAND_NAME} sl newGroup",
        // Each new command in this group needs to be added here
        Subcommands:    []cli.Command{
            NewActionMetaData(),
        },
    }
}
```

4. Create a new file for the action itself. `softlayer-cli/plugin/commands/newGroup/newAction.go`

```go
package newGroup

import (
    "github.com/softlayer/softlayer-go/session"
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
    "github.com/urfave/cli"

    . "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type NewActionCommand struct {
    UI             terminal.UI
    Session        *session.Session
}

func NewNewActionCommand(ui terminal.UI, session *session.Session) (cmd *NewActionCommand) {
    return &NewActionCommand{
        UI:             ui,
        Session:        session,
    }
}

func NewActionMetaData() cli.Command {
    return cli.Command{
        Category: "newGroup",
        Name:     "newAction",
        Description: T("Does some sort of action thing"),
        Usage: T(`${COMMAND_NAME} sl newGroup newAction`),
        Flags: []cli.Flag{
            metadata.OutputFlag(),
        },
    }
}

func (cmd *NewActionCommand) Run(c *cli.Context) error {
    // If you need access to a manager, define it here.
    newGroupManager := managers.NewAccountManager(cmd.Session)
    // The CLI code for the action goes here. 
    return nil
}

```

5. Create a new file for the unit tests `softlayer-cli/plugin/commands/newGroup/newAction_test.go`

```go
package newGroup_test


import (
    "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/softlayer/softlayer-go/session"
    "github.com/urfave/cli"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/newGroup"
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("NewGroup NewAction", func() {
    var (
        fakeUI          *terminal.FakeUI
        cmd             *newGroup.NewActionCommand
        cliCommand      cli.Command
        fakeSession     *session.Session
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        cmd = newGroup.NewNewActionCommand(fakeUI, fakeSession)
        cliCommand = cli.Command{
            Name:   newGroup.NewActionMetaData().Name,
            Description: newGroup.NewActionMetaData().Description,
            Usage:  newGroup.NewActionMetaData().Usage,
            Flags:  newGroup.NewActionMetaData().Flags,
            Action: cmd.Run,
        }
    })

    Describe("NewAction Testing", func() {
        Context("Happy Path", func() {
            It("Runs without issue", func() {
                err := testhelpers.RunCommand(cliCommand)
                Expect(err).NotTo(HaveOccurred())
                outputs := fakeUI.Outputs()
                Expect(outputs).To(ContainSubstring("some string we are testing for"))
            })
        })
    })
})
```

6. Edit `softlayer-cli/plugin/actions.go`
You will need to add a section for the new group of commands you are making. When just adding a new action to an existing group, this step can be skipped.

```go
    // ibmcloud sl newGroup
    newGroupCommands := newGroup.GetCommandAcionBindings(context, ui, session)
    for name, action := range newGroupCommands {
        CommandActionBindings[name] = action
    }

```

7. After adding your command and tests, build the i18n files `sh bin/fixeverything_i18n.sh`
8. Commit your code, push to a branch and submit a pull request

With this new pattern, making new commands should be a lot more straight forward and easier to manage.

## The Refactor

To begin the refactor, I am going to switch the `actions.go` file to the new pattern, so all CommandActionBindings will be in their resptive command group.

From there however there is still a lot of work to do. I will refactor the `dedicatedhost` group of commands as an example here on what needs to be moved where.

**All work done should be starting from the `cliRefactor` branch!**

1. Check that there is a file in your actions' group that matches the group name. For example with `dedicatedhost` group of actions, there should be a `softlayer-cli/plugin/commands/dedicatedhost/dedicatedhost.go` file. That file should contain 
    + `func GetCommandActionBindings()` that lists all the CommandActionBindings (and should NOT use any Constants to define those). These should be cut from the `actions.go` file and moved here.
        * In the return of each CommandActionBinding, it will be just `NewListGuestsCommand(ui, dedicatedhostManager).Run(c)` instead of `dedicatedhost.NewListGuestsCommand(ui, dedicatedhostManager).Run(c)`, since the code is now in the same namespace is the Command definition.
    + A Namespace (`func DedicatedHostNamespace() plugin.Namespace {`) in this case. Taken from metadata
    + Metadata (`DedicatedHostMetaData()`) in this case. Taken from metadata
    + A manager definition in the `GetCommandActionBindings` function if needed.
    + Make sure `GetCommandActionBindings` is spelled properly. Its missing the `T` in action in most cases.
2. Check the action you are working on is following the new pattern
    + `softlayer-cli/plugin/commands/dedicatedhost/list_guest.go`
        * Copy over the DedicatedhostListGuestsMetaData() from metadata
        * If command is using `OutputFlag()` that needs to be changed to `metadata.OutputFlag()` and make sure `"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"` is in the import.
        * Replace any Constants with proper strings
            - `Category:    CMD_DEDICATEDHOST_NAME` -> `Category:    "dedicatedhost",`
            - `Name:        CMD_DEDICATEDHOST_LIST_GUESTS_NAME,` -> `Name:        "list-guests",`
    + `softlayer-cli/plugin/commands/dedicatedhost/create.go`
        * Copy over the DedicatedhostListGuestsMetaData() from metadata
        * Make sure the `ForceFlag()` is changed to `metadata.ForceFlag()`
        * Replace any Constants with proper strings
    + Delete `softlayer-cli/plugin/metadata/dedicatedhost.go` after everything important is copied over.
    + Make sure `softlayer-cli/plugin/actions.go` is importing your CommandActionBindings. Should look like this.

```go
// ibmcloud sl dedicatedhost
dedicatedhostCommands := dedicatedhost.GetCommandActionBindings(context, ui, session)
for name, action := range dedicatedhostCommands {
    CommandActionBindings[name] = action
}
```

3. Update `softlayer-cli/plugin/plugin.go` to use the new NameSpace and MetaData locations
    +   `metadata.DedicatedhostMetaData(),` -> `dedicatedhost.DedicatedhostMetaData(),`
    +   `metadata.DedicatedhostNamespace(),` -> `dedicatedhost.DedicatedhostNamespace(),`
4. Clean up any `plugin\actions.go:57:2: dedicatedhostManager declared but not used` type of errors
5. Make sure you can build the plugin
6. update the Unit tests, as they also reference the metadata.
    + `softlayer-cli/plugin/commands/dedicatedhost/list_guest_test.go` 
        * `metadata.DedicatedhostListGuestsMetaData()` -> `dedicatedhost.DedicatedhostListGuestsMetaData()`
        * Remove the `"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"` import as its not needed here
    * `softlayer-cli/plugin/commands/dedicatedhost/create_test.go` 
        * `metadata.DedicatedhostListGuestsMetaData()` -> `dedicatedhost.DedicatedhostListGuestsMetaData()`
        * Remove the `"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"` import as its not needed here
7. Run the CLI unit tests
    *   `go test -coverprofile=coverage.out github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost`
8. Commit your changes and make a pull request against the `cliRefactor` branch.

## Checklist