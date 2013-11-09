package routers

import (
	"github.com/astaxie/beego"

	"./../models"
)

// SettingsRouter serves user settings.
type SettingsRouter struct {
	baseRouter
}

// Profile implemented user profile settings page.
func (this *SettingsRouter) Profile() {
	this.TplNames = "settings/profile.html"

	// need login
	if this.CheckLoginRedirect() {
		return
	}

	form := models.ProfileForm{Locale: this.Locale}
	form.SetFromUser(&this.user)
	this.SetFormSets(&form)

	formPwd := models.PasswordForm{Locale: this.Locale}
	this.SetFormSets(&formPwd)
}

// ProfileSave implemented save user profile.
func (this *SettingsRouter) ProfileSave() {
	this.TplNames = "settings/profile.html"

	if this.CheckLoginRedirect() {
		return
	}

	action := this.GetString("action")

	if this.IsAjax() {
		switch action {
		case "send-verify-email":
			if this.user.IsActivated {
				this.Data["json"] = false
			} else {
				models.SendActiveMail(this.Locale, &this.user)
				this.Data["json"] = true
			}

			this.ServeJson()
			return
		}
		return
	}

	profileForm := models.ProfileForm{Locale: this.Locale}
	profileForm.SetFromUser(&this.user)

	pwdForm := models.PasswordForm{Locale: this.Locale}

	this.Data["Form"] = profileForm

	switch action {
	case "save-profile":
		if this.ValidFormSets(&profileForm) {
			if err := profileForm.SaveUserProfile(&this.user); err != nil {
				beego.Error("ProfileSave: save-profile", err)
			}
			this.FlashRedirect("/settings/profile", 302, "ProfileSave")
			return
		}

	case "change-password":
		if this.ValidFormSets(&pwdForm) {
			if models.VerifyPassword(pwdForm.PasswordOld, this.user.Password) {
				// verify success and save new password
				if err := models.SaveNewPassword(&this.user, pwdForm.Password); err == nil {
					this.FlashRedirect("/settings/profile", 302, "PasswordSave")
					return
				} else {
					beego.Error("ProfileSave: change-password", err)
				}
			} else {
				this.SetFormError(&pwdForm, "PasswordOld", "Your old password not correct")
			}
		}

	default:
		this.Redirect("/settings/profile", 302)
		return
	}

	if action != "save-profile" {
		this.SetFormSets(&profileForm)
	}
	if action != "change-password" {
		this.SetFormSets(&pwdForm)
	}
}
