package managers_test

import (
	//"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"io/ioutil"
	"os"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Ticket", func() {
	var (
		fakeSLSession *session.Session
		TicketManager managers.TicketManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		TicketManager = managers.NewTicketManager(fakeSLSession)
	})

	Describe("ListTickets", func() {
		Context("ListTicket succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getTickets",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Return tickets", func() {
				tickets, err := TicketManager.ListTickets()
				Expect(err).ToNot(HaveOccurred())

				Expect(*tickets[0].AccountId).To(Equal(int(12345)))
				Expect(*tickets[0].LastEditType).To(Equal("AUTO"))
				Expect(*tickets[0].Id).To(Equal(int(76767688)))

				Expect(*tickets[1].AccountId).To(Equal(int(12345)))
				Expect(*tickets[1].LastEditType).To(Equal("EMPLOYEE"))
				Expect(*tickets[1].Id).To(Equal(int(76767699)))
			})
		})
	})

	Describe("GetTicket", func() {
		Context("GetTicket succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Ticket_getObject",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Return ticket", func() {
				ticket, err := TicketManager.GetTicket(76768)
				Expect(err).ToNot(HaveOccurred())

				Expect(*ticket.Id).To(Equal(int(76768)))
				Expect(*ticket.UpdateCount).To(Equal(uint(3)))

			})
		})
	})

	Describe("GetSummary", func() {
		Context("GetSummary succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getObject",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Return summary", func() {
				summary, err := TicketManager.Summary()
				Expect(err).ToNot(HaveOccurred())

				Expect((*summary).Accounting).To(Equal(uint(1)))
				Expect((*summary).Billing).To(Equal(uint(2)))
				Expect((*summary).Sales).To(Equal(uint(0)))
				Expect((*summary).Support).To(Equal(uint(4)))
			})
		})
	})

	Describe("GetText", func() {
		Context("GetText succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getObject",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Edit Text", func() {
				if os.Getenv("OS") == "Windows_NT"  {
					Skip("Test doesn't work in windows.")
				}
				tmp, err := ioutil.TempFile("", "ibmcloudsl-test")
				Expect(err).ToNot(HaveOccurred())

				tmp.Write([]byte("#!/bin/bash\necho 'Hello, There!' > $1\n"))

				tmp.Close()
				defer os.Remove(tmp.Name())

				os.Chmod(tmp.Name(), 0700)
				os.Setenv("EDITOR", tmp.Name())
				text, err := TicketManager.GetText()

				Expect(err).ToNot(HaveOccurred())
				Expect(text).To(Equal("Hello, There!\n\n\n***POSTED FROM IBMCLOUD SL***"))
			})
		})

		Context("GetText fail", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getObject",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Edit Text getpath notfound", func() {
				if os.Getenv("OS") == "Windows_NT"  {
					Skip("Test doesn't work in windows.")
				}
				os.Setenv("EDITOR", "shouldnotbearealfile")
				text, err := TicketManager.GetText()

				Expect(err).To(HaveOccurred())
				Expect(text).To(Equal(""))
			})

			It("Edit Text getpath found", func() {
				if os.Getenv("OS") == "Windows_NT"  {
					Skip("Test doesn't work in windows.")
				}
				os.Setenv("EDITOR", "ed")
				text, err := TicketManager.GetText()

				Expect(err).To(HaveOccurred())
				Expect(text).To(Equal(""))
			})

		})

	})

	Describe("AddUpdate", func() {
		Context("AddUpdate succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Ticket_addUpdate",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Add update", func() {
				err := TicketManager.AddUpdate(76768, "This is a test")

				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("CreateTicket", func() {
		Context("CreateTicket succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Ticket_createStandardTicket",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Add update", func() {
				var args managers.TicketArguments
				contents := "Example ticket contents"
				title := "Example ticket title"
				args.Content = &contents
				args.Title = &title

				id, err := TicketManager.CreateStandardTicket(&args)

				Expect(err).ToNot(HaveOccurred())
				Expect(*id).To(Equal(int(76768)))
			})
		})
	})

	Describe("AttachDevices", func() {
		Context("AttachDevices succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Ticket_addAttachedHardware",
					"SoftLayer_Ticket_addAttachedVirtualGuest",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Add hardware", func() {
				err := TicketManager.AttachDeviceToTicket(76768, 222222, true)

				Expect(err).ToNot(HaveOccurred())
			})

			It("Add virtual guest", func() {
				err := TicketManager.AttachDeviceToTicket(76768, 222222, false)

				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("GetSubject", func() {
		Context("AddUpdate succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Ticket_Subject_getAllObjects",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Add update", func() {
				_, err := TicketManager.GetSubjects()

				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("RemoveDevice", func() {
		Context("RemoveDevice succ", func() {
			BeforeEach(func() {
				fakeSLSession = testhelpers.NewFakeSoftlayerSession_True()
				TicketManager = managers.NewTicketManager(fakeSLSession)
			})

			It("Remove virtual guest", func() {
				err := TicketManager.RemoveDeviceFromTicket(76768, 222222, false)

				Expect(err).ToNot(HaveOccurred())
			})

			It("Remove hardware", func() {
				err := TicketManager.RemoveDeviceFromTicket(76768, 222222, true)

				Expect(err).ToNot(HaveOccurred())
			})

		})
	})

})
