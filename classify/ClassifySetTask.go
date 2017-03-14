package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifySetTaskResult struct {
	app.Result
	Classify *Classify `json:"classify,omitempty"`
}

type ClassifySetTask struct {
	app.Task
	Id     int64       `json:"id"`
	Alias  string      `json:"alias"`
	Pid    interface{} `json:"pid"`
	Name   interface{} `json:"name"` //名称
	Tags   interface{} `json:"tags"` //搜索标签
	Result ClassifySetTaskResult
}

func (task *ClassifySetTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifySetTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifySetTask) GetClientName() string {
	return "Classify.Set"
}
