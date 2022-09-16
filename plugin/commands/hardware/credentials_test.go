package hardware_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware credentials", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.CredentialsCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		cmd = hardware.NewCredentialsCommand(fakeUI, fakeHardwareManager)
		cliCommand = cli.Command{
			Name:        hardware.HardwareCredentialsMetaData().Name,
			Description: hardware.HardwareCredentialsMetaData().Description,
			Usage:       hardware.HardwareCredentialsMetaData().Usage,
			Flags:       hardware.HardwareCredentialsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware credentials", func() {
		Context("hardware credentials without id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("hardware credentials with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("hardware credentials with server fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})

		Context("hardware credentials with no-complete response", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{Id: sl.Int(1234)}}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to find credentials of hardware server 1234."))
			})
		})

		Context("hardware credentials", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{
						Id: sl.Int(1234),
						OperatingSystem: &datatypes.Software_Component_OperatingSystem{
							Software_Component: datatypes.Software_Component{
								Passwords: []datatypes.Software_Component_Password{
									datatypes.Software_Component_Password{
										Username: sl.String("root"),
										Password: sl.String("MdZYMicl"),
									},
									datatypes.Software_Component_Password{
										Username: sl.String("user1"),
										Password: sl.String("pIzdjMvf3mE"),
									},
								},
							},
						},
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("root"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("MdZYMicl"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("user1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("pIzdjMvf3mE"))
			})
		})
	})
})
