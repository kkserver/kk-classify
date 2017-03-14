package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifyCreateTaskResult struct {
	app.Result
	Classify *Classify `json:"classify,omitempty"`
}

type ClassifyCreateTask struct {
	app.Task
	Pid    int64  `json:"pid"`
	Alias  string `json:"alias"` //别名
	Name   string `json:"name"`  //名称
	Tags   string `json:"tags"`  //搜索标签
	Result ClassifyCreateTaskResult
}

func (task *ClassifyCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifyCreateTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifyCreateTask) GetClientName() string {
	return "Classify.Create"
}
