package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	mocket "github.com/selvatico/go-mocket"
	"github.com/spf13/viper"
)

var (
	DB             *gorm.DB
)

type Database struct {
	Host              string
	User              string
	Password          string
	DBName            string
	Port              int
	ReconnectRetry    int
	ReconnectInterval int64
	DebugMode         bool
}

// LoadDBConfig load database configuration
func LoadDBConfig(name string) Database {
	db := viper.Sub("database." + name)
	conf := Database{
		Host:              db.GetString("host"),
		User:              db.GetString("user"),
		Password:          db.GetString("password"),
		DBName:            db.GetString("db_name"),
		Port:              db.GetInt("port"),
		ReconnectRetry:    db.GetInt("reconnect_retry"),
		ReconnectInterval: db.GetInt64("reconnect_interval"),
		DebugMode:         db.GetBool("debug"),
	}
	return conf
}

// MysqlConnect connect to mysql using config name. return *gorm.DB incstance
func MysqlConnect(configName string) *gorm.DB {
	mysql := LoadDBConfig(configName)
	connection, err := gorm.Open("mysql", mysql.User+":"+mysql.Password+"@tcp("+mysql.Host+":"+strconv.Itoa(mysql.Port)+")/"+mysql.DBName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	if mysql.DebugMode {
		return connection.Debug()
	}

	return connection
}

func MysqlConnectTest(configName string) *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	connection, err := gorm.Open(mocket.DriverName, "connection_string")

	if err != nil {
		//panic(err)
		fmt.Print(err)
	}

	return connection
}

func OpenDbPool() {
	if flag.Lookup("test.v") == nil && !strings.HasSuffix(os.Args[0], ".test") {
		DB = MysqlConnect("mysql")
	} else {
		DB = MysqlConnectTest("mysql")
	}
	pool := viper.Sub("database.mysql.pool")
	DB.DB().SetMaxOpenConns(pool.GetInt("maxOpenConns"))
	DB.DB().SetMaxIdleConns(pool.GetInt("maxIdleConns"))
	DB.DB().SetConnMaxLifetime(pool.GetDuration("maxLifetime") * time.Second)
}