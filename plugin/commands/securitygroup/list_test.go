package securitygroup_test

import (
	"errors"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Securitygroup list", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *securitygroup.ListCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = securitygroup.NewListCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        securitygroup.SecurityGroupListMetaData().Name,
			Description: securitygroup.SecurityGroupListMetaData().Description,
			Usage:       securitygroup.SecurityGroupListMetaData().Usage,
			Flags:       securitygroup.SecurityGroupListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Securitygroup list", func() {
		Context("list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Options for --sortby are: id,name,description,created"))
			})
		})
		Context("list but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSecurityGroupsReturns(nil, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get security groups."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("list zero result", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListSecurityGroupsReturns(nil, nil)
			})
			It("return not found", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No security groups are found."))
			})
		})
		Context("list non-zero result", func() {
			BeforeEach(func() {
				created1, _ := time.Parse(time.RFC3339, "2017-11-01T00:00:00Z")
				created2, _ := time.Parse(time.RFC3339, "2017-11-02T00:00:00Z")
				fakeNetworkManager.ListSecurityGroupsReturns([]datatypes.Network_SecurityGroup{
					datatypes.Network_SecurityGroup{
						Id:          sl.Int(123),
						Name:        sl.String("abc"),
						Description: sl.String("def"),
						CreateDate:  sl.Time(created2),
					},
					datatypes.Network_SecurityGroup{
						Id:          sl.Int(321),
						Name:        sl.String("efr"),
						Description: sl.String("abf"),
						CreateDate:  sl.Time(created1),
					},
				}, nil)
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "123")).To(BeTrue())
				Expect(strings.Contains(results[2], "321")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "abc")).To(BeTrue())
				Expect(strings.Contains(results[2], "efr")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "description")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "abf")).To(BeTrue())
				Expect(strings.Contains(results[2], "def")).To(BeTrue())
			})
			It("return table", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "created")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "2017-11-01T00:00:00Z")).To(BeTrue())
				Expect(strings.Contains(results[2], "2017-11-02T00:00:00Z")).To(BeTrue())
			})
		})
	})
})
