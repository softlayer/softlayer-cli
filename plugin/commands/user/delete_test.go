package user_test

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Delete", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.DeleteCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewDeleteCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserDeleteMataData().Name,
			Description: user.UserDeleteMataData().Description,
			Usage:       user.UserDeleteMataData().Usage,
			Flags:       user.UserDeleteMataData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("user delete", func() {
		Context("user create with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument")).To(BeTrue())
			})
		})
	})

})
