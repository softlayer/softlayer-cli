package testhelpers

import (
	"flag"
	"fmt"
	"log"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	faketerminal "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.com/urfave/cli"
	"github.com/spf13/cobra"
)

func RunCommand(c cli.Command, args ...string) error {

	cli.OsExiter = func(errorCode int) {
		fmt.Println("Fake OS Exit:", errorCode)
	}
	app := cli.NewApp()
	set := flag.NewFlagSet("test "+c.Name, 0)
	for _, f := range c.Flags {
		f.Apply(set)
	}
	err := set.Parse(append([]string{c.Name}, args...))
	if err != nil {
		return err
	}
	ctx := cli.NewContext(app, set, nil)
	return c.Run(ctx)
}



func RunCobraCommand(cmd *cobra.Command, args ...string) error {
	// If we do cmd.SetArgs(args) with no args, Cobra will try to read them from the actual command line
	// which breaks unit tests when using -ginkgo.focus (or other) flags.
	if len(args) == 0 {
		cmd.SetArgs([]string{})	
	} else {
		cmd.SetArgs(args)
	}
	
	
	_, err := cmd.ExecuteC()
	return err
}

// For when you just want to get a fake context, not actually run the command yet
func GetCliContext(name string) *cli.Context {
	app := cli.NewApp()
	set := flag.NewFlagSet("test "+name, 0)

	ctx := cli.NewContext(app, set, nil)
	return ctx
}

// For when you just want to get a fake context, and have the --help flag set
func GetCliContextHelp(name string) *cli.Context {

	app := cli.NewApp()
	set := flag.NewFlagSet("test "+name, 0)
	helpFlag := cli.BoolFlag{
		Name:  "help",
		Usage: "The Help Flag",
	}
	helpFlag.Apply(set)
	err := set.Parse(append([]string{name}, "help"))
	if err != nil {
		fmt.Printf("Error in GetCliContextHelp() "+ err.Error() + "\n")
	}
	ctx := cli.NewContext(app, set, nil)

	return ctx
}

type CMD struct {
	UI terminal.UI
}

func NewCommand(ui terminal.UI) *CMD {
	return &CMD{
		UI: ui,
	}
}

func (cmd *CMD) Run(c *cli.Context) error {
	cmd.UI.Print("command name:", c.Command.Name)
	cmd.UI.Print("command args:", c.Args())
	cmd.UI.Print("flag value:", c.String("f"))
	return nil
}

func main() {
	fakeUI := faketerminal.NewFakeUI()
	cmd := NewCommand(fakeUI)
	cliCommand := cli.Command{
		Name:        "block",
		Description: "test command",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "f",
				Usage: "help",
			},
		},
		Action: cmd.Run,
	}
	err := RunCommand(cliCommand, "123", "-f", "456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fakeUI.Outputs())
}
