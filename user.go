package rboot

import (
	"encoding/json"
	"fmt"
	"sync"
)

// User 包含了用户的ID，名称，类型，用户附加信息和成员列表(比如群组)
// Data 为用户附加信息
type User struct {
	ID   string                 `json:"id"`   // 用户唯一标识
	Name string                 `json:"name"` // 用户名称
	Data map[string]interface{} `json:"data"` // 附加信息
}

func newUser(users map[string]interface{}) (*User, error) {
	data, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	var u *User
	err = json.Unmarshal(data, &u)
	return u, err
}

type cacheUser struct {
	sync.Mutex
	contact map[string]*User
}

func newCache() *cacheUser {
	return &cacheUser{
		contact: make(map[string]*User),
	}
}

func (c *cacheUser) update(u map[string]interface{}) error {
	nu, err := newUser(u)
	if err != nil {
		return err
	}

	if len(nu.ID) <= 0 {
		return fmt.Errorf("bad data: %v", u)
	}

	c.contact[nu.ID] = nu

	return nil
}

func (c *cacheUser) append(us []map[string]interface{}) (invalids []map[string]interface{}) {
	c.Lock()
	defer c.Unlock()

	for _, u := range us {
		err := c.update(u)

		if err != nil {
			invalids = append(invalids, u)
		}
	}

	return invalids
}

func (c *cacheUser) getUser(id string) *User {
	c.Lock()
	defer c.Unlock()

	if user, ok := c.contact[id]; ok {
		return user
	}
	return nil
}

func (c *cacheUser) getUserName(id string) string {
	user := c.getUser(id)
	if user != nil {
		return user.Name
	}

	return id
}

func (c *cacheUser) contacts() map[string]*User {
	return c.contact
}

func (c *cacheUser) clear() {
	c.contact = make(map[string]*User)
}

func (c *cacheUser) delete(uid string) {
	delete(c.contact, uid)
}

// SyncContacts 同步联系人信息，invalids为数据格式错误的用户列表，对应 us
func (bot *Robot) SyncContacts(us []map[string]interface{}) (invalids []map[string]interface{}) {
	return bot.contact.append(us)
}

// GetUser 根据用户ID获取用户信息
func (bot *Robot) GetUser(id string) *User {
	return bot.contact.getUser(id)
}

// GetUserName 根据用户ID获取用户名称，当找不到用户时返回发送的 id
func (bot *Robot) GetUserName(id string) string {
	return bot.contact.getUserName(id)
}

// AllUsers 获取所有用户
func (bot *Robot) AllUsers() []*User {
	var values []*User
	for _, value := range bot.contact.contact {
		values = append(values, value)
	}
	return values
}

// ClearUser 清空所有用户
func (bot *Robot) ClearUser() {
	bot.contact.clear()
}
