package models

import (
	"encoding/hex"
	"fmt"
	// "strings"
	// "time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"

	//"github.com/dchest/uniuri"

	//"database/sql"

	"./../utils"
)

// CanRegistered checks if the e-mail is available.
func CanRegistered(email string) (bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("Email", email)

	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("user").SetCond(cond).Values(&maps, "Email")
	if err != nil {
		return false, err
	}

	e1 := true

	if n > 0 {
		for _, m := range maps {
			if e1 && orm.ToStr(m["Email"]) == email {
				e1 = false
			}
		}
	}

	return e1, nil
}

// check if exist user by email
func HasUser(user *User, email string) bool {
	var err error
	qs := orm.NewOrm()
	user.Email = email
	err = qs.Read(user, "Email")
	if err == nil {
		return true
	}
	return false
}

// return a user salt token
func GetUserSalt() string {
	return utils.GetRandomString(10)
}

// register create user
func RegisterUser(user *User, form RegisterForm) error {
	// use random salt encode password
	salt := GetUserSalt()
	pwd := utils.EncodePassword(form.Password, salt)

	user.Email = form.Email
	//user.NickName = form.UserName

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	// save md5 email value for gravatar
	//user.GrEmail = utils.EncodeMd5(form.Email)

	err := user.Insert()
	if err != nil {
		return err
	}

	return nil
}

// set a new password to user
func SaveNewPassword(user *User, password string) error {
	salt := GetUserSalt()
	user.Password = fmt.Sprintf("%s$%s", salt, utils.EncodePassword(password, salt))
	user.Rands = GetUserSalt()
	return user.Update("Password", "Rands")
}

// login user
func LoginUser(user *User, c *beego.Controller, remember bool) {
	// weird way of beego session regenerate id...
	c.CruSession = beego.GlobalSessions.SessionRegenerateId(c.Ctx.ResponseWriter, c.Ctx.Request)
	c.Ctx.Input.CruSession = c.CruSession

	sess := c.StartSession()
	sess.Set("auth_user_id", user.Id)
}

// logout user
func LogoutUser(c *beego.Controller) {
	sess := c.StartSession()
	sess.Delete("auth_user_id")
	c.DestroySession()
}

// get user if key exist in session
func GetUserFromSession(user *User, sess session.SessionStore) bool {
	if id, ok := sess.Get("auth_user_id").(int64); ok && id > 0 {
		*user = User{Id: id}
		if user.Read() == nil {
			return true
		}
	}

	return false
}

// verify email and password
func VerifyUser(user *User, email, password string) bool {
	// search user by email
	if HasUser(user, email) == false {
		return false
	}

	if VerifyPassword(password, user.Password) {
		// success
		return true
	}
	return false
}

// compare raw password and encoded password
func VerifyPassword(rawPwd, encodedPwd string) bool {

	// split
	var salt, encoded string
	if len(encodedPwd) > 11 {
		salt = encodedPwd[:10]
		encoded = encodedPwd[11:]
	}

	return utils.EncodePassword(rawPwd, salt) == encoded
}

// get user by erify code
func getVerifyUser(user *User, code string) bool {
	if len(code) <= utils.TimeLimitCodeLength {
		return false
	}

	// use tail hex email query user
	hexStr := code[utils.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		user.Email = string(b)
		if user.Read("Email") == nil {
			return true
		}
	}

	return false
}

// verify active code when active account
func VerifyUserActiveCode(user *User, code string) bool {
	days := utils.ActivationCodeLives

	if getVerifyUser(user, code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(user.Id) + user.Email + user.Password + user.Rands

		return utils.VerifyTimeLimitCode(data, days, prefix)
	}

	return false
}

// create a time limit code for user active
func CreateUserActiveCode(user *User, startInf interface{}) string {
	days := utils.ActivationCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.Password + user.Rands
	code := utils.CreateTimeLimitCode(data, days, startInf)

	// add tail hex email
	code += hex.EncodeToString([]byte(user.Email))
	return code
}

// verify code when reset password
func VerifyUserResetPwdCode(user *User, code string) bool {
	days := utils.ResetPwdCodeLives

	if getVerifyUser(user, code) {
		// time limit code
		prefix := code[:utils.TimeLimitCodeLength]
		data := utils.ToStr(user.Id) + user.Email + user.Password + user.Rands + user.Updated.String()

		return utils.VerifyTimeLimitCode(data, days, prefix)
	}

	return false
}

// create a time limit code for user reset password
func CreateUserResetPwdCode(user *User, startInf interface{}) string {
	days := utils.ResetPwdCodeLives
	data := utils.ToStr(user.Id) + user.Email + user.Password + user.Rands + user.Updated.String()
	code := utils.CreateTimeLimitCode(data, days, startInf)

	// add tail hex email
	code += hex.EncodeToString([]byte(user.Email))
	return code
}
