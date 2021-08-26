package main

import (
	"database/sql"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/julez-dev/go2todo/repo/lists"
	"github.com/julez-dev/go2todo/repo/tasks"
	"github.com/julez-dev/go2todo/service"
	"github.com/julez-dev/go2todo/ui"
	_ "modernc.org/sqlite"
)

func chooseStore() *service.Storage {
	storageType := os.Getenv("GO2TODO_STORAGETYPE")

	if storageType == "sql" {
		sql, err := sql.Open("sqlite", os.Getenv("GO2TODO_SQLPATH"))

		if err != nil {
			log.Fatalln("could not open sql file: %w", err)
		}

		// TODO: Close db

		listsDB, err := lists.NewInSQL(sql)

		if err != nil {
			log.Fatalln(err)
		}

		tasksDB, err := tasks.NewInSQL(sql)

		if err != nil {
			log.Fatal(err)
		}

		return service.NewStorage(tasksDB, listsDB)
	}

	taskPath := os.Getenv("GO2TODO_TASKPATH")

	if taskPath == "" {
		taskPath = "tasks.json"
	}

	listPath := os.Getenv("GO2TODO_LISTPATH")

	if listPath == "" {
		listPath = "lists.json"
	}

	taskDB, err := tasks.NewInFile(taskPath)

	if err != nil {
		log.Fatalln(err)
	}

	listDB, err := lists.NewInFile(listPath)

	if err != nil {
		log.Fatalln(err)
	}

	return service.NewStorage(taskDB, listDB)
}

func main() {
	store := chooseStore()
	tea.NewProgram(ui.New(store)).Start()
}
