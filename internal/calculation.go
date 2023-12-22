package internal

import (
	"context"
	"fmt"
	"strconv"
)

func DayNewCalculation() (proteinsNorm, fatsNorm, carbsNorm float64) {
	getUserData, err := ConnectionDB.GetUserData(context.Background())
	if err != nil {
		return
	}

	age, err := strconv.ParseFloat(getUserData["age"], 3)
	height, err := strconv.ParseFloat(getUserData["height"], 3)
	weight, err := strconv.ParseFloat(getUserData["weight"], 3)
	amount, err := strconv.ParseFloat(getUserData["amount"], 3)

	k := 1.2
	fmt.Println(age, height, weight, amount, k)
	kcalNorm := ((10 * weight) + (6.25 * height) - (5 * age)) * k
	proteinsNorm = kcalNorm * 0.23
	fatsNorm = kcalNorm * 0.3
	carbsNorm = kcalNorm * 0.5

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

func CreatePlanForDay() {
	DayNewCalculation()

}
