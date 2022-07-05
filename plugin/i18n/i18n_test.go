package i18n_test

import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

func TestI18N(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "I18N Suite")
}

func prepareConfigForCLI(cliConfigContent string) core_config.Repository {
	ioutil.WriteFile("config.json", []byte(cliConfigContent), 0644)
	ioutil.WriteFile("cf_config.json", []byte(""), 0644)
	return core_config.NewCoreConfigFromPath("cf_config.json", "config.json", func(err error) {
		fmt.Printf("prepareConfigForCLI() Error: %v", err)
	})
}

var xlationMap map[string]string

var _ = Describe("I18NTests", func() {
	coreConfig := prepareConfigForCLI(`{"UsageStatsEnabled": true}`)
	xlationMap = map[string]string {
		"de_DE": "Wiederkehrender Preis",
		"es_ES": "Precio recurrente",
		"fr_FR": "Prix récurrent",
		"it_IT": "Prezzo ricorrente",
		"ja_JP": "定期払い価格",
		"ko_KR": "반복 가격",
		"pt_BR": "Preço recorrente",
		"zh_Hans": "重复出价",
		"zh_Hant": "循環價格",
		"en_US": "Recurring Price",
	}
	Describe("Language Init Tests", func() {
		Context("Tests All Languages", func() {
			for _, language := range i18n.SUPPORTED_LOCALES {
				language := language

				It("Testing " + language, func() {
					coreConfig.SetLocale(language)
					translator := i18n.Init(coreConfig)
					Expect(translator("Recurring Price")).To(Equal(xlationMap[language]))
				})
			}
		})
	})
	Describe("Test ENV Local Lookup", func() {
		oldLang := os.Getenv("LANGUAGE")
		BeforeEach(func() {
			coreConfig.SetLocale("")
		})
		Context("Test loading from ENV variables", func() {
			for _, language := range i18n.SUPPORTED_LOCALES {
				language := language
				envLang := strings.Replace(language, "_", "-", 1)
				It("LANGUAGE=" + envLang, func() {
					os.Setenv("LANGUAGE", envLang)
					translator := i18n.Init(coreConfig)
					Expect(translator("Recurring Price")).To(Equal(xlationMap[language]))
				})
			}
		})
		AfterEach(func() {
			defer os.Setenv("LANGUAGE", oldLang)	
		})
	})
	AfterEach(func() {
		defer os.Remove("config.json")
		defer os.Remove("cf_config.json")
	})
})

