package client

import (
	"os"
	"strconv"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/cgallo/softlayer-cli/plugin/version"
)

const (
	SoftlayerAPIEndpoint              = "SoftlayerApiEndpoint"
	SoftlayerAPIEndpointPublicDefault = "https://api.softlayer.com/rest/v3.1"
	ENV_SL_API_ENDPOINT               = "ENV_SL_API_ENDPOINT"
)

func GetSLApiEndPoint(context plugin.PluginContext) string {
	//get from environment variable
	if os.Getenv(ENV_SL_API_ENDPOINT) != "" {
		return os.Getenv(ENV_SL_API_ENDPOINT)
	}

	//get from plugin config
	apiEndPoint, err := context.PluginConfig().GetString(SoftlayerAPIEndpoint)
	if err == nil && apiEndPoint != "" {
		return apiEndPoint
	}

	//get default value
	return SoftlayerAPIEndpointPublicDefault
}

func GetDebug(trace string) bool {
	debug, err := strconv.ParseBool(trace)
	if err != nil {
		debug = false
	}
	return debug
}

func GetTimeout(timeout int) time.Duration {
	return time.Duration(timeout) * time.Second
}

func NewSoftlayerClientSessionFromConfig(context plugin.PluginContext) (*session.Session, error) {

	token := context.IAMToken()

	transportHandler := &CLIRestTransport{
		Context:       context,
		RestTransport: &session.RestTransport{},
	}
	sess := &session.Session{
		Endpoint:         GetSLApiEndPoint(context),
		Debug:            GetDebug(context.Trace()),
		Timeout:          GetTimeout(context.HTTPTimeout()),
		IAMToken:         token,
		TransportHandler: transportHandler,
	}
	sess.AppendUserAgent(version.UsageAgentHeader)
	return sess, nil
}
