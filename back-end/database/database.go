package database

import (
	"abs/util"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Init() {
	host := util.GodotEnv("DB_HOST")
	auth := util.GodotEnv("DB_AUTH")
	dbName := util.GodotEnv("DB_NAME")

	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI("mongodb+srv://"+auth+"@"+host+"/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
}
