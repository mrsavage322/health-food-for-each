package internal

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

//Модуль должен принимать на вход json в виде {"login": "xyz", password: "123456789asd"}
//Затем сверяет с БД параметры логина и пароля и отдает ответ Success или Password is incorrect
//
//Если все не ок - ошибка, если ок - успешная авторизация и доступ к таблицам и настройкам

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Result string
}

var request Request
var store = sessions.NewCookieStore([]byte("hello-it-is-my-first-project"))

// Авторизация - принимаем на вход json: login - password. Парсим его, хэшируем пароль и проверяем в БД. Ответ - ОК

func getHash(login string, password string) []byte {
	hashPassword := fmt.Sprintf(login + password)
	h := sha256.New()
	h.Write([]byte(hashPassword))
	hash := h.Sum(nil)
	return hash
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Sign In",
	}

	// Если запрос методом POST
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "session")
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword := getHash(request.Login, request.Password)

		pass, err := ConnectionDB.GetAuthData(context.Background(), request.Login, hashPassword)
		if err != nil {
			log.Println(err)
			log.Println(pass)
			//log.Println(hashPassword)
			resp := Response{Result: "Failed!"}
			responseData, err := json.Marshal(resp)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(responseData)
			return
		} else {
			resp := Response{Result: "Success!"}
			responseData, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			session.Values["authenticated"] = true
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			//log.Println(hashPassword)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(responseData)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	}
	// Используем шаблон для отображения страницы
	tmpl, err := template.New("signin").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<form id="signinForm">
				<label for="login">Логин:</label>
				<input type="text" id="login" name="login" required><br>

				<label for="password">Пароль:</label>
				<input type="password" id="password" name="password" required><br>

				<button type="button" onclick="submitForm()">Авторизоваться</button>
			</form>

			<script>
				function submitForm() {
					var login = document.getElementById("login").value;
					var password = document.getElementById("password").value;

					var data = {
						"login": login,
						"password": password
					};

					fetch('/sign_in', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify(data)
					})
					.then(response => response.json())
					.then(data => {
						console.log('Ответ сервера:', data);
						// Обработка ответа от сервера здесь
					})
					.catch(error => {
						console.error('Ошибка при отправке данных:', error);
					});
				}
			</script>
		</body>
		</html>
	`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pageVariables)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Sign Up",
	}

	// Если запрос методом POST
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashPassword := getHash(request.Login, request.Password)

		er := ConnectionDB.SetAuthData(context.Background(), request.Login, hashPassword)
		if er != nil {
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

	// Используем шаблон для отображения страницы
	tmpl, err := template.New("signup").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<form id="signupForm">
				<label for="login">Логин:</label>
				<input type="text" id="login" name="login" required><br>

				<label for="password">Пароль:</label>
				<input type="password" id="password" name="password" required><br>

				<button type="button" onclick="submitForm()">Зарегистрироваться</button>
			</form>

			<script>
				function submitForm() {
					var login = document.getElementById("login").value;
					var password = document.getElementById("password").value;

					var data = {
						"login": login,
						"password": password
					};

					fetch('/sign_up', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify(data)
					})
					.then(response => response.json())
					.then(data => {
						console.log('Ответ сервера:', data);
						// Обработка ответа от сервера здесь
					})
					.catch(error => {
						console.error('Ошибка при отправке данных:', error);
					});
				}
			</script>
		</body>
		</html>
	`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pageVariables)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		w.Write([]byte("Добро пожаловать!"))
	}
}

type PageVariables struct {
	Title string
}
