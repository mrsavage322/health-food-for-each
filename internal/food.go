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

			var isErrorData bool
			p, err := strconv.Atoi(food.Proteins)
			f, err := strconv.Atoi(food.Fats)
			c, err := strconv.Atoi(food.Carbs)

			if len(food.Foodname) > 200 {
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
		if r.Method == http.MethodGet {
			w.Write([]byte("С помощью данного эндпоинта можно добавить продукты. \n" +
				"Для этого отправьте Get запрос со следующими параметрами:\n" +
				"\n" +
				"foodname: название продукта \n" +
				"proteins: количество белков на 100 г\n" +
				"fats: количество жиров на 100 г\n" +
				"carbs: количество углеводов на 100 г\n" +
				"feature: особенность продукта, например: овощ, фрукт, мясо, рыба, орехи, завтрак, крупа\n"))

		}
	}

}

func ShowFood(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {

			getFoodData, err := ConnectionDB.GetUserFood(context.Background())
			if err != nil {
				http.Error(w, "You don't have any products", http.StatusBadRequest)
				return
			}

			for key, value := range getFoodData {
				fmt.Println(key, value)
			}

			//foodname := getFoodData["foodname"]
			//proteins := getFoodData["height"]
			//fats := getFoodData["weight"]
			//carbs := getFoodData["carbs"]
			//feature := getFoodData["gender"]

			//decoder := json.NewDecoder(r.Body)
			//er := decoder.Decode(&food)
			//if er != nil {
			//	http.Error(w, err.Error(), http.StatusBadRequest)
			//	return
			//}

			//resp := Response{Result: "Success!"}
			responseData, err := json.Marshal(getFoodData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(responseData)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Please use GET method for this endpoint"))
		}
	}

}
