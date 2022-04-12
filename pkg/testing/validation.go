package testing

import (
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Valid do the validation for the flags
func (f FlagsValidation) Valid(t *testing.T, flagSet *flag.FlagSet) {
	for i := range f {
		tt := f[i]
		t.Run(tt.Name, func(t *testing.T) {
			targetFlag := flagSet.Lookup(tt.Name)

			assert.NotNil(t, targetFlag)
			assert.Equal(t, tt.Shorthand, targetFlag.Shorthand)

			if tt.UsageIsEmpty {
				assert.Empty(t, targetFlag.Usage)
			} else {
				assert.NotEmpty(t, targetFlag.Usage)
			}
		})
	}
}
