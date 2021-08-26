package order_test

import (
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Place", func() {
	var (
		fakeUI        *terminal.FakeUI
		cmd           *order.PlaceQuoteCommand
		cliCommand    cli.Command
		fakeSLSession *session.Session
		OrderManager  managers.OrderManager
	)
	BeforeEach(func() {

		filenames := []string{"getDatacenters_1",}
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
		OrderManager = managers.NewOrderManager(fakeSLSession)
		fakeUI = terminal.NewFakeUI()
		cmd = order.NewPlaceQuoteCommand(fakeUI, OrderManager, nil)
		cliCommand = cli.Command{
			Name:        metadata.OrderPlaceQuoteMetaData().Name,
			Description: metadata.OrderPlaceQuoteMetaData().Description,
			Usage:       metadata.OrderPlaceQuoteMetaData().Usage,
			Flags:       metadata.OrderPlaceQuoteMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("order create", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--name", "test", "--send-email")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2523413"}))
				})

			})
		}
	})
	Describe("order create", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2523413"}))
				})

			})
		}
	})
})
