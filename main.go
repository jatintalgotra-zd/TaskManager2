package main

import (
	"fmt"
	"net/http"
	"time"

	"TaskManager2/datasource/mysql"
	taskHandler "TaskManager2/handler/task"
	userHandler "TaskManager2/handler/user"
	taskService "TaskManager2/service/task"
	userService "TaskManager2/service/user"
	taskStore "TaskManager2/store/task"
	userStore "TaskManager2/store/user"
)

func main() {
	db, err := mysql.New("root", "root123", "test_db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	taskStr := taskStore.New(db)
	userStr := userStore.New(db)

	userSvc := userService.New(userStr)
	taskSvc := taskService.New(taskStr, userSvc)

	taskHndlr := taskHandler.New(taskSvc)
	userHndlr := userHandler.New(userSvc)

	http.HandleFunc("POST /task", taskHndlr.Post)
	http.HandleFunc("GET /task", taskHndlr.Get)
	http.HandleFunc("GET /task/{id}", taskHndlr.GetByID)
	http.HandleFunc("PUT /task", taskHndlr.Put)
	http.HandleFunc("DELETE /task/{id}", taskHndlr.DeleteByID)

	http.HandleFunc("POST /user", userHndlr.Post)
	http.HandleFunc("GET /user/{id}", userHndlr.GetByID)

	server := http.Server{
		Addr:        ":8080",
		Handler:     http.DefaultServeMux,
		ReadTimeout: 5 * time.Second,
	}

	err2 := server.ListenAndServe()
	if err2 != nil {
		fmt.Println(err)
	}
}
