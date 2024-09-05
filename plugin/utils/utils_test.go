package utils_test

import (
	"sort"
	// "fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/pflag"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers")

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
Describe("Block Sorting Utility Tests", func() {
		var blockVolumes []datatypes.Network_Storage
		BeforeEach(func() {
			fakeSession := testhelpers.NewFakeSoftlayerSession([]string{})
			err := fakeSession.TransportHandler.DoRequest(fakeSession, "SoftLayer_Account", "getNasNetworkStorage", nil, nil, &blockVolumes)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Test utils.VolumeById", func() {
			sort.Sort(utils.VolumeById(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(123))
			Expect(*blockVolumes[1].Id).To(Equal(154123))
			Expect(*blockVolumes[2].Id).To(Equal(4917309))
			Expect(*blockVolumes[3].Id).To(Equal(21021427))
		})
		It("Test utils.VolumeByUsername", func() {
			sort.Sort(utils.VolumeByUsername(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(21021427))
			Expect(*blockVolumes[1].Id).To(Equal(4917309))
			Expect(*blockVolumes[2].Id).To(Equal(154123))
			Expect(*blockVolumes[3].Id).To(Equal(123))
		})
		It("Test utils.VolumeByDatacenter", func() {
			sort.Sort(utils.VolumeByDatacenter(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(4917309))
			Expect(*blockVolumes[1].Id).To(Equal(21021427))
			Expect(*blockVolumes[2].Id).To(Equal(123))
			Expect(*blockVolumes[3].Id).To(Equal(154123))
		})
		It("Test utils.VolumeByStorageType", func() {
			sort.Sort(utils.VolumeByStorageType(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(123))
			Expect(*blockVolumes[1].Id).To(Equal(21021427))
			Expect(*blockVolumes[2].Id).To(Equal(4917309))
			Expect(*blockVolumes[3].Id).To(Equal(154123))
		})
		It("Test utils.VolumeByCapacity", func() {
			sort.Sort(utils.VolumeByCapacity(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(4917309))
			Expect(*blockVolumes[1].Id).To(Equal(21021427))
			Expect(*blockVolumes[2].Id).To(Equal(123))
			Expect(*blockVolumes[3].Id).To(Equal(154123))
		})
		It("Test utils.VolumeByBytesUsed", func() {
			sort.Sort(utils.VolumeByBytesUsed(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(21021427))
			Expect(*blockVolumes[1].Id).To(Equal(123))
			Expect(*blockVolumes[2].Id).To(Equal(154123))
			Expect(*blockVolumes[3].Id).To(Equal(4917309))
		})
		It("Test utils.VolumeByIPAddress", func() {
			sort.Sort(utils.VolumeByIPAddress(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(154123))
			Expect(*blockVolumes[1].Id).To(Equal(123))
			Expect(*blockVolumes[2].Id).To(Equal(21021427))
			Expect(*blockVolumes[3].Id).To(Equal(4917309))
		})
		It("Test utils.VolumeByTxnCount", func() {
			sort.Sort(utils.VolumeByTxnCount(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(4917309))
			Expect(*blockVolumes[1].Id).To(Equal(21021427))
			Expect(*blockVolumes[2].Id).To(Equal(123))
			Expect(*blockVolumes[3].Id).To(Equal(154123))
		})
		It("Test utils.VolumeByCreatedBy", func() {
			sort.Sort(utils.VolumeByCreatedBy(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(21021427))
			Expect(*blockVolumes[1].Id).To(Equal(123))
			Expect(*blockVolumes[2].Id).To(Equal(4917309))
			Expect(*blockVolumes[3].Id).To(Equal(154123))
		})
		It("Test utils.VolumeByMountAddr", func() {
			sort.Sort(utils.VolumeByMountAddr(blockVolumes))
			Expect(len(blockVolumes)).To(Equal(4))
			Expect(*blockVolumes[0].Id).To(Equal(4917309))
			Expect(*blockVolumes[1].Id).To(Equal(21021427))
			Expect(*blockVolumes[2].Id).To(Equal(154123))
			Expect(*blockVolumes[3].Id).To(Equal(123))
		})
	})
})
