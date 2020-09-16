package validate

import (
	"github.com/urfave/cli"
)

func ValidateServer(c *cli.Context) error {
	if !c.IsSet("folder") {
		return cli.NewExitError("folder argument is missing", 2)
	}
	return nil
}
