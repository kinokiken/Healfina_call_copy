package controllers

import (
	"Healfina_call/database"
	"Healfina_call/models"
	"encoding/json"
	"fmt"
	"log"

	beego "github.com/beego/beego/v2/server/web"
)

type DialogController struct {
	beego.Controller
}

// AddDialog
// @Summary Добавить диалог для пользователя
// @Description Добавляет новый диалог с заранее подготовленным планом общения или с последним сохранённым планом.
// @Tags dialog
// @Accept json
// @Produce json
// @Success 200 {object} models.AddDialogResponse "Успешно добавлен новый диалог"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @router /dialog [post]
func (c *DialogController) AddDialog() {
	user_id := AppContext.UserID

	var plan string

	latestSummary, err := database.GetLatestDialogSummary(user_id)
	if err != nil {
		log.Printf("Error fetching latest summary: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to get latest summary: %v", err)}
		c.ServeJSON()
		return
	}

	if latestSummary == "" {
		plan = `План общения:

		1. Знакомство: Первое сообщение: "Рада видеть тебя на нашей первой консультации 😊 Как тебя зовут?". В следующих сообщениях узнай базовую информацию о клиенте, такую как возраст, пол, профессия, его увлечения. Затем, узнайте о проблемах, с которыми столкнулся клиент.
		2. Объяснение КПТ и схематерапии: После того, как вы узнали клиента, объясните, что такое КПТ и схематерапия, как это работает и чего клиент может ожидать от терапии. Кратко объясните, что цель терапии — помочь клиентам понять и изменить долгосрочные эмоциональные проблемы. Также обсудите важность домашних заданий между сессиями терапии. 
		3. Оценка проблемы клиента: Используйте открытые вопросы, чтобы побудить клиента говорить о своих проблемах, симптомах и целях для терапии. Это поможет вам понять модели мышления, поведения и эмоций клиента. Помните о том, что нужно дать клиенту свободно говорить и активно слушать.
		4. Установка целей: Как только у вас будет лучше понимание проблем клиента, работайте вместе, чтобы установить цели лечения. Эти цели должны быть конкретными, измеримыми, достижимыми, актуальными и ограниченными по времени (SMART). Также обсудите, как будет измеряться прогресс.
		5. Первичная когнитивно-поведенческая концептуализация: Начните разработку когнитивно-поведенческой концептуализации проблемы клиента. Это включает в себя определение связей между мыслями, чувствами и поведением клиента. Возможно, вы не завершите это на первой сессии, но вы можете начать строить предварительное понимание, которое включает ключевые черты характера клиента, его механизмы справления со стрессом и преобладающие когнитивные шаблоны. Также обсудите, как прошлые опыты и семейная история могут влиять на текущее состояние.
		6. Введение в навыки КПТ: В зависимости от времени и готовности клиента, вы можете представить основной навык КПТ, такой как когнитивная перестройка или медитация. Дайте краткое объяснение и посмотрите, как они отреагируют.
		7. Домашнее задание: Задайте первое домашнее задание, которое может быть простым заданием по самонаблюдению, например, отслеживанием настроения и связанных с ним мыслей или деятельности, чтобы начать процесс увеличения самосознания.
		8. Конец сессии: Закончите сессию, подводя итоги ключевых моментов и подтверждая домашнее задание.`
	} else {
		plan = latestSummary
	}

	dialogID, err := database.AddDialog(user_id, plan)
	if err != nil {
		log.Printf("Error adding dialog: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to add dialog: %v", err)}
		c.ServeJSON()
		return
	}

	response := models.AddDialogResponse{
		DialogID: dialogID,
	}
	c.Data["json"] = response
	c.ServeJSON()
}

// SetDialogMessages
// @Summary Установить сообщения для диалога
// @Description Обновляет или добавляет сообщения в диалог пользователя.
// @Tags dialog
// @Accept json
// @Produce json
// @Param body body models.SetMessagesRequest true "Сообщения для добавления"
// @Success 200 {object} models.SetMessagesResponse "Сообщения успешно добавлены"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Ошибка сервера"
// @router /dialog/messages [put]
func (c *DialogController) SetDialogMessages() {
	user_id := AppContext.UserID

	var request models.SetMessagesRequest

	body := c.Ctx.Input.RequestBody

	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Error parsing request: %v", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = models.ErrorResponse{Message: "Invalid request body"}
		c.ServeJSON()
		return
	}

	err = database.SetDialogMessages(user_id, request.Messages, &request.DialogID)
	if err != nil {
		log.Printf("Error setting dialog messages: %v", err)
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = models.ErrorResponse{Message: fmt.Sprintf("Failed to set dialog messages: %v", err)}
		c.ServeJSON()
		return
	}

	response := models.SetMessagesResponse{
		Status: "success",
	}
	c.Data["json"] = response
	c.ServeJSON()
}
