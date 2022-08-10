package exercise

import (
	"course/internal/domain"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseUsecase struct {
	db *gorm.DB
}

func NewExerciseUsecase(db *gorm.DB) *ExerciseUsecase {
	return &ExerciseUsecase{db}
}

func (exUsecase ExerciseUsecase) GetExerciseByID(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = exUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "exercise not found",
		})
		return
	}
	c.JSON(200, exercise)
}

func (exUsecase ExerciseUsecase) CalculateUserScore(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = exUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": "exercise not found",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = exUsecase.db.Where("user_id = ?", userID).Find(&answers).Error
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"message": "error when find answers",
		})
		return
	}
	if len(answers) == 0 {
		c.JSON(200, map[string]interface{}{
			"score": 0,
		})
		return
	}

	mapQuestion := make(map[int]domain.Question)
	for _, question := range exercise.Questions {
		mapQuestion[question.ID] = question
	}

	var score int

	for _, answer := range answers {
		if strings.EqualFold(answer.Answer, mapQuestion[answer.QuestionID].CorrectAnswer) {
			score += mapQuestion[answer.QuestionID].Score
		}
	}
	c.JSON(200, map[string]interface{}{
		"score": score,
	})
}

func (exUsecase ExerciseUsecase) CreateExercise(c *gin.Context) {
	var exercise domain.Exercise
	err := c.ShouldBind(&exercise)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	err = exUsecase.db.Create(&exercise).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create exercise",
		})
		return
	}

	c.JSON(201, gin.H{
		"id":          exercise.ID,
		"title":       exercise.Title,
		"description": exercise.Description,
	})
}

func (exUsecase ExerciseUsecase) CreateQuesetion(c *gin.Context) {
	paramID := c.Param("id")
	exercise_id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"message": "invalid exercise id",
		})
		return
	}

	var question domain.Question
	err = c.ShouldBind(&question)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	var exercise domain.Exercise
	err = exUsecase.db.Where("id = ?", exercise_id).Take(&exercise).Error
	if err != nil {
		c.JSON(404, gin.H{
			"message": "exercise not found",
		})
		return
	}

	question.ExerciseID = exercise_id
	userID := int(c.Request.Context().Value("user_id").(float64))
	question.CreatorID = userID
	err = exUsecase.db.Create(&question).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when add question",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "question has been saved",
	})
}

func (exUsecase ExerciseUsecase) CreateAnswer(c *gin.Context) {
	paramID := c.Param("id")
	exercise_id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	paramQuestionID := c.Param("questionId")
	question_id, err := strconv.Atoi(paramQuestionID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid question id",
		})
		return
	}

	var exercise domain.Exercise
	err = exUsecase.db.Where("id = ?", exercise_id).Take(&exercise).Error
	if err != nil {
		c.JSON(404, gin.H{
			"message": "exercise not found",
		})
		return
	}

	var question domain.Question
	err = exUsecase.db.Where("id = ?", question_id).Take(&question).Error
	if err != nil {
		c.JSON(404, gin.H{
			"message": "question not found",
		})
		return
	}

	var answer domain.Answer
	err = c.ShouldBind(&answer)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid answer input",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	answer.ExerciseID = exercise_id
	answer.QuestionID = question_id
	answer.UserID = userID

	err = exUsecase.db.Create(&answer).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create answer",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "answer has been saved",
	})
}
