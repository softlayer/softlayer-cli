package loadbal_test

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
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Load balancer create", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeLBManager *testhelpers.FakeLoadBalancerManager
		cmd           *loadbal.CreateCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewCreateCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        metadata.LoadbalOrderMetadata().Name,
			Description: metadata.LoadbalOrderMetadata().Description,
			Usage:       metadata.LoadbalOrderMetadata().Usage,
			Flags:       metadata.LoadbalOrderMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("create without name", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-n, --name' is required")).To(BeTrue())
		})
	})
	Context("create with wrong subnet id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "-n", "name")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: '-d, --datacenter' is required")).To(BeTrue())
		})
	})
	Context("create without confirmation", func() {
		It("return aborted", func() {
			fakeUI.Inputs("No")
			err := testhelpers.RunCommand(cliCommand, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted"}))
		})
	})
	Context("create with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-f")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("create with correct load balancer ID", func() {
		BeforeEach(func() {
			fakeLBManager.CreateLoadBalancerReturns(datatypes.Container_Product_Order_Receipt{
				OrderId: sl.Int(12345678),
			}, nil)
		})
		It("return no error", func() {
			err := testhelpers.RunCommand(cliCommand, "-n", "1234", "-t", "publictoprivate", "-d", "dal09", "-s", "123", "-f")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
		})
	})
})
