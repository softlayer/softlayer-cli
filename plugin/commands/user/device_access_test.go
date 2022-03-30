package user_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Device access", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cmd             *user.DeviceAccessCommand
		cliCommand      cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		cmd = user.NewDeviceAccessCommand(fakeUI, fakeUserManager)
		cliCommand = cli.Command{
			Name:        user.UserDeviceAccessMetaData().Name,
			Description: user.UserDeviceAccessMetaData().Description,
			Usage:       user.UserDeviceAccessMetaData().Usage,
			Flags:       user.UserDeviceAccessMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("user device-access", func() {
		Context("Return error", func() {
			It("Set command without identifier", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid identifier", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetUserAllowDevicesPermissionsReturns([]datatypes.User_Customer_CustomerPermission_Permission{}, errors.New("Internal Server Error"))
			})
			It("failed get permissions", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get user permissions."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetDedicatedHostsReturns([]datatypes.Virtual_DedicatedHost{}, errors.New("Internal Server Error"))
			})
			It("failed get dedicated hosts", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get dedicated hosts."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetHardwareReturns([]datatypes.Hardware{}, errors.New("Internal Server Error"))
			})
			It("failed get bare metal servers", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get bare metal servers."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeUserManager.GetVirtualGuestsReturns([]datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("failed get virtual servers", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual servers."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerPermissions := []datatypes.User_Customer_CustomerPermission_Permission{
					datatypes.User_Customer_CustomerPermission_Permission{
						KeyName: sl.String("ACCESS_ALL_GUEST"),
						Name:    sl.String("All Guest Access"),
					},
				}
				fakerdedicatedHosts := []datatypes.Virtual_DedicatedHost{
					datatypes.Virtual_DedicatedHost{
						Id:    sl.Int(333333),
						Name:  sl.String("myDedicatedHost"),
						Notes: sl.String("my dedicated notes"),
					},
				}
				fakerHardares := []datatypes.Hardware{
					datatypes.Hardware{
						Id:                       sl.Int(111111),
						FullyQualifiedDomainName: sl.String("hardware.mydomain.com"),
						PrimaryIpAddress:         sl.String("11.11.11.11"),
						PrimaryBackendIpAddress:  sl.String("10.10.10.11"),
						Notes:                    sl.String("my hardware notes"),
					},
				}
				fakerVirtualGuests := []datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:                       sl.Int(222222),
						FullyQualifiedDomainName: sl.String("virtual.mydomain.com"),
						PrimaryIpAddress:         sl.String("12.12.12.12"),
						PrimaryBackendIpAddress:  sl.String("10.10.10.12"),
						Notes:                    sl.String("my virtual guests notes"),
					},
				}
				fakeUserManager.GetUserAllowDevicesPermissionsReturns(fakerPermissions, nil)
				fakeUserManager.GetDedicatedHostsReturns(fakerdedicatedHosts, nil)
				fakeUserManager.GetHardwareReturns(fakerHardares, nil)
				fakeUserManager.GetVirtualGuestsReturns(fakerVirtualGuests, nil)
			})
			It("Set command with a valid user", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACCESS_ALL_GUEST"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("All Guest Access"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("333333"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myDedicatedHost"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("my dedicated notes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hardware.mydomain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("11.11.11.11"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.10.10.11"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("my hardware notes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual.mydomain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12.12.12.12"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("10.10.10.12"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("my virtual guests notes"))
			})
		})

		Context("Return no error", func() {
			It("User set does not have devices and permissions", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
			})
		})
	})
})
