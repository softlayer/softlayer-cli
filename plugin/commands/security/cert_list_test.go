package security_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certiticate List", func() {
	var (
		fakeUI              *terminal.FakeUI
		cliCommand          *security.CertListCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
		fakeSecurityManager *testhelpers.FakeSecurityManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = security.NewCertListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cliCommand.SecurityManager = fakeSecurityManager
	})

	Describe("Certiticate list", func() {
		Context("Certiticate list with wrong status", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--status", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [--status] must be either all, valid or expired."))
			})
		})

		Context("Certiticate list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("Certiticate list but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListCertificatesReturns([]datatypes.Security_Certificate{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list SSL certificates on your account"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Certiticate list with different sortby", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListCertificatesReturns([]datatypes.Security_Certificate{
					datatypes.Security_Certificate{
						Id:           sl.Int(123),
						CommonName:   sl.String("mon"),
						ValidityDays: sl.Int(365),
						Notes:        sl.String("Docker"),
					},
					datatypes.Security_Certificate{
						Id:           sl.Int(110),
						CommonName:   sl.String("nom"),
						ValidityDays: sl.Int(150),
						Notes:        sl.String("Armer"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("110"))
				Expect(results[2]).To(ContainSubstring("123"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "common_name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("mon"))
				Expect(results[2]).To(ContainSubstring("nom"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "days_until_expire")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("150"))
				Expect(results[2]).To(ContainSubstring("365"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "note")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("Armer"))
				Expect(results[2]).To(ContainSubstring("Docker"))
			})
		})
	})
})
