package routers

import (
	"./../utils"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"strings"

	"./../models"
)

// LoginRouter serves login page.
type LoginRouter struct {
	baseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"
	this.SetPageTitleFromKey("pagetitles.user_login")

	loginRedirect := strings.TrimSpace(this.GetString("to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/home"
	}

	// no need login
	if this.CheckLoginRedirect(false, loginRedirect) {
		return
	}

	if len(loginRedirect) > 0 {
		this.Ctx.SetCookie("login_to", loginRedirect, 0, "/")
	}

	form := models.LoginForm{Locale: this.Locale}
	this.SetFormSets(&form)
}

// Login implemented POST user login.
func (this *LoginRouter) Login() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User

	form := models.LoginForm{Locale: this.Locale}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		if this.IsAjax() {
			goto ajaxError
		}
		return
	}

	if models.VerifyUser(&user, form.Email, form.Password) {
		loginRedirect := strings.TrimSpace(this.Ctx.GetCookie("login_to"))
		if utils.IsMatchHost(loginRedirect) == false {
			loginRedirect = "/home"
		} else {
			this.Ctx.SetCookie("login_to", "", -1, "/")
		}

		// login user
		models.LoginUser(&user, &this.Controller, form.Remember)

		this.setLangCookie(i18n.GetLangByIndex(user.Lang))

		if this.IsAjax() {
			this.Data["json"] = map[string]interface{}{
				"success":  true,
				"message":  this.Tr("auth.login_success_ajax"),
				"redirect": loginRedirect,
			}
			this.ServeJson()
			return
		}

		this.Redirect(loginRedirect, 302)
		return
	} else {
		if this.IsAjax() {
			goto ajaxError
		}
	}
	this.Data["Error"] = true
	return

ajaxError:
	this.Data["json"] = map[string]interface{}{
		"success": false,
		"message": this.Tr("auth.login_error_ajax"),
		"once":    this.Data["once_token"],
	}
	this.ServeJson()
}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	models.LogoutUser(&this.Controller)

	// write flash message
	this.FlashWrite("HasLogout", "true")

	this.Redirect("/login", 302)
}

// RegisterRouter serves register page.
type RegisterRouter struct {
	baseRouter
}

// Get implemented Get method for RegisterRouter.
func (this *RegisterRouter) Get() {
	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"
	this.SetPageTitleFromKey("pagetitles.user_register")

	form := models.RegisterForm{Locale: this.Locale}
	this.SetFormSets(&form)
}

// Register implemented Post method for RegisterRouter.
func (this *RegisterRouter) Register() {
	this.Data["IsRegister"] = true
	this.TplNames = "auth/register.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := models.RegisterForm{Locale: this.Locale}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// Create new user.
	user := new(models.User)
	if err := models.RegisterUser(user, form); err == nil {
		models.SendRegisterMail(this.Locale, user)

		// login user
		models.LoginUser(user, &this.Controller, false)

		this.FlashRedirect("/settings/profile", 302, "RegSuccess")

		return

	} else {
		beego.Error("Register: Failed ", err)
	}
}

// Active implemented check Email actice code.
func (this *RegisterRouter) Active() {
	this.TplNames = "auth/active.html"
	this.SetPageTitleFromKey("pagetitles.user_activate_email")

	// no need active
	if this.CheckActiveRedirect(false) {
		return
	}

	code := this.GetString(":code")

	var user models.User

	if models.VerifyUserActiveCode(&user, code) {
		user.IsActivated = true
		user.Rands = models.GetUserSalt()
		if err := user.Update("IsActivated", "Rands", "Updated"); err != nil {
			beego.Error("Active: user Update ", err)
		}
		if this.isLogin {
			this.user = user
		}

		this.Redirect("/active/success", 302)

	} else {
		this.Data["Success"] = false
	}
}

// ActiveSuccess implemented success page when email active code verified.
func (this *RegisterRouter) ActiveSuccess() {
	this.TplNames = "auth/active.html"
	this.SetPageTitleFromKey("pagetitles.user_activate_success")

	this.Data["Success"] = true
}

// ForgotRouter serves login page.
type ForgotRouter struct {
	baseRouter
}

// Get implemented Get method for ForgotRouter.
func (this *ForgotRouter) Get() {
	this.TplNames = "auth/forgot.html"
	this.SetPageTitleFromKey("pagetitles.user_forgot_password")

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	form := models.ForgotForm{Locale: this.Locale}
	this.SetFormSets(&form)
}

// Get implemented Post method for ForgotRouter.
func (this *ForgotRouter) Post() {
	this.TplNames = "auth/forgot.html"

	// no need login
	if this.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	form := models.ForgotForm{Locale: this.Locale, User: &user}
	// valid form and put errors to template context
	if this.ValidFormSets(&form) == false {
		return
	}

	// send reset password email
	models.SendResetPwdMail(this.Locale, &user)

	this.FlashRedirect("/forgot", 302, "SuccessSend")
}

// Reset implemented user password reset.
func (this *ForgotRouter) Reset() {
	this.TplNames = "auth/reset.html"
	this.SetPageTitleFromKey("pagetitles.user_reset_password")

	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if models.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true
		form := models.ResetPwdForm{Locale: this.Locale}
		this.SetFormSets(&form)
	} else {
		this.Data["Success"] = false
	}
}

// Reset implemented user password reset.
func (this *ForgotRouter) ResetPost() {
	this.TplNames = "auth/reset.html"

	code := this.GetString(":code")
	this.Data["Code"] = code

	var user models.User

	if models.VerifyUserResetPwdCode(&user, code) {
		this.Data["Success"] = true

		form := models.ResetPwdForm{Locale: this.Locale}
		if this.ValidFormSets(&form) == false {
			return
		}

		user.IsActivated = true
		if err := models.SaveNewPassword(&user, form.Password); err != nil {
			beego.Error("ResetPost Save New Password: ", err)
		}

		if this.isLogin {
			models.LogoutUser(&this.Controller)
		}

		this.FlashRedirect("/login", 302, "ResetSuccess")

	} else {
		this.Data["Success"] = false
	}
}
