package client

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type CLIRestTransport struct {
	*session.RestTransport
	Context plugin.PluginContext
}

func (r *CLIRestTransport) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	err := r.RestTransport.DoRequest(sess, service, method, args, options, pResult)
	slError, ok := err.(sl.Error)
	if ok {
		if slError.StatusCode == 500 && slError.Exception == "SoftLayer_Exception_Account_Authentication_AccessTokenValidation" {
			newIAMToken, tokenErr := r.Context.RefreshIAMToken()
			if tokenErr != nil {
				return tokenErr
			}
			sess.IAMToken = newIAMToken
			err = r.RestTransport.DoRequest(sess, service, method, args, options, pResult)
		}
	}
	return err
}
