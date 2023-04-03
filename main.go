package main

import (
	"net/http"
	"strconv"
	kawethradb "github.com/Hasan-Kilici/kawethradb"
	"github.com/gin-gonic/gin"
)

type Task struct {
	ID         int
	UserID     int
	Tasks      string
	Taskstatus string
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("src/*.tmpl")
	r.Static("/static", "./static/")

	r.GET("/", func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("ID")
		if err != nil {
			id := kawethradb.Count("./Tasks.csv")
			ctx.SetCookie("ID", strconv.Itoa(id), 36000, "/", "", false, true)
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		userID, _ := strconv.Atoi(cookie)
		results, _ := kawethradb.FindAll("./Tasks.csv", "UserID", userID)
		var tasks []Task
		for _, result := range results {
			taskid, _ := strconv.Atoi(result["ID"])
			task := Task{
				ID:         taskid,
				UserID:     userID,
				Tasks:      result["Tasks"],
				Taskstatus: result["Taskstatus"],
			}
			tasks = append(tasks, task)
		}

		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Tasks":  tasks,
			"UserID": userID,
		})
	})

	r.POST("/addTask/:userID", func(ctx *gin.Context) {
		userID, _ := strconv.Atoi(ctx.Param("userID"))
		taskInput := ctx.PostForm("task")
		task := Task{
			ID:         kawethradb.Count("./Tasks.csv") + 1,
			UserID:     userID,
			Tasks:      taskInput,
			Taskstatus: "Not Finished",
		}
		kawethradb.Insert("./Tasks.csv", task)
		ctx.Redirect(http.StatusFound, "/")
	})

	r.POST("/deleteTask/:TaskID", func(ctx *gin.Context) {
		taskID, _ := strconv.Atoi(ctx.Param("TaskID"))
		kawethradb.DeleteByID("./Tasks.csv", taskID)

		ctx.Redirect(http.StatusFound, "/")
	})

	r.POST("/finishTask/:TaskID", func(ctx *gin.Context) {
		cookie, _ := ctx.Cookie("id")
		taskID := ctx.Param("TaskID")
		userID, _ := strconv.Atoi(cookie)
		newUserID := strconv.Itoa(userID)
		formTask, _ := kawethradb.FindByID("./Tasks.csv", taskID)
		task := formTask["Tasks"]
		updatedTask := []string{taskID, newUserID, task, "Finished"}
		kawethradb.Update("./Tasks.csv", "ID", taskID, updatedTask)
		ctx.Redirect(http.StatusFound, "/")
	})

	r.POST("/unfinishTask/:TaskID", func(ctx *gin.Context) {
		cookie, _ := ctx.Cookie("id")
		taskID := ctx.Param("TaskID")
		userID, _ := strconv.Atoi(cookie)
		newUserID := strconv.Itoa(userID)
		formTask, _ := kawethradb.FindByID("./Tasks.csv", taskID)
		task := formTask["Tasks"]
		updatedTask := []string{taskID, newUserID, task, "Not Finished"}
		kawethradb.Update("./Tasks.csv", "ID", taskID, updatedTask)
		ctx.Redirect(http.StatusFound, "/")
	})

	r.Run()

}

