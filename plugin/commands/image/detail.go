package image

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("IDENTIFIER"),
		Short: T("Get details for an image"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	imageID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Image ID")
	}

	outputFormat := cmd.GetOutputFlag()

	image, err := cmd.ImageManager.GetImage(imageID)
	if err != nil {
		return bmxErr.NewAPIError(T("Failed to get image: {{.ImageID}}.\n", map[string]interface{}{"ImageID": imageID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, image)
	}

	diskspace := 0
	for _, child := range image.Children {
		if child.BlockDevicesDiskSpaceTotal != nil {
			diskspace += int(*child.BlockDevicesDiskSpaceTotal)
		}
	}
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(image.Id))
	table.Add(T("global_identifier"), utils.FormatStringPointer(image.GlobalIdentifier))
	table.Add(T("name"), utils.FormatStringPointer(image.Name))
	if image.Status != nil {
		table.Add(T("status"), utils.FormatStringPointer(image.Status.Name))
	} else {
		table.Add(T("status"), "-")
	}

	table.Add(T("account"), utils.FormatIntPointer(image.AccountId))

	if image.PublicFlag != nil && *image.PublicFlag == 1 {
		table.Add(T("visibility"), T("Public"))
	} else {
		table.Add(T("visibility"), T("Private"))
	}
	if image.ImageType != nil {
		table.Add(T("type"), utils.FormatStringPointer(image.ImageType.Name))
	} else {
		table.Add(T("type"), "-")
	}

	table.Add(T("flex"), utils.FormatBoolPointer(image.FlexImageFlag))
	table.Add(T("note"), utils.FormatStringPointer(image.Note))
	table.Add(T("created"), utils.FormatSLTimePointer(image.CreateDate))
	table.Add(T("disk_space"), utils.B2GB(diskspace))
	table.Add(T("datacenter"), "-------------------------------")
	for _, child := range image.Children {
		transaction := ""

		if child.Transaction != nil && child.Transaction.TransactionStatus != nil {
			transaction = fmt.Sprintf("(%s)", utils.FormatStringPointer(child.Transaction.TransactionStatus.Name))
		}
		message := fmt.Sprintf("%s %s", utils.B2GB(int(*child.BlockDevicesDiskSpaceTotal)), transaction)
		table.Add(utils.FormatStringPointer(child.Datacenter.Name), message)
	}
	table.Print()
	return nil
}
