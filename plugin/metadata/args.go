package metadata

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

// Copied https://github.com/spf13/cobra/blob/a0aadc68eb88af3acad9be0ec18c8b91438fd984/args.go so we can Translate
// these error messages.

// NoArgs returns an error if any args are included.
func NoArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		subs := map[string]interface{}{
			"Arg":     args[0],
			"Command": cmd.CommandPath(),
		}
		return fmt.Errorf(T("unknown command {{.Arg}} for {{.Command}}", subs))
	}
	return nil
}

// OnlyValidArgs returns an error if any args are not in the list of ValidArgs.
func OnlyValidArgs(cmd *cobra.Command, args []string) error {
	if len(cmd.ValidArgs) > 0 {
		// Remove any description that may be included in ValidArgs.
		// A description is following a tab character.
		var validArgs []string
		for _, v := range cmd.ValidArgs {
			validArgs = append(validArgs, strings.Split(v, "\t")[0])
		}

		for _, v := range args {
			if !stringInSlice(v, validArgs) {
				subs := map[string]interface{}{
					"Arg":  v,
					"Path": cmd.CommandPath(),
				}
				return fmt.Errorf(T("invalid argument {{.Arg}} for {{.Path}}", subs))
			}
		}
	}
	return nil
}

// ArbitraryArgs never returns an error.
func ArbitraryArgs(cmd *cobra.Command, args []string) error {
	return nil
}

// MinimumNArgs returns an error if there is not at least N args.
func MinimumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			subs := map[string]interface{}{
				"Limit": n,
				"Args":  len(args),
			}
			return fmt.Errorf(T("requires at least {{.Limit}} arg(s), only received {{.Args}}", subs))
		}
		return nil
	}
}

// MaximumNArgs returns an error if there are more than N args.
func MaximumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			subs := map[string]interface{}{
				"Limit": n,
				"Args":  len(args),
			}
			return fmt.Errorf(T("accepts at most {{.Limit}} arg(s), received {{.Args}}", subs))
		}
		return nil
	}
}

// ExactArgs returns an error if there are not exactly n args.
func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			subs := map[string]interface{}{"Limit": n, "Args": len(args)}
			return fmt.Errorf(T("accepts {{.Limit}} arg(s), received {{.Args}}", subs))
		}
		return nil
	}
}

// These ARGS exist so I can specify the proper translated error message. The default from spf13 doesn't allow for custom
// error messages

// Just one arg
func OneArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(T("Incorrect Usage: This command requires one argument."))
	}
	return nil
}

// Just two arg
func TwoArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(T("Incorrect Usage: This command requires two arguments."))
	}
	return nil
}

// Just three arg
func ThreeArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf(T("Incorrect Usage: This command requires three arguments."))
	}
	return nil
}

// ExactValidArgs returns an error if
// there are not exactly N positional args OR
// there are any positional args that are not in the `ValidArgs` field of `Command`
func ExactValidArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := ExactArgs(n)(cmd, args); err != nil {
			return err
		}
		return OnlyValidArgs(cmd, args)
	}
}

// RangeArgs returns an error if the number of args is not within the expected range.
func RangeArgs(min int, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min || len(args) > max {
			subs := map[string]interface{}{
				"Min":  min,
				"Max":  max,
				"Args": len(args),
			}
			return fmt.Errorf(T("accepts between {{.Min}} and {{.Max}} arg(s), received {{.Args}}", subs))
		}
		return nil
	}
}

// MatchAll allows combining several PositionalArgs to work in concert.
func MatchAll(pargs ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, parg := range pargs {
			if err := parg(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
