package dns_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.com/softlayer/softlayer-go/session"

	. "github.com/onsi/ginkgo"
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
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Zone print with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Unable to find zone abc."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get zone ID from zone name: abc.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Unable to find zone abc.")).To(BeTrue())
			})
		})

		Context("Zone print with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.DumpZoneReturns("", errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to dump content for zone: abc.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				Expect(strings.Contains(results[1], "$ORIGIN dal06.bluemix.ibmcsf.net.")).To(BeTrue())
				Expect(strings.Contains(results[2], "$TTL 900")).To(BeTrue())
				Expect(strings.Contains(results[3], "@ IN SOA ns1.softlayer.com. support.softlayer.com. (")).To(BeTrue())
				Expect(strings.Contains(results[4], "2014121600        ; Serial")).To(BeTrue())
				Expect(strings.Contains(results[5], "7200              ; Refresh")).To(BeTrue())
				Expect(strings.Contains(results[6], "600               ; Retry")).To(BeTrue())
				Expect(strings.Contains(results[7], "1728000           ; Expire")).To(BeTrue())
				Expect(strings.Contains(results[8], "43200)            ; Minimum")).To(BeTrue())
				Expect(strings.Contains(results[10], "@                      900      IN NS    ns1.softlayer.com.")).To(BeTrue())
				Expect(strings.Contains(results[11], "@                      900      IN NS    ns2.softlayer.com.")).To(BeTrue())
				Expect(strings.Contains(results[13], "@                      900      IN MX 10 mail.dal06.bluemix.ibmcsf.net.")).To(BeTrue())
				Expect(strings.Contains(results[15], "txt                    900      IN TXT   bcr01.dal06.bluemix.ibmcsf.net")).To(BeTrue())
				Expect(strings.Contains(results[16], "@                      900      IN A     127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[17], "ftp                    86400    IN A     127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[18], "mail                   86400    IN A     127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[19], "webmail                86400    IN A     127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[20], "www                    86400    IN A     127.0.0.1")).To(BeTrue())
			})
		})
	})
})
