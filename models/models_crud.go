package models

import (
	"github.com/astaxie/beego/orm"
)

func InsertModel(m interface{}) error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func ReadModel(m interface{}, fields []string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func UpdateModel(m interface{}, fields []string) error {
	fields = append(fields, "Updated")
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func DeleteModel(m interface{}) error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

//User
func (m *User) Insert() error {
	m.Rands = GetUserSalt()
	return InsertModel(m)
}

func (m *User) Read(fields ...string) error {
	return ReadModel(m, fields)
}

func (m *User) Update(fields ...string) error {
	return UpdateModel(m, fields)
}

func (m *User) Delete() error {
	return DeleteModel(m)
}
