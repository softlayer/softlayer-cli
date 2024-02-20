package dns_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.com/softlayer/softlayer-go/session"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Zone print", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.ZonePrintCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewZonePrintCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Zone print", func() {
		Context("zone print without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("Zone print with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Unable to find zone abc."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get zone ID from zone name: abc."))
				Expect(err.Error()).To(ContainSubstring("Unable to find zone abc."))
			})
		})

		Context("Zone print with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.DumpZoneReturns("", errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to dump content for zone: abc."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Zone print", func() {
			BeforeEach(func() {
				fakeDNSManager.DumpZoneReturns(
					`
$ORIGIN dal06.bluemix.ibmcsf.net.
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
`, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("$ORIGIN dal06.bluemix.ibmcsf.net."))
				Expect(results[2]).To(ContainSubstring("$TTL 900"))
				Expect(results[3]).To(ContainSubstring("@ IN SOA ns1.softlayer.com. support.softlayer.com. ("))
				Expect(results[4]).To(ContainSubstring("2014121600        ; Serial"))
				Expect(results[5]).To(ContainSubstring("7200              ; Refresh"))
				Expect(results[6]).To(ContainSubstring("600               ; Retry"))
				Expect(results[7]).To(ContainSubstring("1728000           ; Expire"))
				Expect(results[8]).To(ContainSubstring("43200)            ; Minimum"))
				Expect(results[10]).To(ContainSubstring("@                      900      IN NS    ns1.softlayer.com."))
				Expect(results[11]).To(ContainSubstring("@                      900      IN NS    ns2.softlayer.com."))
				Expect(results[13]).To(ContainSubstring("@                      900      IN MX 10 mail.dal06.bluemix.ibmcsf.net."))
				Expect(results[15]).To(ContainSubstring("txt                    900      IN TXT   bcr01.dal06.bluemix.ibmcsf.net"))
				Expect(results[16]).To(ContainSubstring("@                      900      IN A     127.0.0.1"))
				Expect(results[17]).To(ContainSubstring("ftp                    86400    IN A     127.0.0.1"))
				Expect(results[18]).To(ContainSubstring("mail                   86400    IN A     127.0.0.1"))
				Expect(results[19]).To(ContainSubstring("webmail                86400    IN A     127.0.0.1"))
				Expect(results[20]).To(ContainSubstring("www                    86400    IN A     127.0.0.1"))
			})
		})
	})
})
