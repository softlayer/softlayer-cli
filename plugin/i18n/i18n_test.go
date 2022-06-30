package i18n_test

import (
	"os"
	"fmt"
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

var _ = Describe("I18NTests", func() {
	coreConfig := prepareConfigForCLI(`{"UsageStatsEnabled": true}`)
	Describe("Language Init Tests", func() {
		Context("Tests All Languages", func() {
			for _, language := range i18n.SUPPORTED_LOCALES {
				language := language
				It("Testing " + language, func() {
					coreConfig.SetLocale(language)
					translator := i18n.Init(coreConfig)
					Expect(translator("Recurring Price")).To(Equal("Recurring Price"))
				})
			}
		})
	})
	Describe("Test ENV Local Lookup", func() {
		oldLang := os.Getenv("LANGUAGE")
		BeforeEach(func() {
			coreConfig.SetLocale("")
		})

		It("Tests es_ES from ENV", func() {
			os.Setenv("LANGUAGE", "es-ES")
			translator := i18n.Init(coreConfig)
			Expect(translator("# of Active Transactions")).To(Equal("# de transacciones activas"))
		})
		It("Tests de_DE from ENV", func() {
			os.Setenv("LANGUAGE", "de-DE")
			translator := i18n.Init(coreConfig)
			Expect(translator("# of Active Transactions")).To(Equal("Anzahl aktiver Transaktionen"))
		})
		It("Tests ja_JP from ENV", func() {
			os.Setenv("LANGUAGE", "ja_JP")
			translator := i18n.Init(coreConfig)
			Expect(translator("# of Active Transactions")).To(Equal("アクティブなトランザクションの数"))
		})
		It("Tests zh_Hans from ENV", func() {
			os.Setenv("LANGUAGE", "zh-CN")
			translator := i18n.Init(coreConfig)
			Expect(translator("# of Active Transactions")).To(Equal("活动事务数"))
		})
		It("Tests zh_Hant from ENV", func() {
			os.Setenv("LANGUAGE", "zh-TW")
			translator := i18n.Init(coreConfig)
			Expect(translator("# of Active Transactions")).To(Equal("作用中交易數目"))
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

