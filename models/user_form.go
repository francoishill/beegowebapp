package models

import (
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/beego/i18n"

	"./../utils"
)

// Register form
type RegisterForm struct {
	Email      string      `valid:"Required;Email;MaxSize(80)"`
	Password   string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	Locale     i18n.Locale `form:"-"`
}

func (form *RegisterForm) Valid(v *validation.Validation) {

	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "auth.repassword_not_match")
		return
	}

	e1, _ := CanRegistered(form.Email)

	if !e1 {
		v.SetError("Email", "auth.email_already_taken")
	}
}

func (form *RegisterForm) Labels() map[string]string {
	return map[string]string{
		"Email":      "auth.login_email",
		"Password":   "auth.login_password",
		"PasswordRe": "auth.retype_password",
	}
}

func (form *RegisterForm) Helps() map[string]string {
	return map[string]string{
	//"Email": form.Locale.Tr("valid.min_length_is", 5) + ", " + form.Locale.Tr("valid.only_contains", "a-z 0-9 - _"),
	}
}

func (form *RegisterForm) Placeholders() map[string]string {
	return map[string]string{
		"Email":      "auth.plz_enter_email",
		"Password":   "auth.plz_enter_password",
		"PasswordRe": "auth.plz_reenter_password",
	}
}

// Login form
type LoginForm struct {
	Email    string `valid:"Required; Email; MaxSize(240)"`
	Password string `form:"type(password)" valid:"Required"`
	Remember bool
	Locale   i18n.Locale `form:"-"`
}

func (form *LoginForm) Labels() map[string]string {
	return map[string]string{
		"Email":    "auth.email",
		"Password": "auth.login_password",
		"Remember": "auth.login_remember_me",
	}
}

func (form *LoginForm) Valid(v *validation.Validation) {
	valid := validation.Validation{}
	valid.Email(form.Email, "email")
	if valid.HasErrors() {
		v.SetError("Email", form.Locale.Tr("valid.invalid_email"))
	}
}

// Forgot form
type ForgotForm struct {
	Email  string      `valid:"Required;Email;MaxSize(80)"`
	User   *User       `form:"-"`
	Locale i18n.Locale `form:"-"`
}

func (form *ForgotForm) Labels() map[string]string {
	return map[string]string{
		"Email": "auth.login_email",
	}
}

func (form *ForgotForm) Helps() map[string]string {
	return map[string]string{
		"Email": "auth.forgotform_email_help",
	}
}

func (form *ForgotForm) Valid(v *validation.Validation) {
	if HasUser(form.User, form.Email) == false {
		v.SetError("Email", "auth.forgotform_wrong_email")
	}
}

// Reset password form
type ResetPwdForm struct {
	Password   string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	Locale     i18n.Locale `form:"-"`
}

func (form *ResetPwdForm) Valid(v *validation.Validation) {
	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "auth.repassword_not_match")
		return
	}
}

func (form *ResetPwdForm) Labels() map[string]string {
	return map[string]string{
		"PasswordRe": "auth.retype_password",
	}
}

func (form *ResetPwdForm) Placeholders() map[string]string {
	return map[string]string{
		"Password":   "auth.plz_enter_password",
		"PasswordRe": "auth.plz_reenter_password",
	}
}

// Settings Profile form
type ProfileForm struct {
	NickName    string      `valid:"Required;MaxSize(30)"`
	Url         string      `valid:"MaxSize(100)"`
	Info        string      `form:"type(textarea)" valid:"MaxSize(255)"`
	Email       string      `valid:"Required;Email;MaxSize(100)"`
	PublicEmail bool        `valid:""`
	GrEmail     string      `valid:"Required;MaxSize(80)"`
	Lang        int              `form:"type(select);attr(rel,select2)" valid:""`
	LangAdds    SliceStringField `form:"type(select);attr(rel,select2);attr(multiple,multiple)" valid:""`
	Locale      i18n.Locale `form:"-"`
}

func (form *ProfileForm) LangSelectData() [][]string {
	langs := utils.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *ProfileForm) LangAddsSelectData() [][]string {
	langs := utils.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *ProfileForm) Valid(v *validation.Validation) {
	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "Can not be empty")
	}

	if len(form.LangAdds) > 0 {
		adds := make(SliceStringField, 0, len(form.LangAdds))
		for _, l := range form.LangAdds {
			if d, err := utils.StrTo(l).Int(); err == nil {
				if form.Lang == d {
					continue
				}
				if len(i18n.GetLangByIndex(form.Lang)) == 0 {
					v.SetError("Lang", "Can not be empty")
					return
				}
				adds = append(adds, l)
			}
		}
		form.LangAdds = adds
	}
}

func (form *ProfileForm) SetFromUser(user *User) {
	utils.SetFormValues(user, form)
}

func (form *ProfileForm) SaveUserProfile(user *User) error {
	// set md5 value if the value is an email
	if strings.IndexRune(form.GrEmail, '@') != -1 {
		form.GrEmail = utils.EncodeMd5(form.GrEmail)
	}

	changes := utils.FormChanges(user, form)
	if len(changes) > 0 {
		// if email changed then need re-active
		if user.Email != form.Email {
			user.IsActivated = false
			changes = append(changes, "IsActivated")
		}

		utils.SetFormValues(form, user)
		return user.Update(changes...)
	}
	return nil
}

func (form *ProfileForm) Labels() map[string]string {
	return map[string]string{
		"Lang":        "auth.profile_lang",
		"LangAdds":    "auth.profile_lang_additional",
		"NickName":    "model.user_nickname",
		"PublicEmail": "auth.profile_publicemail",
		"GrEmail":     "auth.profile_gremail",
		"Url":         "auth.profile_url",
	}
}

func (form *ProfileForm) Helps() map[string]string {
	return map[string]string{
		"GrEmail": "auth.profile_gremail_help",
		"Info":    form.Locale.Tr("Max-length is %d", 255),
	}
}

func (form *ProfileForm) Placeholders() map[string]string {
	return map[string]string{
		"NickName": "auth.plz_enter_nickname",
		"GrEmail":  "auth.plz_enter_gremail",
		"Url":      "auth.plz_enter_website",
		"Info":     "auth.plz_enter_your_info",
	}
}

// Change password form
type PasswordForm struct {
	PasswordOld string      `form:"type(password)" valid:"Required"`
	Password    string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe  string      `form:"type(password)" valid:"Required;MinSize(4);MaxSize(30)"`
	User        *User  `form:"-"`
	Locale      i18n.Locale `form:"-"`
}

func (form *PasswordForm) Valid(v *validation.Validation) {
	// Check if passwords of two times are same.
	if form.Password != form.PasswordRe {
		v.SetError("PasswordRe", "auth.repassword_not_match")
		return
	}

	if VerifyPassword(form.PasswordOld, form.User.Password) == false {
		v.SetError("PasswordOld", "auth.old_password_wrong")
	}
}

func (form *PasswordForm) Labels() map[string]string {
	return map[string]string{
		"PasswordOld": "auth.old_password",
		"Password":    "auth.new_password",
		"PasswordRe":  "auth.retype_password",
	}
}

func (form *PasswordForm) Placeholders() map[string]string {
	return map[string]string{
		"PasswordOld": "auth.plz_enter_old_password",
		"Password":    "auth.plz_enter_new_password",
		"PasswordRe":  "auth.plz_reenter_password",
	}
}

type UserAdminForm struct {
	Create      bool   `form:"-"`
	Id          int64  `form:"-"`
	Email       string `valid:"Required;Email;MaxSize(100)"`
	PublicEmail bool   ``
	//NickName    string      `valid:"Required;MaxSize(30)"`
	Url  string `valid:"MaxSize(100)"`
	Info string `form:"type(textarea)" valid:"MaxSize(255)"`
	//GrEmail   string      `valid:"Required;MaxSize(80)"`
	Followers   int         ``
	Following   int         ``
	IsAdmin     bool        ``
	IsActivated bool        ``
	IsBanned    bool        ``
	Lang        int              `form:"type(select);attr(rel,select2)" valid:""`
	LangAdds    SliceStringField `form:"type(select);attr(rel,select2);attr(multiple,multiple)" valid:""`
}

func (form *UserAdminForm) LangSelectData() [][]string {
	langs := utils.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *UserAdminForm) LangAddsSelectData() [][]string {
	langs := utils.Langs
	data := make([][]string, 0, len(langs))
	for i, lang := range langs {
		data = append(data, []string{lang, utils.ToStr(i)})
	}
	return data
}

func (form *UserAdminForm) Valid(v *validation.Validation) {
	qs := Users()

	if CheckIsExist(qs, "Email", form.Email, form.Id) {
		v.SetError("Email", "auth.email_already_taken")
	}

	if len(i18n.GetLangByIndex(form.Lang)) == 0 {
		v.SetError("Lang", "Can not be empty")
	}

	if len(form.LangAdds) > 0 {
		adds := make(SliceStringField, 0, len(form.LangAdds))
		for _, l := range form.LangAdds {
			if d, err := utils.StrTo(l).Int(); err == nil {
				if form.Lang == d {
					continue
				}
				if len(i18n.GetLangByIndex(form.Lang)) == 0 {
					v.SetError("Lang", "Can not be empty")
					return
				}
				adds = append(adds, l)
			}
		}
		form.LangAdds = adds
	}
}

func (form *UserAdminForm) Helps() map[string]string {
	return nil
}

func (form *UserAdminForm) Labels() map[string]string {
	return nil
}

func (form *UserAdminForm) SetFromUser(user *User) {
	utils.SetFormValues(user, form)
}

func (form *UserAdminForm) SetToUser(user *User) {
	// set md5 value if the value is an email
	/*if strings.IndexRune(form.GrEmail, '@') != -1 {
		form.GrEmail = utils.EncodeMd5(form.GrEmail)
	}*/

	utils.SetFormValues(form, user)
}
