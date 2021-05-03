package ipsec_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("IPSec update", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeIPSecManager *testhelpers.FakeIPSECManager
		cmd              *ipsec.UpdateCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeIPSecManager = new(testhelpers.FakeIPSECManager)
		cmd = ipsec.NewUpdateCommand(fakeUI, fakeIPSecManager)
		cliCommand = cli.Command{
			Name:        metadata.IpsecUpdateMetaData().Name,
			Description: metadata.IpsecUpdateMetaData().Description,
			Usage:       metadata.IpsecUpdateMetaData().Usage,
			Flags:       metadata.IpsecUpdateMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Context("update without contextID", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
		})
	})
	Context("update with wrong context id", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Invalid input for 'Context ID'. It must be a positive integer.")).To(BeTrue())
		})
	})
	Context("update with wrong phase1-auth", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-a", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -a|--phase1-auth must be either MD5 or SHA1 or SHA256.")).To(BeTrue())
		})
	})
	Context("update with wrong phase1-crypto", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-c", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "-c|--phase1-crypto must be either DES or 3DES or AES128 or AES192 or AES256.")).To(BeTrue())
		})
	})
	Context("update with wrong phase1-dh", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-d", "10")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -d|--phase1-dh must be either 0 or 1 or 2 or 5.")).To(BeTrue())
		})
	})
	Context("update with wrong phase1-key-ttl", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-t", "100")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -t|--phase1-key-ttl must be in range 120-172800.")).To(BeTrue())
		})
	})
	Context("update with wrong phase2-auth", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-u", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -u|--phase2-auth must be either MD5 or SHA1 or SHA256.")).To(BeTrue())
		})
	})
	Context("update with wrong phase2-crypto", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-y", "abc")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -y|--phase2-crypto must be either DES or 3DES or AES128 or AES192 or AES256.")).To(BeTrue())
		})
	})
	Context("update with wrong phase2-dh", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-e", "3")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -e|--phase2-dh must be either 0 or 1 or 2 or 5.")).To(BeTrue())
		})
	})
	Context("update with wrong phase2-forward-secrecy", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-f", "2")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -f|--phase2-forward-secrecy must be either 0 or 1.")).To(BeTrue())
		})
	})
	Context("update with wrong phase2-key-ttl", func() {
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-l", "100")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Incorrect Usage: -l|--phase2-key-ttl must be in range 120-172800.")).To(BeTrue())
		})
	})
	Context("update with server fail", func() {
		BeforeEach(func() {
			fakeIPSecManager.UpdateTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{}, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-l", "150")
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to update IPSec 1234.")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})
	Context("update", func() {
		BeforeEach(func() {
			fakeIPSecManager.UpdateTunnelContextReturns(datatypes.Network_Tunnel_Module_Context{Id: sl.Int(1234)}, nil)
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-a", "SHA256")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-c", "AES256")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-d", "1")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-t", "1000")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-u", "SHA256")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-y", "AES256")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-e", "2")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-f", "1")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
		It("succeed", func() {
			err := testhelpers.RunCommand(cliCommand, "1234", "-l", "150")
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated IPSec 1234."}))
		})
	})
})
