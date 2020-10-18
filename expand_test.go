package strings

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

const (
	what = `
		admin.host: "${SESSIONNAME}${USERNAME}${SESSIONNAME}"
		admin.port: "${NotExistingOne}"
		admin.database: "$incomplete$$${abc${"
		admin.user: "$$$"
		admin.password: "&5m$$$$6WZyT>P$VMrY(N+-+?ZxXHrpy"
`

	prefix  = `${`
	postfix = `}`

	intact = `
		admin.host: "---session------username------session---"
		admin.port: "${NotExistingOne}"
		admin.database: "$incomplete$$${abc${"
		admin.user: "$$$"
		admin.password: "&5m$$$$6WZyT>P$VMrY(N+-+?ZxXHrpy"
`
	empty = `
		admin.host: "---session------username------session---"
		admin.port: ""
		admin.database: "$incomplete$$${abc${"
		admin.user: "$$$"
		admin.password: "&5m$$$$6WZyT>P$VMrY(N+-+?ZxXHrpy"
`
)

func TestOne(t *testing.T) {
	logger := log.WithField("area", "test")

	os.Setenv("SESSIONNAME", "---session---")
	os.Setenv("USERNAME", "---username---")

	back, err := Expand(what, prefix, postfix, EnvironmentResolverIntact, logger)
	if err != nil {
		t.Fatalf("error (%v)", err)
	}
	if back != intact {
		t.Fatalf("failed intact test")
	}

	back, err = Expand(what, prefix, postfix, EnvironmentResolver, logger)
	if err != nil {
		t.Fatalf("error (%v)", err)
	}
	if back != empty {
		t.Fatalf("failed empty test")
	}

	back, err = Expand(what, prefix, postfix, EnvironmentResolverFail, logger)
	if err == nil {
		t.Fatalf("it should've failed")
	}
}
