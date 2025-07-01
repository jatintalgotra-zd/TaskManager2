package main

import (
	"fmt"

	"gofr.dev/pkg/gofr"

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

	app := gofr.New()

	app.GET("/task", taskHndlr.GetAllHandler)
	app.GET("/task/{id}", taskHndlr.GetByIDHandler)
	app.POST("/task", taskHndlr.PostHandler)
	app.PUT("/task/{id}", taskHndlr.PutHandler)
	app.DELETE("/task/{id}", taskHndlr.DeleteHandler)

	app.GET("/user/{id}", userHndlr.GetByIDHandler)
	app.POST("/user", userHndlr.PostHandler)

	app.Run()
}
