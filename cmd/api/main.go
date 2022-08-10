package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"
	"course/internal/user/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	db := database.NewConnDatabase()
	exerciseService := exercise.NewExerciseUsecase(db)
	// repo := repository.NewDatabaseRepo(db)
	mcsrv := repository.NewMcsrvRepo()
	userUsecase := user.NewUserUsecase(mcsrv)

	r.GET("/exercises/:id", middleware.WithTimeout(), middleware.WithJWT(userUsecase), exerciseService.GetExerciseByID)
	r.GET("/exercises/:id/score", middleware.WithJWT(userUsecase), exerciseService.CalculateUserScore)
	r.POST("/exercises", middleware.WithJWT(userUsecase), exerciseService.CreateExercise)
	r.POST("/exercises/:id/questions", middleware.WithJWT(userUsecase), exerciseService.CreateQuesetion)
	r.POST("/exercises/:id/questions/:questionId/answers", middleware.WithJWT(userUsecase), exerciseService.CreateAnswer)

	r.Run(":1234")
}
