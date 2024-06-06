package models

// Importing the required packages
import (
	"errors"
	"time"
)

// Check if any errors are encountered
var ErrNoRecord = errors.New("models: no matching record found")

// Create a new structure to hold the record data of the Table in the database
type TodoDB struct {
	ID      int
	Task    string
	Details string
	Created time.Time
	Expires time.Time
}

// Creating a Task structure Which is used to store Unique id, Name of the task and Details of the task
type Task struct {
	ID      int
	Name    string
	Details string
}

// Creating a Slice which stores values of Task structure
var TaskList []Task
