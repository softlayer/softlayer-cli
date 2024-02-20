package securitygroup_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup create", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cliCommand         *securitygroup.CreateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = securitygroup.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("Securitygroup create", func() {
		Context("create with server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.CreateSecurityGroupReturns(datatypes.Network_SecurityGroup{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create security group with name test."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("create succeed", func() {
			BeforeEach(func() {
				fakeNetworkManager.CreateSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id:   sl.Int(1234),
					Name: sl.String("test"),
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test"))
			})
		})
		Context("create succeed", func() {
			BeforeEach(func() {
				fakeNetworkManager.CreateSecurityGroupReturns(datatypes.Network_SecurityGroup{
					Id:          sl.Int(1234),
					Name:        sl.String("test"),
					Description: sl.String("test-desc"),
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "test", "-d", "test-desc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test-desc"))
			})
		})
	})
})
