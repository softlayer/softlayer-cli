package loadbal_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer create", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.CreateCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager

		fakeLBManager.CreateLoadBalancerReturns(datatypes.Container_Product_Order_Receipt{
			OrderId: sl.Int(123456),
			OrderDetails: &datatypes.Container_Product_Order{
				Prices: []datatypes.Product_Item_Price{
					datatypes.Product_Item_Price{
						Item: &datatypes.Product_Item{
							Description: sl.String("Item Description"),
						},
						HourlyRecurringFee: sl.Float(23.5),
					},
				},
			},
		}, nil)

		fakeLBManager.CreateLoadBalancerVerifyReturns(datatypes.Container_Product_Order{
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Item: &datatypes.Product_Item{
						Description: sl.String("Item Verify Description"),
					},
					HourlyRecurringFee: sl.Float(24.7),
				},
			},
		}, nil)
	})

	Context("create without name", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-n, --name' is required")).To(BeTrue())
		})
	})
	Context("create without datacenter", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "name")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-d, --datacenter' is required")).To(BeTrue())
		})
	})
	Context("create without type", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "name", "-d", "dal09")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-t, --type' is required"))
		})
	})
	Context("create with incorrect type", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "name", "-d", "dal09", "-t", "abcd")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: The value of option '-t, --type' should be PublicToPrivate | PrivateToPrivate | PublicToPublic"))
		})
	})
	Context("create without subnet", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-s, --subnet' is required"))
		})
	})
	Context("create without confirmation", func() {
		It("return aborted", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted"}))
		})
	})
	Context("create with confirmation error", func() {
		It("return error", func() {
			fakeUI.Inputs("123456")
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123")
			Expect(err).To(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
			Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
		})
	})
	Context("create with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-f")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("create load balancer with incorrect sticky", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-c", "3", "--sticky", "abcd", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Value of option '--sticky' should be cookie or source-ip"))
		})
	})
	Context("create with correct load balancer ID", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-c", "3", "--sticky", "cookie", "-f")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Order ID: 123456"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Item               Cost"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Item Description   23.500000"))
		})
	})
	Context("verify create load balancer", func() {
		It("return no error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-c", "3", "--sticky", "source-ip", "--verify")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Item                      Cost"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("Item Verify Description   24.700000"))
		})
	})
	Context("verify create with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerVerifyReturns(datatypes.Container_Product_Order{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-c", "3", "--sticky", "source-ip", "--verify")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to verify load balancer with name 1234 on dal09."))
			Expect(err.Error()).To(ContainSubstring("Internal server error"))
		})
	})
	Context("create load balancer with type publictopublic and subnet different to 0", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "publictopublic", "-d", "dal09", "-s", "123", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -s, --subnet is only available in PublicToPrivate and PrivateToPrivate load balancer type."))
		})
	})
	Context("create load balancer with type privatetoprivate and public subnet", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "-n", "1234", "-t", "privatetoprivate", "-d", "dal09", "-s", "123", "--use-public-subnet", "-f")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --use-public-subnet is only available in PublicToPrivate."))
		})
	})
})
