package smartquery

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/navidrome/navidrome/log"
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestCriteriagw(t *testing.T) {
	log.SetLevel(log.LevelFatal)
	gomega.RegisterFailHandler(Fail)
	RunSpecs(t, "Smart Criteria Suite")
}
