package gopg

import (
	"errors"
	"fmt"
	"time"

	"github.com/joelclouddistrict/gochecklist/log"
	"github.com/joelclouddistrict/gochecklist/services"
	"github.com/joelclouddistrict/gochecklist/store"
	"gopkg.in/pg.v5/orm"
)

// CreateTodo adds a record in database
func (p *GopgStore) CreateTodo(todo *services.TodoMessage) (*services.TodoMessage, error) {
	mTodo := &ModelTodo{}

	mTodo.Unmarshal(todo)
	mTodo.Created = time.Now()
	mTodo.Updated = time.Now()

	_, err := p.Pool.Model(mTodo).
		Returning("*").
		Insert()
	if err != nil {
		log.Err(fmt.Sprintf("Error creating todo (%v): %s", todo, err))
		return nil, err
	}

	err = mTodo.Marshal(todo)

	if err != nil {
		log.Err(fmt.Sprintf("Error creating todo (%v): %s", todo, err))
		return nil, err
	}

	return todo, nil
}

// PatchTodo updates a record in database by ID
func (p *GopgStore) PatchTodo(todo *services.TodoMessage) (*services.TodoMessage, error) {
	mTodo := &ModelTodo{}

	mTodo.Unmarshal(todo)
	mTodo.Updated = time.Now()

	q := p.Pool.Model(mTodo)

	getTodoUpdatedColumns(mTodo, q)

	res, err := q.Returning("*").Update()
	if err != nil {
		log.Err(fmt.Sprintf("Error updating todo (%v): %s", todo, err))
		return nil, err
	}

	if res.RowsAffected() != 1 {
		err = errors.New(fmt.Sprintf("Todo %d not found", todo.Id))
		log.Err(fmt.Sprintf("Error updating todo (%v): %s", todo, err))
		return nil, err
	}

	err = mTodo.Marshal(todo)
	if err != nil {
		log.Err(fmt.Sprintf("Error updating todo (%v): %s", todo, err))
		return nil, err
	}

	return todo, nil
}

// GetTodo retrieves a record in database by ID
func (p *GopgStore) GetTodo(id *services.IdMessage) (*services.TodoMessage, error) {
	mTodo := &ModelTodo{}

	q := p.Pool.Model(mTodo).
		Where("id = ?", id.Id)

	err := q.Select()
	if err != nil {
		log.Err(fmt.Sprintf("Error obtaining todo %d: %s", id.Id, err))
		return nil, err
	}

	todo := &services.TodoMessage{}
	mTodo.Marshal(todo)

	return todo, nil
}

// ListTodos retrieves a list of records applying filters
func (p *GopgStore) ListTodos(filters *services.TodoFilter) (*services.TodoArray, error) {
	var mTodos []ModelTodo

	q := p.Pool.Model(&mTodos).
		Limit(store.PAGINATION_LIMIT)

	applyTodoFilters(filters, q)

	count, err := q.SelectAndCount()
	if err != nil {
		log.Err(fmt.Sprintf("Error obtaining todos: %s %v", err, filters))
		return nil, err
	}

	arr := &services.TodoArray{}

	for _, mTodo := range mTodos {
		todo := &services.TodoMessage{}
		err = mTodo.Marshal(todo)
		if err != nil {
			log.Err(fmt.Sprintf("Error obtaining todos: %s %v", err, filters))
			continue
		}
		arr.Todos = append(arr.Todos, todo)
	}

	arr.Total = int32(count)

	return arr, nil
}

// DeleteTodo removes a record in database by ID
func (p *GopgStore) DeleteTodo(id *services.IdMessage) error {
	mTodo := &ModelTodo{
		Id: id.Id,
	}

	err := p.Pool.Delete(mTodo)
	if err != nil {
		log.Err(fmt.Sprintf("Error deleting todo %d: %s", id.Id, err))
		return err
	}

	return nil
}

// getTodoUpdatedColumns adds to an update query the columns that should be updated based on which
// attributes of the models are set
func getTodoUpdatedColumns(mTodo *ModelTodo, q *orm.Query) {
	if len(mTodo.Task) > 0 {
		q.Column("task")
	}

	q.Column("done")

	if !mTodo.Updated.IsZero() {
		q.Column("updated")
	}
}

func applyTodoFilters(filters *services.TodoFilter, q *orm.Query) {
	if filters == nil {
		return
	}

	if len(filters.Terms) > 0 {
		q.Where("todos.task ILIKE", fmt.Sprintf("%%%s%%", filters.Terms))
	}

	if filters.Done == 1 {
		q.Where("todos.done = ?", true)
	}

	if filters.Done == 2 {
		q.Where("todos.done = ?", false)
	}

	q.Offset(int(filters.Offset))

	if filters.Limit > 0 && filters.Limit <= store.PAGINATION_LIMIT {
		q.Limit(int(filters.Limit))
	}

}
