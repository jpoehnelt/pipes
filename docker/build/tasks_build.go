package build

import (
	"time"

	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerBuildParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build", "parent").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return setup.TL.Pipe.Docker.UseBuildx
		}).
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				DockerBuild(tl).Job(),
				DockerPush(tl).Job(),
			)
		})
}

func DockerBuild(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				t.Pipe.DockerFile.Name,
				t.Pipe.DockerFile.Context,
			)

			// build image
			t.CreateCommand(
				setup.DOCKER_EXE,
				"build",
			).
				Set(func(c *Command[Pipe]) error {
					for _, v := range t.Pipe.DockerImage.BuildArgs {
						c.AppendArgs("--build-arg", v)
					}

					if t.Pipe.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					for _, tag := range t.Pipe.Ctx.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						t.Pipe.DockerFile.Name,
						".",
					)

					c.SetDir(t.Pipe.DockerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					if setup.TL.Pipe.Docker.UseBuildKit {
						t.Log.Debugf("Using Docker BuildKit for the build operation.")

						c.AppendEnvironment(
							map[string]string{
								"DOCKER_BUILDKIT": "1",
							},
						)
					}

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerPush(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("push").
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						setup.DOCKER_EXE,
						"push",
						tag,
					).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Pushing Docker image: %s",
								tag,
							)

							return nil
						}).
						SetRetries(3, false, time.Second*10).
						AddSelfToTheTask()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}