package main

// Importing required packages
import (
	// "fmt"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"todo/pkg/models"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
)

// Id is initialized to make a unique ID for each Task
var Id int

var homeErrors = make(map[string]string)

// var ErrorStruct templateData
// Created Home Function which is a structure which includes the information of Info Log and Error log
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// models.TaskList=append(TaskList, models.OldTaskList)
	// Checking the path is in home or not
	// panic("oops! something went wrong")

	// Geting the task list from the database
	// Flash := app.session.PopString(r ,"flash")
	// fmt.Println(Flash)

	models.TaskList, err = app.Todos.GetFromSql()
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Using String Slice to add files which will give the template.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl"}
	ts, err := template.ParseFiles(files...)
	// Parsing files and checking for any error
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Adding the Information in Info.log file
	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf("At Home")
	if ferr != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// Executing the template with slice of Task and checking for any error
	err = ts.Execute(w, struct {
		Tasks      []models.Task
		Flash      string
		FormData   url.Values
		FormErrors map[string]string
	}{
		Tasks:      models.TaskList,
		Flash:      app.session.PopString(r, "flash"),
		FormErrors: homeErrors,
		FormData:   r.PostForm,
	})

	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

// Created AddList Function which is a structure which includes the information of Info Log and Error log
func (app *application) AddList(w http.ResponseWriter, r *http.Request) {
	// Adding the Id initially as it need to be unique
	Id++
	var FormData url.Values
	var FormErrors map[string]string
	var errors = make(map[string]string)
	// Adding new task by giving the Name and Details in the Forms
	getTask := r.FormValue("TaskName")
	getDetails := r.FormValue("Details")

	if strings.TrimSpace(getTask) == "" && strings.TrimSpace(getDetails) == "" {
		errors["TaskName"] = "This field cannot be blank"
		errors["Details"] = "This field cannot be blank"
		app.session.Put(r, "flash", "Please Give Task and Details!")
	} else {
		if strings.TrimSpace(getTask) == "" {
			errors["TaskName"] = "This field cannot be blank"
			app.session.Put(r, "flash", "Please Give Task!")
		} else if utf8.RuneCountInString(getTask) > 100 {
			errors["TaskName"] = "This field is too long (maximum is 100 characters)"
			app.session.Put(r, "flash", "Too many characters!")
		}

		if strings.TrimSpace(getDetails) == "" {
			errors["Details"] = "This field cannot be blank"
			app.session.Put(r, "flash", "Please Give Details!")
		} else if utf8.RuneCountInString(getDetails) > 10000 {
			errors["Details"] = "This field is too long (maximum is 10000 characters)"
			app.session.Put(r, "flash", "Too many characters!")
		}
	}
	if len(errors) > 0 {
		FormErrors = errors
		FormData = r.PostForm
		homeErrors = errors
		fmt.Println(FormErrors, FormData)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	newTask := models.Task{
		ID:      Id,
		Name:    getTask,
		Details: getDetails,
	}

	// Appending the Task into the String slice
	models.TaskList = append(models.TaskList, newTask)
	// Adding the Information in Info.log file
	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf("Added a Task ID:%d, TaskName: %s, Details: %s", newTask.ID, newTask.Name, newTask.Details)
	if ferr != nil {
		log.Fatal(ferr)
	}
	// Created a new Task with the given Task name and Task details
	task := newTask.Name
	details := newTask.Details
	// Predefined Task Expires in 7 Days
	expires := "7"
	// Insert a new Task into the database and check if any errors have happened
	id, err := app.Todos.Insert(task, details, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Task successfully created!")
	// Making the Global Id is equal to the Database Id
	Id = id
	// Redirecting to the home page by using http.Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Created DeleteList Function which is a structure which includes the information of Info Log and Error log
func (app *application) DeleteList(w http.ResponseWriter, r *http.Request) {
	// Get the value "id" from r.FormValue
	value, _ := strconv.Atoi((r.URL.Query().Get("id")))
	// Checking the id and deleting the task by appending
	for i, val := range models.TaskList {
		if val.ID == value {
			app.Todos.DeleteList(value)
			models.TaskList = append(models.TaskList[:i], models.TaskList[i+1:]...)
			break
		}
	}
	// Adding the Information in Info.log file
	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf("Deleted Task of ID: %d\n", value)
	if ferr != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	app.session.Put(r, "flash", "Task successfully Deleted!")
	// Redirecting to the home page by using http.Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Created updateList Function which is a structure which includes the information of Info Log and Error log
func (app *application) updateList(w http.ResponseWriter, r *http.Request) {
	// Get the value "id" from r.FormValue
	value, _ := strconv.Atoi(r.FormValue("id"))
	// Created the Task and details variables as String
	var task string
	var details string
	// Getting the Task and details from the form
	task = r.FormValue("TaskName")
	details = r.FormValue("Details")
	// Check if the Task and details are Present in the form
	if len(task) != 0 && len(details) != 0 {
		// Updating the Task and details in the Database and checking if any errors there
		err := app.Todos.UpdateList(value, task, details)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// Checking if No value are found in the form
	} else if len(task) == 0 && len(details) == 0 {
		// Redirect to the home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		// Checking if no value are found in the Details form
	} else if len(task) != 0 && len(details) == 0 {
		// Geting the details from the Database for that specific Id
		for _, val := range models.TaskList {
			if val.ID == value {
				details = val.Details
				break
			}
		}
		// Update the Task and details in the database and checking if any errors there
		err := app.Todos.UpdateList(value, task, details)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		// Checking if no value are found in the Name form
	} else if len(task) == 0 && len(details) != 0 {
		// Geting the Task from the Database for that specific Id
		for _, val := range models.TaskList {
			if val.ID == value {
				task = val.Name
				break
			}
		}
		// Update the Task and details in the database and checking if any errors there
		err := app.Todos.UpdateList(value, task, details)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}
	app.session.Put(r, "flash", "Task successfully Updated!")
	// Redirecting to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Created a Function to Open the Database
func openDB(dsn string) (*sql.DB, error) {
	// Opening the Database and checking if any errors there
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user signup form...")
}
func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
