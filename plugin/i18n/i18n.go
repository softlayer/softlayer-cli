package i18n

import (
	"embed"
	"strings"
	"golang.org/x/text/language"
	"encoding/json"
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/Xuanwo/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:embed v2Resources/active.*.json
var LocaleFS embed.FS

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

var localizer = Init()


// Translates a string, with any substitutions needed
// text: string to be translated
// subs: A single map[string]interface{}
func T(text string, subs ...interface{}) string {

	config := &goi18n.LocalizeConfig{
		DefaultMessage: &goi18n.Message{ID: text, Other: text, One: text},
	}
	// Need to use `subs ...interface{}` so that we can have 0 or 1 subs.
	// Should never have 2
	if subs != nil && len(subs) == 1 {
		config.TemplateData = subs[0]
	}

	l_string, err := localizer.Localize(config)
	if err != nil {
		// fmt.Printf("ERROR i18n: %v\n", err.Error())
		// return err.Error()
	}
	return l_string
}

// Sets the localizer, reads local from config/system
func Init() *goi18n.Localizer {
	
	coreConfig := core_config.NewCoreConfig(func(error) {})
	userLocale := coreConfig.Locale()
	locale := supportedLocale(userLocale)
	return InitWithLocale(locale)
}

// Sets the localizer with the proper language
func InitWithLocale(locale string) *goi18n.Localizer {
	
	bundle := goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := bundle.LoadMessageFileFS(LocaleFS, "v2Resources/active.en-US.json")
	if err != nil {
		fmt.Printf("Unable to load language from v2Resources/active.en-US.json\n")
	}	
	if locale != "en_US" {
		lang_file := fmt.Sprintf("v2Resources/active.%s.json", locale)
		_, err = bundle.LoadMessageFileFS(LocaleFS, lang_file)
		if err != nil {
			fmt.Printf("Unable to load language from %s\n", lang_file)
		}
	}
	loc := goi18n.NewLocalizer(bundle, locale)
	return loc
}

// Used for testing and changing the language output dynamically
func SetLocalizer(new_localizer *goi18n.Localizer) {
	localizer = new_localizer
}

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
