package security_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Key add", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeSecurityManager *testhelpers.FakeSecurityManager
		cmd                 *security.KeyAddCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = security.NewKeyAddCommand(fakeUI, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        metadata.SecuritySSHKeyAddMetaData().Name,
			Description: metadata.SecuritySSHKeyAddMetaData().Description,
			Usage:       metadata.SecuritySSHKeyAddMetaData().Usage,
			Flags:       metadata.SecuritySSHKeyAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Key add", func() {
		Context("Key add without label", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Key add without key", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "label")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: either [-k] or [-f|--in-file] is required.")).To(BeTrue())
			})
		})
		Context("Key add with both key and keyfile", func() {
			It("return error", func() {
				keyFile := "/tmp/key.pub"
				err := testhelpers.RunCommand(cliCommand, "label", "-k", "ssh-rsa ndhd", "-f", keyFile)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-k] is not allowed with [-f|--in-file].")).To(BeTrue())
			})
		})
		// Context("Key add with failure to read keyfile", func() {
		// 	It("return error", func() {
		// 		keyFile := "/tmp/key.pub"
		// 		os.Remove(keyFile)
		// 		err := testhelpers.RunCommand(cliCommand, "label", "-f", keyFile)
		// 		Expect(err).To(HaveOccurred())
		// 		fmt.Println("ut error:" + err.Error())
		// 		Expect(strings.Contains(err.Error(), "Failed to read SSH key from file: "+keyFile)).To(BeTrue())
		// 	})
		// })
		Context("Key add with server fails", func() {
			BeforeEach(func() {
				fakeSecurityManager.AddSSHKeyReturns(datatypes.Security_Ssh_Key{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "label", "-k", "ssh-rsa ndhd")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to add SSH key.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("Key add ", func() {
			BeforeEach(func() {
				fakeSecurityManager.AddSSHKeyReturns(datatypes.Security_Ssh_Key{
					Id:          sl.Int(1234),
					Fingerprint: sl.String("37:87:03:ec:cd:b9:7e:fa:63:9c:83:21:d4:35:a4:ed"),
				}, nil)
			})
			// It("return no error", func() {
			// 	keyFile := "/tmp/key.pub"
			// 	ioutil.WriteFile(keyFile, []byte("ssh-rsa ndhd"), 0755)
			// 	err := testhelpers.RunCommand(cliCommand, "label", "-f", keyFile)
			// 	Expect(err).NotTo(HaveOccurred())
			// 	Expect(fakeUI.ContainsOutput("OK")).To(BeTrue())
			// 	Expect(fakeUI.ContainsOutput("SSH key was added: 37:87:03:ec:cd:b9:7e:fa:63:9c:83:21:d4:35:a4:ed")).To(BeTrue())
			// })
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "label", "-k", "ssh-rsa ndhd")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"SSH key was added: 37:87:03:ec:cd:b9:7e:fa:63:9c:83:21:d4:35:a4:ed"}))
			})
		})
	})
})
