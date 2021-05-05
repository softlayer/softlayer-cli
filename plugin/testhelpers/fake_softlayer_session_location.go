package testhelpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type FakeTransportLocationHandler struct {
	FileNames []string
}

func (h FakeTransportLocationHandler) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	// fmt.Println("service:\t", service)
	// fmt.Println("method:\t", method)
	// fmt.Println("filenames:\t", h.FileNames)
	// for _, arg := range args {
	// 	fmt.Println("args:\t", arg)
	// }
	// if options.Id != nil {
	// 	fmt.Println("options-id:\t", *options.Id)
	// }
	// if options.Mask != "" {
	// 	fmt.Println("options-mask:\t", options.Mask)
	// }
	// if options.Filter != "" {
	// 	fmt.Println("options-filter:\t", options.Filter)
	// }

	b, err := readJsonLocationTestFixtures(service, method, h.FileNames)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, pResult)
	//fmt.Println(pResult)
	return err
}

func NewFakeSoftlayerLocationSession(fileNames []string) *session.Session {
	return &session.Session{
		TransportHandler: FakeTransportLocationHandler{fileNames},
	}
}

func readJsonLocationTestFixtures(service string, method string, fileNames []string) ([]byte, error) {
	wd, _ := os.Getwd()
	var fixture string
	if len(fileNames) == 0 {
		fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+".json")
	} else {
		//find the file name that matches the service and method name
		for _, filename := range fileNames {
			//fmt.Println("check file:" + filename)
			nameSegments := strings.Split(filename, "_")
			if nameSegments[0] == "SoftLayer" && nameSegments[1] == "Account" {
				if len(nameSegments) == 3 {
					if service == nameSegments[0]+"_"+nameSegments[1] && method == nameSegments[2] {
						fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+".json")
						break
					}
				} else if len(nameSegments) == 4 {
					if service == nameSegments[0]+"_"+nameSegments[1] && method == nameSegments[2] {
						fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+"_"+nameSegments[3]+".json")
						break
					}
				}
			} else if nameSegments[0] == "SoftLayer" && nameSegments[1] != "Account" {
				if len(nameSegments) == 4 {
					if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2] && method == nameSegments[3] {
						if service == "SoftLayer_Location_Datacenter" {
							fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+"_1.json")

						} else {
							fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+".json")

						}
						break
					}
				} else if len(nameSegments) == 5 {
					if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2] && method == nameSegments[3] {
						fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+"_"+nameSegments[4]+".json")
						break
					} else if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2]+"_"+nameSegments[3] && method == nameSegments[4] {
						fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+".json")
						break
					}
				}
			}
		}
	}
	if fixture != "" {
		//fmt.Println("read file:" + fixture)
		return ioutil.ReadFile(fixture) // #nosec
	}
	files := utils.StringSliceToString(fileNames)
	return nil, errors.New("failed to find test fixture file:serivce=" + service + ",method=" + method + ",files:" + files)
}
