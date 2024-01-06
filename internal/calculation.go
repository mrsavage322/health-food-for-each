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

func CalculateDinner(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			breakfast, err := getDinner()
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

func CalculateLunch(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			breakfast, err := getLunch()
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
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForBreakfast(context.Background())
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

func getDinner() ([]map[string]float64, error) {
	DayNewCalculation()
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForDinner(context.Background())
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

	fourthProduct := getFoodData[3]
	a1 := fourthProduct["proteins"]
	a2 := fourthProduct["fats"]
	a3 := fourthProduct["carbs"]
	aName := getFoodName[3]

	forDinnerProtein := 0.0
	forDinnerFat := 0.0
	forDinnerCarb := 0.0

	var counterFirst float64
	var counterSecond float64
	var counterThird float64
	var counterFourth float64

	for n.ProteinsNorm*0.35*0.9-forDinnerProtein > 1 {
		forDinnerProtein += x1 * 0.01
		forDinnerFat += x2 * 0.01
		forDinnerCarb += x3 * 0.01
		counterFirst += 1
	}
	for n.CarbsNorm*0.2*0.9-forDinnerCarb > 1 {
		forDinnerProtein += y1 * 0.01
		forDinnerFat += y2 * 0.01
		forDinnerCarb += y3 * 0.01
		counterSecond += 1
	}
	for n.FatsNorm*0.25-forDinnerFat > 1 {
		forDinnerProtein += z1 * 0.01
		forDinnerFat += z2 * 0.01
		forDinnerCarb += z3 * 0.01
		counterThird += 1
	}
	for counterFourth < 250 {
		forDinnerProtein += a1 * 0.01
		forDinnerFat += a2 * 0.01
		forDinnerCarb += a3 * 0.01
		counterFourth += 1
	}

	firstProductGram := counterFirst
	secondProductGram := counterSecond
	thirdsProductGram := counterThird
	fourthProductGram := counterFourth

	var foodsData []map[string]float64

	log.Println(firstProductGram, secondProductGram, thirdsProductGram, fourthProductGram)
	log.Println(xName, yName, zName, aName)
	log.Println(n.ProteinsNorm*0.35, n.FatsNorm*0.2, n.CarbsNorm*0.25, counterFirst, counterSecond)

	n.ProteinsNorm -= forDinnerProtein
	n.FatsNorm -= forDinnerFat
	n.CarbsNorm -= forDinnerCarb

	log.Println(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
	foodData := make(map[string]float64)

	foodData[xName] = firstProductGram
	foodData[yName] = secondProductGram
	foodData[zName] = thirdsProductGram
	foodData[aName] = fourthProductGram

	foodsData = append(foodsData, foodData)

	return foodsData, nil
}

func getLunch() ([]map[string]float64, error) {
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

	fourthProduct := getFoodData[3]
	a1 := fourthProduct["proteins"]
	a2 := fourthProduct["fats"]
	a3 := fourthProduct["carbs"]
	aName := getFoodName[3]

	fifthProduct := getFoodData[4]
	b1 := fifthProduct["proteins"]
	b2 := fifthProduct["fats"]
	b3 := fifthProduct["carbs"]
	bName := getFoodName[4]

	sixthProduct := getFoodData[5]
	c1 := fifthProduct["proteins"]
	c2 := fifthProduct["fats"]
	c3 := fifthProduct["carbs"]
	cName := getFoodName[5]

	forDinnerProtein := 0.0
	forDinnerFat := 0.0
	forDinnerCarb := 0.0

	var counterFirst float64
	var counterSecond float64
	var counterThird float64
	var counterFourth float64
	var counterFifth float64
	var counterSixth float64

	for n.ProteinsNorm*0.35*0.9-forDinnerProtein > 1 {
		forDinnerProtein += x1 * 0.01
		forDinnerFat += x2 * 0.01
		forDinnerCarb += x3 * 0.01
		counterFirst += 1
	}
	for n.CarbsNorm*0.2*0.9-forDinnerCarb > 1 {
		forDinnerProtein += y1 * 0.01
		forDinnerFat += y2 * 0.01
		forDinnerCarb += y3 * 0.01
		counterSecond += 1
	}
	for n.FatsNorm*0.25-forDinnerFat > 1 {
		forDinnerProtein += z1 * 0.01
		forDinnerFat += z2 * 0.01
		forDinnerCarb += z3 * 0.01
		counterThird += 1
	}
	for counterFourth < 250 {
		forDinnerProtein += a1 * 0.01
		forDinnerFat += a2 * 0.01
		forDinnerCarb += a3 * 0.01
		counterFourth += 1
	}

	firstProductGram := counterFirst
	secondProductGram := counterSecond
	thirdsProductGram := counterThird
	fourthProductGram := counterFourth

	var foodsData []map[string]float64

	log.Println(firstProductGram, secondProductGram, thirdsProductGram, fourthProductGram)
	log.Println(xName, yName, zName, aName)
	log.Println(n.ProteinsNorm*0.35, n.FatsNorm*0.2, n.CarbsNorm*0.25, counterFirst, counterSecond)

	n.ProteinsNorm -= forDinnerProtein
	n.FatsNorm -= forDinnerFat
	n.CarbsNorm -= forDinnerCarb

	log.Println(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
	foodData := make(map[string]float64)

	foodData[xName] = firstProductGram
	foodData[yName] = secondProductGram
	foodData[zName] = thirdsProductGram
	foodData[aName] = fourthProductGram

	foodsData = append(foodsData, foodData)

	return foodsData, nil
}
