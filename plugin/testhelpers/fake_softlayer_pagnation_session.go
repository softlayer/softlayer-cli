package testhelpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type FakePagnationTransportHandler struct {
	FileNames []string
}

func (h FakePagnationTransportHandler) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
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

	b, err := readPagnationJsonTestFixtures(service, method, h.FileNames, options)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	err = json.Unmarshal(b, pResult)
	//	fmt.Println("=================")

	//fmt.Println(pResult)
	return err
}

func NewFakeSoftlayerPagnationSession(fileNames []string) *session.Session {
	return &session.Session{
		TransportHandler: FakePagnationTransportHandler{fileNames},
	}
}

func readPagnationJsonTestFixtures(service string, method string, fileNames []string, options *sl.Options) ([]byte, error) {
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
						fixture = filepath.Join(wd, "..", "testfixtures", "services", service+"_"+method+".json")
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
		if *options.Offset != metadata.LIMIT {
			return ioutil.ReadFile(fixture) // #nosec
		} else {
			return []byte("[]"), nil
		}
		//fmt.Println("read file:" + fixture)

	}
	files := utils.StringSliceToString(fileNames)
	return nil, errors.New("failed to find test fixture file:serivce=" + service + ",method=" + method + ",files:" + files)
}
