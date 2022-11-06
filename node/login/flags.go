package login

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_LOGIN = "Login"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_LOGIN

	&cli.StringFlag{
		Category: CATEGORY_NODE_LOGIN,
		Name:     "npm.login",
		Usage:    "npm registries to login to. json(slice({ username: string, password: string, registry?: string, useHttps?: boolean }))",
		Required: false,
		EnvVars:  []string{"NPM_LOGIN"},
		Value:    "",
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_NODE_LOGIN,
		Name:        "npm.npmrc_file",
		Usage:       ".npmrc file to use.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC_FILE"},
		Value:       cli.NewStringSlice(".npmrc"),
		Destination: &TL.Pipe.Npm.NpmRcFile,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_LOGIN,
		Name:        "npm.npmrc",
		Usage:       "Pass direct contents of the NPMRC file.",
		Required:    false,
		EnvVars:     []string{"NPM_NPMRC"},
		Value:       "",
		Destination: &TL.Pipe.Npm.NpmRc,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) Job {
	return tl.CreateBasicJob(func() error {
		if err := json.Unmarshal([]byte(tl.CliContext.String("npm.login")), &tl.Pipe.Npm.Login); err != nil {
			return fmt.Errorf("Can not unmarshal Npm registry login credentials: %w", err)
		}

		return nil
	})
}
