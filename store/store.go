package store

import (
	"github.com/joelclouddistrict/gochecklist/services"
)

// Storer es la interfaz que deben implementar los diferentes backends que realizen
// la persistencia del microservicio
type Storer interface {
	// Conecta con el Store
	Dial(options Options) error
	// Status devuelve el estado de la conexión
	Status() (int, string)
	// Close cierra la conexión con el Store
	Close() error
	// TodoManager exposes todo management methods
	TodoManager
}

type TodoManager interface {
	// CreateTodo adds a record in database
	CreateTodo(todo *services.TodoMessage) (*services.TodoMessage, error)
	// GetTodo retrieves a record in database by ID
	GetTodo(id *services.IdMessage) (*services.TodoMessage, error)
	// ListTodos retrieves a list of records applying filters
	ListTodos(filter *services.TodoFilter) (*services.TodoArray, error)
	// DeleteTodo removes a record in database by ID
	DeleteTodo(id *services.IdMessage) error
	// PatchTodo updates a record in database by ID
	PatchTodo(todo *services.TodoMessage) (*services.TodoMessage, error)
}

const (
	PAGINATION_LIMIT = 50
	// This limit is used on endpoints that do not have a pagination limit
	UPPER_LIMIT = 10000
	// DISCONNECTED indicates that there is no connection with the Storer
	DISCONNECTED = iota
	// CONNECTED indicate that the connection with the Storer is up and running
	CONNECTED
)

var (
	// StatusStr is a string representation of the status of the connections with the Storer
	StatusStr = []string{"Disconnected", "Connected"}
)

// Options is a map to hold the database connection options
type Options map[string]interface{}
