package exec

import (
	"time"

	"github.com/soderasen-au/go-common/util"
)

type (
	// App can do a Task for many times, each time of the task-execution is a Job
	Task struct {
		Name        string     `json:"name,omitempty" yaml:"name,omitempty" bson:"name,omitempty"`
		LastJob     *Job       `json:"last_job,omitempty" yaml:"last_job,omitempty" bson:"last_job,omitempty" gorm:"serializer:json"`
		LastSucc    *Job       `json:"last_succ,omitempty" yaml:"last_succ,omitempty" bson:"last_succ,omitempty" gorm:"serializer:json"`
		LastFail    *Job       `json:"last_fail,omitempty" yaml:"last_fail,omitempty" bson:"last_fail,omitempty" gorm:"serializer:json"`
		History     []*Job     `json:"history,omitempty" yaml:"history,omitempty" bson:"history,omitempty" gorm:"serializer:json"`
		NextRun     *time.Time `json:"next_run,omitempty" yaml:"next_run,omitempty" bson:"next_run,omitempty"`
		MaxJobCount int        `json:"max_job_count,omitempty" yaml:"max_job_count,omitempty" bson:"max_job_count,omitempty"`
	}

	// A Job finishes a Task using a Request
	Job struct {
		TaskName   string         `json:"task_name,omitempty" yaml:"task_name,omitempty" bson:"task_name,omitempty" gorm:"index"`
		ReqID      string         `json:"req_id,omitempty" yaml:"req_id,omitempty" bson:"req_id,omitempty" gorm:"index"`
		Status     Status         `json:"status,omitempty" yaml:"status,omitempty" bson:"status,omitempty"`
		StartedBy  string         `json:"started_by,omitempty" yaml:"started_by,omitempty" bson:"started_by,omitempty"`
		StartedAt  time.Time      `json:"started_at,omitempty" yaml:"started_at,omitempty" bson:"started_at,omitempty"`
		FinishedAt *time.Time     `json:"finished_at,omitempty" yaml:"finished_at,omitempty" bson:"finished_at,omitempty"`
		LogFile    string         `json:"log_file,omitempty" yaml:"log_file,omitempty" bson:"log_file,omitempty"`
		Errors     []*util.Result `json:"errors,omitempty" yaml:"errors,omitempty" bson:"errors,omitempty" gorm:"serializer:json"`
	}
)
