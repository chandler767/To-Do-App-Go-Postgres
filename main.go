package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	tpl = template.Must(template.ParseFiles("task-entry.html"))
	db  *gorm.DB
)

type newTask struct {
	gorm.Model
	Task string
}

func newEntry(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		task := req.FormValue("task")
		if task != "" {
			db.Create(&newTask{Task: task}) // Add to database
			err := tpl.ExecuteTemplate(w, "task-entry.html", "test")
			if err != nil {
				http.Error(w, err.Error(), 500)
				log.Fatalln(err)
			}
			return
		}
	}
	err := tpl.ExecuteTemplate(w, "task-entry.html", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func main() {
	var err error
	db, err = gorm.Open("postgres", "host="+os.Getenv("PGHOST")+" port="+os.Getenv("PGPORT")+" user="+os.Getenv("PGUSER")+" dbname="+os.Getenv("PGDATABASE")+" sslmode=disable password="+os.Getenv("PGPASSWORD"))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&newTask{}) // Creates our task table
	http.HandleFunc("/", newEntry)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil))
}
