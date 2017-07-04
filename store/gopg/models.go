package gopg

import (
	"errors"
	"time"

	"github.com/joelclouddistrict/gochecklist/services"
)

type ModelTodo struct {
	TableName struct{} `sql:"todos,alias:todos"`
	Id        int32
	Task      string
	Done      bool `sql:",notnull"`
	Created   time.Time
	Updated   time.Time
}

// Marshal will convert the Todo model data into a corresponding gRPC message
// pointed by v.
// The gRPC message should be the required one. It is not infered.
func (m *ModelTodo) Marshal(v interface{}) error {
	todo, ok := v.(*services.TodoMessage)
	if !ok {
		return errors.New("Value is not of type services.TodoMessage")
	}

	todo.Id = m.Id
	todo.Task = m.Task
	todo.Done = m.Done

	return nil
}

// Unmarshal will convert the data of the gRPC message pointed in v
// and put it in the Todo model data
// The gRPC message should be the required one. It is not infered.
func (m *ModelTodo) Unmarshal(v interface{}) error {
	todo, ok := v.(*services.TodoMessage)
	if !ok {
		return errors.New("Value is not of type services.TodoMessage")
	}

	m.Id = todo.Id
	m.Task = todo.Task
	m.Done = todo.Done

	return nil
}
