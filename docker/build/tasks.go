package build

import (
	"gitlab.kilic.dev/devops/pipes/common/parser"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				job,
				ParseReferences(tl).Job(),
			)
		})
}

func ParseReferences(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init", "references").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.References = parser.ParseGitReferences(t.Pipe.Git.Tag, t.Pipe.Git.Branch)

			t.Log.Debugf("References for environment selection: %v", t.Pipe.Ctx.References)

			return nil
		})
}

func DockerInspect(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("inspect").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.DockerImage.Inspect
		}).
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						setup.DOCKER_EXE,
						"manifest",
						"inspect",
						tag,
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Inspecting Docker image: %s",
								tag,
							)

							return nil
						}).
						AddSelfToTheTask()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
