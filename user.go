package rboot

import (
	"encoding/json"
)

type User struct {
	ID         string                 // 用户唯一标识
	Name       string                 // 用户名称
	Type       string                 // 用户类型
	Data       map[string]interface{} // 附加信息
	MemberList []*User                // 成员列表
}

type cacheUser map[string]*User

func newContact(m map[string]interface{}) (*User, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var u *User
	err = json.Unmarshal(data, &u)
	return u, err
}

// GetUserByUserID 根据用户ID获取用户信息
/*func (bot *Robot) GetUser(id string) *User {
	if user, found := bot.users[id]; found {
		return user
	}
	return nil
}

func (bot *Robot) GetUserName(id string) string {
	if user, found := bot.users[id]; found {
		return user.Name
	}
	return id
}

// AllUsers 获取所有用户
func (bot *Robot) AllUsers() []*User {
	var values []*User
	for _, value := range bot.users {
		values = append(values, value)
	}
	return values
}

// ClearUser 清空所有用户
func (bot *Robot) ClearUser() {
	bot.users = make(map[string]*User)
}*/
