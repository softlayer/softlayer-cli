package client_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	bxmodel "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/models"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/cgallo/softlayer-cli/plugin/client"
)

var _ = Describe("Classic infrastructure CLI Client", func() {
	var (
		context plugin.PluginContext
	)

	BeforeEach(func() {
		os.Setenv("BLUEMIX_HOME", "../testfixtures")
		config := core_config.NewCoreConfig(func(err error) {
			panic(fmt.Sprintf("Configuration error: %v", err))
		})
		context = plugin.InitPluginContext("softlayer")
		config.UnsetAPI()
		config.ClearSession()
		config.SetAccount(bxmodel.Account{})
		cleanSession(context)
	})
	AfterEach(func() {
		os.Setenv("BLUEMIX_HOME", filepath.Join(os.Getenv("HOME"), ".bluemix"))
		cleanSession(context)
	})

	Describe("Create a softlayer session from existing plugin context", func() {
		BeforeEach(func() {
			context.PluginConfig().Set(client.ENV_SL_API_ENDPOINT, client.SoftlayerAPIEndpointPublicDefault)
		})
		It("should return a softlayer session", func() {
			session, err := client.NewSoftlayerClientSessionFromConfig(context)
			Expect(err).ToNot(HaveOccurred())
			Expect(session.Endpoint).To(Equal(client.SoftlayerAPIEndpointPublicDefault))
		})
	})

	Describe("GetSLApiEndPoint", func() {
		Context("//environment variable stores /rest/ endpoint", func() {
			BeforeEach(func() {
				os.Setenv("ENV_SL_API_ENDPOINT", "https://api.softlayer.com/rest/v3.1")
				os.Setenv("BLUEMIX_HOME", "../testfixtures")
				context.PluginConfig().Erase(client.SoftlayerAPIEndpoint)
			})
			It("should return rest endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/rest/v3.1")
			})
			It("should return mobile endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/mobile/v3.1")
			})
		})

		Context("environment variable stores /mobile/ endpoint", func() {
			BeforeEach(func() {
				os.Setenv("ENV_SL_API_ENDPOINT", "https://api.softlayer.com/mobile/v3.1")
				os.Setenv("BLUEMIX_HOME", "../testfixtures")
				context.PluginConfig().Erase(client.SoftlayerAPIEndpoint)
			})
			It("should return rest endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/rest/v3.1")
			})
			It("should return mobile endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/mobile/v3.1")
			})
		})

		Context("context stores /rest/ endpoint", func() {
			BeforeEach(func() {
				context.PluginConfig().Set(client.SoftlayerAPIEndpoint, client.SoftlayerAPIEndpointPublicDefault)
			})
			It("should return rest endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/rest/v3.1")
			})
			It("should return mobile endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/mobile/v3.1")
			})
		})

		Context("context stores /mobile/ endpoint", func() {
			BeforeEach(func() {
				context.PluginConfig().Set(client.SoftlayerAPIEndpoint, client.SoftlayerAPIEndpointPublicDefault)
			})
			It("should return rest endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/rest/v3.1")
			})
			It("should return mobile endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/mobile/v3.1")
			})
		})

		Context("nothing is stored", func() {
			BeforeEach(func() {
				context.PluginConfig().Erase(client.SoftlayerAPIEndpoint)
			})
			It("should return rest endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/rest/v3.1")
			})
			It("should return mobile endpoint", func() {
				apiEndpoint := client.GetSLApiEndPoint(context)
				Expect(apiEndpoint, "https://api.softlayer.com/mobile/v3.1")
			})
		})
	})
})

func cleanSession(context plugin.PluginContext) {
	context.PluginConfig().Erase(client.SoftlayerAPIEndpoint)
}
