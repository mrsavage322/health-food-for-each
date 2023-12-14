package app

type AuthData interface {
	SetAuthData()
	GetAuthData()
}

type FoodData interface {
	SetFoodData()
	GetFoodData()
}

type PersonData struct {
	Username   string
	Password   string
	Cookie     string
	Age        int
	Height     int
	Weight     int
	Sex        int
	MealNubmer int
}
