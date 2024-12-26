package controllers

import (
	"Healfina_call/database"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	fmt.Println(user_id)
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
	data, err := os.ReadFile("static/dist/static/browser/index.html")
	if err != nil {
		c.Ctx.WriteString("Error loading index.html: " + err.Error())
		return
	}
	c.Ctx.Output.ContentType("text/html")
	c.Ctx.Output.Body(data)
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
		// Возвращаем ошибку в JSON
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": "Ошибка при получении профиля"}
		c.ServeJSON()
		return
	}

	// Возвращаем данные о пользователе в JSON
	c.Data["json"] = user
	c.ServeJSON()
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

type RegisterController struct {
	beego.Controller
}

type registerRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

func (c *RegisterController) Register() {
	// Считываем JSON: { "user_id": "...", "password": "..." }
	var req registerRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "invalid json"}
		c.ServeJSON()
		return
	}

	// Преобразуем user_id в int
	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "invalid user_id format"}
		c.ServeJSON()
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	if user == nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = map[string]string{"error": "user not found"}
		c.ServeJSON()
		return
	}

	if user.Password != nil {
		// Значит пароль уже установлен
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "user is already registered"}
		c.ServeJSON()
		return
	}

	if req.Password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "password cannot be empty"}
		c.ServeJSON()
		return
	}

	// Установим пароль (в реальном коде -- хешируем)
	err = database.SetUserPassword(userID, req.Password)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]string{"message": "registered successfully"}
	c.ServeJSON()
}

type AuthController struct {
	beego.Controller
}

type loginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// POST /api/login
func (c *AuthController) Login() {
	var req loginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid JSON"}
		c.ServeJSON()
		return
	}

	if req.UserID == "" || req.Password == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "user_id and password required"}
		c.ServeJSON()
		return
	}

	// Преобразуем user_id в int
	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "invalid user_id format"}
		c.ServeJSON()
		return
	}

	// Ищем в Mongo
	user, err := database.GetUserByID(userID)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}
	if user == nil {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "user not found"}
		c.ServeJSON()
		return
	}

	// Сверка пароля - здесь для примера if req.Password != "123"
	// Лучше хранить хеш
	if req.Password != "123" {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "wrong password"}
		c.ServeJSON()
		return
	}

	// Всё ок: запишем user_id в сессию
	c.SetSession("user_id", user.UserID)
	log.Printf("User %d logged in", user.UserID)
	AppContext.UserID = user.UserID

	c.Data["json"] = map[string]string{"message": "ok"}
	c.ServeJSON()
	c.Redirect("/", 200)
}

// POST /api/logout
func (c *AuthController) Logout() {
	c.DestroySession()
	c.Data["json"] = map[string]string{"message": "logged out"}
	c.ServeJSON()
}

type SpaController struct {
	beego.Controller
}

func (c *SpaController) Get() {
	// Отдаём единую страницу SPA
	// Здесь предполагается, что ваш SPA фронтенд собран и находится в директории static/dist
	c.TplName = "index.html" // Путь к вашему HTML файлу для SPA
}
