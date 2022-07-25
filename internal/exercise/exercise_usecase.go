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
