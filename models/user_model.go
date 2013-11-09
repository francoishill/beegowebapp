package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"./../utils"
)

// global settings name -> value
type Setting struct {
	Id      int64
	Name    string `orm:"unique"`
	Value   string `orm:"type(text)"`
	Updated string `orm:"auto_now"`
}

// main user table
// IsAdmin: user is admininstator
// IsActivated: set active when email is verified
// IsBanned: forbid user login
type User struct {
	Id          int64  `orm:"pk;auto"`
	Email       string `orm:"unique;size(240)" valid:"Required; Email; MaxSize(240)"` //User unique email
	Password    string `orm:"size(240)"`                                              //User password
	IsActivated bool   `orm:"index;default(false)"`
	/*DatabaseName string    `orm:"index;null"` //Allow null as we first insert the record to get the ID and use it as part of Dbname
	DbUsername   string    `orm:"size(240)"`
	DbPassword   string    `orm:"size(240)"`
	//ServerLocation string `orm:"null"` //Can add an url to the server later on*/
	IsAdmin  bool      `orm:"index"`
	IsBanned bool      `orm:"index"`
	Rands    string    `orm:"size(10)"`
	Created  time.Time `orm:"auto_now_add"`
	Updated  time.Time `orm:"auto_now"`
	Lang        int              `orm:"index"`
	LangAdds    SliceStringField `orm:"size(50)"`
}

func (u *User) TableEngine() string {
	return "INNODB"
}

/*func (m *User) RefreshFavTopics() int {
	cnt, err := FollowTopics().Filter("User", m.Id).Count()
	if err == nil {
		m.FavTopics = int(cnt)
		m.Update("FavTopics")
	}
	return m.FavTopics
}*/

func (m *User) String() string {
	return utils.ToStr(m.Id)
}

func (m *User) Link() string {
	return fmt.Sprintf("%su/%s", utils.AppUrl, m.Email)
}

/*func (m *User) AvatarLink() string {
	return fmt.Sprintf("%s%s", utils.AvatarURL, m.GrEmail)
}*/

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user").OrderBy("-Id")
}

/*// user follow
type Follow struct {
	Id         int32
	User       *User `orm:"rel(fk)"`
	FollowUser *User `orm:"rel(fk)"`
	Mutual     bool
	Created    time.Time `orm:"auto_now_add"`
}

func (*Follow) TableUnique() [][]string {
	return [][]string{
		[]string{"User", "FollowUser"},
	}
}

func (m *Follow) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Follow) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func Follows() orm.QuerySeter {
	return orm.NewOrm().QueryTable("follow").OrderBy("-Id")
}*/

func init() {
	orm.RegisterModel(new(Setting), new(User) /*, new(Follow)*/)
}
