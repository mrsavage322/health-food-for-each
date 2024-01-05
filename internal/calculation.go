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
	carbsNorm = kcalNorm * 0.47 / 4

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
//func CreateMealForBreakast() error {
//	proteinsNorm, fatsNorm, carbsNorm := DayNewCalculation()
//	fmt.Println(proteinsNorm, fatsNorm, carbsNorm)
//	getFoodData, err := ConnectionDB.CreateMealForLunch(context.Background())
//	if err != nil {
//		log.Println("Dont get fooddata!")
//		return err
//	}
//
//	//xyz := make(map[string]string)
//
//	//TODO: Переписать циклом
//	firstProduct := getFoodData[0]
//	//x := firstProduct["foodname"]
//	x1 := firstProduct["proteins"]
//	x2 := firstProduct["fats"]
//	x3 := firstProduct["carbs"]
//
//	secondProduct := getFoodData[1]
//	//y := secondProduct["foodname"]
//	y1 := secondProduct["proteins"]
//	y2 := secondProduct["fats"]
//	y3 := secondProduct["carbs"]
//
//	thirdProduct := getFoodData[2]
//	//z := thirdProduct["foodname"]
//	z1 := thirdProduct["proteins"]
//	z2 := thirdProduct["fats"]
//	z3 := thirdProduct["carbs"]
//
//	proteinForMeal = x1 + y1 + z1
//	fatForMeal = x2 + y2 + z2
//	carbsForMeal = x3 + y3 + z3
//
//	resultProtein := getProteinForMeal(proteinsNorm)
//	resultFats := getFatsForMeal(fatsNorm)
//	resultCarbs := getCarbsForMeal(carbsNorm)
//
//	log.Println(proteinsNorm)
//	log.Println(getFoodData)
//
//	//TODO: Вынести кэфы
//	if x1 > y1 && x1 > z1 {
//		//Сколько неужно КБЖУ на порцию каждого продукта
//		log.Println("REAL RESULT:")
//		log.Println(resultProtein, resultProtein*0.7*100/x1, resultProtein*0.15*100/y1, resultProtein*0.15*100/z1)
//		log.Println(resultFats, resultFats*0.7*100/x2, resultFats*0.15*100/y2, resultFats*0.15*100/z2)
//		log.Println(resultCarbs, resultCarbs*0.7*100/x3, resultCarbs*0.15*100/y3, resultCarbs*0.15*100/z3)
//
//		//proteinForMeal = x1*0.7 + y1*0.15 + z1*0.15
//		//fatForMeal = x2*0.7 + y2*0.15 + z2*0.15
//		//carbsForMeal = x3*0.7 + y3*0.15 + z3*0.15
//	} else if y1 > x1 && y1 > z1 {
//
//		log.Println(resultProtein, resultProtein*0.15, resultProtein*0.7, resultProtein*0.15)
//		log.Println(resultFats, resultFats*0.15, resultFats*0.7, resultFats*0.15)
//		log.Println(resultCarbs, resultCarbs*0.15, resultCarbs*0.7, resultCarbs*0.15)
//		//proteinForMeal = x1*0.15 + y1*0.7 + z1*0.15
//		//fatForMeal = x2*0.7 + y2*0.15 + z2*0.15
//		//carbsForMeal = x3*0.7 + y3*0.15 + z3*0.7
//	} else {
//
//		log.Println(resultProtein, resultProtein*0.15, resultProtein*0.15, resultProtein*0.7)
//		log.Println(resultFats, resultFats*0.15, resultFats*0.15, resultFats*0.7)
//		log.Println(resultCarbs, resultCarbs*0.15, resultCarbs*0.15, resultCarbs*0.7)
//		//proteinForMeal = x1*0.15 + y1*0.15 + z1*0.7
//		//fatForMeal = x2*0.15 + y2*0.7 + z2*0.15
//		//carbsForMeal = x3*0.15 + y3*0.15 + z3*0.7
//	}

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

//return nil
//

func CalculateBreakfast(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/sign_in", http.StatusSeeOther)
	} else {
		if r.Method == http.MethodGet {
			err := getBreakfast()
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

//func getProteinForMeal(proteinsNorm float64) float64 {
//	for proteinsNorm*0.35-proteinForMeal > 2 || proteinsNorm*0.35-proteinForMeal < 0 {
//		if proteinsNorm*0.35-proteinForMeal > 2 {
//			proteinForMeal *= 1.5
//			log.Printf("Proteins for meal: %f\n", proteinForMeal)
//			log.Printf("Proteins norm: %f\n", proteinsNorm*0.35)
//		} else if proteinsNorm*0.35-proteinForMeal < 0 {
//			proteinForMeal *= 0.5
//			log.Printf("Proteins for meal: %f\n", proteinForMeal)
//			log.Printf("Proteins norm: %f\n", proteinsNorm*0.35)
//		} else {
//			break
//		}
//	}
//	return proteinForMeal
//}

//func getFatsForMeal(fatsNorm float64) float64 {
//	for fatsNorm*0.35-fatForMeal > 2 || fatsNorm*0.35-fatForMeal < 0 {
//		if fatsNorm*0.35-fatForMeal > 2 {
//			fatForMeal *= 1.5
//			log.Println(fatForMeal)
//			log.Printf("Fats for meal: %f\n", fatForMeal)
//			log.Printf("Fats norm: %f\n", fatsNorm*0.35)
//		} else if fatsNorm*0.35-fatForMeal < 0 {
//			fatForMeal *= 0.5
//			log.Println(fatForMeal)
//			log.Printf("Fats for meal: %f\n", fatForMeal)
//			log.Printf("Fats norm: %f\n", fatsNorm*0.35)
//		} else {
//			break
//		}
//	}
//	return fatForMeal
//}

//func getCarbsForMeal(carbsNorm float64) float64 {
//	for carbsNorm*0.35-carbsForMeal > 2 || carbsNorm*0.35-carbsForMeal < 0 {
//		if carbsNorm*0.35-carbsForMeal > 2 {
//			carbsForMeal *= 1.5
//			log.Printf("Carbs for meal: %f\n", carbsForMeal)
//			log.Printf("Carbs norm: %f\n", carbsNorm*0.35)
//		} else if carbsNorm*0.35-carbsForMeal < 0 {
//			carbsForMeal *= 0.5
//			log.Printf("Carbs for meal: %f\n", carbsForMeal)
//			log.Printf("Carbs norm: %f\n", carbsNorm*0.35)
//		} else {
//			break
//		}
//	}
//	return carbsForMeal
//}

//Расчет приема пищи на обед, ужин

//Расчет приема пищи на перекус

//type Product struct {
//	x float64
//	y float64
//	z float64
//}

//func ForOneMealCalc() (gramP, gramF, gramC float64) {
//	proteinsNorm, fatsNorm, carbsNorm := DayNewCalculation()
//
//	getFoodData, err := ConnectionDB.CreateMealForLunch(context.Background())
//	if err != nil {
//		log.Println("Dont get fooddata!")
//		return 0, 0, 0
//	}
//
//	firstMap := getFoodData[0]
//	firstProduct := Product{firstMap["proteins"], firstMap["fats"], firstMap["carbs"]}
//
//	secondMap := getFoodData[1]
//	secondProduct := Product{secondMap["proteins"], secondMap["fats"], secondMap["carbs"]}
//
//	thirdMap := getFoodData[2]
//	thirdProduct := Product{thirdMap["proteins"], thirdMap["fats"], thirdMap["carbs"]}
//
//	pullX := firstProduct.x + secondProduct.x + thirdProduct.x
//	pullY := firstProduct.y + secondProduct.y + thirdProduct.y
//	pullZ := firstProduct.z + secondProduct.z + thirdProduct.z
//
//	pullProduct := Product{pullX, pullY, pullZ}
//
//	for proteinsNorm*0.35-pullProduct.x > 2 || fatsNorm*0.35-pullProduct.y > 2 || carbsNorm*0.35-pullProduct.z > 2 {
//		resultProduct := Product{}
//		MinusProduct.Minus(firstProduct, pullProduct, resultProduct)
//		if pullProduct.x < 2 && pullProduct.y < 2 {
//			if firstProduct.x < 2 || firstProduct.y < 2 {
//				resultProduct.z = pullZ - firstProduct.z
//			}
//			if secondProduct.x < 2 || secondProduct.y < 2 {
//				resultProduct.z = pullZ - secondProduct.z
//			}
//			if thirdProduct.x < 2 || thirdProduct.y < 2 {
//				resultProduct.z = pullZ - thirdProduct.z
//			}
//		}
//
//	}
//
//	//fmt.Println(firstProduct, secondProduct, thirdProduct)
//	return 0, 0, 0
//}
//
//func (p Product) Minus(product, pullProduct Product) Product {
//	outputProduct := p
//	outputProduct.x = pullProduct.x - product.x*0.01
//	outputProduct.y = pullProduct.y - product.y*0.01
//	outputProduct.z = pullProduct.z - product.z*0.01
//	return outputProduct
//}
//
//type MinusProduct interface {
//	Minus(product, pullProduct Product) Product
//}

//type threeMeal interface {
//	getBreakfast() error
//	getLunch()
//	getDinner()
//}

func getBreakfast() error {
	proteinsNorm, fatsNorm, carbsNorm := DayNewCalculation()
	getFoodData, err := ConnectionDB.CreateMealForLunch(context.Background())
	if err != nil {
		log.Println("Dont get fooddata!")
		return err
	}

	firstProduct := getFoodData[0]
	x1 := firstProduct["proteins"]
	x2 := firstProduct["fats"]
	x3 := firstProduct["carbs"]

	secondProduct := getFoodData[1]
	y1 := secondProduct["proteins"]
	y2 := secondProduct["fats"]
	y3 := secondProduct["carbs"]

	thirdProduct := getFoodData[2]
	z1 := thirdProduct["proteins"]
	z2 := thirdProduct["fats"]
	z3 := thirdProduct["carbs"]

	forBreakfastProtein := 0.0
	forBreakfastFat := 0.0
	forBreakfastCarb := 0.0

	var counterFirst float64
	var counterSecond float64

	log.Println(x1, x2, x3)
	log.Println(y1, y2, y3)
	log.Println(z1, z2, z3)

	for (proteinsNorm*0.25-forBreakfastProtein > 1 && fatsNorm*0.2-forBreakfastFat > 1) && carbsNorm*0.35*0.9-forBreakfastCarb > 1 {
		log.Println(forBreakfastCarb)
		forBreakfastProtein += x1*0.02 + y1*0.01
		forBreakfastFat += x2*0.02 + y2*0.01
		forBreakfastCarb += x3*0.02 + y3*0.01
		counterFirst += 1
	}
	for carbsNorm*0.35-forBreakfastCarb > 2 && counterSecond < 200 {
		log.Println(forBreakfastCarb)
		forBreakfastProtein += z1 * 0.01
		forBreakfastFat += z2 * 0.01
		forBreakfastCarb += z3 * 0.01
		counterSecond += 1
	}
	firstProductGram := counterFirst * 100 * 0.02
	secondProductProductGram := counterFirst
	thirdsProductGram := counterSecond

	log.Println(firstProductGram, secondProductProductGram, thirdsProductGram)
	log.Println(proteinsNorm-forBreakfastProtein, fatsNorm-forBreakfastFat, carbsNorm-forBreakfastCarb, counterFirst, counterSecond)

	return nil
}
