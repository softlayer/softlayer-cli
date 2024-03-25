package i18n

import (
	"path/filepath"
	"strings"
	"golang.org/x/text/language"
	"encoding/json"
	"fmt"
	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/Xuanwo/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	// "github.ibm.com/SoftLayer/softlayer-cli/plugin/resources"
)

const (
	DEFAULT_LOCALE = "en_US"
)

var SUPPORTED_LOCALES = []string{
	"de_DE",
	"en_US",
	"es_ES",
	"fr_FR",
	"it_IT",
	"ja_JP",
	"ko_KR",
	"pt_BR",
	"zh_Hans",
	"zh_Hant",
}

var resourcePath = filepath.Join("plugin", "i18n", "resources")
var localizer = Init()

func GetResourcePath() string {
	return resourcePath
}

func SetResourcePath(path string) {
	resourcePath = path
}

// var T goi18n.TranslateFunc = Init(core_config.NewCoreConfig(func(error) {}))

func T(text string, subs ...interface{}) string {

	// fmt.Printf("SUBS: %v\n", subs)
	message := &goi18n.Message{ID: text, Other: text}
	config := &goi18n.LocalizeConfig{DefaultMessage: message}
	l_string, err := localizer.Localize(config)
	if err != nil {
		fmt.Printf("ERROR i18n\n")
		return err.Error()
	}
	return l_string
}




func Init() *goi18n.Localizer {
	bundle := goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("plugin/i18n/resources/en_US.json")
	bundle.MustLoadMessageFile("plugin/i18n/resources/ja_JP.json")
	loc := goi18n.NewLocalizer(bundle, language.English.String())
	return loc
}

/*
func initWithLocale(locale string) goi18n.TranslateFunc {
	err := loadFromAsset(locale)
	if err != nil {
		locale = DEFAULT_LOCALE
	}
	return goi18n.MustTfunc(locale)
}

func loadFromAsset(locale string) (err error) {
	assetName := locale + ".all.json"
	assetKey := filepath.Join(resourcePath, assetName)
	bytes, err := resources.Asset(assetKey)
	if err != nil {
		return
	}
	err = goi18n.ParseTranslationFileBytes(assetName, bytes)
	return
}

*/
// Tries to determine the system locale, when local isn't set, default to en_US
func DetectLocal() string {
	tag, err := locale.Detect()
	if err != nil {
		return DEFAULT_LOCALE
	}
	// tag is en-US, needs to be en_US
	locale := strings.Replace(tag.String(), "-", "_", 1)
	return locale
}


// Tries to match the system locale with a supported locale, otherwise sets a DEFAULT_LOCALE
func supportedLocale(configLocal string) string {

	// Check if the configLocal matches, this takes precendent
	for _, l := range SUPPORTED_LOCALES {
		if strings.EqualFold(configLocal, l) {
			return l
		}
	}

	// Check if the system has a local that matches
	locale := DetectLocal()
	for _, l := range SUPPORTED_LOCALES {
		if strings.EqualFold(locale, l) {
			return l
		}
	}
	switch strings.ToLower(locale) {
	case "zh_cn", "zh_sg":
		return "zh_Hans"
	case "zh_hk", "zh_tw":
		return "zh_Hant"
	}
	return DEFAULT_LOCALE
}
