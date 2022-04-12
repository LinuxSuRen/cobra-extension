package testing

// FlagValidation contains some fields of a flag to valid
type FlagValidation struct {
	Name         string
	Shorthand    string
	UsageIsEmpty bool
}

// FlagsValidation is an array of the FlagValidations
type FlagsValidation []FlagValidation
