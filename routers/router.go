package routers

import (
	"Healfina_call/controllers"
	"strings"

	httpSwagger "github.com/swaggo/http-swagger"

	beego "github.com/beego/beego/v2/server/web"
	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

var authFilter = func(ctx *beegoCtx.Context) {
	path := ctx.Request.RequestURI
	// Пропускаем некоторые пути
	if path == "/login" || path == "/register" || path == "/api/login" || path == "/api/register" {
		return
	}

	if strings.HasPrefix(path, "/static") {
		// Если пользователь не авторизован, блокируем доступ
		sessUser := ctx.Input.Session("user_id")
		if sessUser == nil {
			ctx.Output.SetStatus(401)
			ctx.Output.JSON(map[string]string{
				"error": "unauthorized",
			}, false, false)
			return
		}
	}

	sessUser := ctx.Input.Session("user_id")
	if sessUser == nil {
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]string{
			"error": "unauthorized",
		}, false, false)
		ctx.ResponseWriter.WriteHeader(401)
		ctx.ResponseWriter.Flush()
		// Завершаем обработку
	}
}

func init() {
	beego.SetStaticPath("/static", "static/dist/static/browser/")
	beego.Router("/api/login", &controllers.AuthController{}, "post:Login")
	beego.Router("/api/register", &controllers.RegisterController{}, "post:Register")
	beego.Router("/logout", &controllers.AuthController{}, "post:Logout")
	beego.Router("/*", &controllers.SpaController{}, "get:Get")
	beego.Router("/", &controllers.MainController{})
	beego.Handler("/swagger/*", httpSwagger.WrapHandler, false)
	beego.Router("/profile", &controllers.ProfileController{}, "get:GetProfile")
	beego.Router("/set_dark_mode", &controllers.MainController{}, "post:SetDarkMode")
	beego.Router("/stream", &controllers.AudioController{}, "get:StreamAudio")
	beego.Router("/dialog", &controllers.DialogController{}, "post:AddDialog")                             //
	beego.Router("/dialog/messages", &controllers.DialogController{}, "put:SetDialogMessages")             //
	beego.Router("/summary", &controllers.SummaryController{}, "get:GetOverallSummary")                    //
	beego.Router("/records/add", &controllers.RecordController{}, "put:AddRecordAfter")                    //
	beego.Router("/records", &controllers.RecordController{}, "get:GetUserRecords")                        //
	beego.Router("/records/update/:record_id", &controllers.RecordController{}, "put:UpdateUserRecord")    //
	beego.Router("/records/delete/:record_id", &controllers.RecordController{}, "delete:DeleteUserRecord") //
	beego.Router("/records/search", &controllers.RecordController{}, "get:SearchUserRecords")

	beego.InsertFilter("/*", beego.BeforeRouter, authFilter)
}
