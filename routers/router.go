package routers

import (
	"Healfina_call/controllers"

	httpSwagger "github.com/swaggo/http-swagger"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
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
}
