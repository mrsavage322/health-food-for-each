package internal

type DataBase interface {
	SetURL
	GetURL
	SaveToFile() error
}
