package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("ImageManager", func() {
	var (
		fakeSLSession *session.Session
		imageManager  managers.ImageManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		imageManager = managers.NewImageManager(fakeSLSession)
	})

	Describe("Get an image", func() {
		Context("Get an image given its ID", func() {
			It("It returns a image", func() {
				image, err := imageManager.GetImage(1335057)
				Expect(err).ToNot(HaveOccurred())
				Expect(*image.Name).To(Equal("testimage"))
				Expect(*image.ImageType.KeyName).To(Equal("SYSTEM"))
				Expect(*image.Status.Name).To(Equal("Active"))
			})
		})
	})

	Describe("Delete an image", func() {
		Context("Delete an image given its ID", func() {
			It("It returns no error", func() {
				err := imageManager.DeleteImage(1335057)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("List private images", func() {

		Context("List private images under current account", func() {
			It("It returns a list of private images", func() {
				images, err := imageManager.ListPrivateImages("", "")
				Expect(err).ToNot(HaveOccurred())
				for _, image := range images {
					Expect(*image.Id).ShouldNot(BeNil())
					Expect(*image.AccountId).ShouldNot(BeNil())
					Expect(*image.Name).ShouldNot(BeNil())
					if image.ImageType == nil || image.ImageType.KeyName == nil {
						//fmt.Println("empty type:%d, %s\n", *image.Id, *image.Name)
					} else {
						//fmt.Println("nonempty type:%d, %s\n", *image.Id, *image.Name)
						Expect(*image.ImageType.KeyName).Should(Or(Equal("SYSTEM"), Equal("DISK_CAPTURE")))
					}
				}
			})
		})
	})

	Describe("List public images", func() {

		Context("List public images", func() {
			It("It returns a list of public images", func() {
				images, err := imageManager.ListPublicImages("", "")
				Expect(err).ToNot(HaveOccurred())
				for _, image := range images {
					Expect(*image.Id).ShouldNot(BeNil())
					Expect(*image.AccountId).ShouldNot(BeNil())
					Expect(*image.Name).ShouldNot(BeNil())
					if image.ImageType == nil || image.ImageType.KeyName == nil {
						//fmt.Println("empty type:%d, %s\n", *image.Id, *image.Name)
					} else {
						//fmt.Println("nonempty type:%d, %s\n", *image.Id, *image.Name)
						Expect(*image.ImageType.KeyName).Should(Or(Equal("SYSTEM"), Equal("DISK_CAPTURE")))
					}
				}
			})
		})
	})

	Describe("Edit an image", func() {
		Context("Edit an image's name, note and tag", func() {
			It("It returns a list of successes and a list of messages", func() {
				successes, errs := imageManager.EditImage(166680, "rename", "mynote", "updatedTags")
				for index, success := range successes {
					Expect(success).Should(BeTrue())
					Expect(errs[index]).ShouldNot(BeNil())
				}
			})
		})
	})
})
