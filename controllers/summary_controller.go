package controllers

import (
	"Healfina_call/database"
	"Healfina_call/models"
	"fmt"
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

type SummaryController struct {
	beego.Controller
}

// GetOverallSummary
// @Summary Получить общий обзор
// @Description Получает общий обзор данных пользователя на основе его ID.
// @Tags summary
// @Accept json
// @Produce json
// @Success 200 {object} models.GetOverallSummaryResponse "Успешный ответ с общим обзором"
// @Failure 500 {object} models.ErrorResponse "Ошибка при получении общего обзора"
// @router /summary [get]
func (c *SummaryController) GetOverallSummary() {
	user_id := AppContext.UserID

	overallSummary, err := database.GetOverallSummary(user_id)
	if err != nil {
		log.Printf("Error getting overall summary: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to get overall summary: %v", err)}
		c.ServeJSON()
		return
	}

	response := models.GetOverallSummaryResponse{OverallSummary: overallSummary}
	c.Data["json"] = response
	c.ServeJSON()
}
