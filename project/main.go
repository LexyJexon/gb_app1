package main

import (
	"encoding/gob"
	"fmt"
	_ "fmt"
	"github.com/BurntSushi/toml"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"project/models"
	_ "project/routers"
)

type (
	Config struct {
		DB struct {
			DBHOST     string
			PORT       int
			DBNAME     string
			DBUSERNAME string
			DBUSERPWD  string
		}
	}
)

func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	_, err = toml.DecodeFile("conf/config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		config.DB.DBUSERNAME,
		config.DB.DBUSERPWD,
		config.DB.DBHOST,
		config.DB.PORT,
		config.DB.DBNAME)

	err = orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		log.Fatal("Ошибка декодирования конфигурации из файла TOML:", err)
	}

	gob.Register(&models.Users{})
}

func main() {
	beego.Run()
}
