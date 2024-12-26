package controllers

import (
	"Healfina_call/database"
	"Healfina_call/models"
	"encoding/json"
	"fmt"
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

type RecordController struct {
	beego.Controller
}

// AddRecordAfter
// @Summary Добавить запись пользователя
// @Description Добавляет новую запись для пользователя после завершения сессии.
// @Tags records
// @Accept json
// @Produce json
// @Param body body models.AddRecordRequest true "Данные для добавления записи"
// @Success 200 {object} models.GenericResponse "Успешное добавление записи"
// @Failure 400 {string} string "Некорректный формат JSON"
// @Failure 500 {string} string "Ошибка сервера"
// @router /records/add [put]
func (c *RecordController) AddRecordAfter() {
	user_id := AppContext.UserID

	var req models.AddRecordRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = models.ErrorResponse{Message: "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	err := database.AddRecordAfter(user_id, req.Record)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to add record: %v", err)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.GenericResponse{Status: "success"}
	c.ServeJSON()
}

// GetUserRecords
// @Summary Получить все записи пользователя
// @Description Получает список всех записей пользователя.
// @Tags records
// @Accept json
// @Produce json
// @Success 200 {array} models.RecordAfter "Список записей пользователя"
// @Failure 500 {string} string "Ошибка сервера"
// @router /records [get]
func (c *RecordController) GetUserRecords() {
	user_id := AppContext.UserID

	records, err := database.GetUserRecords(user_id)
	if err != nil {
		log.Printf("Error fetching user records: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to fetch records: %v", err)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = records
	c.ServeJSON()
}

// UpdateUserRecord
// @Summary Обновить запись пользователя
// @Description Обновляет существующую запись пользователя по ID.
// @Tags records
// @Accept json
// @Produce json
// @Param record_id path string true "ID записи"
// @Param body body models.UpdateRecordRequest true "Данные для обновления записи"
// @Success 200 {object} models.GenericResponse "Успешное обновление записи"
// @Failure 400 {string} string "Некорректный формат JSON или отсутствует ID записи"
// @Failure 500 {string} string "Ошибка сервера"
// router /records/update/:record_id [put]
func (c *RecordController) UpdateUserRecord() {
	user_id := AppContext.UserID
	recordID := c.Ctx.Input.Param(":record_id")
	if recordID == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = models.ErrorResponse{Message: "Record ID is required"}
		c.ServeJSON()
		return
	}

	var req models.UpdateRecordRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = models.ErrorResponse{Message: "Invalid JSON format"}
		c.ServeJSON()
		return
	}

	err := database.UpdateUserRecord(user_id, recordID, req.Record)
	if err != nil {
		log.Printf("Error updating record: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to update record: %v", err)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.GenericResponse{Status: "success"}
	c.ServeJSON()
}

// DeleteUserRecord
// @Summary Удалить запись пользователя
// @Description Удаляет запись пользователя по ID.
// @Tags records
// @Accept json
// @Produce json
// @Param record_id path string true "ID записи"
// @Success 200 {object} models.GenericResponse "Успешное удаление записи"
// @Failure 500 {string} string "Ошибка сервера"
// @router /records/delete/:record_id [delete]
func (c *RecordController) DeleteUserRecord() {
	user_id := AppContext.UserID
	recordID := c.Ctx.Input.Param(":record_id")

	err := database.DeleteUserRecord(user_id, recordID)
	if err != nil {
		log.Printf("Error deleting record: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to delete record: %v", err)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = models.GenericResponse{Status: "success"}
	c.ServeJSON()
}

// Функция для поиска записей по тексту
func (c *RecordController) SearchUserRecords() {
	user_id := AppContext.UserID
	searchQuery := c.GetString("searchQuery") // Получаем строку поиска из запроса

	if searchQuery == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = models.ErrorResponse{Message: "Search query is required"}
		c.ServeJSON()
		return
	}

	records, err := database.SearchUserRecords(user_id, searchQuery)
	if err != nil {
		log.Printf("Error fetching user records: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to fetch records: %v", err)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = records
	c.ServeJSON()
}
