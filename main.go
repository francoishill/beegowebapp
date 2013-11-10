// An open source project for Gopher community.
package main

import (
	"fmt"

	"github.com/astaxie/beego"

	"./routers"
	"./utils"
)

// We have to call a initialize function manully
// because we use `bee bale` to pack static resources
// and we cannot make sure that which init() execute first.
func initialize() {
	utils.LoadConfig()
}

func main() {
	initialize()

	if utils.IsProMode {
		beego.Info("Product mode enabled")
	} else {
		beego.Info("Development mode enabled")
	}
	beego.Info(beego.AppName, utils.APP_VER, utils.AppUrl)

	if !utils.IsProMode {
		beego.SetStaticPath("/static_source", "static_source")
	}

	// Register routers.

	general := new(routers.GeneralRouter)
	beego.Router("/", general, "get:PublicHome")
	beego.Router("/home", general, "get:UserHome")
	beego.Router("/acknowledgements", general, "get:Acknowledgements")

	user := new(routers.UserRouter)
	beego.Router("/u/:email", user, "get:Home")

	login := new(routers.LoginRouter)
	beego.Router("/login", login, "post:Login")
	beego.Router("/logout", login, "get:Logout")

	register := new(routers.RegisterRouter)
	beego.Router("/register", register, "post:Register")
	beego.Router("/active/success", register, "get:ActiveSuccess")
	beego.Router("/active/:code([0-9a-zA-Z]+)", register, "get:Active")

	settings := new(routers.SettingsRouter)
	beego.Router("/settings/profile", settings, "get:Profile;post:ProfileSave")

	forgot := new(routers.ForgotRouter)
	beego.Router("/forgot", forgot)
	beego.Router("/reset/:code([0-9a-zA-Z]+)", forgot, "get:Reset;post:ResetPost")

	adminDashboard := new(routers.AdminDashboardRouter)
	beego.Router("/admin", adminDashboard)

	routes := map[string]beego.ControllerInterface{
		"user": new(routers.UserAdminRouter),
		//Can place more routers to also use generic List/Save/Update/Delete urls below
	}
	for name, router := range routes {
		beego.Router(fmt.Sprintf("/admin/:model(%s)", name), router, "get:List")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id(new)", name), router, "get:Create;post:Save")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)", name), router, "get:Edit;post:Update")
		beego.Router(fmt.Sprintf("/admin/:model(%s)/:id([0-9]+)/:action(delete)", name), router, "get:Confirm;post:Delete")
	}

	// "robot.txt"
	beego.Router("/robot.txt", &routers.RobotRouter{})

	// For all unknown pages.
	beego.Run()
}
