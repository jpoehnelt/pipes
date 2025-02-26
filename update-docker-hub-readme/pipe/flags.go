package pipe

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER_HUB = "DockerHub"
	CATEGORY_README     = "Readme"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_HUB,
		Name:        "docker_hub.username",
		Usage:       "DockerHub username for updating the readme.",
		EnvVars:     []string{"DOCKER_USERNAME"},
		Required:    true,
		Destination: &TL.Pipe.DockerHub.Username,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_HUB,
		Name:        "docker_hub.password",
		Usage:       "DockerHub password for updating the readme.",
		EnvVars:     []string{"DOCKER_PASSWORD"},
		Required:    true,
		Destination: &TL.Pipe.DockerHub.Password,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_HUB,
		Name:        "docker_hub.address",
		Usage:       "HTTP address for the DockerHub compatible service.",
		EnvVars:     []string{"DOCKER_HUB_ADDRESS"},
		Value:       "https://hub.docker.com/v2/repositories",
		Destination: &TL.Pipe.DockerHub.Address,
	},

	&cli.StringFlag{
		Category:    CATEGORY_README,
		Name:        "readme.repository",
		Usage:       "Repository for applying the readme on.",
		EnvVars:     []string{"DOCKER_IMAGE_NAME", "README_REPOSITORY"},
		Required:    false,
		Value:       "",
		Destination: &TL.Pipe.Readme.Repository,
	},

	&cli.StringFlag{
		Category:    CATEGORY_README,
		Name:        "readme.file",
		Usage:       "Readme file for the given repository.",
		EnvVars:     []string{"README_FILE"},
		Value:       "README.md",
		Destination: &TL.Pipe.Readme.File,
		Required:    false,
	},

	&cli.StringFlag{
		Category:    CATEGORY_README,
		Name:        "readme.short_description",
		Usage:       "Short description to display on DockerHub.",
		EnvVars:     []string{"README_DESCRIPTION"},
		Destination: &TL.Pipe.Readme.Description,
		Required:    false,
	},

	&cli.StringFlag{
		Category: CATEGORY_README,
		Name:     "readme.matrix",
		Usage:    "Matrix of multiple README files to update. json([]struct{ repository: string, file: string, description?: string })",
		EnvVars:  []string{"README_MATRIX"},
		Required: false,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	if len(tl.Pipe.Readme.Description) > 100 {
		return fmt.Errorf(
			"Readme short description can only be 100 characters long while you have: %d",
			len(tl.Pipe.Readme.Description),
		)
	}

	if v := tl.CliContext.String("readme.matrix"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.Readme.Matrix); err != nil {
			return fmt.Errorf("Can not unmarshal Readme matrix: %w", err)
		}
	}

	if tl.Pipe.Readme.Repository == "" && len(tl.Pipe.Readme.Matrix) == 0 {
		return fmt.Errorf("You have to either provide a target via Repository or multiple targets through the Matrix.")
	}

	tl.Pipe.Ctx.Readme = make(map[string]ParsedReadme)

	return nil
}
