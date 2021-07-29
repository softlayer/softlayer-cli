package managers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("TagsManager", func() {
	var (
		fakeSLSession 	 *session.Session
		tagsManager    managers.TagsManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		tagsManager = managers.NewTagsManager(fakeSLSession)
	})

	Describe("ListTags", func() {
		Context("TagManager.ListTags()", func() {
			It("Returns no errors", func() {
				_, err := tagsManager.ListTags()
				Expect(err).ToNot(HaveOccurred())
			})
			It("Handles Error", func() {
				fakeHandler := testhelpers.FakeTransportHandler{}
                fakeHandler.AddApiError("SoftLayer_Tag", "getAttachedTagsForCurrentUser", 500, "BAD")
                fakeSLSession := &session.Session{TransportHandler: fakeHandler,}

				tagsManager = managers.NewTagsManager(fakeSLSession)
				_, err := tagsManager.ListTags()
				apiError := err.(sl.Error)
				Expect(err).To(HaveOccurred())
				Expect(apiError.StatusCode).To(Equal(500))
			})
		})
		Context("TagManager.ListEmptyTags()", func() {
			It("Returns no errors", func() {
				_, err := tagsManager.ListEmptyTags()
				Expect(err).ToNot(HaveOccurred())
			})
			It("Handles Error", func() {
				fakeHandler := testhelpers.FakeTransportHandler{}
                fakeHandler.AddApiError("SoftLayer_Tag", "getUnattachedTagsForCurrentUser", 500, "Test Error")
                fakeSLSession := &session.Session{TransportHandler: fakeHandler,}

				tagsManager = managers.NewTagsManager(fakeSLSession)
				_, err := tagsManager.ListEmptyTags()
				apiError := err.(sl.Error)
				Expect(err).To(HaveOccurred())
				Expect(apiError.StatusCode).To(Equal(500))
			})
		})
		Context("TagManager.GetTagByTagName()", func() {
			It("Returns no errors", func() {
				tags, err := tagsManager.GetTagByTagName("test1")
				Expect(err).ToNot(HaveOccurred())
				Expect(*tags[0].Name).To(Equal("test1"))
			})
		})
		Context("TagManager.GetTagReferences()", func() {
			It("Returns no errors", func() {
				tags, err := tagsManager.GetTagReferences(1234)
				Expect(err).ToNot(HaveOccurred())
				Expect(*tags[0].Tag.Name).To(Equal("tag03022020"))
			})
			It("Handles Error", func() {
				fakeHandler := testhelpers.FakeTransportHandler{}
                fakeHandler.AddApiError("SoftLayer_Tag", "getReferences", 500, "Test Error")
                fakeSLSession := &session.Session{TransportHandler: fakeHandler,}
				tagsManager = managers.NewTagsManager(fakeSLSession)
				_, err := tagsManager.GetTagReferences(1234)
				apiError := err.(sl.Error)
				Expect(err).To(HaveOccurred())
				Expect(apiError.StatusCode).To(Equal(500))
			})
		})
		Context("TagManager.DeleteTag()", func() {
			It("Returns no errors", func() {
				tags, err := tagsManager.DeleteTag("test1")
				Expect(err).ToNot(HaveOccurred())
				Expect(tags).To(Equal(true))
			})
		})
		Context("TagManager.ReferenceLookup()", func() {
			It("HARDWARE", func() {
				result := tagsManager.ReferenceLookup("HARDWARE", 1)
				Expect(result).To(Equal("ys1-0-baremetal-uaadb.softlayer.com"))	
			})
			It("GUEST", func() {
				result := tagsManager.ReferenceLookup("GUEST", 1)
				Expect(result).To(Equal("wilma2.wilma.org"))	
			})
			It("TICKET", func() {
				result := tagsManager.ReferenceLookup("TICKET", 1)
				Expect(result).To(Equal("API Question - Test ticket"))	
			})
			It("NETWORK_VLAN_FIREWALL", func() {
				result := tagsManager.ReferenceLookup("NETWORK_VLAN_FIREWALL", 1)
				Expect(result).To(Equal("1.2.3.4"))	
			})
			It("IMAGE_TEMPLATE", func() {
				result := tagsManager.ReferenceLookup("IMAGE_TEMPLATE", 1)
				Expect(result).To(Equal("testimage"))	
			})
			It("APPLICATION_DELIVERY_CONTROLLER", func() {
				result := tagsManager.ReferenceLookup("APPLICATION_DELIVERY_CONTROLLER", 1)
				Expect(result).To(Equal("testNetworkObject"))	
			})
			It("NETWORK_VLAN", func() {
				result := tagsManager.ReferenceLookup("NETWORK_VLAN", 1)
				Expect(result).To(Equal("-"))	
			})
			It("NETWORK_SUBNET", func() {
				result := tagsManager.ReferenceLookup("NETWORK_SUBNET", 1)
				Expect(result).To(Equal("10.40.92.68"))	
			})
			It("DEDICATED_HOST", func() {
				result := tagsManager.ReferenceLookup("DEDICATED_HOST", 1)
				Expect(result).To(Equal("VirtualDedicated name"))	
			})
		})
		Context("TagManager.ReferenceLookup()", func() {
			It("Handle 404 Error", func() {
				apiError := sl.Error{
					StatusCode: 404,
					Message: "Fake 404",
				}
				name := "HARDWARE"
				result := managers.NameCheck(&name, apiError)
				Expect(result).To(Equal("Not Found"))	
			})
			It("Handle Other Error", func() {
				apiError := sl.Error{
					StatusCode: 123,
					Message: "Fake Error",
				}
				name := "HARDWARE"
				result := managers.NameCheck(&name, apiError)
				Expect(result).To(Equal("Fake Error (HTTP 123)"))	
			})
		})
	})
})
