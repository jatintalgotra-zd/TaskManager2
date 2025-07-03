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

	app.GET("/task", taskHndlr.GetAll)
	app.GET("/task/{id}", taskHndlr.GetByID)
	app.POST("/task", taskHndlr.Post)
	app.PUT("/task/{id}", taskHndlr.Put)
	app.DELETE("/task/{id}", taskHndlr.Delete)

	app.GET("/user/{id}", userHndlr.GetByID)
	app.POST("/user", userHndlr.Post)

	app.Run()
}
