package virtual

import (
	"encoding/json"
	"errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	bxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	ImageManager         managers.ImageManager
	Context              plugin.PluginContext
}

func NewCreateCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, imageManager managers.ImageManager, context plugin.PluginContext) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		ImageManager:         imageManager,
		Context:              context,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	virtualGuest := datatypes.Virtual_Guest{}
	var err error
	if c.NumFlags() == 0 {
		confirm, err := cmd.UI.Confirm(T("Please make sure you know all the creation options by running command: '{{.CommandName}} sl vs options'. Continue?",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
		params := make(map[string]interface{})
		params["hostname"], _ = cmd.UI.Ask(T("Hostname: "))
		params["domain"], _ = cmd.UI.Ask(T("Domain: "))
		inputCpu, _ := cmd.UI.Ask(T("Cpu: "))
		cpu, err := strconv.Atoi(inputCpu)
		if err != nil {
			return slErrors.NewInvalidSoftlayerIdInputError("CPU")
		}
		params["cpu"] = cpu
		inputMemory, _ := cmd.UI.Ask(T("Memory: "))
		memory, err := strconv.Atoi(inputMemory)
		if err != nil {
			return slErrors.NewInvalidSoftlayerIdInputError("Memory")
		}
		params["memory"] = memory
		params["datacenter"], _ = cmd.UI.Ask(T("Datacenter: "))
		params["os"], _ = cmd.UI.Ask(T("Operating System Code: "))

		_, err = cmd.VirtualServerManager.GenerateInstanceCreationTemplate(&virtualGuest, params)
		if err != nil {
			return err
		}
	} else {
		//create virtual server with customized parameters
		params, err := verifyParams(cmd.ImageManager, c)
		if err != nil {
			return err
		}
		_, err = cmd.VirtualServerManager.GenerateInstanceCreationTemplate(&virtualGuest, params)
		if err != nil {
			return err
		}
	}

	//do export
	if c.IsSet("export") {
		content, err := json.Marshal(virtualGuest)
		if err != nil {
			return cli.NewExitError(T("Failed to marshal virtual server template.\n")+err.Error(), 1)
		}
		export := c.String("export")
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(export, content, 0644)
		if err != nil {
			return cli.NewExitError(T("Failed to write virtual server template file to: {{.Template}}.",
				map[string]interface{}{"Template": export})+err.Error(), 1)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Virtual server template is exported to: {{.Template}}.", map[string]interface{}{"Template": export}))
		return nil
	}

	//do test
	if c.IsSet("test") {
		_, err := cmd.VirtualServerManager.VerifyInstanceCreation(virtualGuest)
		if err != nil {
			return cli.NewExitError(T("Failed to verify virtual server creation.\n")+err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order is correct."))
		return nil
	}

	quantity := c.Int("quantity")
	if quantity <= 0 {
		return bxErr.NewInvalidUsageError(T("The value of option '--quantity' should be greater or equal to 1."))
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	var multiErrors []error

	if c.IsSet("u") || c.IsSet("F") {
		var userData string
		if c.IsSet("u") {
			userData = c.String("u")
		}
		if c.IsSet("F") {
			userfile := c.String("F")
			content, err := ioutil.ReadFile(userfile) // #nosec
			if err != nil {
				newError := errors.New(T("Failed to read user data from file: {{.File}}.", map[string]interface{}{"File": userfile}))
				multiErrors = append(multiErrors, newError)
			}
			userData = string(content)
		}
		virtualGuest.UserData = []datatypes.Virtual_Guest_Attribute{datatypes.Virtual_Guest_Attribute{Value: &userData}}
	}

	var virtualGuests []datatypes.Virtual_Guest
	if quantity > 1 {
		virtualGuests, err = cmd.CreateMutliVSIWithSameConfig(virtualGuest, quantity)
		if err != nil {
			return cli.NewExitError(T("Failed to create multi virtual server instances.\n")+err.Error(), 2)
		}
	} else {
		virtualGuest, err = cmd.VirtualServerManager.CreateInstance(&virtualGuest)
		if err != nil {
			return cli.NewExitError(T("Failed to create virtual server instance.\n")+err.Error(), 2)
		}
		virtualGuests = append(virtualGuests, virtualGuest)
	}
	if c.IsSet("tag") || c.IsSet("g") {
		tagString := utils.StringSliceToString(c.StringSlice("tag"))
		for _, vs := range virtualGuests {
			err := cmd.VirtualServerManager.SetTags(*vs.Id, tagString)
			if err != nil {
				newError := errors.New(T("Failed to update the tag of virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": *vs.Id}) + err.Error())
				multiErrors = append(multiErrors, newError)
			}
		}
	}

	if len(virtualGuests) == 1 {
		cmd.printVirtualGuest(virtualGuests[0], c, &multiErrors)
	} else {
		cmd.printVirtualGuests(virtualGuests)
	}

	if len(multiErrors) > 0 {
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
	}
	return nil
}

func (cmd *CreateCommand) CreateMutliVSIWithSameConfig(virtualGuest datatypes.Virtual_Guest, quantity int) ([]datatypes.Virtual_Guest, error) {
	var virtualGuests []datatypes.Virtual_Guest
	var virtualGuestTmp datatypes.Virtual_Guest
	for i := 0; i < quantity; i++ {
		virtualGuestTmp = virtualGuest
		if i != 0 {
			hostName := *virtualGuestTmp.Hostname + "-" + strconv.Itoa(i)
			virtualGuestTmp.Hostname = &hostName
		}
		virtualGuests = append(virtualGuests, virtualGuestTmp)
	}
	virtualGuests, err := cmd.VirtualServerManager.CreateInstances(virtualGuests)
	if err != nil {
		return []datatypes.Virtual_Guest{}, cli.NewExitError(T("Failed to create multi virtual server instances.\n")+err.Error(), 2)
	}
	return virtualGuests, nil
}
func (cmd *CreateCommand) printVirtualGuest(virtualGuest datatypes.Virtual_Guest, c *cli.Context, multiErrors *[]error) {
	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add(T("ID"), utils.FormatIntPointer(virtualGuest.Id))
	table.Add(T("FQDN"), utils.FormatStringPointer(virtualGuest.FullyQualifiedDomainName))
	table.Add(T("Created"), utils.FormatSLTimePointer(virtualGuest.CreateDate))
	table.Add(T("GUID"), utils.FormatStringPointer(virtualGuest.GlobalIdentifier))
	table.Add(T("Placement Group ID"), utils.FormatIntPointer(virtualGuest.PlacementGroupId))

	//do wait
	if c.IsSet("wait") {
		until := time.Now().Add(time.Duration(c.Int("wait")) * time.Second)
		ready, _, err := cmd.VirtualServerManager.InstanceIsReady(*virtualGuest.Id, until)
		if err != nil {
			table.Add(T("ready"), "-")
			newError := errors.New(T("Failed to get ready status of virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": *virtualGuest.Id}) + err.Error())
			(*multiErrors) = append((*multiErrors), newError)
		} else {
			table.Add(T("ready"), strconv.FormatBool(ready))
		}
	}
	table.Print()
}

func (cmd *CreateCommand) printVirtualGuests(virtualGuests []datatypes.Virtual_Guest) {
	table := cmd.UI.Table([]string{T("ID"), T("HostName"), T("GUID"), T("Placement Group ID"), T("Created")})
	for _, vm := range virtualGuests {
		table.Add(utils.FormatIntPointer(vm.Id), utils.FormatStringPointer(vm.FullyQualifiedDomainName), utils.FormatStringPointer(vm.GlobalIdentifier), utils.FormatIntPointer(vm.PlacementGroupId), utils.FormatSLTimePointer(vm.CreateDate))
	}
	table.Print()
}

func verifyParams(imageManager managers.ImageManager, c *cli.Context) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if c.IsSet("flavor") {
		if c.IsSet("c") {
			return nil, bxErr.NewExclusiveFlagsError("[-c|--cpu]", "[--flavor]")
		}
		if c.IsSet("m") {
			return nil, bxErr.NewExclusiveFlagsError("[-m|--memory]", "[--flavor]")
		}
		if c.IsSet("dedicated") {
			return nil, bxErr.NewExclusiveFlagsError("[--dedicated]", "[--flavor]")
		}
		if c.IsSet("host-id") {
			return nil, bxErr.NewExclusiveFlagsError("[--host-id]", "[--flavor]")
		}
		params["flavor"] = c.String("flavor")
	}
	if c.IsSet("o") && c.IsSet("image") {
		return nil, bxErr.NewExclusiveFlagsError("[-o|--os]", "[--image]")
	}

	if c.IsSet("u") && c.IsSet("F") {
		return nil, bxErr.NewExclusiveFlagsError("[-u|--userdata]", "[-F|--userfile]")
	}

	if c.IsSet("t") {
		params["template"] = c.String("t")
	}

	if c.IsSet("like") {
		params["like"] = c.Int("like")
	}

	if c.IsSet("H") {
		params["hostname"] = c.String("H")
	}

	if c.IsSet("D") {
		params["domain"] = c.String("D")
	}

	if c.IsSet("c") {
		params["cpu"] = c.Int("c")
	}

	if c.IsSet("m") {
		params["memory"] = c.Int("m")
	}

	if c.IsSet("d") {
		params["datacenter"] = c.String("d")
	}

	if c.IsSet("o") {
		params["os"] = c.String("o")
	}

	if c.IsSet("image") {
		image, err := imageManager.GetImage(c.Int("image"))
		if err != nil {
			return nil, err
		}
		if image.GlobalIdentifier != nil {
			params["image"] = *image.GlobalIdentifier
		}
	}

	if !c.IsSet("billing") {
		params["billing"] = true
	} else {
		billing := c.String("billing")
		if billing == "hourly" {
			params["billing"] = true
		} else if billing == "monthly" {
			params["billing"] = false
		} else {
			return nil, bxErr.NewInvalidUsageError(T("[--billing] billing rate must be either hourly or monthly."))
		}
	}

	if c.IsSet("dedicated") {
		params["dedicated"] = true
	} else {
		params["dedicated"] = false
	}

	if c.IsSet("host-id") {
		params["host-id"] = c.Int("host-id")
		params["dedicated"] = true
	}

	if c.IsSet("private") {
		params["private"] = true
	} else {
		params["private"] = false
	}

	if c.IsSet("san") {
		params["san"] = true
	}

	if c.IsSet("i") {
		params["i"] = c.String("i")
	}

	if c.IsSet("key") || c.IsSet("k") {
		params["sshkeys"] = c.IntSlice("k")
	}

	if c.IsSet("disk") {
		params["disks"] = c.IntSlice("disk")
	}

	if c.IsSet("n") {
		params["network"] = c.Int("n")
	}

	if c.IsSet("vlan-public") {
		params["vlan-public"] = c.Int("vlan-public")
	}

	if c.IsSet("vlan-private") {
		params["vlan-private"] = c.Int("vlan-private")
	}

	if c.IsSet("subnet-public") {
		if !c.IsSet("vlan-public") {
			return nil, bxErr.NewMissingInputError("--vlan-public")
		}
		params["subnet-public"] = c.Int("subnet-public")
	}

	if c.IsSet("subnet-private") {
		if !c.IsSet("vlan-private") {
			return nil, bxErr.NewMissingInputError("--vlan-private")
		}
		params["subnet-private"] = c.Int("subnet-private")
	}

	if c.IsSet("S") || c.IsSet("public-security-group") {
		params["public-security-group"] = c.IntSlice("public-security-group")
	}

	if c.IsSet("s") || c.IsSet("private-security-group") {
		params["private-security-group"] = c.IntSlice("private-security-group")
	}

	if c.IsSet("boot-mode") {
		if c.String("boot-mode") != "HVM" && c.String("boot-mode") != "PV" {
			return nil, bxErr.NewInvalidUsageError("--boot-mode should be HVM | PV")
		}
		params["boot-mode"] = c.String("boot-mode")
	}

	if c.IsSet("transient") {
		params["transient"] = c.Bool("transient")
	}

	if c.IsSet("placement-group-id") {
		params["placement-group-id"] = c.Int("placement-group-id")
	}
	return params, nil
}

func VSCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "create",
		Description: T("Create virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs create [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413
	This command orders a virtual server instance with hostname is myvsi, domain is ibm.com, 4 cpu cores, 4096M memory, located at datacenter: dal10,
	operation system is UBUNTU 16 64 bits, 2 disks, one is 100G, the other is 1000G, and placed at public vlan with ID 413.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --test
	This command tests whether the order is valid with above options before the order is actually placed.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --export ~/myvsi.txt
	This command exports above options to a file: myvsi.txt under user home directory for later use.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN [required]"),
			},
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Number of CPU cores [required]"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Memory in megabytes [required]"),
			},
			cli.StringFlag{
				Name:  "flavor",
				Usage: T("Public Virtual Server flavor key name"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname [required]"),
			},
			cli.StringFlag{
				Name:  "o,os",
				Usage: T("OS install code. Tip: you can specify <OS>_LATEST"),
			},
			cli.IntFlag{
				Name:  "image",
				Usage: T("Image ID. See: '${COMMAND_NAME} sl image list' for reference"),
			},
			cli.StringFlag{
				Name:  "billing",
				Usage: T("Billing rate. Default is: hourly. Options are: hourly, monthly"),
			},
			cli.BoolFlag{
				Name:  "dedicated",
				Usage: T("Create a dedicated Virtual Server (Private Node)"),
			},
			cli.IntFlag{
				Name:  "host-id",
				Usage: T("Host Id to provision a Dedicated Virtual Server onto"),
			},
			cli.BoolFlag{
				Name:  "san",
				Usage: T("Use SAN storage instead of local disk"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not actually create the virtual server"),
			},
			cli.StringFlag{
				Name:  "export",
				Usage: T("Exports options to a template file"),
			},
			cli.StringFlag{
				Name:  "i,postinstall",
				Usage: T("Post-install script to download"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("The IDs of the SSH keys to add to the root user (multiple occurrence permitted)"),
			},
			cli.IntSliceFlag{
				Name:  "disk",
				Usage: T("Disk sizes (multiple occurrence permitted)"),
			},
			cli.BoolFlag{
				Name:  "private",
				Usage: T("Forces the virtual server to only have access the private network"),
			},
			cli.StringFlag{
				Name:  "like",
				Usage: T("Use the configuration from an existing virtual server"),
			},
			cli.IntFlag{
				Name:  "n,network",
				Usage: T("Network port speed in Mbps"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Tags to add to the instance (multiple occurrence permitted)"),
			},
			cli.StringFlag{
				Name:  "t,template",
				Usage: T("A template file that defaults the command-line options"),
			},
			cli.StringFlag{
				Name:  "u,userdata",
				Usage: T("User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "F,userfile",
				Usage: T("Read userdata from file"),
			},
			cli.StringFlag{
				Name:  "vlan-public",
				Usage: T("The ID of the public VLAN on which you want the virtual server placed"),
			},
			cli.StringFlag{
				Name:  "vlan-private",
				Usage: T("The ID of the private VLAN on which you want the virtual server placed"),
			},
			cli.IntSliceFlag{
				Name:  "S,public-security-group",
				Usage: T("Security group ID to associate with the public interface (multiple occurrence permitted)"),
			},
			cli.IntSliceFlag{
				Name:  "s,private-security-group",
				Usage: T("Security group ID to associate with the private interface (multiple occurrence permitted)"),
			},
			cli.IntFlag{
				Name:  "wait",
				Usage: T("Wait until the virtual server is finished provisioning for up to X seconds before returning. It's not compatible with option --quantity"),
			},
			cli.IntFlag{
				Name:  "placement-group-id",
				Usage: T("Placement Group Id to order this guest on."),
			},
			cli.StringFlag{
				Name:  "boot-mode",
				Usage: T("Specify the mode to boot the OS in. Supported modes are HVM and PV."),
			},
			cli.IntFlag{
				Name:  "subnet-public",
				Usage: T("The ID of the public SUBNET on which you want the virtual server placed"),
			},
			cli.IntFlag{
				Name:  "subnet-private",
				Usage: T("The ID of the private SUBNET on which you want the virtual server placed"),
			},
			cli.BoolFlag{
				Name:  "transient",
				Usage: T("Create a transient virtual server"),
			},
			cli.IntFlag{
				Name:  "quantity",
				Usage: T("The quantity of virtual server be created. It should be greater or equal to 1. This value defaults to 1."),
				Value: 1,
			},
			metadata.ForceFlag(),
		},
	}
}
