package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

var _ = Describe("Utility Tests", func() {
	DescribeTable("FormatBoolPointerToYN Tests",
		func(input bool, expected string) {
			var result string
			// Special case for passing in a nil reference
			// go otherwise treans nil as false if it is assigned
			if expected == "-" {
				result = utils.FormatBoolPointerToYN(nil)
			} else {
				result = utils.FormatBoolPointerToYN(&input)
			}

			Expect(result).To(Equal(expected))
		},
		Entry("Nil Input", nil, "No"),
		Entry("True Input", true, "Yes"),
		Entry("False Input", false, "No"),
	)
	DescribeTable("ResolveVirtualGuestId Tests",
		func(input string, expected int, expectedErr string) {
			result, err := utils.ResolveVirtualGuestId(input)
			if expectedErr != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(expectedErr))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
			Expect(result).To(Equal(expected))
		},
		Entry("Number", "99", 99, nil),
		Entry("Large Number", "30000000", 30000000, nil),
		Entry("String", "NinteyNine", 0, "strconv.Atoi: parsing \"NinteyNine\": invalid syntax"),
		Entry("Nil", nil, 0, "strconv.Atoi: parsing \"\": invalid syntax"),
	)
})

/*
func TestSliceInSlice(t *testing.T) {
	source := []string{"id", "hostnamdde"}
	defaultColumns := []string{"id", "hostname", "domain", "cpu", "memory", "primary_ip", "backend_ip", "datacenter", "action"}
	optionalColumns := []string{"guid", "power_state", "created_by", "tags"}
	target := append(defaultColumns, optionalColumns...)

	exist, idx := SliceInSlice(source, target)
	Equal(t, false, exist)
	Equal(t, 1, idx)
}

func TestStringSliceToString(t *testing.T) {
	sclice := []string{"aaa", "bbb"}
	result := StringSliceToString(sclice)
	Equal(t, "aaa,bbb", result)
}

func TestTagRefsToString(t *testing.T) {
	tags := []datatypes.Tag_Reference{
		datatypes.Tag_Reference{
			Tag: &datatypes.Tag{
				Name: sl.String("aaa"),
			},
		},
		datatypes.Tag_Reference{
			Tag: &datatypes.Tag{
				Name: sl.String("bbb"),
			},
		},
	}
	result := TagRefsToString(tags)
	Equal(t, "aaa,bbb", result)
}

func TestBool2Int(t *testing.T) {
	value1 := true
	result1 := Bool2Int(value1)
	Equal(t, result1, 1)
	value2 := false
	result2 := Bool2Int(value2)
	Equal(t, result2, 0)
}

func TestStructToMap(t *testing.T) {
	access := Access{
		ID:               "1",
		Name:             "2",
		Type:             "3",
		PrivateIPAddress: "4",
		SourceSubnet:     "5",
		HostIQN:          "6",
		UserName:         "7",
		Password:         "8",
		AllowedHostID:    "9",
	}
	accessMap, err := StructToMap(access)
	Equal(t, err, nil)
	Equal(t, accessMap, map[string]string{"id": "1", "name": "2", "type": "3", "private_ip_address": "4", "source_subnet": "5", "host_iqn": "6", "username": "7", "password": "8", "allowed_host_id": "9"})
}

func TestStructToMap1(t *testing.T) {
	access := Access{
		ID:               "a",
		Name:             "b",
		Type:             "c",
		PrivateIPAddress: "d",
		SourceSubnet:     "e",
		HostIQN:          "f",
		UserName:         "g",
		Password:         "h",
		AllowedHostID:    "i",
	}
	accessMap, err := StructToMap(access)
	Equal(t, err, nil)
	Equal(t, accessMap, map[string]string{"id": "a", "name": "b", "type": "c", "private_ip_address": "d", "source_subnet": "e", "host_iqn": "f", "username": "g", "password": "h", "allowed_host_id": "i"})
}

*/
