package gogshook

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type Commits struct {
	Message string
	Url     string
}

type Res struct {
	Commits []Commits

	Repository struct {
		Name      string
		UpdatedAt time.Time
	}
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("read body err, %v\n", err)
			return
		}
		defer r.Body.Close()

		var res Res

		err = json.Unmarshal(body, &res)
		if err != nil {
			logrus.Errorf("unmarshal json err, %v\n", err)
			return
		}

		fmt.Println(res)

	} else {
		w.Write([]byte("down"))
	}

}
