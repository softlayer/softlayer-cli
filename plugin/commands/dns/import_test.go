package dns_test

import (
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("DNS Import", func() {

	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.ImportCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewImportCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        metadata.DnsImportMetaData().Name,
			Description: metadata.DnsImportMetaData().Description,
			Usage:       metadata.DnsImportMetaData().Usage,
			Flags:       metadata.DnsImportMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("DNS import", func() {
		Context("DNS import without file", func() {
			It("without any argument", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("with an inexist file", func() {
				err := testhelpers.RunCommand(cliCommand, "not-exist.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to read file: not-exist.txt."))
				Expect(err.Error()).To(ContainSubstring("open not-exist.txt: The system cannot find the file specified."))
			})
		})

		Context("DNS send a file import", func() {

			dirFile := os.TempDir() + "file.txt"

			It("send a empty file", func() {
				file, _ := os.Create(dirFile)
				err := testhelpers.RunCommand(cliCommand, file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse file."))
				Expect(err.Error()).To(ContainSubstring("Unable to parse zone from BIND file."))
			})

			content := `$ORIGIN`
			It("no send a TTL in the file", func() {
				file, _ := os.Create(dirFile)
				file.WriteString(content)
				err := testhelpers.RunCommand(cliCommand, file.Name())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to parse file."))
				Expect(err.Error()).To(ContainSubstring("dns: not a TTL: \"$ORIGIN\" at line: 1:7"))
			})

			content_ := `$ORIGIN dal06.bluemix.ibmcsf.net.
$TTL 900
@ IN SOA ns1.softlayer.com. support.softlayer.com. (
					2014121600        ; Serial
					7200              ; Refresh
					600               ; Retry
					1728000           ; Expire
					43200)            ; Minimum

@                      900      IN NS    ns1.softlayer.com.
@                      900      IN NS    ns2.softlayer.com.

@                      900      IN MX 10 mail.dal06.bluemix.ibmcsf.net.

txt                    900      IN TXT   bcr01.dal06.bluemix.ibmcsf.net
@                      900      IN A     127.0.0.1
ftp                    86400    IN A     127.0.0.1
mail                   86400    IN A     127.0.0.1
webmail                86400    IN A     127.0.0.1
www                    86400    IN A     127.0.0.1
`
			It("send a good file with --dry-run argument", func() {
				file, _ := os.Create(dirFile)
				file.WriteString(content_)
				err := testhelpers.RunCommand(cliCommand, file.Name(), "--dry-run")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})

			content__ := `$ORIGIN dal06.bluemix.ibmcsf.net.
$TTL 900
@ IN SOA ns1.softlayer.com. support.softlayer.com. (
					2014121600        ; Serial
					7200              ; Refresh
					600               ; Retry
					1728000           ; Expire
					43200)            ; Minimum
@                      900      IN NS    ns1.softlayer.com.
@                      900      IN NS    ns2.softlayer.com.
`
			It("send a good file", func() {
				file, _ := os.Create(dirFile)
				file.WriteString(content__)
				err := testhelpers.RunCommand(cliCommand, file.Name())
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Zone  was created."))
			})
			os.Remove(dirFile)
		})
	})
})
