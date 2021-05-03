package testhelpers

import (
	"flag"
	"fmt"
	"log"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	faketerminal "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.com/urfave/cli"
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
