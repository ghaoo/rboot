package wework

import (
	"log"
	"os"
)

func (w *wework) getContacts() []map[string]interface{} {
	// 获取所有部门
	client := w.client.WithSecret(os.Getenv("WORKWX_CONTACT_SECRET"))
	depts, err := client.ListDepartment()
	if err != nil {
		log.Printf("list department err: %v", err)
		return nil
	}

	// 获取部门下的所有成员名单
	var contacts = make([]map[string]interface{}, 0)
	for _, dept := range depts {
		users, err := client.SimpleListUser(dept.ID)
		if err != nil {
			log.Printf("list department %s user err: %v", dept.Name, err)
		}
		for _, user := range users {
			contacts = append(contacts, map[string]interface{}{
				"id":   user.UserID,
				"name": user.Name,
			})
		}
	}

	return contacts
}
