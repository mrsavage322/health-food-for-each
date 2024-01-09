package internal

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Коэф физ активности
const k = 1.2

// Количество приемов пищи
var mealAmount int

type Norm struct {
	ProteinsNorm float64
	FatsNorm     float64
	CarbsNorm    float64
}

var n Norm
var def Norm

// Расчет КБЖУ на день, возращаем proteins, fats, carbs
func DayNewCalculation() (float64, float64, float64) {
	getUserData, err := ConnectionDB.GetUserData(context.Background())
	if err != nil {
		return 0, 0, 0
	}

	age, err := strconv.ParseFloat(getUserData["age"], 3)
	height, err := strconv.ParseFloat(getUserData["height"], 3)
	weight, err := strconv.ParseFloat(getUserData["weight"], 3)
	mealAmount, err = strconv.Atoi(getUserData["amount"])
	gender := getUserData["gender"]

	log.Println(age, height, weight, mealAmount, k)

	var kcalNorm float64
	if gender == "M" {
		kcalNorm = (10 * weight) + (6.25 * height) - (5 * age) + 5
	} else {
		kcalNorm = (10 * weight) + (6.25 * height) - (5 * age) - 161
	}

	return kcalNorm * 0.23 / 4, kcalNorm * 0.3 / 9, kcalNorm * 0.47 / 4
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
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForBreakfast(context.Background())
	if err != nil {
		log.Println("Dont get food data!")
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
	firstProductGram := counterFirst * 2
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
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForDinner(context.Background())
	if err != nil {
		log.Println("Dont get food data!")
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
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForLunch(context.Background())
	if err != nil {
		log.Println("Dont get food data!")
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

	//fourthProduct := getFoodData[3]
	//a1 := fourthProduct["proteins"]
	//a2 := fourthProduct["fats"]
	//a3 := fourthProduct["carbs"]
	aName := getFoodName[3]

	fifthProduct := getFoodData[4]
	b1 := fifthProduct["proteins"]
	b2 := fifthProduct["fats"]
	b3 := fifthProduct["carbs"]
	bName := getFoodName[4]

	sixthProduct := getFoodData[5]
	c1 := sixthProduct["proteins"]
	c2 := sixthProduct["fats"]
	c3 := sixthProduct["carbs"]
	cName := getFoodName[5]

	forLunchProtein := 0.0
	forLunchFat := 0.0
	forLunchCarb := 0.0

	var counterFirst float64
	var counterSecond float64
	var counterThird float64
	var counterFourth float64
	var counterFifth float64
	//var counterSixth float64

	log.Println("Attention!", n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)

	for counterFirst < 120 {
		forLunchProtein += c1 * 0.2
		forLunchFat += c2 * 0.2
		forLunchCarb += c3 * 0.20
		counterFirst += 20
	}

	for counterSecond < 200 {
		forLunchProtein += b1 * 0.5
		forLunchFat += b2 * 0.5
		forLunchCarb += b3 * 0.5
		counterSecond += 50
	}

	for n.ProteinsNorm-forLunchProtein > 10 {
		log.Println(n.ProteinsNorm)
		forLunchProtein += x1 * 0.01
		forLunchFat += x2 * 0.01
		forLunchCarb += x3 * 0.01
		counterThird += 1
	}

	for n.CarbsNorm-forLunchCarb > 5 {
		forLunchProtein += y1 * 0.01
		forLunchFat += y2 * 0.01
		forLunchCarb += y3 * 0.01
		counterFourth += 1
	}

	for n.FatsNorm-forLunchFat > 0 {
		forLunchProtein += z1 * 0.01
		forLunchFat += z2 * 0.01
		forLunchCarb += z3 * 0.01
		counterFifth += 1
	}

	firstProductGram := counterThird
	secondProductGram := counterFourth
	thirdsProductGram := counterFifth
	fourthProductGram := counterFifth
	fifthProductGram := counterSecond
	sixthProductGram := counterFirst

	var foodsData []map[string]float64

	//log.Println(firstProductGram, secondProductGram, thirdsProductGram, fourthProductGram)
	log.Println(xName, yName, zName, aName)
	log.Println(n.ProteinsNorm*0.35, n.FatsNorm*0.2, n.CarbsNorm*0.25, counterFirst, counterSecond)

	n.ProteinsNorm -= forLunchProtein
	n.FatsNorm -= forLunchFat
	n.CarbsNorm -= forLunchCarb

	//var ResultProtein = n.ProteinsNorm
	//var ResultFat = n.FatsNorm
	//var ResultCarb = n.CarbsNorm

	log.Println(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
	foodData := make(map[string]float64)

	foodData[xName] = firstProductGram
	foodData[yName] = secondProductGram
	foodData[zName] = thirdsProductGram
	foodData[aName] = fourthProductGram
	foodData[bName] = fifthProductGram
	foodData[cName] = sixthProductGram

	foodsData = append(foodsData, foodData)

	return foodsData, nil
}

func CheckResult(p float64, f float64, c float64) bool {
	if p > def.ProteinsNorm*1.1 || p < def.ProteinsNorm*(-1.1) {
		return false
	} else if f > def.FatsNorm*1.05 || f < def.FatsNorm*(-1.05) {
		return false
	} else if c > def.CarbsNorm*1.05 || c > def.CarbsNorm*(-1.05) {
		return false
	} else {
		return true
	}
}

func CalculateDay(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			CalculateDayNotHandler()
			responseData, err := json.Marshal(dayMeal)
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

var dayMeal [][]map[string]float64

// TODO: rename
func CalculateDayNotHandler() [][]map[string]float64 {
	n.ProteinsNorm, n.FatsNorm, n.CarbsNorm = DayNewCalculation()
	def.ProteinsNorm, def.FatsNorm, def.CarbsNorm = DayNewCalculation()
	dayMeal = nil
	log.Println(mealAmount)
	if mealAmount == 3 {
		breakfast, err := getBreakfast()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		dinner, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		lunch, err := getLunch()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		correctCalc := CheckResult(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
		if correctCalc == true {
			dayMeal = append(dayMeal, breakfast, lunch, dinner)
		} else {
			CalculateDayNotHandler()
		}

	}
	if mealAmount == 4 {
		breakfast, err := getBreakfast()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		dinner, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		secondDinner, err := getSecondDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		lunch, err := getLunch()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		correctCalc := CheckResult(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
		if correctCalc == true {
			dayMeal = append(dayMeal, breakfast, lunch, secondDinner, dinner)
		} else {
			CalculateDayNotHandler()
		}
	}
	if mealAmount == 5 {
		breakfast, err := getBreakfast()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		dinner, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		dinnerSecond, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		secondDinner, err := getSecondDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		lunch, err := getLunch()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		correctCalc := CheckResult(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
		if correctCalc == true {
			dayMeal = append(dayMeal, breakfast, lunch, dinnerSecond, secondDinner, dinner)
		} else {
			CalculateDayNotHandler()
		}
	}
	if mealAmount == 6 {
		breakfast, err := getBreakfast()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		dinner, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}
		secondDinner, err := getSecondDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		dinnerSecond, err := getDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		secondDinnerSecond, err := getSecondDinner()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		lunch, err := getLunch()
		if err != nil {
			log.Println("Failed to create meal for day")
			return nil
		}

		correctCalc := CheckResult(n.ProteinsNorm, n.FatsNorm, n.CarbsNorm)
		if correctCalc == true {
			dayMeal = append(dayMeal, breakfast, lunch, dinnerSecond, secondDinnerSecond, secondDinner, dinner)
		} else {
			CalculateDayNotHandler()
		}
	}

	return dayMeal
}

func CalculateWeekNotHandler() [][][]map[string]float64 {
	var weekMeal [][][]map[string]float64
	weekMeal = nil
	//week := [7]string{"Monday","Tuesday","Wednesday","Thursday","Friday","Saturday","Sunday"}
	//weekMap := []map[string]string
	//first := weekMap[0]
	//first["Day"] = first["Monday"]
	//weekMap[0]["Day"] = weekMap[0]["Monday"]
	//weekMap[1]["Day"] = weekMap[1]["Tuesday"]
	//weekMap[2]["Day"] = weekMap[2]["Wednesday"]
	//weekMap[3]["Day"] = weekMap[3]["Thursday"]
	//weekMap[4]["Day"] = weekMap[4]["Friday"]
	//weekMap[5]["Day"] = weekMap[5]["Saturday"]
	//weekMap[6]["Day"] = weekMap[6]["Sunday"]

	for i := 0; i < 7; i++ {
		weekDay := CalculateDayNotHandler()
		weekMeal = append(weekMeal, weekDay)
	}
	return weekMeal
}

func CalculateWeek(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			weekMeal := CalculateWeekNotHandler()
			responseData, err := json.Marshal(weekMeal)
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

func getSecondDinner() ([]map[string]float64, error) {
	getFoodData, getFoodName, err := ConnectionDB.CreateMealForSecondDinner(context.Background())
	if err != nil {
		log.Println("Dont get food data!")
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

	forDinnerProtein := 0.0
	forDinnerFat := 0.0
	forDinnerCarb := 0.0

	var counterFirst float64
	var counterSecond float64

	for n.ProteinsNorm*0.15-forDinnerProtein > 1 {
		forDinnerProtein += x1 * 0.01
		forDinnerFat += x2 * 0.01
		forDinnerCarb += x3 * 0.01
		counterFirst += 1
	}
	for counterSecond < 125 {
		forDinnerProtein += y1*0.02 + z1*0.01
		forDinnerFat += y2*0.02 + z2*0.01
		forDinnerCarb += y3*0.02 + z3*0.01
		counterSecond += 1
	}

	firstProductGram := counterFirst
	secondProductGram := counterSecond * 2
	thirdsProductGram := counterSecond

	var foodsData []map[string]float64

	n.ProteinsNorm -= forDinnerProtein
	n.FatsNorm -= forDinnerFat
	n.CarbsNorm -= forDinnerCarb

	foodData := make(map[string]float64)

	foodData[xName] = firstProductGram
	foodData[yName] = secondProductGram
	foodData[zName] = thirdsProductGram

	foodsData = append(foodsData, foodData)

	return foodsData, nil
}
