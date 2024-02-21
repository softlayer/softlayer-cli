package loadbal_test

import (
	"errors"
	"time"

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

var _ = Describe("Load balancer L7 policies", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *loadbal.L7PolicyListCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeLBManager *testhelpers.FakeLoadBalancerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = loadbal.NewL7PolicyListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeLBManager = new(testhelpers.FakeLoadBalancerManager)
		cliCommand.LoadBalancerManager = fakeLBManager
	})

	id := "1234"
	name := "test"
	uuid := "uuid-123456"
	var action string
	createdTime := "2017-01-03T00:00:00Z"
	created, _ := time.Parse(time.RFC3339, createdTime)

	Context("list without protocol-id", func() {
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("--protocol-id"))
		})
	})

	Context("list with server fails", func() {
		BeforeEach(func() {
			fakeLBManager.GetL7PoliciesReturns(nil, errors.New("Internal Server Error"))
		})
		It("return error", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to get l7 policies"))
		})
	})

	Context("list not policies", func() {
		BeforeEach(func() {
			fakeLBManager.GetL7PoliciesReturns(nil, nil)
		})
		It("return L7 empty policies", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring("No l7 policies was found."))

		})
	})

	Context("list policy REJECT", func() {
		BeforeEach(func() {
			action = "REJECT"
			fakeLBManager.GetL7PoliciesReturns([]datatypes.Network_LBaaS_L7Policy{
				{
					Id:         sl.Int(1234),
					Name:       &name,
					Uuid:       &uuid,
					Action:     &action,
					CreateDate: sl.Time(created),
				},
			}, nil)
		})
		It("return L7 policy REJECT", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(id))
			Expect(fakeUI.Outputs()).To(ContainSubstring(name))
			Expect(fakeUI.Outputs()).To(ContainSubstring(uuid))
			Expect(fakeUI.Outputs()).To(ContainSubstring(createdTime))
		})
	})

	Context("list policy REDIRECT_POOL", func() {
		BeforeEach(func() {
			action = "REDIRECT_POOL"
			fakeLBManager.GetL7PoliciesReturns([]datatypes.Network_LBaaS_L7Policy{
				{
					Id:         sl.Int(1234),
					Name:       &name,
					Uuid:       &uuid,
					Action:     &action,
					CreateDate: sl.Time(created),
				},
			}, nil)
		})
		It("return L7 policy REDIRECT_POOL", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(id))
			Expect(fakeUI.Outputs()).To(ContainSubstring(name))
			Expect(fakeUI.Outputs()).To(ContainSubstring(uuid))
			Expect(fakeUI.Outputs()).To(ContainSubstring(createdTime))
		})
	})

	Context("list policy REDIRECT_URL", func() {
		BeforeEach(func() {
			action = "REDIRECT_URL"
			fakeLBManager.GetL7PoliciesReturns([]datatypes.Network_LBaaS_L7Policy{
				{
					Id:         sl.Int(1234),
					Name:       &name,
					Uuid:       &uuid,
					Action:     &action,
					CreateDate: sl.Time(created),
				},
			}, nil)
		})
		It("return L7 policy REDIRECT_URL", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(id))
			Expect(fakeUI.Outputs()).To(ContainSubstring(name))
			Expect(fakeUI.Outputs()).To(ContainSubstring(uuid))
			Expect(fakeUI.Outputs()).To(ContainSubstring(createdTime))
		})
	})

	Context("list policy REDIRECT_HTTPS", func() {
		BeforeEach(func() {
			action = "REDIRECT_HTTPS"
			fakeLBManager.GetL7PoliciesReturns([]datatypes.Network_LBaaS_L7Policy{
				{
					Id:         sl.Int(1234),
					Name:       &name,
					Uuid:       &uuid,
					Action:     &action,
					CreateDate: sl.Time(created),
				},
			}, nil)
		})
		It("return L7 policy REDIRECT_HTTPS", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(id))
			Expect(fakeUI.Outputs()).To(ContainSubstring(name))
			Expect(fakeUI.Outputs()).To(ContainSubstring(uuid))
			Expect(fakeUI.Outputs()).To(ContainSubstring(createdTime))
		})
	})

	Context("list policy REDIRECT_URL", func() {
		BeforeEach(func() {
			action = "REDIRECT_POOL"
			fakeLBManager.GetL7PoliciesReturns([]datatypes.Network_LBaaS_L7Policy{
				{
					Id:         sl.Int(1234),
					Name:       &name,
					Uuid:       &uuid,
					Action:     &action,
					CreateDate: sl.Time(created),
				},
			}, nil)
		})
		It("return L7 policy", func() {
			err := testhelpers.RunCobraCommand(cliCommand.Command, "--protocol-id", "1234")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstring(id))
			Expect(fakeUI.Outputs()).To(ContainSubstring(name))
			Expect(fakeUI.Outputs()).To(ContainSubstring(uuid))
			Expect(fakeUI.Outputs()).To(ContainSubstring(createdTime))
		})
	})
})
