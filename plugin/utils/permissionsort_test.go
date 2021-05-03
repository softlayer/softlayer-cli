package utils_test

import (
	"sort"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	. "github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

var string1 = "a"
var string2 = "b"
var string3 = "c"

var permissionsBykeyName = []datatypes.User_Customer_CustomerPermission_Permission{
	datatypes.User_Customer_CustomerPermission_Permission{
		KeyName: &string2,
	},
	datatypes.User_Customer_CustomerPermission_Permission{
		KeyName: &string1,
	},
	datatypes.User_Customer_CustomerPermission_Permission{
		KeyName: &string3,
	},
}

var _ = Describe("Permissionsort", func() {
	It("", func() {
		Expect(strings.Contains(*permissionsBykeyName[1].KeyName, "a")).To(BeTrue())
		Expect(strings.Contains(*permissionsBykeyName[0].KeyName, "b")).To(BeTrue())
		Expect(strings.Contains(*permissionsBykeyName[2].KeyName, "c")).To(BeTrue())
		sort.Sort(PermissionsBykeyName(permissionsBykeyName))
		Expect(strings.Contains(*permissionsBykeyName[0].KeyName, "a")).To(BeTrue())
		Expect(strings.Contains(*permissionsBykeyName[1].KeyName, "b")).To(BeTrue())
		Expect(strings.Contains(*permissionsBykeyName[2].KeyName, "c")).To(BeTrue())
	})
})
