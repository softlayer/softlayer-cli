package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
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
	Describe("normalizeQuietFlag Tests", func() {
		It("Test quite => quiet", func() {
			flagSet := pflag.NewFlagSet("testSet", 0)
			var fakeBoolVar bool
			flagSet.BoolVarP(&fakeBoolVar, "quiet", "q", false, "test")
			result := utils.NormalizeQuietFlag(flagSet, "quite")
			Expect(string(result)).To(Equal("quiet"))
		})
		It("Test nothing else is changed", func() {
			flagSet := pflag.NewFlagSet("testSet", 0)
			var fakeBoolVar bool
			flagSet.BoolVarP(&fakeBoolVar, "quiet", "q", false, "test")
			result := utils.NormalizeQuietFlag(flagSet, "asd")
			Expect(string(result)).To(Equal("asd"))
		})
	})
})
