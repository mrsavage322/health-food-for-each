package app

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Result string
}

var Request UserRequest
var store = sessions.NewCookieStore([]byte("hello-it-is-my-first-project"))

// Функция получает на вход логин и пароль, возвращает массив байт (хэш)
func getHash(login string, password string) []byte {
	hashPassword := fmt.Sprintf(login + password)
	h := sha256.New()
	h.Write([]byte(hashPassword))
	hash := h.Sum(nil)
	return hash
}

// Хэндлер авторизации принимает json: login - password. Парсит его, хэширует пароль и проверяет в БД
func SignIn(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "session")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&Request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword := getHash(Request.Login, Request.Password)

		err = ConnectionDB.GetAuthData(r.Context(), Request.Login, hashPassword)
		if err != nil {
			resp := Response{Result: "Failed!"}
			responseData, err := json.Marshal(resp)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(responseData)
			return
		}

		resp := Response{Result: "Success!"}
		responseData, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	}

}

// Хэндлер регистрации пользователя.
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&Request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword := getHash(Request.Login, Request.Password)

		ok := ConnectionDB.SetAuthData(r.Context(), Request.Login, hashPassword)
		if ok != true {
			log.Println("Login already exists")
			resp := Response{Result: "Login already exists!"}
			responseData, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

}

func MainPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "You are not auth", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Добро пожаловать!\n Список эндпоитнов для работы с данной программой: \n" +
		"/ - Приветствие и список доступных эндпоитнов\n" +
		"/sign_in - Авторизация\n" +
		"/sign_up - Регистрация\n" +
		"/food/add - Добавить продукт\n" +
		"/settings - Настройки пользователя\n" +
		"/calc/day - Получить питание на один день\n" +
		"/calc/weel - Получить питание на неделю\n"))

}
