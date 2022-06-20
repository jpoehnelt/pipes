package pipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	StepsResponse       GLApiSuccessfulStepsResponse
	Steps               []Step
	Client              *http.Client
	DownloadedArtifacts []DownloadedArtifact
	JobNames            []string
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.Steps = []Step{}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {

			reqUrl := fmt.Sprintf(
				"%s/projects/%s/pipelines/%s/jobs/?scope=success",
				t.Pipe.Gitlab.ApiUrl,
				t.Pipe.Gitlab.ParentProjectId,
				t.Pipe.Gitlab.ParentPipelineId,
			)

			t.Log.Debugf("Pipeline steps API request URL: %s", reqUrl)

			req, err := http.NewRequest(
				http.MethodGet,
				reqUrl,
				nil,
			)

			if err != nil {
				return err
			}

			req.Header.Set("PRIVATE-TOKEN", t.Pipe.Gitlab.Token)

			t.Pipe.Ctx.Client = &http.Client{}

			res, err := t.Pipe.Ctx.Client.Do(req)

			if err != nil {
				return err
			}

			defer res.Body.Close()

			if err = ParseGLApiResponseCode(t, reqUrl, res.StatusCode); err != nil {
				return err
			}

			decoder := json.NewDecoder(res.Body)

			if err = decoder.Decode(&t.Pipe.Ctx.StepsResponse); err != nil {
				return err
			}

			return nil
		})
}

func DiscoverArtifacts(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("discover").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.JobNames = strings.Split(t.Pipe.Gitlab.DownloadArtifacts, "|")

			for _, step := range t.Pipe.Ctx.JobNames {
				func(step string) {
					t.CreateSubtask("").
						Set(func(t *Task[Pipe]) error {
							found := false

							for _, v := range t.Pipe.Ctx.StepsResponse {
								if v.Name == step {
									t.Log.Debugf(
										"Adding step artifacts: %s with id %d", step, v.ID,
									)

									t.Pipe.Ctx.Steps = append(t.Pipe.Ctx.Steps, Step{id: v.ID, name: step})

									found = true
								}
							}

							if !found {
								t.Log.Errorf(
									"Job with name is not found so artifacts are not downloaded: %s ",
									step,
								)

								var available []string = []string{}

								for _, v := range t.Pipe.Ctx.StepsResponse {
									available = append(available, v.Name)
								}

								return fmt.Errorf("Available steps are: %s", strings.Join(available, ", "))
							}

							return nil
						}).
						AddSelfToParent(func(pt, st *Task[Pipe]) {
							pt.ExtendSubtask(func(j Job) Job {
								return tl.JobParallel(j, st.Job())
							})
						})
				}(step)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func DownloadArtifacts(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("download").
		Set(func(t *Task[Pipe]) error {
			url := fmt.Sprintf(
				"%s/projects/%s/jobs/%s/artifacts/",
				t.Pipe.Gitlab.ApiUrl,
				t.Pipe.Gitlab.ParentProjectId,
				"%d",
			)

			for _, step := range t.Pipe.Ctx.Steps {
				func(step Step) {
					t.CreateSubtask(fmt.Sprintf("download:%s", step.name)).
						Set(func(t *Task[Pipe]) error {
							t.Log.Debugf(
								"Will download artifact with parent job: %s -> %d",
								step.name,
								step.id,
							)

							path, err := DownloadArtifact(t, fmt.Sprintf(url, step.id))

							if err != nil {
								return fmt.Errorf(
									"Can not download artifacts from stage: %s -> %d with error: %s",
									step.name,
									step.id,
									err,
								)
							}

							t.Pipe.Ctx.DownloadedArtifacts = append(
								t.Pipe.Ctx.DownloadedArtifacts,
								DownloadedArtifact{name: step.name, path: path},
							)

							return nil
						}).AddSelfToParent(func(pt, st *Task[Pipe]) {
						pt.ExtendSubtask(func(j Job) Job {
							return tl.JobParallel(j, st.Job())
						})
					})
				}(step)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func UnarchiveArtifacts(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("unarchive").
		Set(func(t *Task[Pipe]) error {
			for _, artifact := range t.Pipe.Ctx.DownloadedArtifacts {
				func(artifact DownloadedArtifact) {
					t.CreateSubtask(fmt.Sprintf("unarchive:%s", artifact.name)).
						Set(func(t *Task[Pipe]) error {
							t.Log.Debugf(
								"Decompressing artifact: %s -> %s",
								artifact.name,
								artifact.path,
							)

							t.CreateCommand(
								"unzip",
								"-o",
								artifact.path,
								"-d",
								"./",
							).AddSelfToTheTask()

							return nil
						}).ShouldRunAfter(func(t *Task[Pipe]) error {
						return t.RunCommandJobAsJobParallel()
					}).AddSelfToParent(func(pt, st *Task[Pipe]) {
						pt.ExtendSubtask(func(j Job) Job {
							return tl.JobParallel(j, st.Job())
						})
					})
				}(artifact)
			}

			return nil
		}).ShouldRunAfter(func(t *Task[Pipe]) error {
		return t.RunCommandJobAsJobParallel()
	})
}
