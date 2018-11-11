package trainsim

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db gorm.DB

type DbConfig struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

func (c DbConfig) GetPostgres() string {
	connectUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", c.Host, c.Port, c.User, c.DbName, c.Password)
	return connectUrl
}

//db.AutoMigrate(
func NewDb() *gorm.DB {
	dbconf := DbConfig{
		Host:   "localhost",
		Port:   "5432",
		DbName: "trains",
	}
	ReadJson("db.json", &dbconf)
	db, err := gorm.Open("postgres", dbconf.GetPostgres())
	if err != nil {
		panic("cannot open connection: " + err.Error())
	}
	db.LogMode(false)
	return db
}
