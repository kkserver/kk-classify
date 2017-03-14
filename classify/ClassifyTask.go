package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifyTaskResult struct {
	app.Result
	Classify *Classify `json:"classify,omitempty"`
}

type ClassifyTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result ClassifyTaskResult
}

func (task *ClassifyTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifyTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifyTask) GetClientName() string {
	return "Classify.Get"
}
