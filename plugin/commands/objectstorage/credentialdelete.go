package objectstorage

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CredentialDeleteCommand struct {
	UI                   terminal.UI
	ObjectStorageManager managers.ObjectStorageManager
}

func NewCredentialDeleteCommand(ui terminal.UI, objectStorageManager managers.ObjectStorageManager) (cmd *CredentialDeleteCommand) {
	return &CredentialDeleteCommand{
		UI:                   ui,
		ObjectStorageManager: objectStorageManager,
	}
}

func CredentialDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    "object-storage",
		Name:        "credential-delete",
		Description: T("Delete the credential of an Object Storage Account."),
		Usage: T(`${COMMAND_NAME} sl object-storage credential-delete IDENTIFIER [OPTIONS]

Examples:
	${COMMAND_NAME} sl object-storage credential-delete ObjectStorageID --credential_id CredentialID
	${COMMAND_NAME} sl object-storage credential-delete 123456 --credential_id 654321`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "credential_id",
				Usage: T("This is the credential id associated with the volume. [REQUIRED]"),
			},
		},
	}
}

func (cmd *CredentialDeleteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	StorageID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Storage ID")
	}

	if !c.IsSet("credential_id") {
		return slErr.NewMissingInputError("--credential_id")
	}

	CredentialID := c.Int("credential_id")

	err = cmd.ObjectStorageManager.DeleteCredential(StorageID, CredentialID)
	if err != nil {
		if strings.Contains(err.Error(), slErr.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find object-storage with ID: {{.StorageID}}.\n", map[string]interface{}{"StorageID": StorageID})+err.Error(), 0)
		}
		if strings.Contains(err.Error(), "ObjectNotFound") {
			return cli.NewExitError(T("Unable to find credential with ID: {{.CredentialID}}.\n", map[string]interface{}{"CredentialID": CredentialID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to delete credential: {{.StorageID}}.\n", map[string]interface{}{"StorageID": StorageID})+err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Credential: {{.CredentialID}} was deleted.", map[string]interface{}{"CredentialID": CredentialID}))
	return nil
}
