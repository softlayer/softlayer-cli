// Code generated by counterfeiter. DO NOT EDIT.
package testhelpers

import (
	"sync"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type FakeDedicatedHostManager struct {
	CancelGuestsStub        func(int) ([]managers.StatusInfo, error)
	cancelGuestsMutex       sync.RWMutex
	cancelGuestsArgsForCall []struct {
		arg1 int
	}
	cancelGuestsReturns struct {
		result1 []managers.StatusInfo
		result2 error
	}
	cancelGuestsReturnsOnCall map[int]struct {
		result1 []managers.StatusInfo
		result2 error
	}
	DeleteHostStub        func(int) error
	deleteHostMutex       sync.RWMutex
	deleteHostArgsForCall []struct {
		arg1 int
	}
	deleteHostReturns struct {
		result1 error
	}
	deleteHostReturnsOnCall map[int]struct {
		result1 error
	}
	GenerateOrderTemplateStub        func(string, string, string, string, string, int) (datatypes.Container_Product_Order_Virtual_DedicatedHost, error)
	generateOrderTemplateMutex       sync.RWMutex
	generateOrderTemplateArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 int
	}
	generateOrderTemplateReturns struct {
		result1 datatypes.Container_Product_Order_Virtual_DedicatedHost
		result2 error
	}
	generateOrderTemplateReturnsOnCall map[int]struct {
		result1 datatypes.Container_Product_Order_Virtual_DedicatedHost
		result2 error
	}
	GetCreateOptionsStub        func(datatypes.Product_Package) map[string]map[string]string
	getCreateOptionsMutex       sync.RWMutex
	getCreateOptionsArgsForCall []struct {
		arg1 datatypes.Product_Package
	}
	getCreateOptionsReturns struct {
		result1 map[string]map[string]string
	}
	getCreateOptionsReturnsOnCall map[int]struct {
		result1 map[string]map[string]string
	}
	GetInstanceStub        func(int, string) (datatypes.Virtual_DedicatedHost, error)
	getInstanceMutex       sync.RWMutex
	getInstanceArgsForCall []struct {
		arg1 int
		arg2 string
	}
	getInstanceReturns struct {
		result1 datatypes.Virtual_DedicatedHost
		result2 error
	}
	getInstanceReturnsOnCall map[int]struct {
		result1 datatypes.Virtual_DedicatedHost
		result2 error
	}
	GetPackageStub        func() (datatypes.Product_Package, error)
	getPackageMutex       sync.RWMutex
	getPackageArgsForCall []struct {
	}
	getPackageReturns struct {
		result1 datatypes.Product_Package
		result2 error
	}
	getPackageReturnsOnCall map[int]struct {
		result1 datatypes.Product_Package
		result2 error
	}
	GetVlansOptionsStub        func(string, string, datatypes.Product_Package) ([]datatypes.Network_Vlan, error)
	getVlansOptionsMutex       sync.RWMutex
	getVlansOptionsArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 datatypes.Product_Package
	}
	getVlansOptionsReturns struct {
		result1 []datatypes.Network_Vlan
		result2 error
	}
	getVlansOptionsReturnsOnCall map[int]struct {
		result1 []datatypes.Network_Vlan
		result2 error
	}
	ListDedicatedHostStub        func(string, string, string, int) ([]datatypes.Virtual_DedicatedHost, error)
	listDedicatedHostMutex       sync.RWMutex
	listDedicatedHostArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 int
	}
	listDedicatedHostReturns struct {
		result1 []datatypes.Virtual_DedicatedHost
		result2 error
	}
	listDedicatedHostReturnsOnCall map[int]struct {
		result1 []datatypes.Virtual_DedicatedHost
		result2 error
	}
	ListGuestsStub        func(int, int, string, string, int, []string, string) ([]datatypes.Virtual_Guest, error)
	listGuestsMutex       sync.RWMutex
	listGuestsArgsForCall []struct {
		arg1 int
		arg2 int
		arg3 string
		arg4 string
		arg5 int
		arg6 []string
		arg7 string
	}
	listGuestsReturns struct {
		result1 []datatypes.Virtual_Guest
		result2 error
	}
	listGuestsReturnsOnCall map[int]struct {
		result1 []datatypes.Virtual_Guest
		result2 error
	}
	OrderInstanceStub        func(datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order_Receipt, error)
	orderInstanceMutex       sync.RWMutex
	orderInstanceArgsForCall []struct {
		arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost
	}
	orderInstanceReturns struct {
		result1 datatypes.Container_Product_Order_Receipt
		result2 error
	}
	orderInstanceReturnsOnCall map[int]struct {
		result1 datatypes.Container_Product_Order_Receipt
		result2 error
	}
	VerifyInstanceCreationStub        func(datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order, error)
	verifyInstanceCreationMutex       sync.RWMutex
	verifyInstanceCreationArgsForCall []struct {
		arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost
	}
	verifyInstanceCreationReturns struct {
		result1 datatypes.Container_Product_Order
		result2 error
	}
	verifyInstanceCreationReturnsOnCall map[int]struct {
		result1 datatypes.Container_Product_Order
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDedicatedHostManager) CancelGuests(arg1 int) ([]managers.StatusInfo, error) {
	fake.cancelGuestsMutex.Lock()
	ret, specificReturn := fake.cancelGuestsReturnsOnCall[len(fake.cancelGuestsArgsForCall)]
	fake.cancelGuestsArgsForCall = append(fake.cancelGuestsArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.CancelGuestsStub
	fakeReturns := fake.cancelGuestsReturns
	fake.recordInvocation("CancelGuests", []interface{}{arg1})
	fake.cancelGuestsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) CancelGuestsCallCount() int {
	fake.cancelGuestsMutex.RLock()
	defer fake.cancelGuestsMutex.RUnlock()
	return len(fake.cancelGuestsArgsForCall)
}

func (fake *FakeDedicatedHostManager) CancelGuestsCalls(stub func(int) ([]managers.StatusInfo, error)) {
	fake.cancelGuestsMutex.Lock()
	defer fake.cancelGuestsMutex.Unlock()
	fake.CancelGuestsStub = stub
}

func (fake *FakeDedicatedHostManager) CancelGuestsArgsForCall(i int) int {
	fake.cancelGuestsMutex.RLock()
	defer fake.cancelGuestsMutex.RUnlock()
	argsForCall := fake.cancelGuestsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDedicatedHostManager) CancelGuestsReturns(result1 []managers.StatusInfo, result2 error) {
	fake.cancelGuestsMutex.Lock()
	defer fake.cancelGuestsMutex.Unlock()
	fake.CancelGuestsStub = nil
	fake.cancelGuestsReturns = struct {
		result1 []managers.StatusInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) CancelGuestsReturnsOnCall(i int, result1 []managers.StatusInfo, result2 error) {
	fake.cancelGuestsMutex.Lock()
	defer fake.cancelGuestsMutex.Unlock()
	fake.CancelGuestsStub = nil
	if fake.cancelGuestsReturnsOnCall == nil {
		fake.cancelGuestsReturnsOnCall = make(map[int]struct {
			result1 []managers.StatusInfo
			result2 error
		})
	}
	fake.cancelGuestsReturnsOnCall[i] = struct {
		result1 []managers.StatusInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) DeleteHost(arg1 int) error {
	fake.deleteHostMutex.Lock()
	ret, specificReturn := fake.deleteHostReturnsOnCall[len(fake.deleteHostArgsForCall)]
	fake.deleteHostArgsForCall = append(fake.deleteHostArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.DeleteHostStub
	fakeReturns := fake.deleteHostReturns
	fake.recordInvocation("DeleteHost", []interface{}{arg1})
	fake.deleteHostMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDedicatedHostManager) DeleteHostCallCount() int {
	fake.deleteHostMutex.RLock()
	defer fake.deleteHostMutex.RUnlock()
	return len(fake.deleteHostArgsForCall)
}

func (fake *FakeDedicatedHostManager) DeleteHostCalls(stub func(int) error) {
	fake.deleteHostMutex.Lock()
	defer fake.deleteHostMutex.Unlock()
	fake.DeleteHostStub = stub
}

func (fake *FakeDedicatedHostManager) DeleteHostArgsForCall(i int) int {
	fake.deleteHostMutex.RLock()
	defer fake.deleteHostMutex.RUnlock()
	argsForCall := fake.deleteHostArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDedicatedHostManager) DeleteHostReturns(result1 error) {
	fake.deleteHostMutex.Lock()
	defer fake.deleteHostMutex.Unlock()
	fake.DeleteHostStub = nil
	fake.deleteHostReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDedicatedHostManager) DeleteHostReturnsOnCall(i int, result1 error) {
	fake.deleteHostMutex.Lock()
	defer fake.deleteHostMutex.Unlock()
	fake.DeleteHostStub = nil
	if fake.deleteHostReturnsOnCall == nil {
		fake.deleteHostReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteHostReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplate(arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 int) (datatypes.Container_Product_Order_Virtual_DedicatedHost, error) {
	fake.generateOrderTemplateMutex.Lock()
	ret, specificReturn := fake.generateOrderTemplateReturnsOnCall[len(fake.generateOrderTemplateArgsForCall)]
	fake.generateOrderTemplateArgsForCall = append(fake.generateOrderTemplateArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 int
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.GenerateOrderTemplateStub
	fakeReturns := fake.generateOrderTemplateReturns
	fake.recordInvocation("GenerateOrderTemplate", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.generateOrderTemplateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplateCallCount() int {
	fake.generateOrderTemplateMutex.RLock()
	defer fake.generateOrderTemplateMutex.RUnlock()
	return len(fake.generateOrderTemplateArgsForCall)
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplateCalls(stub func(string, string, string, string, string, int) (datatypes.Container_Product_Order_Virtual_DedicatedHost, error)) {
	fake.generateOrderTemplateMutex.Lock()
	defer fake.generateOrderTemplateMutex.Unlock()
	fake.GenerateOrderTemplateStub = stub
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplateArgsForCall(i int) (string, string, string, string, string, int) {
	fake.generateOrderTemplateMutex.RLock()
	defer fake.generateOrderTemplateMutex.RUnlock()
	argsForCall := fake.generateOrderTemplateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplateReturns(result1 datatypes.Container_Product_Order_Virtual_DedicatedHost, result2 error) {
	fake.generateOrderTemplateMutex.Lock()
	defer fake.generateOrderTemplateMutex.Unlock()
	fake.GenerateOrderTemplateStub = nil
	fake.generateOrderTemplateReturns = struct {
		result1 datatypes.Container_Product_Order_Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GenerateOrderTemplateReturnsOnCall(i int, result1 datatypes.Container_Product_Order_Virtual_DedicatedHost, result2 error) {
	fake.generateOrderTemplateMutex.Lock()
	defer fake.generateOrderTemplateMutex.Unlock()
	fake.GenerateOrderTemplateStub = nil
	if fake.generateOrderTemplateReturnsOnCall == nil {
		fake.generateOrderTemplateReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Product_Order_Virtual_DedicatedHost
			result2 error
		})
	}
	fake.generateOrderTemplateReturnsOnCall[i] = struct {
		result1 datatypes.Container_Product_Order_Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetCreateOptions(arg1 datatypes.Product_Package) map[string]map[string]string {
	fake.getCreateOptionsMutex.Lock()
	ret, specificReturn := fake.getCreateOptionsReturnsOnCall[len(fake.getCreateOptionsArgsForCall)]
	fake.getCreateOptionsArgsForCall = append(fake.getCreateOptionsArgsForCall, struct {
		arg1 datatypes.Product_Package
	}{arg1})
	stub := fake.GetCreateOptionsStub
	fakeReturns := fake.getCreateOptionsReturns
	fake.recordInvocation("GetCreateOptions", []interface{}{arg1})
	fake.getCreateOptionsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeDedicatedHostManager) GetCreateOptionsCallCount() int {
	fake.getCreateOptionsMutex.RLock()
	defer fake.getCreateOptionsMutex.RUnlock()
	return len(fake.getCreateOptionsArgsForCall)
}

func (fake *FakeDedicatedHostManager) GetCreateOptionsCalls(stub func(datatypes.Product_Package) map[string]map[string]string) {
	fake.getCreateOptionsMutex.Lock()
	defer fake.getCreateOptionsMutex.Unlock()
	fake.GetCreateOptionsStub = stub
}

func (fake *FakeDedicatedHostManager) GetCreateOptionsArgsForCall(i int) datatypes.Product_Package {
	fake.getCreateOptionsMutex.RLock()
	defer fake.getCreateOptionsMutex.RUnlock()
	argsForCall := fake.getCreateOptionsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDedicatedHostManager) GetCreateOptionsReturns(result1 map[string]map[string]string) {
	fake.getCreateOptionsMutex.Lock()
	defer fake.getCreateOptionsMutex.Unlock()
	fake.GetCreateOptionsStub = nil
	fake.getCreateOptionsReturns = struct {
		result1 map[string]map[string]string
	}{result1}
}

func (fake *FakeDedicatedHostManager) GetCreateOptionsReturnsOnCall(i int, result1 map[string]map[string]string) {
	fake.getCreateOptionsMutex.Lock()
	defer fake.getCreateOptionsMutex.Unlock()
	fake.GetCreateOptionsStub = nil
	if fake.getCreateOptionsReturnsOnCall == nil {
		fake.getCreateOptionsReturnsOnCall = make(map[int]struct {
			result1 map[string]map[string]string
		})
	}
	fake.getCreateOptionsReturnsOnCall[i] = struct {
		result1 map[string]map[string]string
	}{result1}
}

func (fake *FakeDedicatedHostManager) GetInstance(arg1 int, arg2 string) (datatypes.Virtual_DedicatedHost, error) {
	fake.getInstanceMutex.Lock()
	ret, specificReturn := fake.getInstanceReturnsOnCall[len(fake.getInstanceArgsForCall)]
	fake.getInstanceArgsForCall = append(fake.getInstanceArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	stub := fake.GetInstanceStub
	fakeReturns := fake.getInstanceReturns
	fake.recordInvocation("GetInstance", []interface{}{arg1, arg2})
	fake.getInstanceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) GetInstanceCallCount() int {
	fake.getInstanceMutex.RLock()
	defer fake.getInstanceMutex.RUnlock()
	return len(fake.getInstanceArgsForCall)
}

func (fake *FakeDedicatedHostManager) GetInstanceCalls(stub func(int, string) (datatypes.Virtual_DedicatedHost, error)) {
	fake.getInstanceMutex.Lock()
	defer fake.getInstanceMutex.Unlock()
	fake.GetInstanceStub = stub
}

func (fake *FakeDedicatedHostManager) GetInstanceArgsForCall(i int) (int, string) {
	fake.getInstanceMutex.RLock()
	defer fake.getInstanceMutex.RUnlock()
	argsForCall := fake.getInstanceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeDedicatedHostManager) GetInstanceReturns(result1 datatypes.Virtual_DedicatedHost, result2 error) {
	fake.getInstanceMutex.Lock()
	defer fake.getInstanceMutex.Unlock()
	fake.GetInstanceStub = nil
	fake.getInstanceReturns = struct {
		result1 datatypes.Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetInstanceReturnsOnCall(i int, result1 datatypes.Virtual_DedicatedHost, result2 error) {
	fake.getInstanceMutex.Lock()
	defer fake.getInstanceMutex.Unlock()
	fake.GetInstanceStub = nil
	if fake.getInstanceReturnsOnCall == nil {
		fake.getInstanceReturnsOnCall = make(map[int]struct {
			result1 datatypes.Virtual_DedicatedHost
			result2 error
		})
	}
	fake.getInstanceReturnsOnCall[i] = struct {
		result1 datatypes.Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetPackage() (datatypes.Product_Package, error) {
	fake.getPackageMutex.Lock()
	ret, specificReturn := fake.getPackageReturnsOnCall[len(fake.getPackageArgsForCall)]
	fake.getPackageArgsForCall = append(fake.getPackageArgsForCall, struct {
	}{})
	stub := fake.GetPackageStub
	fakeReturns := fake.getPackageReturns
	fake.recordInvocation("GetPackage", []interface{}{})
	fake.getPackageMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) GetPackageCallCount() int {
	fake.getPackageMutex.RLock()
	defer fake.getPackageMutex.RUnlock()
	return len(fake.getPackageArgsForCall)
}

func (fake *FakeDedicatedHostManager) GetPackageCalls(stub func() (datatypes.Product_Package, error)) {
	fake.getPackageMutex.Lock()
	defer fake.getPackageMutex.Unlock()
	fake.GetPackageStub = stub
}

func (fake *FakeDedicatedHostManager) GetPackageReturns(result1 datatypes.Product_Package, result2 error) {
	fake.getPackageMutex.Lock()
	defer fake.getPackageMutex.Unlock()
	fake.GetPackageStub = nil
	fake.getPackageReturns = struct {
		result1 datatypes.Product_Package
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetPackageReturnsOnCall(i int, result1 datatypes.Product_Package, result2 error) {
	fake.getPackageMutex.Lock()
	defer fake.getPackageMutex.Unlock()
	fake.GetPackageStub = nil
	if fake.getPackageReturnsOnCall == nil {
		fake.getPackageReturnsOnCall = make(map[int]struct {
			result1 datatypes.Product_Package
			result2 error
		})
	}
	fake.getPackageReturnsOnCall[i] = struct {
		result1 datatypes.Product_Package
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetVlansOptions(arg1 string, arg2 string, arg3 datatypes.Product_Package) ([]datatypes.Network_Vlan, error) {
	fake.getVlansOptionsMutex.Lock()
	ret, specificReturn := fake.getVlansOptionsReturnsOnCall[len(fake.getVlansOptionsArgsForCall)]
	fake.getVlansOptionsArgsForCall = append(fake.getVlansOptionsArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 datatypes.Product_Package
	}{arg1, arg2, arg3})
	stub := fake.GetVlansOptionsStub
	fakeReturns := fake.getVlansOptionsReturns
	fake.recordInvocation("GetVlansOptions", []interface{}{arg1, arg2, arg3})
	fake.getVlansOptionsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) GetVlansOptionsCallCount() int {
	fake.getVlansOptionsMutex.RLock()
	defer fake.getVlansOptionsMutex.RUnlock()
	return len(fake.getVlansOptionsArgsForCall)
}

func (fake *FakeDedicatedHostManager) GetVlansOptionsCalls(stub func(string, string, datatypes.Product_Package) ([]datatypes.Network_Vlan, error)) {
	fake.getVlansOptionsMutex.Lock()
	defer fake.getVlansOptionsMutex.Unlock()
	fake.GetVlansOptionsStub = stub
}

func (fake *FakeDedicatedHostManager) GetVlansOptionsArgsForCall(i int) (string, string, datatypes.Product_Package) {
	fake.getVlansOptionsMutex.RLock()
	defer fake.getVlansOptionsMutex.RUnlock()
	argsForCall := fake.getVlansOptionsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeDedicatedHostManager) GetVlansOptionsReturns(result1 []datatypes.Network_Vlan, result2 error) {
	fake.getVlansOptionsMutex.Lock()
	defer fake.getVlansOptionsMutex.Unlock()
	fake.GetVlansOptionsStub = nil
	fake.getVlansOptionsReturns = struct {
		result1 []datatypes.Network_Vlan
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) GetVlansOptionsReturnsOnCall(i int, result1 []datatypes.Network_Vlan, result2 error) {
	fake.getVlansOptionsMutex.Lock()
	defer fake.getVlansOptionsMutex.Unlock()
	fake.GetVlansOptionsStub = nil
	if fake.getVlansOptionsReturnsOnCall == nil {
		fake.getVlansOptionsReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Network_Vlan
			result2 error
		})
	}
	fake.getVlansOptionsReturnsOnCall[i] = struct {
		result1 []datatypes.Network_Vlan
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) ListDedicatedHost(arg1 string, arg2 string, arg3 string, arg4 int) ([]datatypes.Virtual_DedicatedHost, error) {
	fake.listDedicatedHostMutex.Lock()
	ret, specificReturn := fake.listDedicatedHostReturnsOnCall[len(fake.listDedicatedHostArgsForCall)]
	fake.listDedicatedHostArgsForCall = append(fake.listDedicatedHostArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 int
	}{arg1, arg2, arg3, arg4})
	stub := fake.ListDedicatedHostStub
	fakeReturns := fake.listDedicatedHostReturns
	fake.recordInvocation("ListDedicatedHost", []interface{}{arg1, arg2, arg3, arg4})
	fake.listDedicatedHostMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) ListDedicatedHostCallCount() int {
	fake.listDedicatedHostMutex.RLock()
	defer fake.listDedicatedHostMutex.RUnlock()
	return len(fake.listDedicatedHostArgsForCall)
}

func (fake *FakeDedicatedHostManager) ListDedicatedHostCalls(stub func(string, string, string, int) ([]datatypes.Virtual_DedicatedHost, error)) {
	fake.listDedicatedHostMutex.Lock()
	defer fake.listDedicatedHostMutex.Unlock()
	fake.ListDedicatedHostStub = stub
}

func (fake *FakeDedicatedHostManager) ListDedicatedHostArgsForCall(i int) (string, string, string, int) {
	fake.listDedicatedHostMutex.RLock()
	defer fake.listDedicatedHostMutex.RUnlock()
	argsForCall := fake.listDedicatedHostArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeDedicatedHostManager) ListDedicatedHostReturns(result1 []datatypes.Virtual_DedicatedHost, result2 error) {
	fake.listDedicatedHostMutex.Lock()
	defer fake.listDedicatedHostMutex.Unlock()
	fake.ListDedicatedHostStub = nil
	fake.listDedicatedHostReturns = struct {
		result1 []datatypes.Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) ListDedicatedHostReturnsOnCall(i int, result1 []datatypes.Virtual_DedicatedHost, result2 error) {
	fake.listDedicatedHostMutex.Lock()
	defer fake.listDedicatedHostMutex.Unlock()
	fake.ListDedicatedHostStub = nil
	if fake.listDedicatedHostReturnsOnCall == nil {
		fake.listDedicatedHostReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Virtual_DedicatedHost
			result2 error
		})
	}
	fake.listDedicatedHostReturnsOnCall[i] = struct {
		result1 []datatypes.Virtual_DedicatedHost
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) ListGuests(arg1 int, arg2 int, arg3 string, arg4 string, arg5 int, arg6 []string, arg7 string) ([]datatypes.Virtual_Guest, error) {
	var arg6Copy []string
	if arg6 != nil {
		arg6Copy = make([]string, len(arg6))
		copy(arg6Copy, arg6)
	}
	fake.listGuestsMutex.Lock()
	ret, specificReturn := fake.listGuestsReturnsOnCall[len(fake.listGuestsArgsForCall)]
	fake.listGuestsArgsForCall = append(fake.listGuestsArgsForCall, struct {
		arg1 int
		arg2 int
		arg3 string
		arg4 string
		arg5 int
		arg6 []string
		arg7 string
	}{arg1, arg2, arg3, arg4, arg5, arg6Copy, arg7})
	stub := fake.ListGuestsStub
	fakeReturns := fake.listGuestsReturns
	fake.recordInvocation("ListGuests", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6Copy, arg7})
	fake.listGuestsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) ListGuestsCallCount() int {
	fake.listGuestsMutex.RLock()
	defer fake.listGuestsMutex.RUnlock()
	return len(fake.listGuestsArgsForCall)
}

func (fake *FakeDedicatedHostManager) ListGuestsCalls(stub func(int, int, string, string, int, []string, string) ([]datatypes.Virtual_Guest, error)) {
	fake.listGuestsMutex.Lock()
	defer fake.listGuestsMutex.Unlock()
	fake.ListGuestsStub = stub
}

func (fake *FakeDedicatedHostManager) ListGuestsArgsForCall(i int) (int, int, string, string, int, []string, string) {
	fake.listGuestsMutex.RLock()
	defer fake.listGuestsMutex.RUnlock()
	argsForCall := fake.listGuestsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7
}

func (fake *FakeDedicatedHostManager) ListGuestsReturns(result1 []datatypes.Virtual_Guest, result2 error) {
	fake.listGuestsMutex.Lock()
	defer fake.listGuestsMutex.Unlock()
	fake.ListGuestsStub = nil
	fake.listGuestsReturns = struct {
		result1 []datatypes.Virtual_Guest
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) ListGuestsReturnsOnCall(i int, result1 []datatypes.Virtual_Guest, result2 error) {
	fake.listGuestsMutex.Lock()
	defer fake.listGuestsMutex.Unlock()
	fake.ListGuestsStub = nil
	if fake.listGuestsReturnsOnCall == nil {
		fake.listGuestsReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Virtual_Guest
			result2 error
		})
	}
	fake.listGuestsReturnsOnCall[i] = struct {
		result1 []datatypes.Virtual_Guest
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) OrderInstance(arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order_Receipt, error) {
	fake.orderInstanceMutex.Lock()
	ret, specificReturn := fake.orderInstanceReturnsOnCall[len(fake.orderInstanceArgsForCall)]
	fake.orderInstanceArgsForCall = append(fake.orderInstanceArgsForCall, struct {
		arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost
	}{arg1})
	stub := fake.OrderInstanceStub
	fakeReturns := fake.orderInstanceReturns
	fake.recordInvocation("OrderInstance", []interface{}{arg1})
	fake.orderInstanceMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) OrderInstanceCallCount() int {
	fake.orderInstanceMutex.RLock()
	defer fake.orderInstanceMutex.RUnlock()
	return len(fake.orderInstanceArgsForCall)
}

func (fake *FakeDedicatedHostManager) OrderInstanceCalls(stub func(datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order_Receipt, error)) {
	fake.orderInstanceMutex.Lock()
	defer fake.orderInstanceMutex.Unlock()
	fake.OrderInstanceStub = stub
}

func (fake *FakeDedicatedHostManager) OrderInstanceArgsForCall(i int) datatypes.Container_Product_Order_Virtual_DedicatedHost {
	fake.orderInstanceMutex.RLock()
	defer fake.orderInstanceMutex.RUnlock()
	argsForCall := fake.orderInstanceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDedicatedHostManager) OrderInstanceReturns(result1 datatypes.Container_Product_Order_Receipt, result2 error) {
	fake.orderInstanceMutex.Lock()
	defer fake.orderInstanceMutex.Unlock()
	fake.OrderInstanceStub = nil
	fake.orderInstanceReturns = struct {
		result1 datatypes.Container_Product_Order_Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) OrderInstanceReturnsOnCall(i int, result1 datatypes.Container_Product_Order_Receipt, result2 error) {
	fake.orderInstanceMutex.Lock()
	defer fake.orderInstanceMutex.Unlock()
	fake.OrderInstanceStub = nil
	if fake.orderInstanceReturnsOnCall == nil {
		fake.orderInstanceReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Product_Order_Receipt
			result2 error
		})
	}
	fake.orderInstanceReturnsOnCall[i] = struct {
		result1 datatypes.Container_Product_Order_Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreation(arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order, error) {
	fake.verifyInstanceCreationMutex.Lock()
	ret, specificReturn := fake.verifyInstanceCreationReturnsOnCall[len(fake.verifyInstanceCreationArgsForCall)]
	fake.verifyInstanceCreationArgsForCall = append(fake.verifyInstanceCreationArgsForCall, struct {
		arg1 datatypes.Container_Product_Order_Virtual_DedicatedHost
	}{arg1})
	stub := fake.VerifyInstanceCreationStub
	fakeReturns := fake.verifyInstanceCreationReturns
	fake.recordInvocation("VerifyInstanceCreation", []interface{}{arg1})
	fake.verifyInstanceCreationMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreationCallCount() int {
	fake.verifyInstanceCreationMutex.RLock()
	defer fake.verifyInstanceCreationMutex.RUnlock()
	return len(fake.verifyInstanceCreationArgsForCall)
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreationCalls(stub func(datatypes.Container_Product_Order_Virtual_DedicatedHost) (datatypes.Container_Product_Order, error)) {
	fake.verifyInstanceCreationMutex.Lock()
	defer fake.verifyInstanceCreationMutex.Unlock()
	fake.VerifyInstanceCreationStub = stub
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreationArgsForCall(i int) datatypes.Container_Product_Order_Virtual_DedicatedHost {
	fake.verifyInstanceCreationMutex.RLock()
	defer fake.verifyInstanceCreationMutex.RUnlock()
	argsForCall := fake.verifyInstanceCreationArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreationReturns(result1 datatypes.Container_Product_Order, result2 error) {
	fake.verifyInstanceCreationMutex.Lock()
	defer fake.verifyInstanceCreationMutex.Unlock()
	fake.VerifyInstanceCreationStub = nil
	fake.verifyInstanceCreationReturns = struct {
		result1 datatypes.Container_Product_Order
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) VerifyInstanceCreationReturnsOnCall(i int, result1 datatypes.Container_Product_Order, result2 error) {
	fake.verifyInstanceCreationMutex.Lock()
	defer fake.verifyInstanceCreationMutex.Unlock()
	fake.VerifyInstanceCreationStub = nil
	if fake.verifyInstanceCreationReturnsOnCall == nil {
		fake.verifyInstanceCreationReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Product_Order
			result2 error
		})
	}
	fake.verifyInstanceCreationReturnsOnCall[i] = struct {
		result1 datatypes.Container_Product_Order
		result2 error
	}{result1, result2}
}

func (fake *FakeDedicatedHostManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cancelGuestsMutex.RLock()
	defer fake.cancelGuestsMutex.RUnlock()
	fake.deleteHostMutex.RLock()
	defer fake.deleteHostMutex.RUnlock()
	fake.generateOrderTemplateMutex.RLock()
	defer fake.generateOrderTemplateMutex.RUnlock()
	fake.getCreateOptionsMutex.RLock()
	defer fake.getCreateOptionsMutex.RUnlock()
	fake.getInstanceMutex.RLock()
	defer fake.getInstanceMutex.RUnlock()
	fake.getPackageMutex.RLock()
	defer fake.getPackageMutex.RUnlock()
	fake.getVlansOptionsMutex.RLock()
	defer fake.getVlansOptionsMutex.RUnlock()
	fake.listDedicatedHostMutex.RLock()
	defer fake.listDedicatedHostMutex.RUnlock()
	fake.listGuestsMutex.RLock()
	defer fake.listGuestsMutex.RUnlock()
	fake.orderInstanceMutex.RLock()
	defer fake.orderInstanceMutex.RUnlock()
	fake.verifyInstanceCreationMutex.RLock()
	defer fake.verifyInstanceCreationMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDedicatedHostManager) recordInvocation(key string, args []interface{}) {
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

var _ managers.DedicatedHostManager = new(FakeDedicatedHostManager)
