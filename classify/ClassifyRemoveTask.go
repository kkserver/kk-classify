package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifyRemoveTaskResult struct {
	app.Result
}

type ClassifyRemoveTask struct {
	app.Task
	Id     int64       `json:"id"`
	Pid    interface{} `json:"pid"`
	Result ClassifyRemoveTaskResult
}

func (task *ClassifyRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifyRemoveTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifyRemoveTask) GetClientName() string {
	return "Classify.Remove"
}
