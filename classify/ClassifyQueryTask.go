package classify

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ClassifyQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type ClassifyQueryTaskResult struct {
	app.Result
	Counter   *ClassifyQueryCounter `json:"counter,omitempty"`
	Classifys []Classify            `json:"classifys,omitempty"`
}

type ClassifyQueryTask struct {
	app.Task
	Id        int64       `json:"id"`
	Pid       interface{} `json:"pid"`
	Alias     string      `json:"alias"`
	Keyword   string      `json:"q"`
	PageIndex int         `json:"p"`
	PageSize  int         `json:"size"`
	Counter   bool        `json:"counter"`
	Result    ClassifyQueryTaskResult
}

func (task *ClassifyQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *ClassifyQueryTask) GetInhertType() string {
	return "classify"
}

func (task *ClassifyQueryTask) GetClientName() string {
	return "Classify.Query"
}
