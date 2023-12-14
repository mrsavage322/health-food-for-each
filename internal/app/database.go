package app

type Authenticator interface {
	SetAuthData()
	GetAuthData()
}

type FoodAction interface {
	SetFoodData(data FoodData) error
	GetFoodData(data FoodData) error
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

type FoodData struct {
	Kcal     int
	Proteins int
	Fats     int
	Carbs    int
	Feature  string
}

type FoodDataForPerson struct {
	Kcal     int
	Proteins int
	Fats     int
	Carbs    int
	Feature  string
	IsLovely bool
}

func SetFoodData(data FoodData) error {
	//Запрос к базе данных для добавления продукта
	return nil
}

func GetFoodData(data FoodData) error {
	//Запрос к базе данных для получения продукта продукта
	return nil
}
