package user

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ApikeyCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Add         bool
	Remove      bool
	Refresh     bool
}

func NewApikeyCommand(sl *metadata.SoftlayerCommand) (cmd *ApikeyCommand) {
	thisCmd := &ApikeyCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "apikey " + T("IDENTIFIER"),
		Short: T("Allows to create, remove or refresh user's API authentication key"),
		Long:  T("Each user can only have a single API key."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Add, "add", false, T("Create an user's API authentication key"))
	cobraCmd.Flags().BoolVar(&thisCmd.Remove, "remove", false, T("Remove an user's API authentication key"))
	cobraCmd.Flags().BoolVar(&thisCmd.Refresh, "refresh", false, T("Refresh an user's API authentication key"))

	cobraCmd.MarkFlagsMutuallyExclusive("add", "remove", "refresh")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ApikeyCommand) Run(args []string) error {
	if !cmd.Add && !cmd.Remove && !cmd.Refresh {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if cmd.Remove || cmd.Refresh {
		_, err := cmd.UserManager.RemoveApiAuthenticationKey(userId)
		if err != nil {
			return errors.NewAPIError(T("Failed to remove user's API authentication key"), err.Error(), 2)
		}

		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully removed user's API authentication key"))
	}

	if cmd.Refresh || cmd.Add {
		apiAuthenticationKey, err := cmd.UserManager.AddApiAuthenticationKey(userId)
		if err != nil {
			return errors.NewAPIError(T("Failed to add user's API authentication key"), err.Error(), 2)
		}
		i18nsubs := map[string]interface{}{"apiAuthenticationKey": apiAuthenticationKey}
		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully added. New API Authentication Key: {{.apiAuthenticationKey}}", i18nsubs))
	}

	return nil
}
