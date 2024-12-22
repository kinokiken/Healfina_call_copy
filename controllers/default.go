package controllers

import (
	"Healfina_call/database"
	"fmt"
	"log"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// const secretKey = "8F7c862Cl56B7X2sRH3U5cfX7w61v87p"

type MainController struct {
	beego.Controller
}

type ProfileController struct {
	beego.Controller
}

// Get
// @Summary Получить главную страницу
// @Description Возвращает главную страницу с пользовательскими данными и настройками.
// @Tags main
// @Accept json
// @Produce html
// @Success 200 {object} map[string]interface{} "Успешный ответ с данными пользователя"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @router / [get]
func (c *MainController) Get() {

	// user_id := c.GetString("user_id")
	// timestamp := c.GetString("timestamp")
	// signature := c.GetString("signature")

	// if user_id == "" || timestamp == "" || signature == "" {
	// 	c.Ctx.WriteString("Missing parameters")
	// 	return
	// }

	// // Проверка подписи
	// if !verifySignature(user_id, timestamp, signature) {
	// 	c.Ctx.WriteString("Invalid signature")
	// 	return
	// }

	// ts, err := strconv.ParseInt(timestamp, 10, 64)
	// if err != nil || time.Since(time.Unix(ts, 0)) > 60*time.Minute {
	// 	c.Ctx.WriteString("Token expired")
	// 	return
	// }

	var user_id int

	userID := "685751542"

	user_id, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Некорректный userID: %v", err)
		c.Ctx.WriteString("Некорректный идентификатор пользователя")
		return
	}

	AppContext.UserID = user_id
	c.Ctx.SetCookie("user_id", fmt.Sprintf("%d", user_id), 3600)

	user, err := database.GetUserByID(user_id)
	if err != nil {
		c.Data["UserData"] = map[string]string{"error": "Ошибка при получении пользователя"}
	} else if user == nil {
		c.Data["UserData"] = map[string]string{"error": "Пользователь не найден"}
	} else {
		c.Data["UserData"] = user
	}

	darkMode, ok := c.GetSession("dark_mode").(bool)
	if !ok {
		darkMode = false
	}
	log.Println("Dark mode:", darkMode)

	c.Data["DarkMode"] = darkMode
	c.TplName = "index.html"
}

// GetProfile
// @Summary Получить профиль пользователя
// @Description Возвращает данные профиля текущего пользователя.
// @Tags profile
// @Accept json
// @Produce html
// @Success 200 {object} map[string]interface{} "Данные пользователя"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @router /profile [get]
func (c *ProfileController) GetProfile() {
	user_id := AppContext.UserID

	user, err := database.GetUserByID(user_id)
	if err != nil {
		log.Printf("Ошибка при получении профиля: %v", err)
		c.Ctx.WriteString("Ошибка при получении профиля")
		return
	}

	c.Data["User"] = user
	c.TplName = "profile.html"
}

// SetDarkMode
// @Summary Изменить режим отображения
// @Description Переключает темный режим интерфейса для пользователя.
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Состояние режима отображения"
// @Failure 500 {string} string "Ошибка сервера"
// @router /set_dark_mode [post]
func (c *MainController) SetDarkMode() {
	darkMode, ok := c.GetSession("dark_mode").(bool)
	if !ok {
		darkMode = false
	}

	newDarkMode := !darkMode
	c.SetSession("dark_mode", newDarkMode)

	c.Data["json"] = map[string]interface{}{
		"darkMode": newDarkMode,
	}
	c.ServeJSON()
}

// func verifySignature(user_id, timestamp, signature string) bool {
// 	mac := hmac.New(sha256.New, []byte(secretKey))
// 	mac.Write([]byte(fmt.Sprintf("%s:%s", user_id, timestamp)))
// 	expectedSignature := hex.EncodeToString(mac.Sum(nil))
// 	return hmac.Equal([]byte(expectedSignature), []byte(signature))
// }
