package run

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	NodeCommand struct {
		Script string
		Cwd    string `validate:"dir"`
	}

	Pipe struct {
		Ctx

		NodeCommand
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				RunNodeScript(tl).Job(),
			)
		})
}
