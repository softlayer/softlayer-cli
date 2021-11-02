package tags_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"testing"
)

var FakeTags = []datatypes.Tag{
	datatypes.Tag{
		Id:             sl.Int(1234),
		Name:           sl.String("TEST TAG"),
		ReferenceCount: sl.Uint(1),
	},
}

var FakeTagReference = []datatypes.Tag_Reference{
	datatypes.Tag_Reference{
		Id:              sl.Int(1111),
		ResourceTableId: sl.Int(22222),
		TagType: &datatypes.Tag_Type{
			Description: sl.String("Test Tag"),
			KeyName:     sl.String("HARDWARE"),
		},
	},
}

func TestManagers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tags Suite")
}
