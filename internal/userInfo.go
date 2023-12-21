package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserData struct {
	Login  string `json:"login"`
	Age    string `json:"age"`
	Height string `json:"height"`
	Weight string `json:"weight"`
	Amount string `json:"amount"`
}

var userData UserData

func Settings(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&userData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println(request.Login, userData.Age, userData.Height, userData.Weight, userData.Amount)

		}
	}
}
