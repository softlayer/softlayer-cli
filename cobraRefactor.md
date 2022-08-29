# Cobra Refactor

This document will serve as a guide for replacing the `github.com/urfave/cli` library with `github.com/spf13/cobra`

The cobra library is what the main ibmcloud CLI uses, and has a lot more features in terms of supporting flag/argument options.

Along with changing the library, I'd also like to enforce that each command take in as arguments the terminal UI, and session parameters, and nothing else. Each CLI command should create an instance of a softlayer-cli manager if needed.

## Environment Setup

We will be working from the [`cobraCommands`](https://github.ibm.com/SoftLayer/softlayer-cli/tree/cobraCommands/plugin) branch for this project. Name your branch something like `cobra<commandGroup>` or something similar. cobra and the commandGroup should be in the branch name, and if needed a number or some other marker.

```bash
git checkout cobraCommands
git checkout -b cobraVirtual
```

## Refactoring Pattern

This example will cover refactoring the `sl account` commands, the same patterns should be used in the other command groups as well.

### Refactor the command group

`GetCommandActionBindings()` This function setup each of the command group's commands. In the new pattern, we will use a new function called `SetupCobraCommands()` the  `Metadata` function will no longer be needed.

#### account.go

Comment out `GetCommandActionBindings` as it will be used to keep track of which commands need to be added. 

The basic template will look something like this. 

```go
// Add these imports
import (
    "github.com/spf13/cobra" // Cobra for new CLI
    "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata" // has the struct SoftlayerCommand
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
    cobraCmd := &cobra.Command{
        Use: "account", // AccountMetaData()->Name
        Short: T("Classic infrastructure Account commands"), // AccountMetaData()->Description
        Long: "${COMMAND_NAME} sl account", // AccountMetaData()->USage
        RunE: nil,
    }
    // Add actual commands in this section.
    cobraCmd.AddCommand(NewBandwidthPoolsCommand(sl).Command)
    cobraCmd.AddCommand(NewBandwidthPoolsDetailCommand(sl).Command)
    cobraCmd.AddCommand(NewBillingItemsCommand(sl).Command)
    return cobraCmd 
}
```


### Refactor the command

Each command has a few things that need to be updated.

1. `type <Command>Command struct {}` which defines a few properties for the command. Important to this refactor is that each paramter be added as a property to this struct.
2. `func New<Command>Command` return type will change to \*cobra.Command and this function will create an instance of the cobraCommand itself. This means we no longer need a <Command>MetaData() function.
3. `func Run()` this function will require quite a few changes, especially in instances where paramters/arguments are referenced.

#### bandwidth_pools.go

1. Add the cobra import, and metadata import (if it doens't exist already)
```go
import "github.com/spf13/cobra"
import "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
```

2. Update pattern and add any flags/options to the `BandwidthPoolsCommand` struct. This command only uses the `OutputFlag` flag... which is included in SoftlayerCommand, so we dont need to do anything else.

```go
type BandwidthPoolsCommand struct {
    *metadata.SoftlayerCommand
    Command *cobra.Command
    AccountManager managers.AccountManager
    //Flag1 int
    //Flag2 string
}
```

`*metadata.SoftlayerCommand  // this format makes BandwidthPoolsCommand inherit the properties from SoftlayerCommand, so it has access to UI and Session `
`Command *cobra.Command  // all commands will have this, the reference to the actual cobra.Command`

Because BandwidthPoolsCommand inherits from metadata.SoftlayerCommand, BandwidthPoolsCommand has access to BandwidthPoolsCommand.UI, BandwidthPoolsCommand.Session, and BandwidthPoolsCommand.OutputFlag.

This struct should also have a reference for the manager the command needs.
Make sure to list any flags/options here as well.


3. Update `NewBandwidthPoolsCommand` to be the new pattern. `thisCmd` will be a reference to the BandwidthPoolsCommand struct, and `cobraCmd` will be the actual Cobra command instance. `RunE` will have the reference to the actual `thisCmd.Run(args)` method. The Command Struct should also have a variable for each manager the command needs (usually just one though)

```go
func NewBandwidthPoolsCommand(sl *metadata.SoftlayerCommand) *BandwidthPoolsCommand {
    thisCmd := &BandwidthPoolsCommand{ //Update this line
        SoftlayerCommand: sl,
        AccountManager: managers.NewAccountManager(sl.Session),
    }
    cobraCmd := &cobra.Command{
        // The first 'word' in the Use line is the command name. Anything after that will show up in the help text
        Use: "THE-COMMAND-NAME-HERE",  // if a command takes arguments, add them here in ex: + T("IDENTIFIER")
        Short: T("Your translated short description goes here"), // Updates this from metadata
        Long: "",  // Remove this if the Usage from the old command is just basic information about how to run it. The Long description should be for examples, detailed information about the command.
        Args: metadata.NoArgs, // Make sure this accepts the correct number of args
        RunE: func(cmd *cobra.Command, args []string) error {
            return thisCmd.Run(args)
        },
    }
    // Add any flags here
    // cobraCmd.Flags().IntVar(&thisCmd.Init, "init", 0, T("Init parameter"))
    thisCmd.Command = cobraCmd

    return thisCmd
}
```
Make sure to remove all these comment lines if you are copy/pasting this bit. They exist to remind you which lines need attention.
Options for `Args`: https://github.com/spf13/cobra/blob/main/user_guide.md#positional-and-custom-arguments
I had to re-define these in the metadata package though (/plugin/commands/sl/args.go) so I could translate the error messages. Basically just copy/paste though. Just do `metadata.MaximumNArgs(1)` instead of `cobra.MaximumNArgs(1)` or whatever.

In the `Long` description, We no longer need filler like `${COMMAND_NAME} sl account bandwidth-pools` as that will be added to the help text automatically now. Its ok to keep anything that are EXAMPLE though.


4. Delete `BandwidthPoolsMetaData()`
5. Update `Run()` to follow new pattern

```go
func (cmd *BandwidthPoolsCommand) Run(args []string) error {
    pools, err := cmd.AccountManager.GetBandwidthPools()
    if err != nil {
        return err
    }

    outputFormat := cmd.GetOutputFlag()
    // Rest of the function was unchanged
    return nil
}

```

Since we can enforce a number of args in the `NewBandwidthPoolsCommand` definition, the following, and similar lines are no longer needed (if you see them, not all commands have them)
```go
// This section can be removed,  `Args: metadata.OneArgs,` is the same thing.
    if c.NArg() != 1 {
        return errors.NewInvalidUsageError(T("This command requires one argument."))
    }

// Replace this
// bandwidthPoolId, err := strconv.Atoi(c.Args()[0])
// With
bandwidthPoolId, err := strconv.Atoi(args[0])
```

```go
/* Replace this outputFormat style
    outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
    if err != nil {
        return err
    }
*/
    outputFormat := cmd.GetOutputFlag()
```

Update any cli.NewExitError messages
```go
// Replace this
//return cli.NewExitError(T("Failed to get Bandwidth Pool.\n")+err.Error(), 2)
// with
return errors.NewAPIError(T("Failed to get Bandwidth Pool."), err.Error(), 2)
```
`errors.NewAPIError`  is something I just added so make printing API error easier. It takes 3 arguments, the error you want to tell the user, the API error, and an error code (which isn't used for anything at the moment). It will automatically put a "\n" between the 2 messages for you. If removing the \n from the translated string, make sure to remove it from the `id`s of all i18n files.


6. Update SetupCobraCommands in account.go to add this command, and remove it from `GetCommandACtionBindings`

```go
func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
    cobraCmd := &cobra.Command{
        Use: "account",
        Short: T("Classic infrastructure Account commands"),
        Long: "${COMMAND_NAME} sl account",
        RunE: nil,
    }
    cobraCmd.AddCommand(NewBandwidthPoolsCommand(sl))
    return cobraCmd 
}
```

7. Check if it builds, remove any unused imports.

### Refactor the action file

### Refactor Unit Tests

Quick Replaces / copy-paste

Test setup, just need to change `cliCommand` and the accountManager
```go
    var (
        fakeUI              *terminal.FakeUI
        cliCommand          *account.{{WHATEVER}}Command
        fakeSession         *session.Session
        slCommand           *metadata.SoftlayerCommand
        fakeAccountManager  *testhelpers.FakeAccountManager
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        fakeAccountManager = new(testhelpers.FakeAccountManager)
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = account.New{{WHATEVER}}Command(slCommand)
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
        cliCommand.AccountManager = fakeAccountManager
    })
```
`RunCommand(cliCommand   ---->   RunCobraCommand(cliCommand.Command`

#### account_test.go

Most of these tests can be removed. We will just need to test `SetupCobraCommands` and `AccountNamespace`

since were not using the actionBindings anymore, I think its ok we don't test to make sure commands don't get accidently removed.
```go
var _ = Describe("Test account.GetCommandActionBindings()", func() {
    fakeUI := terminal.NewFakeUI()
    fakeSession := testhelpers.NewFakeSoftlayerSession(nil)
    slMeta := &metadata.SoftlayerCommand{fakeUI, fakeSession, ""}
    Context("New commands testable", func() {
        accountCommands := account.SetupCobraCommands(slMeta)
        Expect(accountCommands.Name()).To(Equal("account"))
    })
    Context("Account Namespace", func() {
        It("Account Name Space", func() {
            Expect(account.AccountNamespace().ParentName).To(ContainSubstring("sl"))
            Expect(account.AccountNamespace().Name).To(ContainSubstring("account"))
            Expect(account.AccountNamespace().Description).To(ContainSubstring("Classic infrastructure Account"))
        })
    })
})
```


#### bandwidth_pools_test.go
With these tests, you'll need to add a var `slCommand *metadata.SoftlayerCommand` and initialize it with fakeUI and fakeSession like this: `slCommand = &metadata.SoftlayerCommand{fakeUI, fakeSession, ""}`, then use that to create an instance of the command you want to test, `cmd = account.NewBandwidthPoolsCommand(slCommand)`

Replace `RunCommand` with `RunCobraCommand`
 
```go
var _ = Describe("Account Bandwidth-Pools", func() {
    var (
        fakeUI          *terminal.FakeUI
        cliCommand      *account.BandwidthPoolsCommand
        fakeSession     *session.Session
        slCommand       *metadata.SoftlayerCommand
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = account.NewBandwidthPoolsCommand(slCommand)
        // Only needed if your testing json output for this command.
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
    })

    Describe("Bandwidth-Pools Testing", func() {
        Context("Happy Path", func() {
            It("Runs without issue", func() {
                err := testhelpers.RunCobraCommand(cliCommand.Command)
                Expect(err).NotTo(HaveOccurred())
                outputs := fakeUI.Outputs()
                Expect(outputs).To(ContainSubstring("3361 GB      7.13 GB         7.70 GB"))
            })
            It("Outputs JSON", func() {
                err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=JSON")
                Expect(err).NotTo(HaveOccurred())
                outputs := fakeUI.Outputs()
                Expect(outputs).To(ContainSubstring("\"amountIn\": 7.54252,"))
    
            })
        })
    })
})
```

Useful Regex 

```regex
Expect\(strings.Contains\(err.Error\(\), \"(.*)\"\)\)\.To\(BeTrue\(\)\)
Expect(err.Error()).To(ContainSubstring("\1"))
```

```
RunCommand(cliCommand
RunCobraCommand(cliCommand.Command
```


#### bandwidth_pools_details_test.go

For tests that need a fake manager, use this pattern.


```go
var _ = Describe("account bandwidth_pools_details", func() {
    var (
        fakeUI              *terminal.FakeUI
        cliCommand          *account.BandwidthPoolsDetailCommand
        fakeSession         *session.Session
        slCommand           *metadata.SoftlayerCommand
        fakeAccountManager  *testhelpers.FakeAccountManager
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        fakeAccountManager = new(testhelpers.FakeAccountManager)
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = account.NewBandwidthPoolsDetailCommand(slCommand)
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
        cliCommand.AccountManager = fakeAccountManager
    })
    Describe("account bandwidth_pools_details", func() {
        Context("Return error", func() {
            BeforeEach(func() {
                fakeAccountManager.GetBandwidthPoolDetailReturns(datatypes.Network_Bandwidth_Version1_Allotment{}, errors.New("Failed to get Bandwidth Pool."))
            })
            It("Failed Bandwidth Pool", func() {
                err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
                Expect(err).To(HaveOccurred())
                Expect(err.Error()).To(ContainSubstring("Failed to get Bandwidth Pool."))
            })
        })
```


## Pull Requests

Smaller pull requests are easier to check, so I would prefer making pull requests contain AT MOST 1 group of commands, and at least 1 actual command. Ideally limit your pull requests to about a days worth of work, and I'll try to get them all merged in at the start of the day.

`cobraCommands` is the branch we will be merging into while we get everything ready.