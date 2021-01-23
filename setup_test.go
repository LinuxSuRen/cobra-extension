package pkg

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"testing"
)

func TestAll(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	reporter := reporters.NewJUnitReporter("test.xml")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "", []ginkgo.Reporter{reporter})
}
