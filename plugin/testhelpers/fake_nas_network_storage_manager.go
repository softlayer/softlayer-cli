// Code generated by counterfeiter. DO NOT EDIT.
package testhelpers

import (
	"sync"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type FakeNasNetworkStorageManager struct {
	GetNasNetworkStorageStub        func(int, string) (datatypes.Network_Storage, error)
	getNasNetworkStorageMutex       sync.RWMutex
	getNasNetworkStorageArgsForCall []struct {
		arg1 int
		arg2 string
	}
	getNasNetworkStorageReturns struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	getNasNetworkStorageReturnsOnCall map[int]struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	ListNasNetworkStoragesStub        func(string) ([]datatypes.Network_Storage, error)
	listNasNetworkStoragesMutex       sync.RWMutex
	listNasNetworkStoragesArgsForCall []struct {
		arg1 string
	}
	listNasNetworkStoragesReturns struct {
		result1 []datatypes.Network_Storage
		result2 error
	}
	listNasNetworkStoragesReturnsOnCall map[int]struct {
		result1 []datatypes.Network_Storage
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorage(arg1 int, arg2 string) (datatypes.Network_Storage, error) {
	fake.getNasNetworkStorageMutex.Lock()
	ret, specificReturn := fake.getNasNetworkStorageReturnsOnCall[len(fake.getNasNetworkStorageArgsForCall)]
	fake.getNasNetworkStorageArgsForCall = append(fake.getNasNetworkStorageArgsForCall, struct {
		arg1 int
		arg2 string
	}{arg1, arg2})
	stub := fake.GetNasNetworkStorageStub
	fakeReturns := fake.getNasNetworkStorageReturns
	fake.recordInvocation("GetNasNetworkStorage", []interface{}{arg1, arg2})
	fake.getNasNetworkStorageMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorageCallCount() int {
	fake.getNasNetworkStorageMutex.RLock()
	defer fake.getNasNetworkStorageMutex.RUnlock()
	return len(fake.getNasNetworkStorageArgsForCall)
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorageCalls(stub func(int, string) (datatypes.Network_Storage, error)) {
	fake.getNasNetworkStorageMutex.Lock()
	defer fake.getNasNetworkStorageMutex.Unlock()
	fake.GetNasNetworkStorageStub = stub
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorageArgsForCall(i int) (int, string) {
	fake.getNasNetworkStorageMutex.RLock()
	defer fake.getNasNetworkStorageMutex.RUnlock()
	argsForCall := fake.getNasNetworkStorageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorageReturns(result1 datatypes.Network_Storage, result2 error) {
	fake.getNasNetworkStorageMutex.Lock()
	defer fake.getNasNetworkStorageMutex.Unlock()
	fake.GetNasNetworkStorageStub = nil
	fake.getNasNetworkStorageReturns = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *FakeNasNetworkStorageManager) GetNasNetworkStorageReturnsOnCall(i int, result1 datatypes.Network_Storage, result2 error) {
	fake.getNasNetworkStorageMutex.Lock()
	defer fake.getNasNetworkStorageMutex.Unlock()
	fake.GetNasNetworkStorageStub = nil
	if fake.getNasNetworkStorageReturnsOnCall == nil {
		fake.getNasNetworkStorageReturnsOnCall = make(map[int]struct {
			result1 datatypes.Network_Storage
			result2 error
		})
	}
	fake.getNasNetworkStorageReturnsOnCall[i] = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStorages(arg1 string) ([]datatypes.Network_Storage, error) {
	fake.listNasNetworkStoragesMutex.Lock()
	ret, specificReturn := fake.listNasNetworkStoragesReturnsOnCall[len(fake.listNasNetworkStoragesArgsForCall)]
	fake.listNasNetworkStoragesArgsForCall = append(fake.listNasNetworkStoragesArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.ListNasNetworkStoragesStub
	fakeReturns := fake.listNasNetworkStoragesReturns
	fake.recordInvocation("ListNasNetworkStorages", []interface{}{arg1})
	fake.listNasNetworkStoragesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStoragesCallCount() int {
	fake.listNasNetworkStoragesMutex.RLock()
	defer fake.listNasNetworkStoragesMutex.RUnlock()
	return len(fake.listNasNetworkStoragesArgsForCall)
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStoragesCalls(stub func(string) ([]datatypes.Network_Storage, error)) {
	fake.listNasNetworkStoragesMutex.Lock()
	defer fake.listNasNetworkStoragesMutex.Unlock()
	fake.ListNasNetworkStoragesStub = stub
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStoragesArgsForCall(i int) string {
	fake.listNasNetworkStoragesMutex.RLock()
	defer fake.listNasNetworkStoragesMutex.RUnlock()
	argsForCall := fake.listNasNetworkStoragesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStoragesReturns(result1 []datatypes.Network_Storage, result2 error) {
	fake.listNasNetworkStoragesMutex.Lock()
	defer fake.listNasNetworkStoragesMutex.Unlock()
	fake.ListNasNetworkStoragesStub = nil
	fake.listNasNetworkStoragesReturns = struct {
		result1 []datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *FakeNasNetworkStorageManager) ListNasNetworkStoragesReturnsOnCall(i int, result1 []datatypes.Network_Storage, result2 error) {
	fake.listNasNetworkStoragesMutex.Lock()
	defer fake.listNasNetworkStoragesMutex.Unlock()
	fake.ListNasNetworkStoragesStub = nil
	if fake.listNasNetworkStoragesReturnsOnCall == nil {
		fake.listNasNetworkStoragesReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Network_Storage
			result2 error
		})
	}
	fake.listNasNetworkStoragesReturnsOnCall[i] = struct {
		result1 []datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *FakeNasNetworkStorageManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getNasNetworkStorageMutex.RLock()
	defer fake.getNasNetworkStorageMutex.RUnlock()
	fake.listNasNetworkStoragesMutex.RLock()
	defer fake.listNasNetworkStoragesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNasNetworkStorageManager) recordInvocation(key string, args []interface{}) {
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

var _ managers.NasNetworkStorageManager = new(FakeNasNetworkStorageManager)