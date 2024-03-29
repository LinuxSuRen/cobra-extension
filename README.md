[![](https://goreportcard.com/badge/linuxsuren/cobra-extension)](https://goreportcard.com/report/linuxsuren/cobra-extension)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/linuxsuren/cobra-extension)
[![Contributors](https://img.shields.io/github/contributors/linuxsuren/cobra-extension.svg)](https://github.com/linuxsuren/cobra-extension/graphs/contributors)
[![GitHub release](https://img.shields.io/github/release/linuxsuren/cobra-extension.svg?label=release)](https://github.com/linuxsuren/cobra-extension/releases/latest)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/linuxsuren/cobra-extension)
[![HitCount](http://hits.dwyl.com/linuxsuren/cobra-extension.svg)](http://hits.dwyl.com/linuxsuren/cobra-extension)

This project aims to provide an easy way to let you writing a plugin for your CLI project. And it based on [cobra](https://github.com/spf13/cobra).

## Get started

`go get github.com/linuxsuren/cobra-extension`

## Friendly to the flags test

You can add some tests for the flags quickly, for instance:

```go
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
```
