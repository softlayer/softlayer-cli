package virtual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/datatypes"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	Command              *cobra.Command
	VirtualServerManager managers.VirtualServerManager
	ImageManager         managers.ImageManager

	Dedicated      bool
	Private        bool
	San            bool
	Test           bool
	Transient      bool
	Force          bool
	Disk           []int
	Key            []int
	PriSecGroup    []int
	PubSecGroup    []int
	HostId         int
	Image          int
	Like           int
	PlacementGroup int
	Quantity       int
	SubnetPrivate  int
	SubnetPublic   int
	VlanPrivate    int
	VlanPublic     int
	Wait           int
	CPU            int
	Memory         int
	Network        int
	Tag            []string
	Billing        string
	BootMode       string
	Export         string
	Flavor         string
	Datacenter     string
	Domain         string
	Hostname       string
	Os             string
	PostInstall    string
	Template       string
	Userdata       string
	Userfile       string
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) (cmd *CreateCommand) {
	thisCmd := &CreateCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
		ImageManager:         managers.NewImageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create virtual server instance"),
		Long: T(`${COMMAND_NAME} sl vs create [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413
	This command orders a virtual server instance with hostname is myvsi, domain is ibm.com, 4 cpu cores, 4096M memory, located at datacenter: dal10,
	operation system is UBUNTU 16 64 bits, 2 disks, one is 100G, the other is 1000G, and placed at public vlan with ID 413.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --test
	This command tests whether the order is valid with above options before the order is actually placed.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --export ~/myvsi.txt
	This command exports above options to a file: myvsi.txt under user home directory for later use.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVar(&thisCmd.Dedicated, "dedicated", false, T("Create a dedicated Virtual Server (Private Node)"))
	cobraCmd.Flags().BoolVar(&thisCmd.Private, "private", false, T("Forces the virtual server to only have access the private network"))
	cobraCmd.Flags().BoolVar(&thisCmd.San, "san", false, T("Use SAN storage instead of local disk"))
	cobraCmd.Flags().BoolVar(&thisCmd.Test, "test", false, T("Do not actually create the virtual server"))
	cobraCmd.Flags().BoolVar(&thisCmd.Transient, "transient", false, T("Create a transient virtual server"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	cobraCmd.Flags().IntSliceVar(&thisCmd.Disk, "disk", []int{}, T("Disk sizes (multiple occurrence permitted)"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Key, "key", "k", []int{}, T("The IDs of the SSH keys to add to the root user (multiple occurrence permitted)"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.PriSecGroup, "private-security-group", "s", []int{}, T("Security group ID to associate with the private interface (multiple occurrence permitted)"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.PubSecGroup, "public-security-group", "S", []int{}, T("Security group ID to associate with the public interface (multiple occurrence permitted)"))
	cobraCmd.Flags().IntVar(&thisCmd.HostId, "host-id", 0, T("Host Id to provision a Dedicated Virtual Server onto"))
	cobraCmd.Flags().IntVar(&thisCmd.Image, "image", 0, T("Image ID. See: '${COMMAND_NAME} sl image list' for reference"))
	cobraCmd.Flags().IntVar(&thisCmd.Like, "like", 0, T("Use the configuration from an existing virtual server"))
	cobraCmd.Flags().IntVar(&thisCmd.PlacementGroup, "placement-group-id", 0, T("Placement Group Id to order this guest on."))
	cobraCmd.Flags().IntVar(&thisCmd.Quantity, "quantity", 1, T("The quantity of virtual server be created. It should be greater or equal to 1. This value defaults to 1."))
	cobraCmd.Flags().IntVar(&thisCmd.SubnetPrivate, "subnet-private", 0, T("The ID of the private SUBNET on which you want the virtual server placed"))
	cobraCmd.Flags().IntVar(&thisCmd.SubnetPublic, "subnet-public", 0, T("The ID of the public SUBNET on which you want the virtual server placed"))
	cobraCmd.Flags().IntVar(&thisCmd.VlanPrivate, "vlan-private", 0, T("The ID of the private VLAN on which you want the virtual server placed"))
	cobraCmd.Flags().IntVar(&thisCmd.VlanPublic, "vlan-public", 0, T("The ID of the public VLAN on which you want the virtual server placed"))
	cobraCmd.Flags().IntVar(&thisCmd.Wait, "wait", 0, T("Wait until the virtual server is finished provisioning for up to X seconds before returning. It's not compatible with option --quantity"))
	cobraCmd.Flags().IntVarP(&thisCmd.CPU, "cpu", "c", 0, T("Number of CPU cores [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Memory, "memory", "m", 0, T("Memory in megabytes [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Network, "network", "n", 0, T("Network port speed in Mbps"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "g", []string{}, T("Tags to add to the instance (multiple occurrence permitted)"))
	cobraCmd.Flags().StringVar(&thisCmd.Billing, "billing", "hourly", T("Billing rate. Default is: hourly. Options are: hourly, monthly"))
	cobraCmd.Flags().StringVar(&thisCmd.BootMode, "boot-mode", "", T("Specify the mode to boot the OS in. Supported modes are HVM and PV."))
	cobraCmd.Flags().StringVar(&thisCmd.Export, "export", "", T("Exports options to a template file"))
	cobraCmd.Flags().StringVar(&thisCmd.Flavor, "flavor", "", T("Public Virtual Server flavor key name"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Datacenter shortname [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Domain portion of the FQDN [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Host portion of the FQDN [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Os, "os", "o", "", T("OS install code. Tip: you can specify <OS>_LATEST"))
	cobraCmd.Flags().StringVarP(&thisCmd.PostInstall, "postinstall", "i", "", T("Post-install script to download"))
	cobraCmd.Flags().StringVarP(&thisCmd.Template, "template", "t", "", T("A template file that defaults the command-line options"))
	cobraCmd.Flags().StringVarP(&thisCmd.Userdata, "userdata", "u", "", T("User defined metadata string"))
	cobraCmd.Flags().StringVarP(&thisCmd.Userfile, "userfile", "F", "", T("Read userdata from file"))
	return thisCmd
}

func (cmd *CreateCommand) CheckRequiredOptions() bool {
	if cmd.Hostname == "" {
		return false
	} else if cmd.Domain == "" {
		return false
	} else if cmd.Datacenter == "" {
		return false
	} else if cmd.Os == "" && cmd.Image == 0 {
		return false
	}
	return true
}

func (cmd *CreateCommand) Run(args []string) error {
	virtualGuest := datatypes.Virtual_Guest{}
	var err error
	params, err := cmd.verifyParams()
	if err != nil {
		return err
	}
	if !cmd.CheckRequiredOptions() {
		confirm, err := cmd.UI.Confirm(T("Please make sure you know all the creation options by running command: '{{.CommandName}} sl vs options'. Continue?",
			map[string]interface{}{"CommandName": "ibmcloud"}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}

		if params["hostname"] == "" {
			params := make(map[string]interface{})
			params["hostname"], _ = cmd.UI.Ask(T("Hostname") + ": "))
			// return slErrors.NewMissingInputError("[-H|--hostname]")
		}

		if params["domain"] == "" {
			params["domain"], _ = cmd.UI.Ask(T("Domain") + ": "))
		}

		if params["cpu"] == 0 && params["flavor"] == "" {
			inputCpu, _ := cmd.UI.Ask(T("Cpu: "))
			cpu, err := strconv.Atoi(inputCpu)
			if err != nil {
				return slErrors.NewInvalidSoftlayerIdInputError("CPU")
			}
			params["cpu"] = cpu
		}

		if params["memory"] == 0 && params["flavor"] == "" {
			inputMemory, _ := cmd.UI.Ask(T("Memory") + ": "))
			memory, err := strconv.Atoi(inputMemory)
			if err != nil {
				return slErrors.NewInvalidSoftlayerIdInputError("Memory")
			}
			if memory <= 0 {
				return slErrors.NewInvalidUsageError(T("either [-m|--memory] or [--flavor] is required."))
			}
			params["memory"] = memory
		}

		if params["datacenter"] == "" {
			params["datacenter"], _ = cmd.UI.Ask(T("Datacenter") + ": "))
		}

		if params["os"] == "" {
			params["os"], _ = cmd.UI.Ask(T("Operating System Code") + ": "))
		}

		_, err = cmd.VirtualServerManager.GenerateInstanceCreationTemplate(&virtualGuest, params)
		if err != nil {
			return err
		}
	} else {
		//create virtual server with customized parameters
		_, err = cmd.VirtualServerManager.GenerateInstanceCreationTemplate(&virtualGuest, params)
		if err != nil {
			return err
		}
	}

	//do export
	if cmd.Export != "" {
		content, err := json.Marshal(virtualGuest)
		if err != nil {
			return slErrors.NewAPIError(T("Failed to marshal virtual server template.\n"), err.Error(), 1)
		}
		export := cmd.Export
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(export, content, 0644)
		if err != nil {
			return slErrors.NewAPIError(T("Failed to write virtual server template file to: {{.Template}}.",
				map[string]interface{}{"Template": export}), err.Error(), 1)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Virtual server template is exported to: {{.Template}}.", map[string]interface{}{"Template": export}))
		return nil
	}

	//do test
	if cmd.Test {
		_, err := cmd.VirtualServerManager.VerifyInstanceCreation(virtualGuest)
		if err != nil {
			return slErrors.NewAPIError(T("Failed to verify virtual server creation.\n"), err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order is correct."))
		return nil
	}

	quantity := cmd.Quantity
	if quantity <= 0 {
		return slErrors.NewInvalidUsageError(T("The value of option '--quantity' should be greater or equal to 1."))
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	var multiErrors []error

	if cmd.Userdata != "" || cmd.Userfile != "" {
		var userData string
		if cmd.Userdata != "" {
			userData = cmd.Userdata
		}
		if cmd.Userfile != "" {
			userfile := cmd.Userfile
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
			return slErrors.NewAPIError(T("Failed to create multi virtual server instances.\n"), err.Error(), 2)
		}
	} else {
		virtualGuest, err = cmd.VirtualServerManager.CreateInstance(&virtualGuest)
		if err != nil {
			return slErrors.NewAPIError(T("Failed to create virtual server instance.\n"), err.Error(), 2)
		}
		virtualGuests = append(virtualGuests, virtualGuest)
	}
	if len(cmd.Tag) > 0 {
		tagString := utils.StringSliceToString(cmd.Tag)
		for _, vs := range virtualGuests {
			err := cmd.VirtualServerManager.SetTags(*vs.Id, tagString)
			if err != nil {
				newError := errors.New(T("Failed to update the tag of virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": *vs.Id}) + err.Error())
				multiErrors = append(multiErrors, newError)
			}
		}
	}

	if len(virtualGuests) == 1 {
		cmd.printVirtualGuest(virtualGuests[0], &multiErrors)
	} else {
		cmd.printVirtualGuests(virtualGuests)
	}

	if len(multiErrors) > 0 {
		return slErrors.CollapseErrors(multiErrors)
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
		return []datatypes.Virtual_Guest{}, slErrors.NewAPIError(T("Failed to create multi virtual server instances.\n"), err.Error(), 2)
	}
	return virtualGuests, nil
}

func (cmd *CreateCommand) printVirtualGuest(virtualGuest datatypes.Virtual_Guest, multiErrors *[]error) {
	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add(T("ID"), utils.FormatIntPointer(virtualGuest.Id))
	table.Add(T("FQDN"), utils.FormatStringPointer(virtualGuest.FullyQualifiedDomainName))
	table.Add(T("Created"), utils.FormatSLTimePointer(virtualGuest.CreateDate))
	table.Add(T("GUID"), utils.FormatStringPointer(virtualGuest.GlobalIdentifier))
	table.Add(T("Placement Group ID"), utils.FormatIntPointer(virtualGuest.PlacementGroupId))

	//do wait
	if cmd.Wait > 0 {
		until := time.Now().Add(time.Duration(cmd.Wait) * time.Second)
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
	table := cmd.UI.Table([]string{T("ID"), T("Hostname"), T("GUID"), T("Placement Group ID"), T("Created")})
	for _, vm := range virtualGuests {
		table.Add(utils.FormatIntPointer(vm.Id), utils.FormatStringPointer(vm.FullyQualifiedDomainName), utils.FormatStringPointer(vm.GlobalIdentifier), utils.FormatIntPointer(vm.PlacementGroupId), utils.FormatSLTimePointer(vm.CreateDate))
	}
	table.Print()
}

func (cmd *CreateCommand) verifyParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})

	if cmd.Flavor != "" {
		if cmd.CPU != 0 {
			fmt.Printf("Returning an error....\n")
			return nil, slErrors.NewExclusiveFlagsError("[-c|--cpu]", "[--flavor]")
		}
		if cmd.Memory != 0 {
			return nil, slErrors.NewExclusiveFlagsError("[-m|--memory]", "[--flavor]")
		}
		if cmd.Dedicated {
			return nil, slErrors.NewExclusiveFlagsError("[--dedicated]", "[--flavor]")
		}
		if cmd.HostId != 0 {
			return nil, slErrors.NewExclusiveFlagsError("[--host-id]", "[--flavor]")
		}
		params["flavor"] = cmd.Flavor
	}
	if cmd.Os != "" && cmd.Image != 0 {
		return nil, slErrors.NewExclusiveFlagsError("[-o|--os]", "[--image]")
	}

	if cmd.Userdata != "" && cmd.Userfile != "" {
		return nil, slErrors.NewExclusiveFlagsError("[-u|--userdata]", "[-F|--userfile]")
	}

	if cmd.Template != "" {
		params["template"] = cmd.Template
	}

	if cmd.Like != 0 {
		params["like"] = cmd.Like
	}

	if cmd.Hostname != "" {
		params["hostname"] = cmd.Hostname
	}

	if cmd.Domain != "" {
		params["domain"] = cmd.Domain
	}

	if cmd.CPU != 0 {
		params["cpu"] = cmd.CPU
	}

	if cmd.Memory != 0 {
		params["memory"] = cmd.Memory
	}

	if cmd.Datacenter != "" {
		params["datacenter"] = cmd.Datacenter
	}

	if cmd.Os != "" {
		params["os"] = cmd.Os
	}

	if cmd.Image != 0 {
		image, err := cmd.ImageManager.GetImage(cmd.Image)
		if err != nil {
			return nil, err
		}
		if image.GlobalIdentifier != nil {
			params["image"] = *image.GlobalIdentifier
		}
	}

	billing := cmd.Billing
	if billing == "hourly" {
		params["billing"] = true
	} else if billing == "monthly" {
		params["billing"] = false
	} else {
		return nil, slErrors.NewInvalidUsageError(T("[--billing] billing rate must be either hourly or monthly."))
	}

	params["dedicated"] = cmd.Dedicated

	if cmd.HostId != 0 {
		params["host-id"] = cmd.HostId
		params["dedicated"] = true
	}

	params["private"] = cmd.Private

	if cmd.San {
		params["san"] = true
	}

	if cmd.PostInstall != "" {
		params["i"] = cmd.PostInstall
	}

	if len(cmd.Key) > 0 {
		params["sshkeys"] = cmd.Key
	}

	if len(cmd.Disk) > 0 {
		params["disks"] = cmd.Disk
	}

	if cmd.Network > 0 {
		params["network"] = cmd.Network
	}

	if cmd.VlanPublic != 0 {
		params["vlan-public"] = cmd.VlanPublic
	}

	if cmd.VlanPrivate != 0 {
		params["vlan-private"] = cmd.VlanPrivate
	}

	if cmd.SubnetPublic != 0 {
		if cmd.VlanPublic == 0 {
			return nil, slErrors.NewMissingInputError("--vlan-public")
		}
		params["subnet-public"] = cmd.SubnetPublic
	}

	if cmd.SubnetPrivate != 0 {
		if cmd.VlanPrivate == 0 {
			return nil, slErrors.NewMissingInputError("--vlan-private")
		}
		params["subnet-private"] = cmd.SubnetPrivate
	}

	if len(cmd.PubSecGroup) > 0 {
		fmt.Printf("PubSecGroup(%v): %v\n", cmd.Domain, cmd.PubSecGroup)
		params["public-security-group"] = cmd.PubSecGroup
	}

	if len(cmd.PriSecGroup) > 0 {
		fmt.Printf("PriSecGroup: %v\n", cmd.PriSecGroup)
		params["private-security-group"] = cmd.PriSecGroup
	}

	if cmd.BootMode != "" {
		if cmd.BootMode != "HVM" && cmd.BootMode != "PV" {
			return nil, slErrors.NewInvalidUsageError("--boot-mode should be HVM | PV")
		}
		params["boot-mode"] = cmd.BootMode
	}

	if cmd.Transient {
		params["transient"] = cmd.Transient
	}

	if cmd.PlacementGroup != 0 {
		params["placement-group-id"] = cmd.PlacementGroup
	}
	return params, nil
}
