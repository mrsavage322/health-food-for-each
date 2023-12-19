package internal

import (
	"context"
	"encoding/json"
	"fmt"
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

func AddFood(w http.ResponseWriter, r *http.Request) {
	//pageVariables := PageVariables{
	//	Title: "Add food",
	//}

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

			p, _ := strconv.Atoi(food.Proteins)
			f, _ := strconv.Atoi(food.Fats)
			c, _ := strconv.Atoi(food.Carbs)

			// Починить запрос в БД
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

			//foodname := r.FormValue("foodname")
			//proteins := r.FormValue("proteins")
			//fats := r.FormValue("fats")
			//carbs := r.FormValue("carbs")
			//feature := r.FormValue("feature")
			//
			//p, _ := strconv.Atoi(proteins)
			//f, _ := strconv.Atoi(fats)
			//c, _ := strconv.Atoi(carbs)

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
}
