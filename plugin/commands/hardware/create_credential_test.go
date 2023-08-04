package hardware_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hadware create-credential", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.CreateCredentialCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewCreateCredentialCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware create-credential", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage : This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde", "--username=myusername", "--password=password1234", "--software=ubuntu")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--username=myusername", "--password=password1234", "--software=ubuntu", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set command without required options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`required flag(s) "password", "software", "username" not set`))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetHardwareReturns(datatypes.Hardware_Server{}, errors.New("Failed to get hardware server"))
			})
			It("Failed get hardware sensor data", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--username=myusername", "--password=password1234", "--software=ubuntu")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerHardware := datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{
						SoftwareComponents: []datatypes.Software_Component{
							datatypes.Software_Component{
								Id: sl.Int(111111),
								SoftwareLicense: &datatypes.Software_License{
									SoftwareDescription: &datatypes.Software_Description{
										Name: sl.String("Ubuntu"),
									},
								},
							},
						},
					},
				}
				fakeHardwareManager.GetHardwareReturns(fakerHardware, nil)
			})
			It("Failed find software", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--username=myusername", "--password=password1234", `--software="Red Hat"`)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Software not found"))
			})

			It("Failed create Software Credential", func() {
				fakeHardwareManager.CreateSoftwareCredentialReturns(datatypes.Software_Component_Password{}, errors.New("Failed to create Software Credential"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--username=myusername", "--password=password1234", "--software=ubuntu")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create Software Credential"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerHardware := datatypes.Hardware_Server{
					Hardware: datatypes.Hardware{
						SoftwareComponents: []datatypes.Software_Component{
							datatypes.Software_Component{
								Id: sl.Int(111111),
								SoftwareLicense: &datatypes.Software_License{
									SoftwareDescription: &datatypes.Software_Description{
										Name: sl.String("Ubuntu"),
									},
								},
							},
						},
					},
				}
				fakerCredential := datatypes.Software_Component_Password{
					Id:         sl.Int(222222),
					CreateDate: sl.Time(created),
					Username:   sl.String("myusername"),
					Password:   sl.String("password1234"),
					Notes:      sl.String("mynote"),
				}
				fakeHardwareManager.CreateSoftwareCredentialReturns(fakerCredential, nil)
				fakeHardwareManager.GetHardwareReturns(fakerHardware, nil)
			})
			It("display credential", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--username=myusername", "--password=password1234", "--software=ubuntu", "--notes=mynote")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("myusername"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("password1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("mynote"))
			})
		})
	})
})
