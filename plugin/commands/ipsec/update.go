package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type UpdateCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewUpdateCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *UpdateCommand) {
	return &UpdateCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *UpdateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	if c.IsSet("a") {
		phase1Auth := c.String("a")
		if phase1Auth != "MD5" && phase1Auth != "SHA1" && phase1Auth != "SHA256" {
			return errors.NewInvalidUsageError(T("-a|--phase1-auth must be either MD5 or SHA1 or SHA256."))
		}
	}
	if c.IsSet("c") {
		phase1Crypto := c.String("c")
		if phase1Crypto != "DES" && phase1Crypto != "3DES" && phase1Crypto != "AES128" && phase1Crypto != "AES192" && phase1Crypto != "AES256" {
			return errors.NewInvalidUsageError(T("-c|--phase1-crypto must be either DES or 3DES or AES128 or AES192 or AES256."))
		}
	}
	if c.IsSet("d") {
		phase1Dh := c.Int("d")
		if phase1Dh != 0 && phase1Dh != 1 && phase1Dh != 2 && phase1Dh != 5 {
			return errors.NewInvalidUsageError(T("-d|--phase1-dh must be either 0 or 1 or 2 or 5."))
		}
	}
	if c.IsSet("t") {
		phase1KeyLife := c.Int("t")
		if phase1KeyLife < 120 || phase1KeyLife > 172800 {
			return errors.NewInvalidUsageError(T("-t|--phase1-key-ttl must be in range 120-172800."))
		}
	}
	if c.IsSet("u") {
		phase2Auth := c.String("u")
		if phase2Auth != "MD5" && phase2Auth != "SHA1" && phase2Auth != "SHA256" {
			return errors.NewInvalidUsageError(T("-u|--phase2-auth must be either MD5 or SHA1 or SHA256."))
		}
	}
	if c.IsSet("y") {
		phase2Crypto := c.String("y")
		if phase2Crypto != "DES" && phase2Crypto != "3DES" && phase2Crypto != "AES128" && phase2Crypto != "AES192" && phase2Crypto != "AES256" {
			return errors.NewInvalidUsageError(T("-y|--phase2-crypto must be either DES or 3DES or AES128 or AES192 or AES256."))
		}
	}
	if c.IsSet("e") {
		phase2Dh := c.Int("e")
		if phase2Dh != 0 && phase2Dh != 1 && phase2Dh != 2 && phase2Dh != 5 {
			return errors.NewInvalidUsageError(T("-e|--phase2-dh must be either 0 or 1 or 2 or 5."))
		}
	}
	if c.IsSet("f") {
		phase2ForwardSecrecy := c.Int("f")
		if phase2ForwardSecrecy != 0 && phase2ForwardSecrecy != 1 {
			return errors.NewInvalidUsageError(T("-f|--phase2-forward-secrecy must be either 0 or 1."))
		}
	}
	if c.IsSet("l") {
		phase2KeyLife := c.Int("l")
		if phase2KeyLife < 120 || phase2KeyLife > 172800 {
			return errors.NewInvalidUsageError(T("-l|--phase2-key-ttl must be in range 120-172800."))
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.IPSECManager.UpdateTunnelContext(contextId,
		c.String("n"),
		c.String("r"),
		c.String("k"),
		c.String("a"),
		c.String("c"),
		c.Int("d"),
		c.Int("t"),
		c.String("u"),
		c.String("y"),
		c.Int("e"),
		c.Int("f"),
		c.Int("l"))
	if err != nil {
		return cli.NewExitError(T("Failed to update IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Updated IPSec {{.ContextID}}.", map[string]interface{}{"ContextID": contextId}))
	return nil
}
