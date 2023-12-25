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
var proteinForMeal float64
var fatForMeal float64
var carbsForMeal float64

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
	//x := firstProduct["foodname"]
	x1 := firstProduct["proteins"]
	x2 := firstProduct["fats"]
	x3 := firstProduct["carbs"]

	secondProduct := getFoodData[1]
	//y := secondProduct["foodname"]
	y1 := secondProduct["proteins"]
	y2 := secondProduct["fats"]
	y3 := secondProduct["carbs"]

	thirdProduct := getFoodData[2]
	//z := thirdProduct["foodname"]
	z1 := thirdProduct["proteins"]
	z2 := thirdProduct["fats"]
	z3 := thirdProduct["carbs"]

	proteinForMeal = x1 + y1 + z1
	fatForMeal = x2 + y2 + z2
	carbsForMeal = x3 + y3 + z3

	resultProtein := getProteinForMeal(proteinsNorm)
	resultFats := getFatsForMeal(fatsNorm)
	resultCarbs := getCarbsForMeal(carbsNorm)

	log.Println(proteinsNorm)
	log.Println(getFoodData)

	//TODO: Вынести кэфы
	if x1 > y1 && x1 > z1 {
		log.Println(resultProtein, resultProtein*0.7, resultProtein*0.15, resultProtein*0.15)
		log.Println(resultFats, resultFats*0.7, resultFats*0.15, resultFats*0.15)
		log.Println(resultCarbs, resultCarbs*0.7, resultCarbs*0.15, resultCarbs*0.15)

		//proteinForMeal = x1*0.7 + y1*0.15 + z1*0.15
		//fatForMeal = x2*0.7 + y2*0.15 + z2*0.15
		//carbsForMeal = x3*0.7 + y3*0.15 + z3*0.15
	} else if y1 > x1 && y1 > z1 {

		log.Println(resultProtein, resultProtein*0.15, resultProtein*0.7, resultProtein*0.15)
		log.Println(resultFats, resultFats*0.15, resultFats*0.7, resultFats*0.15)
		log.Println(resultCarbs, resultCarbs*0.15, resultCarbs*0.7, resultCarbs*0.15)
		//proteinForMeal = x1*0.15 + y1*0.7 + z1*0.15
		//fatForMeal = x2*0.7 + y2*0.15 + z2*0.15
		//carbsForMeal = x3*0.7 + y3*0.15 + z3*0.7
	} else {

		log.Println(resultProtein, resultProtein*0.15, resultProtein*0.15, resultProtein*0.7)
		log.Println(resultFats, resultFats*0.15, resultFats*0.15, resultFats*0.7)
		log.Println(resultCarbs, resultCarbs*0.15, resultCarbs*0.15, resultCarbs*0.7)
		//proteinForMeal = x1*0.15 + y1*0.15 + z1*0.7
		//fatForMeal = x2*0.15 + y2*0.7 + z2*0.15
		//carbsForMeal = x3*0.15 + y3*0.15 + z3*0.7
	}

	//if proteinsNorm-proteinForMeal > 2 {
	//	//Пересчитывает умнажая на 2
	//	proteinForMeal *= 2
	//} else if proteinsNorm-proteinForMeal < 0 {
	//	proteinForMeal /= 2
	//} else {
	//	return proteinForMeal
	//}

	//log.Println(resultProtein, resultProtein*0.7, resultProtein*0.15, resultProtein*0.15)
	//log.Println(resultProtein, resultProtein*0.7, resultProtein*0.15, resultProtein*0.15)
	//log.Println(resultProtein, resultProtein*0.7, resultProtein*0.15, resultProtein*0.15)

	//log.Println(resultFats, resultFats*0.7, resultFats*0.15, resultFats*0.15)
	//log.Println(resultFats, resultFats*0.7, resultFats*0.15, resultFats*0.15)
	//log.Println(resultFats, resultFats*0.7, resultFats*0.15, resultFats*0.15)

	//log.Println(resultCarbs, resultCarbs*0.7, resultCarbs*0.15, resultCarbs*0.15)
	//log.Println(resultCarbs, resultCarbs*0.7, resultCarbs*0.15, resultCarbs*0.15)
	//log.Println(resultCarbs, resultCarbs*0.7, resultCarbs*0.15, resultCarbs*0.15)

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

func getProteinForMeal(proteinsNorm float64) float64 {
	for proteinsNorm*0.35-proteinForMeal > 2 || proteinsNorm*0.35-proteinForMeal < 0 {
		if proteinsNorm*0.35-proteinForMeal > 2 {
			proteinForMeal *= 1.5
			log.Printf("Proteins for meal: %f\n", proteinForMeal)
			log.Printf("Proteins norm: %f\n", proteinsNorm*0.35)
		} else if proteinsNorm*0.35-proteinForMeal < 0 {
			proteinForMeal *= 0.5
			log.Printf("Proteins for meal: %f\n", proteinForMeal)
			log.Printf("Proteins norm: %f\n", proteinsNorm*0.35)
		} else {
			break
		}
	}
	return proteinForMeal
}

func getFatsForMeal(fatsNorm float64) float64 {
	for fatsNorm*0.35-fatForMeal > 2 || fatsNorm*0.35-fatForMeal < 0 {
		if fatsNorm*0.35-fatForMeal > 2 {
			fatForMeal *= 1.5
			log.Println(fatForMeal)
			log.Printf("Fats for meal: %f\n", fatForMeal)
			log.Printf("Fats norm: %f\n", fatsNorm*0.35)
		} else if fatsNorm*0.35-fatForMeal < 0 {
			fatForMeal *= 0.5
			log.Println(fatForMeal)
			log.Printf("Fats for meal: %f\n", fatForMeal)
			log.Printf("Fats norm: %f\n", fatsNorm*0.35)
		} else {
			break
		}
	}
	return fatForMeal
}

func getCarbsForMeal(carbsNorm float64) float64 {
	for carbsNorm*0.35-carbsForMeal > 2 || carbsNorm*0.35-carbsForMeal < 0 {
		if carbsNorm*0.35-carbsForMeal > 2 {
			carbsForMeal *= 1.5
			log.Printf("Carbs for meal: %f\n", carbsForMeal)
			log.Printf("Carbs norm: %f\n", carbsNorm*0.35)
		} else if carbsNorm*0.35-carbsForMeal < 0 {
			carbsForMeal *= 0.5
			log.Printf("Carbs for meal: %f\n", carbsForMeal)
			log.Printf("Carbs norm: %f\n", carbsNorm*0.35)
		} else {
			break
		}
	}
	return carbsForMeal
}

//Расчет приема пищи на обед, ужин

//Расчет приема пищи на перекус
