package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifyBatchCreateTaskResult struct {
	app.Result
	Classifys []Classify `json:"classifys,omitempty"`
}

type ClassifyBatchCreateTask struct {
	app.Task
	Pid    int64  `json:"id"`
	Alias  string `json:"alias"`
	Names  string `json:"names"` //名称,名称,名称,名称
	Result ClassifyBatchCreateTaskResult
}

func (task *ClassifyBatchCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifyBatchCreateTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifyBatchCreateTask) GetClientName() string {
	return "Classify.BatchCreate"
}
