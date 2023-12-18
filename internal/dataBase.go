package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Auth interface {
	SetAuthData(ctx context.Context, login string, pass []byte) error
	GetAuthData(ctx context.Context, login string, pass []byte) (string, error)
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
		fmt.Println("We have a error!", err)
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

func (d *DBConnect) SetAuthData(ctx context.Context, login string, pass []byte) error {
	err := d.pool.QueryRow(ctx, `
		INSERT INTO authuser (login, password)
		VALUES ($1, $2)
		ON CONFLICT (login)
		DO UPDATE SET id = 1 
	`, login, pass)

	if err != nil {
		log.Println("Have a problem :", err)
		return nil
	}
	log.Println("User was created")
	return nil
}

func (d *DBConnect) GetAuthData(ctx context.Context, login string, pass []byte) ([]byte, error) {
	var userData []byte
	err := d.pool.QueryRow(ctx, "SELECT password FROM AuthUser WHERE login = $1 AND password = $2", login, pass).Scan(&userData)
	if err != nil {
		log.Println("User not exist or password was entered incorrect!")
		return nil, err
	}
	return userData, err
}

func (s *DBConnect) CreateTable() error {
	_, err := s.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS authuser (
            id SERIAL PRIMARY KEY,
            login VARCHAR UNIQUE NOT NULL,
            password BYTEA NOT NULL                               
        );
    `)
	if err != nil {
		return err
	}
	return nil
}
