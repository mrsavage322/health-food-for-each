package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const k = 1.2

var newAmount int
var n Norm

type Norm struct {
	ProteinsNorm float64
	FatsNorm     float64
	CarbsNorm    float64
}

// Расчет КБЖУ на день
func DayNewCalculation() {
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

	n.ProteinsNorm = kcalNorm * 0.23 / 4
	n.FatsNorm = kcalNorm * 0.3 / 9
	n.CarbsNorm = kcalNorm * 0.47 / 4

	return
}

func CalculateBreakfast(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			breakfast, err := getBreakfast()
			if err != nil {
				log.Println("Failed to create meal for breakfast")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			responseData, err := json.Marshal(breakfast)
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

func getBreakfast() ([]map[string]float64, error) {
	DayNewCalculation()
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForLunch(context.Background())
	if err != nil {
		log.Println("Dont get fooddata!")
		return nil, err
	}

	firstProduct := getFoodData[0]
	x1 := firstProduct["proteins"]
	x2 := firstProduct["fats"]
	x3 := firstProduct["carbs"]
	xName := getFoodName[0]

	secondProduct := getFoodData[1]
	y1 := secondProduct["proteins"]
	y2 := secondProduct["fats"]
	y3 := secondProduct["carbs"]
	yName := getFoodName[1]

	thirdProduct := getFoodData[2]
	z1 := thirdProduct["proteins"]
	z2 := thirdProduct["fats"]
	z3 := thirdProduct["carbs"]
	zName := getFoodName[2]

	forBreakfastProtein := 0.0
	forBreakfastFat := 0.0
	forBreakfastCarb := 0.0

	var counterFirst float64
	var counterSecond float64

	for (n.ProteinsNorm*0.25-forBreakfastProtein > 1 && n.FatsNorm*0.2-forBreakfastFat > 1) && n.CarbsNorm*0.35*0.9-forBreakfastCarb > 1 {
		log.Println(forBreakfastCarb)
		forBreakfastProtein += x1*0.02 + y1*0.01
		forBreakfastFat += x2*0.02 + y2*0.01
		forBreakfastCarb += x3*0.02 + y3*0.01
		counterFirst += 1
	}
	for n.CarbsNorm*0.35-forBreakfastCarb > 2 && counterSecond < 200 {
		log.Println(forBreakfastCarb)
		forBreakfastProtein += z1 * 0.01
		forBreakfastFat += z2 * 0.01
		forBreakfastCarb += z3 * 0.01
		counterSecond += 1
	}
	firstProductGram := counterFirst * 100 * 0.02
	secondProductProductGram := counterFirst
	thirdsProductGram := counterSecond

	var foodsData []map[string]float64

	log.Println(firstProductGram, secondProductProductGram, thirdsProductGram)
	log.Println(xName, yName, zName)
	log.Println(n.ProteinsNorm*0.25-forBreakfastProtein, n.FatsNorm*0.2-forBreakfastFat, n.CarbsNorm*0.35-forBreakfastCarb, counterFirst, counterSecond)

	n.ProteinsNorm -= forBreakfastProtein
	n.FatsNorm -= forBreakfastFat
	n.CarbsNorm -= forBreakfastCarb

	log.Println(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
	foodData := make(map[string]float64)

	foodData[xName] = firstProductGram
	foodData[yName] = secondProductProductGram
	foodData[zName] = thirdsProductGram

	foodsData = append(foodsData, foodData)

	return foodsData, nil
}
