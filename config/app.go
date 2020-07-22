package config

import (
	"fmt"
	"go/build"
	"os"

	"github.com/spf13/viper"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

var (
	// AppPath application path
	AppPath string
)

func init() {
	// set config based on env
	loadEnvVars()

	OpenDbPool()
}

func loadEnvVars() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Bind OS environment variable
	viper.SetEnvPrefix("goteach")
	viper.BindEnv("env")
	viper.BindEnv("app_path") // bind GOTEACH_APP_PATH varibale

	if viper.Get("env") == "development" {
		viper.SetConfigName("dev")
		dir, _ := os.Getwd()
		AppPath = dir
	} else if viper.Get("env") == "testing" {
		viper.SetConfigName("testing")
		AppPath = viper.GetString("app_path")
	} else {
		viper.SetConfigName("config")
		dir, _ := os.Getwd()
		AppPath = dir
	}

	viper.SetConfigType("json")
	viper.AddConfigPath(AppPath + "/cfg")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func loadSQLite() {
	db, _ := sql.Open("sqlite3", AppPath+"/"+viper.GetString("database.sqlite.db_name")+".db")
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS queue (id INTEGER PRIMARY KEY, data TEXT)")
	stmt.Exec()
}
