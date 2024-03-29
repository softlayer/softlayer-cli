// Code generated by counterfeiter. DO NOT EDIT.
package testhelpers

import (
	"sync"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type FakeCdnManager struct {
	CreateCdnStub        func(string, string, string, int, int, string, string, string, string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	createCdnMutex       sync.RWMutex
	createCdnArgsForCall []struct {
		arg1  string
		arg2  string
		arg3  string
		arg4  int
		arg5  int
		arg6  string
		arg7  string
		arg8  string
		arg9  string
		arg10 string
	}
	createCdnReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	createCdnReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	DeleteCDNStub        func(string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	deleteCDNMutex       sync.RWMutex
	deleteCDNArgsForCall []struct {
		arg1 string
	}
	deleteCDNReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	deleteCDNReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	EditCDNStub        func(int, string, int, int, string, string, string, string, string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	editCDNMutex       sync.RWMutex
	editCDNArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 int
		arg4 int
		arg5 string
		arg6 string
		arg7 string
		arg8 string
		arg9 string
	}
	editCDNReturns struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	editCDNReturnsOnCall map[int]struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	GetDetailCDNStub        func(int, string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	getDetailCDNMutex       sync.RWMutex
	getDetailCDNArgsForCall []struct {
		arg1 int
		arg2 string
	}
	getDetailCDNReturns struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	getDetailCDNReturnsOnCall map[int]struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	GetNetworkCdnMarketplaceConfigurationMappingStub        func() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	getNetworkCdnMarketplaceConfigurationMappingMutex       sync.RWMutex
	getNetworkCdnMarketplaceConfigurationMappingArgsForCall []struct {
	}
	getNetworkCdnMarketplaceConfigurationMappingReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	getNetworkCdnMarketplaceConfigurationMappingReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}
	GetOriginsStub        func(string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error)
	getOriginsMutex       sync.RWMutex
	getOriginsArgsForCall []struct {
		arg1 string
	}
	getOriginsReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}
	getOriginsReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}
	GetUsageMetricsStub        func(int, int, string) (datatypes.Container_Network_CdnMarketplace_Metrics, error)
	getUsageMetricsMutex       sync.RWMutex
	getUsageMetricsArgsForCall []struct {
		arg1 int
		arg2 int
		arg3 string
	}
	getUsageMetricsReturns struct {
		result1 datatypes.Container_Network_CdnMarketplace_Metrics
		result2 error
	}
	getUsageMetricsReturnsOnCall map[int]struct {
		result1 datatypes.Container_Network_CdnMarketplace_Metrics
		result2 error
	}
	OriginAddCdnStub        func(string, string, string, string, string, int, int, string, string, string, bool, bool, string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error)
	originAddCdnMutex       sync.RWMutex
	originAddCdnArgsForCall []struct {
		arg1  string
		arg2  string
		arg3  string
		arg4  string
		arg5  string
		arg6  int
		arg7  int
		arg8  string
		arg9  string
		arg10 string
		arg11 bool
		arg12 bool
		arg13 string
		arg14 string
	}
	originAddCdnReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}
	originAddCdnReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}
	PurgeStub        func(string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge, error)
	purgeMutex       sync.RWMutex
	purgeArgsForCall []struct {
		arg1 string
		arg2 string
	}
	purgeReturns struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge
		result2 error
	}
	purgeReturnsOnCall map[int]struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge
		result2 error
	}
	RemoveOriginStub        func(string, string) (string, error)
	removeOriginMutex       sync.RWMutex
	removeOriginArgsForCall []struct {
		arg1 string
		arg2 string
	}
	removeOriginReturns struct {
		result1 string
		result2 error
	}
	removeOriginReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCdnManager) CreateCdn(arg1 string, arg2 string, arg3 string, arg4 int, arg5 int, arg6 string, arg7 string, arg8 string, arg9 string, arg10 string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	fake.createCdnMutex.Lock()
	ret, specificReturn := fake.createCdnReturnsOnCall[len(fake.createCdnArgsForCall)]
	fake.createCdnArgsForCall = append(fake.createCdnArgsForCall, struct {
		arg1  string
		arg2  string
		arg3  string
		arg4  int
		arg5  int
		arg6  string
		arg7  string
		arg8  string
		arg9  string
		arg10 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10})
	stub := fake.CreateCdnStub
	fakeReturns := fake.createCdnReturns
	fake.recordInvocation("CreateCdn", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10})
	fake.createCdnMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) CreateCdnCallCount() int {
	fake.createCdnMutex.RLock()
	defer fake.createCdnMutex.RUnlock()
	return len(fake.createCdnArgsForCall)
}

func (fake *FakeCdnManager) CreateCdnCalls(stub func(string, string, string, int, int, string, string, string, string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)) {
	fake.createCdnMutex.Lock()
	defer fake.createCdnMutex.Unlock()
	fake.CreateCdnStub = stub
}

func (fake *FakeCdnManager) CreateCdnArgsForCall(i int) (string, string, string, int, int, string, string, string, string, string) {
	fake.createCdnMutex.RLock()
	defer fake.createCdnMutex.RUnlock()
	argsForCall := fake.createCdnArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9, argsForCall.arg10
}

func (fake *FakeCdnManager) CreateCdnReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.createCdnMutex.Lock()
	defer fake.createCdnMutex.Unlock()
	fake.CreateCdnStub = nil
	fake.createCdnReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) CreateCdnReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.createCdnMutex.Lock()
	defer fake.createCdnMutex.Unlock()
	fake.CreateCdnStub = nil
	if fake.createCdnReturnsOnCall == nil {
		fake.createCdnReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
			result2 error
		})
	}
	fake.createCdnReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) DeleteCDN(arg1 string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	fake.deleteCDNMutex.Lock()
	ret, specificReturn := fake.deleteCDNReturnsOnCall[len(fake.deleteCDNArgsForCall)]
	fake.deleteCDNArgsForCall = append(fake.deleteCDNArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DeleteCDNStub
	fakeReturns := fake.deleteCDNReturns
	fake.recordInvocation("DeleteCDN", []interface{}{arg1})
	fake.deleteCDNMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) DeleteCDNCallCount() int {
	fake.deleteCDNMutex.RLock()
	defer fake.deleteCDNMutex.RUnlock()
	return len(fake.deleteCDNArgsForCall)
}

func (fake *FakeCdnManager) DeleteCDNCalls(stub func(string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)) {
	fake.deleteCDNMutex.Lock()
	defer fake.deleteCDNMutex.Unlock()
	fake.DeleteCDNStub = stub
}

func (fake *FakeCdnManager) DeleteCDNArgsForCall(i int) string {
	fake.deleteCDNMutex.RLock()
	defer fake.deleteCDNMutex.RUnlock()
	argsForCall := fake.deleteCDNArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCdnManager) DeleteCDNReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.deleteCDNMutex.Lock()
	defer fake.deleteCDNMutex.Unlock()
	fake.DeleteCDNStub = nil
	fake.deleteCDNReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) DeleteCDNReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.deleteCDNMutex.Lock()
	defer fake.deleteCDNMutex.Unlock()
	fake.DeleteCDNStub = nil
	if fake.deleteCDNReturnsOnCall == nil {
		fake.deleteCDNReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
			result2 error
		})
	}
	fake.deleteCDNReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) EditCDN(arg1 int, arg2 string, arg3 int, arg4 int, arg5 string, arg6 string, arg7 string, arg8 string, arg9 string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	fake.editCDNMutex.Lock()
	ret, specificReturn := fake.editCDNReturnsOnCall[len(fake.editCDNArgsForCall)]
	fake.editCDNArgsForCall = append(fake.editCDNArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 int
		arg4 int
		arg5 string
		arg6 string
		arg7 string
		arg8 string
		arg9 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9})
	stub := fake.EditCDNStub
	fakeReturns := fake.editCDNReturns
	fake.recordInvocation("EditCDN", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9})
	fake.editCDNMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) EditCDNCallCount() int {
	fake.editCDNMutex.RLock()
	defer fake.editCDNMutex.RUnlock()
	return len(fake.editCDNArgsForCall)
}

func (fake *FakeCdnManager) EditCDNCalls(stub func(int, string, int, int, string, string, string, string, string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)) {
	fake.editCDNMutex.Lock()
	defer fake.editCDNMutex.Unlock()
	fake.EditCDNStub = stub
}

func (fake *FakeCdnManager) EditCDNArgsForCall(i int) (int, string, int, int, string, string, string, string, string) {
	fake.editCDNMutex.RLock()
	defer fake.editCDNMutex.RUnlock()
	argsForCall := fake.editCDNArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9
}

func (fake *FakeCdnManager) EditCDNReturns(result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.editCDNMutex.Lock()
	defer fake.editCDNMutex.Unlock()
	fake.EditCDNStub = nil
	fake.editCDNReturns = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) EditCDNReturnsOnCall(i int, result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.editCDNMutex.Lock()
	defer fake.editCDNMutex.Unlock()
	fake.EditCDNStub = nil
	if fake.editCDNReturnsOnCall == nil {
		fake.editCDNReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
			result2 error
		})
	}
	fake.editCDNReturnsOnCall[i] = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetDetailCDN(arg1 int, arg2 string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	fake.getDetailCDNMutex.Lock()
	ret, specificReturn := fake.getDetailCDNReturnsOnCall[len(fake.getDetailCDNArgsForCall)]
	fake.getDetailCDNArgsForCall = append(fake.getDetailCDNArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	stub := fake.GetDetailCDNStub
	fakeReturns := fake.getDetailCDNReturns
	fake.recordInvocation("GetDetailCDN", []interface{}{arg1, arg2})
	fake.getDetailCDNMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) GetDetailCDNCallCount() int {
	fake.getDetailCDNMutex.RLock()
	defer fake.getDetailCDNMutex.RUnlock()
	return len(fake.getDetailCDNArgsForCall)
}

func (fake *FakeCdnManager) GetDetailCDNCalls(stub func(int, string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)) {
	fake.getDetailCDNMutex.Lock()
	defer fake.getDetailCDNMutex.Unlock()
	fake.GetDetailCDNStub = stub
}

func (fake *FakeCdnManager) GetDetailCDNArgsForCall(i int) (int, string) {
	fake.getDetailCDNMutex.RLock()
	defer fake.getDetailCDNMutex.RUnlock()
	argsForCall := fake.getDetailCDNArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCdnManager) GetDetailCDNReturns(result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.getDetailCDNMutex.Lock()
	defer fake.getDetailCDNMutex.Unlock()
	fake.GetDetailCDNStub = nil
	fake.getDetailCDNReturns = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetDetailCDNReturnsOnCall(i int, result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.getDetailCDNMutex.Lock()
	defer fake.getDetailCDNMutex.Unlock()
	fake.GetDetailCDNStub = nil
	if fake.getDetailCDNReturnsOnCall == nil {
		fake.getDetailCDNReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
			result2 error
		})
	}
	fake.getDetailCDNReturnsOnCall[i] = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetNetworkCdnMarketplaceConfigurationMapping() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Lock()
	ret, specificReturn := fake.getNetworkCdnMarketplaceConfigurationMappingReturnsOnCall[len(fake.getNetworkCdnMarketplaceConfigurationMappingArgsForCall)]
	fake.getNetworkCdnMarketplaceConfigurationMappingArgsForCall = append(fake.getNetworkCdnMarketplaceConfigurationMappingArgsForCall, struct {
	}{})
	stub := fake.GetNetworkCdnMarketplaceConfigurationMappingStub
	fakeReturns := fake.getNetworkCdnMarketplaceConfigurationMappingReturns
	fake.recordInvocation("GetNetworkCdnMarketplaceConfigurationMapping", []interface{}{})
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) GetNetworkCdnMarketplaceConfigurationMappingCallCount() int {
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.RLock()
	defer fake.getNetworkCdnMarketplaceConfigurationMappingMutex.RUnlock()
	return len(fake.getNetworkCdnMarketplaceConfigurationMappingArgsForCall)
}

func (fake *FakeCdnManager) GetNetworkCdnMarketplaceConfigurationMappingCalls(stub func() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)) {
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Lock()
	defer fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Unlock()
	fake.GetNetworkCdnMarketplaceConfigurationMappingStub = stub
}

func (fake *FakeCdnManager) GetNetworkCdnMarketplaceConfigurationMappingReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Lock()
	defer fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Unlock()
	fake.GetNetworkCdnMarketplaceConfigurationMappingStub = nil
	fake.getNetworkCdnMarketplaceConfigurationMappingReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetNetworkCdnMarketplaceConfigurationMappingReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, result2 error) {
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Lock()
	defer fake.getNetworkCdnMarketplaceConfigurationMappingMutex.Unlock()
	fake.GetNetworkCdnMarketplaceConfigurationMappingStub = nil
	if fake.getNetworkCdnMarketplaceConfigurationMappingReturnsOnCall == nil {
		fake.getNetworkCdnMarketplaceConfigurationMappingReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
			result2 error
		})
	}
	fake.getNetworkCdnMarketplaceConfigurationMappingReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetOrigins(arg1 string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error) {
	fake.getOriginsMutex.Lock()
	ret, specificReturn := fake.getOriginsReturnsOnCall[len(fake.getOriginsArgsForCall)]
	fake.getOriginsArgsForCall = append(fake.getOriginsArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetOriginsStub
	fakeReturns := fake.getOriginsReturns
	fake.recordInvocation("GetOrigins", []interface{}{arg1})
	fake.getOriginsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) GetOriginsCallCount() int {
	fake.getOriginsMutex.RLock()
	defer fake.getOriginsMutex.RUnlock()
	return len(fake.getOriginsArgsForCall)
}

func (fake *FakeCdnManager) GetOriginsCalls(stub func(string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error)) {
	fake.getOriginsMutex.Lock()
	defer fake.getOriginsMutex.Unlock()
	fake.GetOriginsStub = stub
}

func (fake *FakeCdnManager) GetOriginsArgsForCall(i int) string {
	fake.getOriginsMutex.RLock()
	defer fake.getOriginsMutex.RUnlock()
	argsForCall := fake.getOriginsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeCdnManager) GetOriginsReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, result2 error) {
	fake.getOriginsMutex.Lock()
	defer fake.getOriginsMutex.Unlock()
	fake.GetOriginsStub = nil
	fake.getOriginsReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetOriginsReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, result2 error) {
	fake.getOriginsMutex.Lock()
	defer fake.getOriginsMutex.Unlock()
	fake.GetOriginsStub = nil
	if fake.getOriginsReturnsOnCall == nil {
		fake.getOriginsReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
			result2 error
		})
	}
	fake.getOriginsReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetUsageMetrics(arg1 int, arg2 int, arg3 string) (datatypes.Container_Network_CdnMarketplace_Metrics, error) {
	fake.getUsageMetricsMutex.Lock()
	ret, specificReturn := fake.getUsageMetricsReturnsOnCall[len(fake.getUsageMetricsArgsForCall)]
	fake.getUsageMetricsArgsForCall = append(fake.getUsageMetricsArgsForCall, struct {
		arg1 int
		arg2 int
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetUsageMetricsStub
	fakeReturns := fake.getUsageMetricsReturns
	fake.recordInvocation("GetUsageMetrics", []interface{}{arg1, arg2, arg3})
	fake.getUsageMetricsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) GetUsageMetricsCallCount() int {
	fake.getUsageMetricsMutex.RLock()
	defer fake.getUsageMetricsMutex.RUnlock()
	return len(fake.getUsageMetricsArgsForCall)
}

func (fake *FakeCdnManager) GetUsageMetricsCalls(stub func(int, int, string) (datatypes.Container_Network_CdnMarketplace_Metrics, error)) {
	fake.getUsageMetricsMutex.Lock()
	defer fake.getUsageMetricsMutex.Unlock()
	fake.GetUsageMetricsStub = stub
}

func (fake *FakeCdnManager) GetUsageMetricsArgsForCall(i int) (int, int, string) {
	fake.getUsageMetricsMutex.RLock()
	defer fake.getUsageMetricsMutex.RUnlock()
	argsForCall := fake.getUsageMetricsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeCdnManager) GetUsageMetricsReturns(result1 datatypes.Container_Network_CdnMarketplace_Metrics, result2 error) {
	fake.getUsageMetricsMutex.Lock()
	defer fake.getUsageMetricsMutex.Unlock()
	fake.GetUsageMetricsStub = nil
	fake.getUsageMetricsReturns = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Metrics
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) GetUsageMetricsReturnsOnCall(i int, result1 datatypes.Container_Network_CdnMarketplace_Metrics, result2 error) {
	fake.getUsageMetricsMutex.Lock()
	defer fake.getUsageMetricsMutex.Unlock()
	fake.GetUsageMetricsStub = nil
	if fake.getUsageMetricsReturnsOnCall == nil {
		fake.getUsageMetricsReturnsOnCall = make(map[int]struct {
			result1 datatypes.Container_Network_CdnMarketplace_Metrics
			result2 error
		})
	}
	fake.getUsageMetricsReturnsOnCall[i] = struct {
		result1 datatypes.Container_Network_CdnMarketplace_Metrics
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) OriginAddCdn(arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 int, arg7 int, arg8 string, arg9 string, arg10 string, arg11 bool, arg12 bool, arg13 string, arg14 string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error) {
	fake.originAddCdnMutex.Lock()
	ret, specificReturn := fake.originAddCdnReturnsOnCall[len(fake.originAddCdnArgsForCall)]
	fake.originAddCdnArgsForCall = append(fake.originAddCdnArgsForCall, struct {
		arg1  string
		arg2  string
		arg3  string
		arg4  string
		arg5  string
		arg6  int
		arg7  int
		arg8  string
		arg9  string
		arg10 string
		arg11 bool
		arg12 bool
		arg13 string
		arg14 string
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14})
	stub := fake.OriginAddCdnStub
	fakeReturns := fake.originAddCdnReturns
	fake.recordInvocation("OriginAddCdn", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14})
	fake.originAddCdnMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) OriginAddCdnCallCount() int {
	fake.originAddCdnMutex.RLock()
	defer fake.originAddCdnMutex.RUnlock()
	return len(fake.originAddCdnArgsForCall)
}

func (fake *FakeCdnManager) OriginAddCdnCalls(stub func(string, string, string, string, string, int, int, string, string, string, bool, bool, string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error)) {
	fake.originAddCdnMutex.Lock()
	defer fake.originAddCdnMutex.Unlock()
	fake.OriginAddCdnStub = stub
}

func (fake *FakeCdnManager) OriginAddCdnArgsForCall(i int) (string, string, string, string, string, int, int, string, string, string, bool, bool, string, string) {
	fake.originAddCdnMutex.RLock()
	defer fake.originAddCdnMutex.RUnlock()
	argsForCall := fake.originAddCdnArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9, argsForCall.arg10, argsForCall.arg11, argsForCall.arg12, argsForCall.arg13, argsForCall.arg14
}

func (fake *FakeCdnManager) OriginAddCdnReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, result2 error) {
	fake.originAddCdnMutex.Lock()
	defer fake.originAddCdnMutex.Unlock()
	fake.OriginAddCdnStub = nil
	fake.originAddCdnReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) OriginAddCdnReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, result2 error) {
	fake.originAddCdnMutex.Lock()
	defer fake.originAddCdnMutex.Unlock()
	fake.OriginAddCdnStub = nil
	if fake.originAddCdnReturnsOnCall == nil {
		fake.originAddCdnReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
			result2 error
		})
	}
	fake.originAddCdnReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) Purge(arg1 string, arg2 string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge, error) {
	fake.purgeMutex.Lock()
	ret, specificReturn := fake.purgeReturnsOnCall[len(fake.purgeArgsForCall)]
	fake.purgeArgsForCall = append(fake.purgeArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.PurgeStub
	fakeReturns := fake.purgeReturns
	fake.recordInvocation("Purge", []interface{}{arg1, arg2})
	fake.purgeMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) PurgeCallCount() int {
	fake.purgeMutex.RLock()
	defer fake.purgeMutex.RUnlock()
	return len(fake.purgeArgsForCall)
}

func (fake *FakeCdnManager) PurgeCalls(stub func(string, string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge, error)) {
	fake.purgeMutex.Lock()
	defer fake.purgeMutex.Unlock()
	fake.PurgeStub = stub
}

func (fake *FakeCdnManager) PurgeArgsForCall(i int) (string, string) {
	fake.purgeMutex.RLock()
	defer fake.purgeMutex.RUnlock()
	argsForCall := fake.purgeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCdnManager) PurgeReturns(result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge, result2 error) {
	fake.purgeMutex.Lock()
	defer fake.purgeMutex.Unlock()
	fake.PurgeStub = nil
	fake.purgeReturns = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) PurgeReturnsOnCall(i int, result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge, result2 error) {
	fake.purgeMutex.Lock()
	defer fake.purgeMutex.Unlock()
	fake.PurgeStub = nil
	if fake.purgeReturnsOnCall == nil {
		fake.purgeReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge
			result2 error
		})
	}
	fake.purgeReturnsOnCall[i] = struct {
		result1 []datatypes.Container_Network_CdnMarketplace_Configuration_Cache_Purge
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) RemoveOrigin(arg1 string, arg2 string) (string, error) {
	fake.removeOriginMutex.Lock()
	ret, specificReturn := fake.removeOriginReturnsOnCall[len(fake.removeOriginArgsForCall)]
	fake.removeOriginArgsForCall = append(fake.removeOriginArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.RemoveOriginStub
	fakeReturns := fake.removeOriginReturns
	fake.recordInvocation("RemoveOrigin", []interface{}{arg1, arg2})
	fake.removeOriginMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCdnManager) RemoveOriginCallCount() int {
	fake.removeOriginMutex.RLock()
	defer fake.removeOriginMutex.RUnlock()
	return len(fake.removeOriginArgsForCall)
}

func (fake *FakeCdnManager) RemoveOriginCalls(stub func(string, string) (string, error)) {
	fake.removeOriginMutex.Lock()
	defer fake.removeOriginMutex.Unlock()
	fake.RemoveOriginStub = stub
}

func (fake *FakeCdnManager) RemoveOriginArgsForCall(i int) (string, string) {
	fake.removeOriginMutex.RLock()
	defer fake.removeOriginMutex.RUnlock()
	argsForCall := fake.removeOriginArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeCdnManager) RemoveOriginReturns(result1 string, result2 error) {
	fake.removeOriginMutex.Lock()
	defer fake.removeOriginMutex.Unlock()
	fake.RemoveOriginStub = nil
	fake.removeOriginReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) RemoveOriginReturnsOnCall(i int, result1 string, result2 error) {
	fake.removeOriginMutex.Lock()
	defer fake.removeOriginMutex.Unlock()
	fake.RemoveOriginStub = nil
	if fake.removeOriginReturnsOnCall == nil {
		fake.removeOriginReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.removeOriginReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeCdnManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createCdnMutex.RLock()
	defer fake.createCdnMutex.RUnlock()
	fake.deleteCDNMutex.RLock()
	defer fake.deleteCDNMutex.RUnlock()
	fake.editCDNMutex.RLock()
	defer fake.editCDNMutex.RUnlock()
	fake.getDetailCDNMutex.RLock()
	defer fake.getDetailCDNMutex.RUnlock()
	fake.getNetworkCdnMarketplaceConfigurationMappingMutex.RLock()
	defer fake.getNetworkCdnMarketplaceConfigurationMappingMutex.RUnlock()
	fake.getOriginsMutex.RLock()
	defer fake.getOriginsMutex.RUnlock()
	fake.getUsageMetricsMutex.RLock()
	defer fake.getUsageMetricsMutex.RUnlock()
	fake.originAddCdnMutex.RLock()
	defer fake.originAddCdnMutex.RUnlock()
	fake.purgeMutex.RLock()
	defer fake.purgeMutex.RUnlock()
	fake.removeOriginMutex.RLock()
	defer fake.removeOriginMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCdnManager) recordInvocation(key string, args []interface{}) {
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

var _ managers.CdnManager = new(FakeCdnManager)
