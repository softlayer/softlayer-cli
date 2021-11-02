package i18n_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"testing"
)

func TestI18N(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "I18N Suite")
}

var _ = Describe("I18NTests", func() {
	coreConfig := core_config.NewCoreConfig(func(error) {})
	detector := new(testhelpers.FakeDetector)
	Describe("Language Init Tests", func() {
		Context("Tests All Languages", func() {
			for _, language := range i18n.SUPPORTED_LOCALES {
				language := language
				It("Testing "+language, func() {
					detector.DetectLocaleReturns(language)
					translator := i18n.Init(coreConfig, detector)
					Expect(translator("Recurring Price")).To(Equal("Recurring Price"))
				})
			}

		})
	})
})
