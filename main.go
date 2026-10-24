package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sabudaye/openproject-cli/cmd"
	"github.com/sabudaye/openproject-cli/components/configuration"
	openerrors "github.com/sabudaye/openproject-cli/components/errors"
)

var (
	version = "current"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cliVersion := configuration.BuildCliVersion(version, commit, date)

	if err := cmd.Execute(cliVersion); err != nil {
		if !errors.Is(err, openerrors.ErrHandled) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
