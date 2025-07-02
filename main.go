package main

import (
	"gofr.dev/pkg/gofr"

	taskHandler "TaskManager2/handler/task"
	userHandler "TaskManager2/handler/user"
	"TaskManager2/migrations"
	taskService "TaskManager2/service/task"
	userService "TaskManager2/service/user"
	taskStore "TaskManager2/store/task"
	userStore "TaskManager2/store/user"
)

func main() {
	taskStr := taskStore.New()
	userStr := userStore.New()

	userSvc := userService.New(userStr)
	taskSvc := taskService.New(taskStr, userSvc)

	taskHndlr := taskHandler.New(taskSvc)
	userHndlr := userHandler.New(userSvc)

	app := gofr.New()

	app.Migrate(migrations.All())

	app.GET("/task", taskHndlr.GetAllHandler)
	app.GET("/task/{id}", taskHndlr.GetByIDHandler)
	app.POST("/task", taskHndlr.PostHandler)
	app.PUT("/task/{id}", taskHndlr.PutHandler)
	app.DELETE("/task/{id}", taskHndlr.DeleteHandler)

	app.GET("/user/{id}", userHndlr.GetByIDHandler)
	app.POST("/user", userHndlr.PostHandler)

	app.Run()
}
