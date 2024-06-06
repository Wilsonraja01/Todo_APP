package mysql

// Importing the required packages
import (
	"database/sql"
	"fmt"
	"todo/pkg/models"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

// This will insert a new Task into the database.
func (m *TodoModel) Insert(task, details, expires string) (int, error) {
	// Given the SQL query for insert into the database.
	stmt := `INSERT INTO Todo (task, details, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Execute the statement and checking if any error is returned
	result, err := m.DB.Exec(stmt, task, details, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	//  Return the Id of the task
	return int(id), nil

}

// This will Fetches the task details from the database
func (m *TodoModel) GetFromSql() ([]models.Task, error) {
	// Given the SQL query for get from the database
	stmt := `SELECT * FROM Todo`
	// Execute the statement and checking if any error is returned
	rows, err := m.DB.Query(stmt)
	if err != nil {
		fmt.Printf("Error at stmt")
		return nil, err
	}
	// Created a Slice of struct Task
	var TaskListInstance []models.Task
	// Iterating over each rows
	for rows.Next() {
		// Created a Task instance for each row
		s := &models.TodoDB{}
		// Scanning the rows and checking if any error is returned
		err = rows.Scan(&s.ID, &s.Task, &s.Details, &s.Created, &s.Expires)
		if err != nil {
			fmt.Printf("Error at inner s")
			return nil, err
		}
		// Created a Task instance for each row
		oldTask := models.Task{
			ID:      s.ID,
			Name:    s.Task,
			Details: s.Details,
		}
		// Appending the oldTask to the TaskListInstance slice
		TaskListInstance = append(TaskListInstance, oldTask)
	}
	// Check if any errors happened
	if err = rows.Err(); err != nil {
		fmt.Printf("Error at final")
		return nil, err
	}
	// If everything went OK then return the TaskListInstance slice.
	return TaskListInstance, nil
}

// This will delete a specific Task based on its id.
func (m *TodoModel) DeleteList(id int) error {
	// Given the SQL query to delete the specific id and its row from the database
	stmt := `DELETE FROM Todo WHERE id =?`
	// Execute the statement and checking if any error is returned
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// This will Update a specific Task based on its id.
func (m *TodoModel) UpdateList(id int, task string, details string) error {
	// Given the SQL query to Update the specific id and its row to the database
	stmt := `UPDATE Todo SET task=?,details=? WHERE id=?`
	// Execute the statement and checking if any error is returned
	_, err := m.DB.Exec(stmt, task, details, id)
	if err != nil {
		return err
	}
	return nil
}
