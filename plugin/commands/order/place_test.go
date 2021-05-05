package order_test

import (
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
)

var _ = Describe("Place", func() {
	var (
		fakeUI        *terminal.FakeUI
		cmd           *order.PlaceCommand
		cliCommand    cli.Command
		fakeSLSession *session.Session
		OrderManager  managers.OrderManager
	)
	BeforeEach(func() {

		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		OrderManager = managers.NewOrderManager(fakeSLSession)
		fakeUI = terminal.NewFakeUI()
		cmd = order.NewPlaceCommand(fakeUI, OrderManager, nil)
		cliCommand = cli.Command{
			Name:        metadata.OrderPlaceMetaData().Name,
			Description: metadata.OrderPlaceMetaData().Description,
			Usage:       metadata.OrderPlaceMetaData().Usage,
			Flags:       metadata.OrderPlaceMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("order verify", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "--verify")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4_PORTABLE_PUBLIC_IP_ADDRESSES"}))
				})

			})
		}
	})
	Describe("order create", func() {
		for k, _ := range order.TYPEMAP {
			Context("successfully"+k, func() {

				k := k
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand, "CLOUD_SERVER", "dal13", "EVAULT_100_GB,CITRIX_VDC", "--complex-type", k, "-f")
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"11493593"}))
				})

			})
		}
	})
})
