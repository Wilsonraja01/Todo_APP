package main

// Importing the required packages
import (
	"bytes"
	"flag"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions" 
	"log"
	"net/http"
	"os"
	"todo/pkg/models/mysql"
)

// Created a application which is a structure
// which includes the information of Info Log and Error log
// and contains information about the database
type application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	session *sessions.Session
	Todos    *mysql.TodoModel
}

type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}


	
// Created infoLog which will update the information
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Here f will open the info.log file and it will write the infoLog details ad ferr has the error
var f, ferr = os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE, 0666)

// Initalied the errorLog which has the information of the log of date,time and error details
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// Initialised er which is used to update the error in the error.log file
var er, err = os.OpenFile("./error.log", os.O_RDWR|os.O_CREATE, 0666)

// Main function
func main() {
	// Address of the HTTP server in Default is http://localhost:4000
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Initialised the MySQL server
	dsn := flag.String("dsn", "root:root@/todoApp?parseTime=true", "MySQL data source name")

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret")
	// Parsing Both Address and Database
	flag.Parse()
	// Check if MySQL server is running
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	// Creating a instance of application
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	app := &application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		session: session,
		Todos:    &mysql.TodoModel{DB: db},
	}

	// If any error happen in f it will give an Error
	if ferr != nil {
		log.Fatal(ferr)
	}
	// Log the information to info.log file with Date and time
	infoLog = log.New(f, "INFO\t", log.Ldate|log.Ltime)
	infoLog.Printf("starting server on %s", *addr)
	defer f.Close()
	// If any error happen in er it will give an Error
	if err != nil {
		log.Fatal(err)
	}
	// Log the error to error.log file with Date and time
	errorLog = log.New(er, "ERROR\t", log.Ldate|log.Ltime)

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Call the ListenAndServe() method in the struct.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
