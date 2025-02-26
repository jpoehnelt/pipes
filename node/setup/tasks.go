package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func SetupPackageManager(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.PackageManager = PackageManager{
				Exe:      t.Pipe.Node.PackageManager,
				Commands: PackageManagers[t.Pipe.Node.PackageManager],
			}

			t.Log.Infof("Using package manager: %s", t.Pipe.Node.PackageManager)

			return nil
		})
}

func NodeVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"node",
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command[Pipe]) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch node.js version.")

						return nil
					}

					t.Log.Infof("node.js version: %s", stream[0])

					return nil
				}).
				AddSelfToTheTask()

			t.CreateCommand(
				t.Pipe.Ctx.PackageManager.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs(t.Pipe.Ctx.PackageManager.Commands.Version...)

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command[Pipe]) error {
					stream := c.GetCombinedStream()

					if len(stream) == 0 {
						t.Log.Debugln("Can not fetch package manager version.")

						return nil
					}

					t.Log.Infof("%s version: v%s", t.Pipe.Ctx.PackageManager.Exe, stream[0])

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
