package testing

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestFlagsValidation_Valid(t *testing.T) {
	boolFlag := true
	emptyFlag := true
	cmd := cobra.Command{}
	cmd.Flags().BoolVarP(&boolFlag, "test", "t", false, "usage test")
	cmd.Flags().BoolVarP(&emptyFlag, "empty", "", false, "")

	flags := FlagsValidation{{
		Name:      "test",
		Shorthand: "t",
	}, {
		Name:         "empty",
		UsageIsEmpty: true,
	}}
	flags.Valid(t, cmd.Flags())
}
