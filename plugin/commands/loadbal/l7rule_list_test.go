package loadbal_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer edit policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.L7RuleListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7RuleListCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7RuleListMetadata().Name,
			Description: loadbal.LoadbalL7RuleListMetadata().Description,
			Usage:       loadbal.LoadbalL7RuleListMetadata().Usage,
			Flags:       loadbal.LoadbalL7RuleListMetadata().Flags,
			Action:      cmd.Run,
		}

		fakeLBManager.ListL7RuleReturns([]datatypes.Network_LBaaS_L7Rule{
			datatypes.Network_LBaaS_L7Rule{
				Id:             sl.Int(123),
				Uuid:           sl.String("abc123"),
				Type:           sl.String("Type"),
				ComparisonType: sl.String("ComparisonType"),
				Value:          sl.String("value"),
				Key:            sl.String("abcd1234"),
				Invert:         sl.Int(5),
			},
			datatypes.Network_LBaaS_L7Rule{
				Id:             sl.Int(1234),
				Uuid:           sl.String("abc1234"),
				Type:           sl.String("Type2"),
				ComparisonType: sl.String("ComparisonType2"),
				Value:          sl.String("value2"),
				Key:            sl.String("abcd12345"),
				Invert:         sl.Int(6),
			},
		}, nil)
	})

	Describe("l7 rule list", func() {
		Context("l7 rule list, missing arguments error", func() {
			It("policy-id is required", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-id' is required"))
			})
		})

		Context("l7 rule list", func() {
			It("with all attributes", func() {
				err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID     UUID      Type    Compare Type      Value    Key         Invert"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123    abc123    Type    ComparisonType    value    abcd1234    5"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234   abc1234   Type2   ComparisonType2   value2   abcd12345   6"))
			})
		})

		Context("errors", func() {
			It("Failed to list l7 rule", func() {
				fakeLBManager.ListL7RuleReturns([]datatypes.Network_LBaaS_L7Rule{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get l7 rules: Internal server error"))
			})
			It("No l7 rules was found", func() {
				fakeLBManager.ListL7RuleReturns([]datatypes.Network_LBaaS_L7Rule{}, nil)
				err := testhelpers.RunCommand(cliCommand, "--policy-id", "123")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No l7 rules was found"))
			})
		})
	})
})
