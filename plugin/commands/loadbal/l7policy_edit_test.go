package loadbal_test

import (
	"errors"
	"strings"
	"time"

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
		cmd           *loadbal.L7PolicyEditCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cmd = loadbal.NewL7PolicyEditCommand(fakeUI, fakeLBManager)
		cliCommand = cli.Command{
			Name:        loadbal.LoadbalL7PolicyEditMetadata().Name,
			Description: loadbal.LoadbalL7PolicyEditMetadata().Description,
			Usage:       loadbal.LoadbalL7PolicyEditMetadata().Usage,
			Flags:       loadbal.LoadbalL7PolicyEditMetadata().Flags,
			Action:      cmd.Run,
		}
	})

	Context("CLI Usage Errors", func() {
		It("Error No policy-id", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--policy-id"))
		})
		It("No valid action", func() {
			command := "--policy-id 12345 -n test-name -a unknown-action"
			command_args := strings.Fields(command)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(
				ContainSubstring("-a, --action should be REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error invalid usage for REJECT", func() {
			command := "--policy-id 12345 -n test-name -a REJECT -r REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is only available with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
		It("Error No --redirect", func() {
			command := "--policy-id 12345 -n test-name -a REDIRECT_URL"
			command_args := strings.Fields(command)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"-r, --redirect is required with action REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"))
		})
	})

	Context("API Error", func() {
		It("Handles API Error", func() {
			command := "--policy-id 12345 -n test-name -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to edit l7 policy"))
		})
		It("Handles API Error", func() {
			command := "--policy-id 12345 -n test-name -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.GetL7PolicyReturns(datatypes.Network_LBaaS_L7Policy{}, errors.New("SL_API_ERROR"))
			err := testhelpers.RunCommand(cliCommand, command_args...)
			println(err)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get l7 policy"))
		})
	})

	Context("CLI Usage", func() {
		It("REJECT", func() {
			command := "--policy-id 12345 --n test-reject -a REJECT"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_POOL", func() {
			command := "--policy-id 12345 --n test-pool -a REDIRECT_POOL -r uuid-pool"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_URL", func() {
			command := "--policy-id 12345 --n test-url -a REDIRECT_URL -r http://example.com"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("REDIRECT_HTTPS", func() {
			command := "--policy-id 12345 --n test-https -a REDIRECT_HTTPS -r uuid-https-protocol"
			command_args := strings.Fields(command)
			fakeLBManager.EditL7PolicyReturns(datatypes.Network_LBaaS_LoadBalancer{}, nil)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).NotTo(HaveOccurred())
		})
		It("priority 1 REDIRECT_POOL", func() {
			command := "--policy-id 12345 --n test-https -a REDIRECT_POOL -r uuid-https-protocol -p 1"
			command_args := strings.Fields(command)
			createDateTest, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
			fakeLBManager.GetL7PolicyReturns(datatypes.Network_LBaaS_L7Policy{
				Id:          sl.Int(123),
				Uuid:        sl.String("abc123"),
				Name:        sl.String("policy name"),
				Action:      sl.String("REDIRECT_HTTPS"),
				RedirectUrl: sl.String("/"),
				Priority:    sl.Int(1),
				CreateDate:  sl.Time(createDateTest),
			}, nil)
			err := testhelpers.RunCommand(cliCommand, command_args...)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("L7 policy edited"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("ID    UUID     Name          Action           Redirect   Priority   Create Date"))
			Expect(fakeUI.Outputs()).To(ContainSubstring("123   abc123   policy name   REDIRECT_HTTPS   /          1          2022-02-01T00:00:00Z"))
		})
	})
})
