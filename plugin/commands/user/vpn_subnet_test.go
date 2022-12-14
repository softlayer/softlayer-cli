package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("user vpn-subnet", func() {
	var (
		fakeUI          *terminal.FakeUI
		cliCommand      *user.VpnSubnetCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
		fakeUserManager *testhelpers.FakeUserManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeUserManager = new(testhelpers.FakeUserManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewVpnSubnetCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager
	})

	Describe("user vpn-subnet", func() {

		Context("Return error", func() {
			It("Set command without Arguments", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires two arguments"))
			})

			It("Set command with an invalid user Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "222222", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'User ID'. It must be a positive integer."))
			})

			It("Set command with an invalid subnet Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "abcde", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Subnet ID'. It must be a positive integer."))
			})

			It("Set without any option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires --add or --remove option"))
			})

			It("Set with both options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--add", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--add', '--remove' are exclusive"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeUserManager.CreateUserVpnOverrideReturns(true, nil)
				fakeUserManager.UpdateVpnUserReturns(true, nil)
			})
			It("Add subnet", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--add")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully added subnet access for user."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.CreateUserVpnOverrideReturns(false, errors.New("Failed to create user vpn override"))
			})
			It("Failed to create user vpn override", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create user vpn override"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.CreateUserVpnOverrideReturns(true, nil)
				fakeUserManager.UpdateVpnUserReturns(false, errors.New("Override created, but unable to update VPN user"))
			})
			It("Failed to update VPN user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Override created, but unable to update VPN user"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetOverridesReturns([]datatypes.Network_Service_Vpn_Overrides{}, errors.New("Failed to get user vpn overrides"))
			})
			It("Failed to get overrides", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user vpn overrides"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetOverridesReturns([]datatypes.Network_Service_Vpn_Overrides{}, nil)
			})
			It("Subnet is not assigned", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is not assigned"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerOverrides := []datatypes.Network_Service_Vpn_Overrides{
					datatypes.Network_Service_Vpn_Overrides{
						SubnetId: sl.Int(222222),
						Id:       sl.Int(123456),
					},
				}
				fakeUserManager.GetOverridesReturns(fakerOverrides, nil)
				fakeUserManager.DeleteUserVpnOverrideReturns(false, errors.New("Failed to delete user vpn override"))
			})
			It("Failed to delete override", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to delete user vpn override"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerOverrides := []datatypes.Network_Service_Vpn_Overrides{
					datatypes.Network_Service_Vpn_Overrides{
						SubnetId: sl.Int(222222),
						Id:       sl.Int(123456),
					},
				}
				fakeUserManager.GetOverridesReturns(fakerOverrides, nil)
				fakeUserManager.DeleteUserVpnOverrideReturns(true, nil)
				fakeUserManager.UpdateVpnUserReturns(false, errors.New("Failed to update VPN user"))
			})
			It("Failed to update VPN user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--remove")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update VPN user"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerOverrides := []datatypes.Network_Service_Vpn_Overrides{
					datatypes.Network_Service_Vpn_Overrides{
						SubnetId: sl.Int(222222),
						Id:       sl.Int(123456),
					},
				}
				fakeUserManager.GetOverridesReturns(fakerOverrides, nil)
				fakeUserManager.DeleteUserVpnOverrideReturns(true, nil)
				fakeUserManager.UpdateVpnUserReturns(true, nil)
			})
			It("Remove subnet", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "111111", "222222", "--remove")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully removed subnet access for user."))
			})
		})
	})
})
