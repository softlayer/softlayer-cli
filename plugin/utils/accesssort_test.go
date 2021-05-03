package utils

import (
	"sort"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestAccessByID(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].ID, "j")
	sort.Sort(AccessByID(accessList))
	Equal(t, accessList[0].ID, "a")
}

func TestAccessByName(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].Name, "k")
	sort.Sort(AccessByName(accessList))
	Equal(t, accessList[0].Name, "b")
}

func TestAccessByType(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].Type, "l")
	sort.Sort(AccessByType(accessList))
	Equal(t, accessList[0].Type, "c")
}

func TestAccessByIP(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].PrivateIPAddress, "m")
	sort.Sort(AccessByPrivateIPAddress(accessList))
	Equal(t, accessList[0].PrivateIPAddress, "d")
}

func TestAccessBySourceSubnet(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].SourceSubnet, "n")
	sort.Sort(AccessBySourceSubnet(accessList))
	Equal(t, accessList[0].SourceSubnet, "e")
}

func TestAccessByHostIQN(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].HostIQN, "o")
	sort.Sort(AccessByHostIQN(accessList))
	Equal(t, accessList[0].HostIQN, "f")
}

func TestAccessByUserName(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].UserName, "p")
	sort.Sort(AccessByUserName(accessList))
	Equal(t, accessList[0].UserName, "g")
}

func TestAccessByPassword(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].Password, "q")
	sort.Sort(AccessByPassword(accessList))
	Equal(t, accessList[0].Password, "h")
}

func TestAccessByAllowedHostID(t *testing.T) {
	accessList := []Access{
		Access{
			ID:               "j",
			Name:             "k",
			Type:             "l",
			PrivateIPAddress: "m",
			SourceSubnet:     "n",
			HostIQN:          "o",
			UserName:         "p",
			Password:         "q",
			AllowedHostID:    "r",
		},
		Access{
			ID:               "a",
			Name:             "b",
			Type:             "c",
			PrivateIPAddress: "d",
			SourceSubnet:     "e",
			HostIQN:          "f",
			UserName:         "g",
			Password:         "h",
			AllowedHostID:    "i",
		},
	}
	Equal(t, accessList[0].AllowedHostID, "r")
	sort.Sort(AccessByAllowedHostID(accessList))
	Equal(t, accessList[0].AllowedHostID, "i")
}
