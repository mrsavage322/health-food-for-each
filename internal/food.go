package internal

import (
	"html/template"
	"net/http"
)

func AddFood(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Add food",
	}

	// Проверяем метод запроса
	if r.Method == http.MethodPost {
		// Получаем значения из формы
		proteins := r.FormValue("proteins")
		fats := r.FormValue("fats")
		carbs := r.FormValue("carbs")
		feature := r.FormValue("feature")

		// Ваша логика обработки логина и пароля
		// Здесь вы можете добавить проверки, хеширование пароля и т. д.

		// Пример: просто выводим в консоль
		println("Proteins:", proteins)
		println("Fats:", fats)
		println("Carbs:", carbs)
		println("Feature:", feature)

	}

	// Используем шаблон для отображения страницы
	tmpl, err := template.New("index").Parse(`
			<!DOCTYPE html>
				<html>
				<head>
					<title>{{.Title}}</title>
				</head>
				<body>
					<h1>{{.Title}}</h1>
					<form method="post" action="/">
						<label for="username">Proteins:</label>
						<input type="text" id="proteins" name="proteins" required><br>
						<label for="fats">Fats:</label>
						<input type="text" id="fats" name="fats" required><br>
						<label for="carbs">Carbs:</label>
						<input type="text" id="carbs" name="carbs" required><br>
						<label for="feature">Feature:</label>
						<input type="text" id="feature" name="feature" required><br>


						<input type="submit" value="Submit">
					</form>
				</body>
				</html>
			`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, pageVariables)
}
