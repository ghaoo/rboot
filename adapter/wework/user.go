package wework

import (
	"encoding/json"
	"github.com/ghaoo/wxwork"
	"log"
	"os"
)

// 获取所有联系人
// 如果需要联系人信息，可保存到缓存
func (w *wework) getContacts() map[string]wxwork.User {
	// 获取所有部门
	client := w.client.WithSecret(os.Getenv("WORKWX_CONTACT_SECRET"))
	depts, err := client.ListDepartment()
	if err != nil {
		log.Printf("list department err: %v", err)
		return nil
	}

	// 获取部门下的所有成员名单
	var contacts = make(map[string]wxwork.User, 0)
	for _, dept := range depts {
		users, err := client.SimpleListUser(dept.ID)
		if err != nil {
			log.Printf("list department %s user err: %v", dept.Name, err)
		}
		for _, user := range users {
			contacts[user.UserID] = user
		}
	}

	return contacts
}

func (w *wework) storeUsers() {
	users := w.getContacts()

	for id, user := range users {
		bUser, _ := json.Marshal(user)
		w.bot.Store("user", id, bUser)
	}
}
