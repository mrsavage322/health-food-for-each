package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth interface {
	SetAuthData()
	GetAuthData()
}

type PersonData struct {
	Username string
	Password string
}

type DBConnect struct {
	pool *pgxpool.Pool
	error
}

// TODO: инит бд один раз
func DataBaseConnection(connection string) DBConnect {
	pool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		fmt.Println("We have a error!")
	}
	//defer pool.Close()

	dbConnect := &DBConnect{
		pool:  pool,
		error: err,
	}

	if err := dbConnect.CreateTable(); err != nil {
		fmt.Println("We have a error!")
	}
	return *dbConnect
}

//type Authenticator interface {
//	SetAuthData()
//	GetAuthData()
//}
//
//type FoodAction interface {
//	SetFoodData(data FoodData) error
//	GetFoodData(data FoodData) error
//}

//type PersonData struct {
//	Username   string
//	Password   string
//	Cookie     string
//	Age        int
//	Height     int
//	Weight     int
//	Sex        int
//	MealNubmer int
//}
//
//type FoodData struct {
//	Kcal     int
//	Proteins int
//	Fats     int
//	Carbs    int
//	Feature  string
//}
//
//type FoodDataForPerson struct {
//	Kcal     int
//	Proteins int
//	Fats     int
//	Carbs    int
//	Feature  string
//	IsLovely bool
//}

func (d *DBConnect) SetAuthData(ctx context.Context, login string, pass []byte) (string, error) {
	err := d.pool.QueryRow(ctx, `
		INSERT INTO kbgu (login, password)
		VALUES ($1, $2)
		ON CONFLICT (login)
		DO UPDATE SET id = 1 
	`, login, pass)

	if err != nil {
		return "User already exist!", nil
	}
	return "User was created!", nil
}

func (d *DBConnect) GetAuthData(ctx context.Context, login string, pass []byte) (string, error) {
	var userData string
	err := d.pool.QueryRow(ctx, "SELECT login, password FROM kbgu.AuthUser WHERE login = $1 AND password = $2", login, pass).Scan(&userData)
	if err != nil {
		return "User not exist or password was entered incorrect!", err
	}
	return userData, err
}

func (s *DBConnect) CreateTable() error {
	_, err := s.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS AuthUser (
            id SERIAL PRIMARY KEY,
            login VARCHAR UNIQUE NOT NULL,
            password VARCHAR NOT NULL                               
        );
    `)
	if err != nil {
		return err
	}
	return nil
}
