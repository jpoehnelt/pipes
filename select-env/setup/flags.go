package setup

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_ENVIRONMENT = "Environment"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(flags.GitFlagsDestination{
	GitBranch: &TL.Pipe.Git.Branch,
	GitTag:    &TL.Pipe.Git.Tag,
}), []cli.Flag{

	// CATEGORY_ENVIRONMENT

	&cli.StringFlag{
		Category: CATEGORY_ENVIRONMENT,
		Name:     "environment.conditions",
		Usage:    `Regex pattern to select an environment. Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. json([]{ condition: RegExp, environment: string })`,
		Required: true,
		EnvVars:  []string{"ENVIRONMENT_CONDITIONS"},
		Value:    `[ { "condition": "^tags/v?\\d.\\d.\\d$", "environment": "production" }, { "condition": "^tags/v?\\d.\\d.\\d-.*\\.\\d$", "environment": "stage" }, { "condition" :"^heads/main$", "environment": "develop" }, { "condition": "^heads/master$", "environment": "develop" } ]`,
		Action: func(ctx *cli.Context, s string) error {
			// setup selection of environment conditions
			if err := json.Unmarshal([]byte(s), &TL.Pipe.Conditions); err != nil {
				return fmt.Errorf("Can not unmarshal environment conditions: %w", err)
			}

			return nil
		},
	},

	&cli.BoolFlag{
		Category:    CATEGORY_ENVIRONMENT,
		Name:        "environment.fail-on-no-reference",
		Usage:       "Whether to fail on missing environment references.",
		Required:    false,
		EnvVars:     []string{"ENVIRONMENT_FAIL_ON_NO_REFERENCE"},
		Value:       true,
		Destination: &TL.Pipe.Environment.FailOnNoReference,
	},
})
