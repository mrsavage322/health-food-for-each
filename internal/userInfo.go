package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type UserData struct {
	Gender string `json:"gender"`
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

			var isErrorData bool

			gender := userData.Gender
			age, err := strconv.Atoi(userData.Age)
			height, err := strconv.Atoi(userData.Height)
			weight, err := strconv.Atoi(userData.Weight)
			amount, err := strconv.Atoi(userData.Amount)

			if gender != "M" && gender != "F" {
				isErrorData = true
				http.Error(w, "Have a problem with input data - we need a correct gender: M or F!", http.StatusNotAcceptable)
			}
			if age > 125 || age < 12 || err != nil {
				isErrorData = true
				http.Error(w, "Have a problem with input data - proteins >100!", http.StatusNotAcceptable)
			}
			if height > 251 || height < 63 || err != nil {
				isErrorData = true
				http.Error(w, "Have a problem with input data - fats >100!", http.StatusNotAcceptable)
			}
			if weight > 635 || weight < 20 || err != nil {
				isErrorData = true
				http.Error(w, "Have a problem with input data - carbs >100!", http.StatusNotAcceptable)
			}
			if amount > 6 || amount < 3 || err != nil {
				isErrorData = true
				http.Error(w, "Have a problem with input data - amount must been 2,3,4,5,6!", http.StatusNotAcceptable)
			}

			if isErrorData {
				http.Error(w, "Have a problem with input data", http.StatusNotAcceptable)
			} else {
				er := ConnectionDB.SetUserData(context.Background(), request.Login, gender, age, height, weight, amount)
				if er != nil {
					log.Println("Have a problem with input data")
					resp := Response{Result: "We have a problem with input data!"}
					responseData, er := json.Marshal(resp)
					if er != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					log.Println(age, gender, height, weight, amount)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(responseData)
					return
				}

				resp := Response{Result: "Success!"}
				responseData, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(responseData)
				return
			}
		} else if r.Method == http.MethodGet {
			getUserData, err := ConnectionDB.GetUserData(context.Background())
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else if len(getUserData) == 0 {
				http.Error(w, "Empty!", http.StatusNoContent)
				return
			}
			var responseUserData []UserData
			resp := UserData{
				Gender: getUserData["gender"],
				Age:    getUserData["age"],
				Height: getUserData["height"],
				Weight: getUserData["weight"],
				Amount: getUserData["amount"],
			}
			response := append(responseUserData, resp)

			responseData, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Данный эндпоинт отображает и обновляет настройки пользователя. \n" +
				"Чтобы обновить параметры, отправьте методом Get следующие параметры в виде json: \n" +
				"\n" +
				"gender: Ваш пол, M или F\n" +
				"age: Ваш возраст\n" +
				"height: Ваш рост в см\n" +
				"weight: Ваш вес в кг\n" +
				"amount: количество приемов пищи в день от 3 до 6\n" +
				"\n" +
				"Текущие параметры: \n" +
				"\n"))
			w.Write(responseData)
		}
	}
}
