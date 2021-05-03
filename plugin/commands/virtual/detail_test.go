package virtual_test

import (
	"errors"
	"strings"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.DetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewDetailCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSDetailMataData().Name,
			Description: metadata.VSDetailMataData().Description,
			Usage:       metadata.VSDetailMataData().Usage,
			Flags:       metadata.VSDetailMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS detail", func() {
		Context("VS detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS detail with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS detail with correct VS ID ", func() {
			created, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
			modified, _ := time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(
					datatypes.Virtual_Guest{
						Id:                       sl.Int(1234),
						GlobalIdentifier:         sl.String("rthtoshfkthr"),
						Hostname:                 sl.String("vs-abc"),
						Domain:                   sl.String("wilma.com"),
						FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
						Status: &datatypes.Virtual_Guest_Status{
							Name: sl.String("Provisioning"),
						},
						PowerState: &datatypes.Virtual_Guest_Power_State{
							Name: sl.String("PowerOn"),
						},
						ActiveTransaction: &datatypes.Provisioning_Version1_Transaction{
							TransactionStatus: &datatypes.Provisioning_Version1_Transaction_Status{
								Name: sl.String("Provisioning"),
							},
						},
						Datacenter: &datatypes.Location{
							Name: sl.String("dal10"),
						},
						OperatingSystem: &datatypes.Software_Component_OperatingSystem{
							Software_Component: datatypes.Software_Component{
								SoftwareLicense: &datatypes.Software_License{
									SoftwareDescription: &datatypes.Software_Description{
										Name:    sl.String("CentOS"),
										Version: sl.String("6.0"),
									},
								},
								Passwords: []datatypes.Software_Component_Password{
									datatypes.Software_Component_Password{
										Username: sl.String("root"),
										Password: sl.String("password4root"),
									},
								},
							},
						},
						MaxCpu:                       sl.Int(8),
						MaxMemory:                    sl.Int(4096),
						PrimaryIpAddress:             sl.String("9.9.9.9"),
						PrimaryBackendIpAddress:      sl.String("1.1.1.1"),
						PrivateNetworkOnlyFlag:       sl.Bool(false),
						DedicatedAccountHostOnlyFlag: sl.Bool(false),
						CreateDate:                   sl.Time(created),
						ModifyDate:                   sl.Time(modified),
						BillingItem: &datatypes.Billing_Item_Virtual_Guest{
							Billing_Item: datatypes.Billing_Item{
								OrderItem: &datatypes.Billing_Order_Item{
									Order: &datatypes.Billing_Order{
										UserRecord: &datatypes.User_Customer{
											Username: sl.String("wilmawang"),
										},
									},
								},
								RecurringFee:                    sl.Float(1000.00),
								NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
							},
						},
						Notes: sl.String("mynotes"),
						TagReferences: []datatypes.Tag_Reference{
							datatypes.Tag_Reference{
								Tag: &datatypes.Tag{
									Name: sl.String("tag1"),
								},
							},
							datatypes.Tag_Reference{
								Tag: &datatypes.Tag{
									Name: sl.String("tag2"),
								},
							},
						},
						NetworkVlans: []datatypes.Network_Vlan{
							datatypes.Network_Vlan{
								Id:           sl.Int(678),
								VlanNumber:   sl.Int(50),
								NetworkSpace: sl.String("PRIMARY"),
							},
						},
					}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"rthtoshfkthr"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Provisioning"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PowerOn"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"CentOS"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"6.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4096"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.1.1.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-25T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilmawang"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"mynotes"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"tag1,tag2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"50"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"root"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"password4root"}))
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"1000.00"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--passwords", "--price")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"rthtoshfkthr"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Provisioning"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PowerOn"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"CentOS"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"6.0"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"8"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"4096"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9.9.9.9"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1.1.1.1"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-25T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilmawang"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"mynotes"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"tag1,tag2"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"50"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"PRIMARY"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"root"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"password4root"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1000.00"}))
			})
		})
	})
})
