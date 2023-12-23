package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const k = 1.2

var newAmount int

// Расчет КБЖУ на день
func DayNewCalculation() (proteinsNorm, fatsNorm, carbsNorm float64) {
	getUserData, err := ConnectionDB.GetUserData(context.Background())
	if err != nil {
		return
	}

	age, err := strconv.ParseFloat(getUserData["age"], 3)
	height, err := strconv.ParseFloat(getUserData["height"], 3)
	weight, err := strconv.ParseFloat(getUserData["weight"], 3)
	newAmount, err = strconv.Atoi(getUserData["amount"])
	gender := getUserData["gender"]

	log.Println(age, height, weight, newAmount, k)

	var kcalNorm float64
	if gender == "M" {
		kcalNorm = (10 * weight) + (6.25 * height) - (5 * age) + 5
	} else {
		kcalNorm = (10 * weight) + (6.25 * height) - (5 * age) - 161
	}

	proteinsNorm = kcalNorm * 0.23 / 4
	fatsNorm = kcalNorm * 0.3 / 9
	carbsNorm = kcalNorm * 0.5 / 4

	return proteinsNorm, fatsNorm, carbsNorm

	//set := UserData{
	//	Age:    getUserData["age"],
	//	Height: getUserData["height"],
	//	Weight: getUserData["weight"],
	//	Amount: getUserData["amount"],
	//}
	//fmt.Println(set)
	//response := append(responseUserData, resp)
	//basic man = (10 * weight) + (6.25 * height) - (5 * age) + 5 * koef
	//basic woman = (10 * weight) + (6.25 * height) - (5 * age) - 161 * koef
	//koef = 1.15
	// prot = 1.4 fat = 0.8 carbs = 2.8
	//protein norm = basic * 0.23  fats * 0.3  carbs * 0.47

}

// Расчет приема пищи
func CreateMealForBreakast() error {
	var proteinForMeal float64
	proteinsNorm, fatsNorm, carbsNorm := DayNewCalculation()
	fmt.Println(proteinsNorm, fatsNorm, carbsNorm)
	getFoodData, err := ConnectionDB.CreateMealForLunch(context.Background())
	if err != nil {
		log.Println("Dont get fooddata!")
		return err
	}

	//xyz := make(map[string]string)

	//TODO: Переписать циклом
	firstProduct := getFoodData[0]
	x := firstProduct["foodname"]
	x1 := firstProduct["proteins"]
	x2 := firstProduct["fats"]
	x3 := firstProduct["carbs"]

	secondProduct := getFoodData[1]
	y := secondProduct["foodname"]
	y1 := secondProduct["proteins"]
	y2 := secondProduct["fats"]
	y3 := secondProduct["carbs"]

	thirdProduct := getFoodData[2]
	z := thirdProduct["foodname"]
	z1 := thirdProduct["proteins"]
	z2 := thirdProduct["fats"]
	z3 := thirdProduct["carbs"]

	//TODO: Вынести кэфы
	if x1 > y1 && x1 > z1 {
		proteinForMeal = x1*0.7 + y1*0.15 + z*0.15
	} else if y1 > x2 && y1 > z1 {
		proteinForMeal = x1*0.15 + y1*0.7 + z*0.15
	} else {
		proteinForMeal = x1*0.15 + y1*0.15 + z*0.7
	}

	if proteinsNorm-proteinForMeal > 2 {
		//Пересчитывает умнажая на 2
	} else if proteinsNorm-proteinForMeal < 0 {
		//Пересчитывает деля на 2
	} else {
		//WIN
	}

	log.Println(getFoodData)
	return nil
}

func CalculateBreakfast(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			err := CreateMealForBreakast()
			if err != nil {
				log.Println("Failed to create meal for breakfast")
				http.Error(w, err.Error(), http.StatusInternalServerError)
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

//Расчет приема пищи на обед, ужин

//Расчет приема пищи на перекус
