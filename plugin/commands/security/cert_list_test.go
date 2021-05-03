package security_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/security"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Certiticate List", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.CertListCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewCertListCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSLCertListMetaData().Name,
			Description: metadata.SecuritySSLCertListMetaData().Description,
			Usage:       metadata.SecuritySSLCertListMetaData().Usage,
			Flags:       metadata.SecuritySSLCertListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Certiticate list", func() {
		Context("Certiticate list with wrong status", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--status", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--status] must be either all, valid or expired.")).To(BeTrue())
			})
		})

		Context("Certiticate list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abc is not supported.")).To(BeTrue())
			})
		})

		Context("Certiticate list but server API call fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.ListCertificatesReturns([]datatypes.Security_Certificate{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list SSL certificates on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "110")).To(BeTrue())
				Expect(strings.Contains(results[2], "123")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "common_name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "mon")).To(BeTrue())
				Expect(strings.Contains(results[2], "nom")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "days_until_expire")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "150")).To(BeTrue())
				Expect(strings.Contains(results[2], "365")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "note")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "Armer")).To(BeTrue())
				Expect(strings.Contains(results[2], "Docker")).To(BeTrue())
			})
		})
	})
})
