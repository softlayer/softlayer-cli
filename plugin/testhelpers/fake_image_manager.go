// Code generated by counterfeiter. DO NOT EDIT.
package testhelpers

import (
	"sync"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type FakeImageManager struct {
	DeleteImageStub        func(int) error
	deleteImageMutex       sync.RWMutex
	deleteImageArgsForCall []struct {
		arg1 int
	}
	deleteImageReturns struct {
		result1 error
	}
	deleteImageReturnsOnCall map[int]struct {
		result1 error
	}
	EditImageStub        func(int, string, string, string) ([]bool, []string)
	editImageMutex       sync.RWMutex
	editImageArgsForCall []struct {
		arg1 int
		arg2 string
		arg3 string
		arg4 string
	}
	editImageReturns struct {
		result1 []bool
		result2 []string
	}
	editImageReturnsOnCall map[int]struct {
		result1 []bool
		result2 []string
	}
	ExportImageStub        func(int, datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (bool, error)
	exportImageMutex       sync.RWMutex
	exportImageArgsForCall []struct {
		arg1 int
		arg2 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration
	}
	exportImageReturns struct {
		result1 bool
		result2 error
	}
	exportImageReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	GetImageStub        func(int) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	getImageMutex       sync.RWMutex
	getImageArgsForCall []struct {
		arg1 int
	}
	getImageReturns struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	getImageReturnsOnCall map[int]struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	ImportImageStub        func(datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	importImageMutex       sync.RWMutex
	importImageArgsForCall []struct {
		arg1 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration
	}
	importImageReturns struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	importImageReturnsOnCall map[int]struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	ListPrivateImagesStub        func(string, string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	listPrivateImagesMutex       sync.RWMutex
	listPrivateImagesArgsForCall []struct {
		arg1 string
		arg2 string
	}
	listPrivateImagesReturns struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	listPrivateImagesReturnsOnCall map[int]struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	ListPublicImagesStub        func(string, string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)
	listPublicImagesMutex       sync.RWMutex
	listPublicImagesArgsForCall []struct {
		arg1 string
		arg2 string
	}
	listPublicImagesReturns struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	listPublicImagesReturnsOnCall map[int]struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImageManager) DeleteImage(arg1 int) error {
	fake.deleteImageMutex.Lock()
	ret, specificReturn := fake.deleteImageReturnsOnCall[len(fake.deleteImageArgsForCall)]
	fake.deleteImageArgsForCall = append(fake.deleteImageArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("DeleteImage", []interface{}{arg1})
	fake.deleteImageMutex.Unlock()
	if fake.DeleteImageStub != nil {
		return fake.DeleteImageStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteImageReturns
	return fakeReturns.result1
}

func (fake *FakeImageManager) DeleteImageCallCount() int {
	fake.deleteImageMutex.RLock()
	defer fake.deleteImageMutex.RUnlock()
	return len(fake.deleteImageArgsForCall)
}

func (fake *FakeImageManager) DeleteImageCalls(stub func(int) error) {
	fake.deleteImageMutex.Lock()
	defer fake.deleteImageMutex.Unlock()
	fake.DeleteImageStub = stub
}

func (fake *FakeImageManager) DeleteImageArgsForCall(i int) int {
	fake.deleteImageMutex.RLock()
	defer fake.deleteImageMutex.RUnlock()
	argsForCall := fake.deleteImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImageManager) DeleteImageReturns(result1 error) {
	fake.deleteImageMutex.Lock()
	defer fake.deleteImageMutex.Unlock()
	fake.DeleteImageStub = nil
	fake.deleteImageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeImageManager) DeleteImageReturnsOnCall(i int, result1 error) {
	fake.deleteImageMutex.Lock()
	defer fake.deleteImageMutex.Unlock()
	fake.DeleteImageStub = nil
	if fake.deleteImageReturnsOnCall == nil {
		fake.deleteImageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteImageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeImageManager) EditImage(arg1 int, arg2 string, arg3 string, arg4 string) ([]bool, []string) {
	fake.editImageMutex.Lock()
	ret, specificReturn := fake.editImageReturnsOnCall[len(fake.editImageArgsForCall)]
	fake.editImageArgsForCall = append(fake.editImageArgsForCall, struct {
		arg1 int
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("EditImage", []interface{}{arg1, arg2, arg3, arg4})
	fake.editImageMutex.Unlock()
	if fake.EditImageStub != nil {
		return fake.EditImageStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.editImageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) EditImageCallCount() int {
	fake.editImageMutex.RLock()
	defer fake.editImageMutex.RUnlock()
	return len(fake.editImageArgsForCall)
}

func (fake *FakeImageManager) EditImageCalls(stub func(int, string, string, string) ([]bool, []string)) {
	fake.editImageMutex.Lock()
	defer fake.editImageMutex.Unlock()
	fake.EditImageStub = stub
}

func (fake *FakeImageManager) EditImageArgsForCall(i int) (int, string, string, string) {
	fake.editImageMutex.RLock()
	defer fake.editImageMutex.RUnlock()
	argsForCall := fake.editImageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeImageManager) EditImageReturns(result1 []bool, result2 []string) {
	fake.editImageMutex.Lock()
	defer fake.editImageMutex.Unlock()
	fake.EditImageStub = nil
	fake.editImageReturns = struct {
		result1 []bool
		result2 []string
	}{result1, result2}
}

func (fake *FakeImageManager) EditImageReturnsOnCall(i int, result1 []bool, result2 []string) {
	fake.editImageMutex.Lock()
	defer fake.editImageMutex.Unlock()
	fake.EditImageStub = nil
	if fake.editImageReturnsOnCall == nil {
		fake.editImageReturnsOnCall = make(map[int]struct {
			result1 []bool
			result2 []string
		})
	}
	fake.editImageReturnsOnCall[i] = struct {
		result1 []bool
		result2 []string
	}{result1, result2}
}

func (fake *FakeImageManager) ExportImage(arg1 int, arg2 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (bool, error) {
	fake.exportImageMutex.Lock()
	ret, specificReturn := fake.exportImageReturnsOnCall[len(fake.exportImageArgsForCall)]
	fake.exportImageArgsForCall = append(fake.exportImageArgsForCall, struct {
		arg1 int
		arg2 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration
	}{arg1, arg2})
	fake.recordInvocation("ExportImage", []interface{}{arg1, arg2})
	fake.exportImageMutex.Unlock()
	if fake.ExportImageStub != nil {
		return fake.ExportImageStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.exportImageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) ExportImageCallCount() int {
	fake.exportImageMutex.RLock()
	defer fake.exportImageMutex.RUnlock()
	return len(fake.exportImageArgsForCall)
}

func (fake *FakeImageManager) ExportImageCalls(stub func(int, datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (bool, error)) {
	fake.exportImageMutex.Lock()
	defer fake.exportImageMutex.Unlock()
	fake.ExportImageStub = stub
}

func (fake *FakeImageManager) ExportImageArgsForCall(i int) (int, datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) {
	fake.exportImageMutex.RLock()
	defer fake.exportImageMutex.RUnlock()
	argsForCall := fake.exportImageArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImageManager) ExportImageReturns(result1 bool, result2 error) {
	fake.exportImageMutex.Lock()
	defer fake.exportImageMutex.Unlock()
	fake.ExportImageStub = nil
	fake.exportImageReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ExportImageReturnsOnCall(i int, result1 bool, result2 error) {
	fake.exportImageMutex.Lock()
	defer fake.exportImageMutex.Unlock()
	fake.ExportImageStub = nil
	if fake.exportImageReturnsOnCall == nil {
		fake.exportImageReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.exportImageReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) GetImage(arg1 int) (datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	fake.getImageMutex.Lock()
	ret, specificReturn := fake.getImageReturnsOnCall[len(fake.getImageArgsForCall)]
	fake.getImageArgsForCall = append(fake.getImageArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("GetImage", []interface{}{arg1})
	fake.getImageMutex.Unlock()
	if fake.GetImageStub != nil {
		return fake.GetImageStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getImageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) GetImageCallCount() int {
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	return len(fake.getImageArgsForCall)
}

func (fake *FakeImageManager) GetImageCalls(stub func(int) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = stub
}

func (fake *FakeImageManager) GetImageArgsForCall(i int) int {
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	argsForCall := fake.getImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImageManager) GetImageReturns(result1 datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = nil
	fake.getImageReturns = struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) GetImageReturnsOnCall(i int, result1 datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.getImageMutex.Lock()
	defer fake.getImageMutex.Unlock()
	fake.GetImageStub = nil
	if fake.getImageReturnsOnCall == nil {
		fake.getImageReturnsOnCall = make(map[int]struct {
			result1 datatypes.Virtual_Guest_Block_Device_Template_Group
			result2 error
		})
	}
	fake.getImageReturnsOnCall[i] = struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ImportImage(arg1 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	fake.importImageMutex.Lock()
	ret, specificReturn := fake.importImageReturnsOnCall[len(fake.importImageArgsForCall)]
	fake.importImageArgsForCall = append(fake.importImageArgsForCall, struct {
		arg1 datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration
	}{arg1})
	fake.recordInvocation("ImportImage", []interface{}{arg1})
	fake.importImageMutex.Unlock()
	if fake.ImportImageStub != nil {
		return fake.ImportImageStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.importImageReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) ImportImageCallCount() int {
	fake.importImageMutex.RLock()
	defer fake.importImageMutex.RUnlock()
	return len(fake.importImageArgsForCall)
}

func (fake *FakeImageManager) ImportImageCalls(stub func(datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (datatypes.Virtual_Guest_Block_Device_Template_Group, error)) {
	fake.importImageMutex.Lock()
	defer fake.importImageMutex.Unlock()
	fake.ImportImageStub = stub
}

func (fake *FakeImageManager) ImportImageArgsForCall(i int) datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration {
	fake.importImageMutex.RLock()
	defer fake.importImageMutex.RUnlock()
	argsForCall := fake.importImageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImageManager) ImportImageReturns(result1 datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.importImageMutex.Lock()
	defer fake.importImageMutex.Unlock()
	fake.ImportImageStub = nil
	fake.importImageReturns = struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ImportImageReturnsOnCall(i int, result1 datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.importImageMutex.Lock()
	defer fake.importImageMutex.Unlock()
	fake.ImportImageStub = nil
	if fake.importImageReturnsOnCall == nil {
		fake.importImageReturnsOnCall = make(map[int]struct {
			result1 datatypes.Virtual_Guest_Block_Device_Template_Group
			result2 error
		})
	}
	fake.importImageReturnsOnCall[i] = struct {
		result1 datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ListPrivateImages(arg1 string, arg2 string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	fake.listPrivateImagesMutex.Lock()
	ret, specificReturn := fake.listPrivateImagesReturnsOnCall[len(fake.listPrivateImagesArgsForCall)]
	fake.listPrivateImagesArgsForCall = append(fake.listPrivateImagesArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ListPrivateImages", []interface{}{arg1, arg2})
	fake.listPrivateImagesMutex.Unlock()
	if fake.ListPrivateImagesStub != nil {
		return fake.ListPrivateImagesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listPrivateImagesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) ListPrivateImagesCallCount() int {
	fake.listPrivateImagesMutex.RLock()
	defer fake.listPrivateImagesMutex.RUnlock()
	return len(fake.listPrivateImagesArgsForCall)
}

func (fake *FakeImageManager) ListPrivateImagesCalls(stub func(string, string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)) {
	fake.listPrivateImagesMutex.Lock()
	defer fake.listPrivateImagesMutex.Unlock()
	fake.ListPrivateImagesStub = stub
}

func (fake *FakeImageManager) ListPrivateImagesArgsForCall(i int) (string, string) {
	fake.listPrivateImagesMutex.RLock()
	defer fake.listPrivateImagesMutex.RUnlock()
	argsForCall := fake.listPrivateImagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImageManager) ListPrivateImagesReturns(result1 []datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.listPrivateImagesMutex.Lock()
	defer fake.listPrivateImagesMutex.Unlock()
	fake.ListPrivateImagesStub = nil
	fake.listPrivateImagesReturns = struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ListPrivateImagesReturnsOnCall(i int, result1 []datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.listPrivateImagesMutex.Lock()
	defer fake.listPrivateImagesMutex.Unlock()
	fake.ListPrivateImagesStub = nil
	if fake.listPrivateImagesReturnsOnCall == nil {
		fake.listPrivateImagesReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
			result2 error
		})
	}
	fake.listPrivateImagesReturnsOnCall[i] = struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ListPublicImages(arg1 string, arg2 string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error) {
	fake.listPublicImagesMutex.Lock()
	ret, specificReturn := fake.listPublicImagesReturnsOnCall[len(fake.listPublicImagesArgsForCall)]
	fake.listPublicImagesArgsForCall = append(fake.listPublicImagesArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ListPublicImages", []interface{}{arg1, arg2})
	fake.listPublicImagesMutex.Unlock()
	if fake.ListPublicImagesStub != nil {
		return fake.ListPublicImagesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listPublicImagesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImageManager) ListPublicImagesCallCount() int {
	fake.listPublicImagesMutex.RLock()
	defer fake.listPublicImagesMutex.RUnlock()
	return len(fake.listPublicImagesArgsForCall)
}

func (fake *FakeImageManager) ListPublicImagesCalls(stub func(string, string) ([]datatypes.Virtual_Guest_Block_Device_Template_Group, error)) {
	fake.listPublicImagesMutex.Lock()
	defer fake.listPublicImagesMutex.Unlock()
	fake.ListPublicImagesStub = stub
}

func (fake *FakeImageManager) ListPublicImagesArgsForCall(i int) (string, string) {
	fake.listPublicImagesMutex.RLock()
	defer fake.listPublicImagesMutex.RUnlock()
	argsForCall := fake.listPublicImagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImageManager) ListPublicImagesReturns(result1 []datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.listPublicImagesMutex.Lock()
	defer fake.listPublicImagesMutex.Unlock()
	fake.ListPublicImagesStub = nil
	fake.listPublicImagesReturns = struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) ListPublicImagesReturnsOnCall(i int, result1 []datatypes.Virtual_Guest_Block_Device_Template_Group, result2 error) {
	fake.listPublicImagesMutex.Lock()
	defer fake.listPublicImagesMutex.Unlock()
	fake.ListPublicImagesStub = nil
	if fake.listPublicImagesReturnsOnCall == nil {
		fake.listPublicImagesReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
			result2 error
		})
	}
	fake.listPublicImagesReturnsOnCall[i] = struct {
		result1 []datatypes.Virtual_Guest_Block_Device_Template_Group
		result2 error
	}{result1, result2}
}

func (fake *FakeImageManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteImageMutex.RLock()
	defer fake.deleteImageMutex.RUnlock()
	fake.editImageMutex.RLock()
	defer fake.editImageMutex.RUnlock()
	fake.exportImageMutex.RLock()
	defer fake.exportImageMutex.RUnlock()
	fake.getImageMutex.RLock()
	defer fake.getImageMutex.RUnlock()
	fake.importImageMutex.RLock()
	defer fake.importImageMutex.RUnlock()
	fake.listPrivateImagesMutex.RLock()
	defer fake.listPrivateImagesMutex.RUnlock()
	fake.listPublicImagesMutex.RLock()
	defer fake.listPublicImagesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeImageManager) recordInvocation(key string, args []interface{}) {
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

var _ managers.ImageManager = new(FakeImageManager)