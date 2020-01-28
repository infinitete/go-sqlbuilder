package ref

import (
	"testing"
)

type User struct {
	Id   int    `json:"id" field:"id"`
	Name string `json:"name" field:"username"`
	Pass string `json:"pass" field:"password"`
}

var user = User{
	Id:   1,
	Name: "测试用户",
	Pass: "123456",
}

var et, _ = NewEntity("user", user)

func TestEntity_Insert(t *testing.T) {

	sql, values := et.Insert()

	t.Logf("%s\n", sql)
	t.Logf("%#v\n", values)
}

func TestEntity_FindById(t *testing.T) {
	sql, _ := et.FindBy([]string{"Id"}, AND)
	t.Logf("%s\n", sql)

	sql, _ = et.FindBy([]string{"Id", "Name"}, AND)
	t.Logf("%s\n", sql)

	sql, _ = et.FindBy([]string{"Id", "Id"}, OR)
	t.Logf("%s\n", sql)

	sql, err := et.DeleteBy([]string{}, AND)
	t.Logf("%s - %s\n", sql, err)

	sql = et.Update()
	t.Logf("%s\n", sql)

	sql = et.Update("Id", "Name")
	t.Logf("%s\n", sql)
}
