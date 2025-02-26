package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DefaultTask(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("default").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"echo",
			).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs("hello")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
