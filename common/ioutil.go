package common

// Printer for print the info
type Printer interface {
	PrintErr(i ...interface{})
	Println(i ...interface{})
	Printf(format string, i ...interface{})
}
