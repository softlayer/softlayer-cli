package testhelpers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"sync"
)

type FakeLicensesManager struct {
	CreateLicensesOptionsStub        func() ([]datatypes.Product_Item, error)
	createLicensesOptionsMutex       sync.RWMutex
	createLicensesOptionsArgsForCall []struct{}
	createLicensesOptionsReturns     struct {
		result1 []datatypes.Product_Item
		result2 error
	}
	createLicensesOptionsReturnsOnCall map[int]struct {
		result1 []datatypes.Product_Item
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
	
}

func (fake *FakeLicensesManager) CreateLicensesOptions() ([]datatypes.Product_Item, error) {
	fake.createLicensesOptionsMutex.Lock()
	ret, specificReturn := fake.createLicensesOptionsReturnsOnCall[len(fake.createLicensesOptionsArgsForCall)]
	fake.createLicensesOptionsArgsForCall = append(fake.createLicensesOptionsArgsForCall, struct{}{})
	fake.recordInvocation("CreateLicensesOptions", []interface{}{})
	fake.createLicensesOptionsMutex.Unlock()
	if fake.CreateLicensesOptionsStub != nil {
		return fake.CreateLicensesOptionsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createLicensesOptionsReturns.result1, fake.createLicensesOptionsReturns.result2
}

func (fake *FakeLicensesManager) CreateLicensesOptionsCallCount() int {
	fake.createLicensesOptionsMutex.RLock()
	defer fake.createLicensesOptionsMutex.RUnlock()
	return len(fake.createLicensesOptionsArgsForCall)
}

func (fake *FakeLicensesManager) CreateLicensesOptionsReturns(result1 []datatypes.Product_Item, result2 error) {
	fake.CreateLicensesOptionsStub = nil
	fake.createLicensesOptionsReturns = struct {
		result1 []datatypes.Product_Item
		result2 error
	}{result1, result2}
}

func (fake *FakeLicensesManager) CreateLicensesOptionsReturnsOnCall(i int, result1 []datatypes.Product_Item, result2 error) {
	fake.CreateLicensesOptionsStub = nil
	if fake.createLicensesOptionsReturnsOnCall == nil {
		fake.createLicensesOptionsReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Product_Item
			result2 error
		})
	}
	fake.createLicensesOptionsReturnsOnCall[i] = struct {
		result1 []datatypes.Product_Item
		result2 error
	}{result1, result2}
}

func (fake *FakeLicensesManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createLicensesOptionsMutex.RLock()
	defer fake.createLicensesOptionsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLicensesManager) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ managers.LicensesManager = new(FakeLicensesManager)