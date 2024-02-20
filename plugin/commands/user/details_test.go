package user_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var testUser datatypes.User_Customer
var _ = Describe("Detail", func() {
	var (
		fakeUI          *terminal.FakeUI
		fakeUserManager *testhelpers.FakeUserManager
		cliCommand      *user.DetailsCommand
		fakeSession     *session.Session
		slCommand       *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeUserManager = new(testhelpers.FakeUserManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = user.NewDetailsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.UserManager = fakeUserManager

		created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")

		testUser = datatypes.User_Customer{
			Id:       sl.Int(5555),
			Username: sl.String("ATestUser"),
			ApiAuthenticationKeys: []datatypes.User_Customer_ApiAuthentication{datatypes.User_Customer_ApiAuthentication{
				AuthenticationKey: sl.String("StringKeyAuthentication"),
			}},
			FirstName:             sl.String("Name"),
			LastName:              sl.String("LastName"),
			Email:                 sl.String("user@email.com"),
			OpenIdConnectUserName: sl.String("123456"),
			Address1:              sl.String("addres with number N 123"),
			CompanyName:           sl.String("NameCompany"),
			CreateDate:            sl.Time(created),
			OfficePhone:           sl.String("123456789"),
			PptpVpnAllowedFlag:    sl.Bool(true),
			SslVpnAllowedFlag:     sl.Bool(true),
			Parent: &datatypes.User_Customer{
				Username: sl.String("ParentName"),
			},
			UserStatus: &datatypes.User_Customer_Status{
				Name: sl.String("ACTIVE"),
			},
			SuccessfulLogins: []datatypes.User_Customer_Access_Authentication{
				datatypes.User_Customer_Access_Authentication{
					CreateDate: sl.Time(created),
					IpAddress:  sl.String("1.1.1.1"),
				},
			},
			UnsuccessfulLogins: []datatypes.User_Customer_Access_Authentication{
				datatypes.User_Customer_Access_Authentication{
					CreateDate: sl.Time(created),
					IpAddress:  sl.String("2.2.2.2"),
				},
			},
			DedicatedHosts: []datatypes.Virtual_DedicatedHost{
				datatypes.Virtual_DedicatedHost{
					Id:             sl.Int(123456),
					Name:           sl.String("dedicatedHostName"),
					CpuCount:       sl.Int(50),
					MemoryCapacity: sl.Int(1000),
					DiskCapacity:   sl.Int(2000),
					CreateDate:     sl.Time(created),
				},
				datatypes.Virtual_DedicatedHost{
					Id:             sl.Int(1234567),
					Name:           sl.String("dedicatedHostName2"),
					CpuCount:       sl.Int(60),
					MemoryCapacity: sl.Int(1100),
					DiskCapacity:   sl.Int(2100),
					CreateDate:     sl.Time(created),
				}},
			Hardware: []datatypes.Hardware{
				datatypes.Hardware{
					Id:                       sl.Int(12345678),
					FullyQualifiedDomainName: sl.String("domain.test.com"),
					PrimaryIpAddress:         sl.String("10.10.10.10"),
					PrimaryBackendIpAddress:  sl.String("11.11.11.11"),
					ProvisionDate:            sl.Time(created),
				},
				datatypes.Hardware{
					Id:                       sl.Int(123456789),
					FullyQualifiedDomainName: sl.String("domain2.test.com"),
					PrimaryIpAddress:         sl.String("20.20.20.20"),
					PrimaryBackendIpAddress:  sl.String("21.21.21.21"),
					ProvisionDate:            sl.Time(created),
				}},
			VirtualGuests: []datatypes.Virtual_Guest{
				datatypes.Virtual_Guest{
					Id:                       sl.Int(654321),
					FullyQualifiedDomainName: sl.String("virtualtest.domain.com"),
					PrimaryIpAddress:         sl.String("30.30.30.30"),
					PrimaryBackendIpAddress:  sl.String("31.31.31.31"),
					ProvisionDate:            sl.Time(created),
				},
				datatypes.Virtual_Guest{
					Id:                       sl.Int(7654321),
					FullyQualifiedDomainName: sl.String("virtualtest2.domain.com"),
					PrimaryIpAddress:         sl.String("40.40.40.40"),
					PrimaryBackendIpAddress:  sl.String("41.41.41.41"),
					ProvisionDate:            sl.Time(created),
				}},
		}

		testPermissions := []datatypes.User_Customer_CustomerPermission_Permission{
			datatypes.User_Customer_CustomerPermission_Permission{
				KeyName: sl.String("KEY_PERMISSION_1"),
				Name:    sl.String("Permission 1"),
			},
			datatypes.User_Customer_CustomerPermission_Permission{
				KeyName: sl.String("KEY_PERMISSION_2"),
				Name:    sl.String("Permission 2"),
			},
		}

		testLogins := []datatypes.User_Customer_Access_Authentication{
			datatypes.User_Customer_Access_Authentication{
				CreateDate:  sl.Time(created),
				IpAddress:   sl.String("50.50.50.50"),
				SuccessFlag: sl.Bool(true),
			},
			datatypes.User_Customer_Access_Authentication{
				CreateDate:  sl.Time(created),
				IpAddress:   sl.String("60.60.60.60"),
				SuccessFlag: sl.Bool(false),
			},
		}

		testEvents := []datatypes.Event_Log{
			datatypes.Event_Log{
				EventCreateDate: sl.Time(created),
				EventName:       sl.String("Login Successful"),
				IpAddress:       sl.String("70.70.70.70"),
				Label:           sl.String("test@test.com"),
				Username:        sl.String("test_test.com"),
			},
			datatypes.Event_Log{
				EventCreateDate: sl.Time(created),
				EventName:       sl.String("IAM Token validation successful"),
				IpAddress:       sl.String("80.80.80.80"),
				Label:           sl.String("test2@test.com"),
				Username:        sl.String("test2_test.com"),
			},
		}

		fakeUserManager.GetUserReturns(testUser, nil)
		fakeUserManager.GetUserPermissionsReturns(testPermissions, nil)
		fakeUserManager.GetLoginsReturns(testLogins, nil)
		fakeUserManager.GetEventsReturns(testEvents, nil)
	})

	Describe("user detail", func() {
		Context("user detail with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("user detail with letters like parameter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: User ID should be a number."))
			})
		})

		Context("user detail error", func() {
			It("return error", func() {
				fakeUserManager.GetUserReturns(datatypes.User_Customer{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show user detail."))
			})
		})

		Context("user detail error with permissions", func() {
			It("return error", func() {
				fakeUserManager.GetUserPermissionsReturns([]datatypes.User_Customer_CustomerPermission_Permission{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--permissions")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show user permissions."))
			})
		})

		Context("user detail error with logins", func() {
			It("return error", func() {
				fakeUserManager.GetLoginsReturns([]datatypes.User_Customer_Access_Authentication{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--logins")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show login history."))
			})
		})

		Context("user detail error with events", func() {
			It("return error", func() {
				fakeUserManager.GetEventsReturns([]datatypes.Event_Log{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--events")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show event log."))
			})
		})
		Context("user detail error with events", func() {
			It("return error", func() {
				fakeUserManager.GetEventsReturns([]datatypes.Event_Log{}, errors.New("Internal server error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--events")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show event log."))
			})
		})
		Context("Error getting hardware", func() {
			It("return error", func() {
				fakeUserManager.GetUserReturnsOnCall(0, testUser, nil)
				fakeUserManager.GetUserReturnsOnCall(1, datatypes.User_Customer{}, errors.New("BAD HARDWARE"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--hardware")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show hardware."))
			})
		})
		Context("Error getting virtual", func() {
			It("return error", func() {
				fakeUserManager.GetUserReturnsOnCall(0, testUser, nil)
				fakeUserManager.GetUserReturnsOnCall(1, datatypes.User_Customer{}, errors.New("BAD VIRTUAL"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--virtual")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to show virual server."))
			})
		})
		Context("user detail with correct id", func() {
			It("return a user", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
			})
		})

		Context("user detail with correct id and apikey", func() {
			It("return a user with apikey", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--keys")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("APIKEY              StringKeyAuthentication"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
			})
		})

		Context("user detail with correct id and permissions", func() {
			It("return a user with permissions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--permissions")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("keyName            name"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("KEY_PERMISSION_1   Permission 1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("KEY_PERMISSION_2   Permission 2"))

			})
		})

		Context("user detail with correct id and hardware", func() {
			It("return a user with hardware", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--hardware")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID        Name                 Cpus   Memory   Disk   Created                Dedicated Access"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456    dedicatedHostName    50     1000     2000   2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234567   dedicatedHostName2   60     1100     2100   2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345678    domain.test.com    10.10.10.10         11.11.11.11          2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789   domain2.test.com   20.20.20.20         21.21.21.21          2017-11-08T00:00:00Z"))
			})
		})

		Context("user detail with correct id and virtual", func() {
			It("return a user with virtual", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--virtual")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("654321    virtualtest.domain.com    30.30.30.30         31.31.31.31          2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("7654321   virtualtest2.domain.com   40.40.40.40         41.41.41.41          2017-11-08T00:00:00Z"))
			})
		})

		Context("user detail with correct id and logins", func() {
			It("return a user with logins", func() {
				fmt.Println("**")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--logins")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Date                   IP Address    Successful Login?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z   50.50.50.50   true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z   60.60.60.60   false"))
			})
		})

		Context("user detail with correct id and events", func() {
			It("return a user with events", func() {
				fmt.Println("**")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--events")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("name                value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                  5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Username            ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Email               user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("OpenID              123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Address             addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Company             NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created             2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Phone Number        123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Parent User         ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Status              ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("PPTP VPN            true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SSL VPN             true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Login          2017-11-08T00:00:00Z From: 1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last Failed Login   2017-11-08T00:00:00Z From: 2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Date                   Type                              IP Address    Label            Username"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z   Login Successful                  70.70.70.70   test@test.com    test_test.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z   IAM Token validation successful   80.80.80.80   test2@test.com   test2_test.com"))
			})
		})

		Context("user detail with correct id and without apikey", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")

				testUser = datatypes.User_Customer{
					Id:                    sl.Int(5555),
					Username:              sl.String("ATestUser"),
					ApiAuthenticationKeys: []datatypes.User_Customer_ApiAuthentication{},
					FirstName:             sl.String("Name"),
					LastName:              sl.String("LastName"),
					Email:                 sl.String("user@email.com"),
					OpenIdConnectUserName: sl.String("123456"),
					Address1:              sl.String("addres with number N 123"),
					CompanyName:           sl.String("NameCompany"),
					CreateDate:            sl.Time(created),
					OfficePhone:           sl.String("123456789"),
					PptpVpnAllowedFlag:    sl.Bool(true),
					SslVpnAllowedFlag:     sl.Bool(true),
					Parent: &datatypes.User_Customer{
						Username: sl.String("ParentName"),
					},
					UserStatus: &datatypes.User_Customer_Status{
						Name: sl.String("ACTIVE"),
					},
				}
				fakeUserManager.GetUserReturns(testUser, nil)
			})

			It("return a user without apikey", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "5555", "--keys")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("5555"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ATestUser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("No"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name LastName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("user@email.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("addres with number N 123 - - - - -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("NameCompany"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ParentName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ACTIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
			})
		})
	})
})
