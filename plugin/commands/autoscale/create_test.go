package autoscale_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale create", func() {
	var (
		fakeUI               *terminal.FakeUI
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
		cmd                  *autoscale.CreateCommand
		cliCommand           cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		cmd = autoscale.NewCreateCommand(fakeUI, fakeAutoScaleManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleCreateMetaData().Name,
			Description: autoscale.AutoScaleCreateMetaData().Description,
			Usage:       autoscale.AutoScaleCreateMetaData().Usage,
			Flags:       autoscale.AutoScaleCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale create", func() {

		Context("Return error", func() {
			It("Set command without options", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`Required flags "name, cooldown, min, max, regional, os, datacenter, hostname, domain, cpu, memory, termination-policy, disk" not set`))
			})

			It("Set command with an invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetDatacenterByNameReturns([]datatypes.Location{}, errors.New("Failed to get Datacenters."))
			})

			It("Failed get datacenter", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Datacenters."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerdadacenters := []datatypes.Location{
					datatypes.Location{
						Id: sl.Int(111111),
					},
				}
				fakeAutoScaleManager.GetDatacenterByNameReturns(fakerdadacenters, nil)
				fakeAutoScaleManager.CreateScaleGroupReturns(datatypes.Scale_Group{}, errors.New("Failed to create Auto Scale Group."))
			})
			It("Failed create scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create Auto Scale Group."))
			})

			It("Set without --policy-amount", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "--policy-relative=ABSOLUTE")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-amount' is required"))
			})

			It("Set without --policy-name", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "--policy-relative=ABSOLUTE", "--policy-amount=1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-name' is required"))
			})

			It("Set without --policy-relative", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "--policy-name=policy", "--policy-amount=1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--policy-relative' is required"))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerdadacenters := []datatypes.Location{
					datatypes.Location{
						Id: sl.Int(111111),
					},
				}
				fakeAutoScaleManager.GetDatacenterByNameReturns(fakerdadacenters, nil)
				fakeUI.Inputs("abcde")
			})
			It("Cancel with invalid input", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("input must be 'y', 'n', 'yes' or 'no'"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakerdadacenters := []datatypes.Location{
					datatypes.Location{
						Id: sl.Int(111111),
					},
				}
				fakerScaleGroup := datatypes.Scale_Group{
					Id:         sl.Int(222222),
					CreateDate: sl.Time(created),
					Name:       sl.String("testcreate"),
					VirtualGuestMembers: []datatypes.Scale_Member_Virtual_Guest{
						datatypes.Scale_Member_Virtual_Guest{
							VirtualGuest: &datatypes.Virtual_Guest{
								Id:       sl.Int(333333),
								Domain:   sl.String("test.com"),
								Hostname: sl.String("testcreatehostname"),
							},
						},
					},
				}
				fakeAutoScaleManager.GetDatacenterByNameReturns(fakerdadacenters, nil)
				fakeAutoScaleManager.CreateScaleGroupReturns(fakerScaleGroup, nil)
			})

			It("Create scale group without confirmation", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate9", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25", "", "--disk=30", "--userdata=Centos", "--policy-relative=ABSOLUTE",
					"--policy-name=policytest", "--policy-amount=3", "--postinstall=http://test.com", "--key=2153554", "--key=2153250", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2016-12-25T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testcreate"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("333333"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("testcreatehostname"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerdadacenters := []datatypes.Location{
					datatypes.Location{
						Id: sl.Int(111111),
					},
				}
				fakeAutoScaleManager.GetDatacenterByNameReturns(fakerdadacenters, nil)
				fakeUI.Inputs("n")
			})

			It("Cancel", func() {
				err := testhelpers.RunCommand(cliCommand, "--name=testcreate2", "--datacenter=ams01", "--domain=test.com",
					"--hostname=testcreatehostname", "--cooldown=3600", "--min=1", "--max=3", "--regional=142", "--termination-policy=2",
					"-os=CENTOS_7_64", "--cpu=2", "--memory=1024", "--disk=25")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
	})
})
