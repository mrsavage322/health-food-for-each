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

	if err := dbConnect.CreateAuthTable(); err != nil {
		fmt.Println("We have a error!")
	}
	if err := dbConnect.CreateFoodTable(); err != nil {
		fmt.Println("We have a error!")
	}

	return *dbConnect
}

//	type Authenticator interface {
//		SetAuthData()
//		GetAuthData()
//	}
type FoodAction interface {
	SetFoodData(ctx context.Context, foodname, proteins, fats, carbs int, feature string) error
	//GetFoodData(data FoodData) error
}

type UserAction interface {
	SetUserData(ctx context.Context, foodname, proteins, fats, carbs int, feature string) error
	GetUserData(ctx context.Context) (map[string]string, error)
	//Get(ctx context.Context, age, height, weight int) error
}

func (d *DBConnect) SetAuthData(ctx context.Context, login string, pass []byte) error {
	err := d.pool.QueryRow(ctx, `
		INSERT INTO userdata (login, password)
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
	err := d.pool.QueryRow(ctx, "SELECT password FROM userdata WHERE login = $1 AND password = $2", login, pass).Scan(&userData)
	if err != nil {
		log.Println("User not exist or password was entered incorrect!")
		return nil, err
	}
	return userData, err
}

func (d *DBConnect) SetFoodData(ctx context.Context, foodname string, proteins, fats, carbs int, feature string) error {
	d.pool.QueryRow(ctx, `
		INSERT INTO food (foodname, proteins, fats, carbs, feature)
		VALUES ($1, $2, $3, $4, $5)
	`, foodname, proteins, fats, carbs, feature)

	log.Println("Food added")
	return nil
}

func (s *DBConnect) CreateAuthTable() error {
	_, err := s.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS userdata (
            id SERIAL PRIMARY KEY,
            login VARCHAR UNIQUE NOT NULL,
            password BYTEA NOT NULL,
            age INT,
            height INT,
            weight INT,
            amount INT                       
        );
    `)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBConnect) CreateFoodTable() error {
	_, err := s.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS food (
            id SERIAL PRIMARY KEY,
            foodname VARCHAR UNIQUE,
            proteins INT NOT NULL,
            fats INT NOT NULL,
            carbs INT NOT NULL,
            feature VARCHAR NOT NULL,
            isLoved BOOL,
            login VARCHAR UNIQUE
        );
    `)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBConnect) SetUserData(ctx context.Context, login string, age, height, weight, amount int) error {
	d.pool.QueryRow(ctx, `
		UPDATE userdata SET age=$2, height=$3, weight=$4, amount=$5
		WHERE login=$1
	`, login, age, height, weight, amount)

	log.Println("User data update!")
	return nil
}

func (d *DBConnect) GetUserData(ctx context.Context) (map[string]string, error) {
	row := d.pool.QueryRow(ctx, "SELECT age, height, weight, amount FROM userdata WHERE login = $1", request.Login)
	userConfig := make(map[string]string)
	var age, height, weight, amount string
	err := row.Scan(&age, &height, &weight, &amount)
	if err != nil {
		return nil, err
	}
	userConfig["age"] = age
	userConfig["height"] = height
	userConfig["weight"] = weight
	userConfig["amount"] = amount

	return userConfig, nil
}
