package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// TODO: Дописать логику расчета нормы БЖУ
func DayCalculation(data map[string]string) (map[string]string, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/settings", nil)
	if err != nil {
		log.Println("Have a problem with GET request", err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return nil, nil
}

func DayNewCalculation() (map[string]string, error) {
	getUserData, err := ConnectionDB.GetUserData(context.Background())
	if err != nil {
		return nil, err
	}

	//set := UserData{
	//	Age:    getUserData["age"],
	//	Height: getUserData["height"],
	//	Weight: getUserData["weight"],
	//	Amount: getUserData["amount"],
	//}
	//fmt.Println(set)
	//response := append(responseUserData, resp)

}
