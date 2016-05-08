package revel_gorp

import (
	"github.com/revel/revel"
	"github.com/u007/go_config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"strings"
	"net/url"
	"fmt"
	"time"
)

var config_file      = "conf/db.conf"
var config           *go_config.IniConfigLoader
var mode        		 = revel.RunMode
var local_zone, local_offset = GetTimeZone()
var time_zone   		 = local_zone //revel.Config.StringDefault("time_zone", local_zone)
var DBMap 							 *gorp.DbMap

func GetTimeZone() (name string, offset int) {
	return time.Now().In(time.Local).Zone()
}

func DatabaseDriver() (string) {
	return config.String(mode, "driver", "")
}

func DatabaseConnectionString() (string, error) {
	var err error
	// super gotcha!!!!!! if you declared outside variable, using same name
	//	locally with := is considered a different variable
	config, err = go_config.NewConfigLoader("ini", config_file)
	if (err != nil) {
    return "", err
  }
  needed := []string{"driver", "user", "host", "encoding", "db", "pass", "connection_pool"}
	// Debug("user: %s", config.String(mode, "user", ""))
	// Debug("driver: %s", config.String(mode, "driver", ""))
  if !CheckRequired(needed...){
    return "", fmt.Errorf("Required configuration missing: %s in %s", strings.Join(needed, ", "), config_file)
  }
	if time_zone == "" {
		return "", fmt.Errorf("Required time_zone in conf/app.conf")
	}
  port := config.Int(mode, "port", 3306)
	connection_string := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=%s",
    config.String(mode, "user", ""), config.String(mode, "pass", ""),
    config.String(mode, "host", ""), port, config.String(mode, "db", ""),
    config.String(mode, "encoding", ""),
    url.QueryEscape(time_zone))
	return connection_string, nil
}

func InitDatabase() (*gorp.DbMap, error) {
	mode = revel.RunMode
	time_zone = revel.Config.StringDefault("time_zone", local_zone)
	Debug("Mode: %s, Time zone: %s, local timezone: %s, offset: %d", mode, time_zone, local_zone, local_offset)
	connection_string, err  := DatabaseConnectionString()
	if (err != nil) {
		return nil, err
	}
	driver := DatabaseDriver()
	// Debug("Driver: %s, Connection: %s", driver, connection_string)

	if (driver == "mysql") {
		db, err := sql.Open(driver,
			connection_string)

		connection_pool := config.Int(mode, "connection_pool", 5)
		if err == nil {
			db.SetMaxOpenConns(connection_pool)
		}

		DBMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

		// defer db.Close()
		return DBMap, err
	}

  return nil, fmt.Errorf("Unsurppoted Database driver %s", driver)
}
//
// func LogValidationErrors(log_prefix string, valid *validation.Validation) {
// 	if valid.HasErrors() {
//     for _, err := range valid.Errors {
// 			Error("[ %s ]Validation %s: %s", log_prefix, err.Key, err.Message)
//     }
//   }
// }

func CheckRequired(args ...string) bool {
  for _, name := range args {
    if (config.String(mode, name, "") == "") {
			err := fmt.Errorf("[ ERROR ] env: %s, %s required in %s", mode, name, config_file)
			Error(err.Error())
			fmt.Println(err.Error())
      return false
    }
  }
  return true
}

const PREFIX = "[ ORM ] "

func Debug(format string, v... interface{}) {
	revel.INFO.Printf(PREFIX + format, v...)
}
func Warning(format string, v... interface{}) {
	revel.WARN.Printf(PREFIX + format, v...)
}
func Error(format string, v... interface{}) {
	revel.ERROR.Printf(PREFIX + format, v...)
	// debug.PrintStack()
}
