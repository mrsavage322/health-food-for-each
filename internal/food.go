package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type FoodData struct {
	Foodname string `json:"foodname"`
	Proteins string `json:"proteins"`
	Fats     string `json:"fats"`
	Carbs    string `json:"carbs"`
	Feature  string `json:"feature"`
}

var food FoodData
var isErrorData bool

func AddFood(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Add food",
	}

	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&food)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println(food.Foodname, food.Proteins, food.Fats, food.Carbs, food.Feature)

			p, err := strconv.Atoi(food.Proteins)
			f, err := strconv.Atoi(food.Fats)
			c, err := strconv.Atoi(food.Carbs)

			if len(food.Foodname) > 50 {
				isErrorData = true
				log.Println("Have a problem with input data!")
				http.Error(w, "Have a problem with input data - too long!", http.StatusNotAcceptable)
			}
			if p > 100 || err != nil {
				isErrorData = true
				log.Println("Have a problem with input data!")
				http.Error(w, "Have a problem with input data - proteins >100!", http.StatusNotAcceptable)
			}
			if f > 100 || err != nil {
				isErrorData = true
				log.Println("Have a problem with input data!")
				http.Error(w, "Have a problem with input data - fats >100!", http.StatusNotAcceptable)
			}
			if c > 100 || err != nil {
				isErrorData = true
				log.Println("Have a problem with input data!")
				http.Error(w, "Have a problem with input data - carbs >100!", http.StatusNotAcceptable)
			}

			feature := [7]string{"мясо", "перекус", "овощ", "фрукт", "орехи", "крупа", "рыба"}
			var isFeature = false
			for _, v := range feature {
				if food.Feature == v {
					isFeature = true
					break
				}
			}
			if !isFeature {
				log.Println("Incorrect feature")
				incorrectFeature := fmt.Sprintf("Incorrect feature, you can use: %s", feature)
				http.Error(w, incorrectFeature, http.StatusNotAcceptable)
			}

			if isErrorData || !isFeature {
				http.Error(w, "Have a problem with input data", http.StatusNotAcceptable)
			} else {
				er := ConnectionDB.SetFoodData(context.Background(), food.Foodname, p, f, c, food.Feature)
				if er != nil {
					log.Println("Have a problem with input data")
					resp := Response{Result: "We have a problem with input data!"}
					responseData, er := json.Marshal(resp)
					if er != nil {
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
				return
			}
		}

	}
	// Проверяем метод запроса

	// Используем шаблон для отображения страницы
	//tmpl, err := template.New("index").Parse(`
	//		<!DOCTYPE html>
	//			<html>
	//			<head>
	//				<title>{{.Title}}</title>
	//			</head>
	//			<body>
	//				<h1>{{.Title}}</h1>
	//				<form method="post" action="/">
	//					<label for="foodname">Foodname:</label>
	//					<input type="text" id="foodname" name="foodname" required><br>
	//					<label for="proteins">Proteins:</label>
	//					<input type="text" id="proteins" name="proteins" required><br>
	//					<label for="fats">Fats:</label>
	//					<input type="text" id="fats" name="fats" required><br>
	//					<label for="carbs">Carbs:</label>
	//					<input type="text" id="carbs" name="carbs" required><br>
	//					<label for="feature">Feature:</label>
	//					<input type="text" id="feature" name="feature" required><br>
	//
	//
	//					<input type="submit" value="Submit">
	//				</form>
	//			</body>
	//			</html>
	//		`)
	//
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//tmpl.Execute(w, pageVariables)
	//TODO: убрать вывод в консоль html
	tmpl, err := template.New("add").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<form id="signupForm">
				<label for="foodname">foodname:</label>
				<input type="text" id="foodname" name="foodname" required><br>

				<label for="proteins">proteins:</label>
				<input type="text" id="proteins" name="proteins" required><br>

				<label for="fats">fats:</label>
				<input type="text" id="fats" name="fats" required><br>

				<label for="carbs">carbs:</label>
				<input type="text" id="carbs" name="carbs" required><br>

				<label for="feature">feature:</label>
				<input type="text" id="feature" name="feature" required><br>

				<button type="button" onclick="submitForm()">Добавить!</button>
			</form>

			<script>
				function submitForm() {
					var foodname = document.getElementById("foodname").value;
					var proteins = document.getElementById("proteins").value;
					var fats = document.getElementById("fats").value;
					var carbs = document.getElementById("carbs").value;
					var feature = document.getElementById("feature").value;

					var data = {
						"foodname": foodname,
						"proteins": proteins,
						"fats": fats,
						"carbs": carbs,
						"feature": feature

					};

					fetch('/add', {
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
