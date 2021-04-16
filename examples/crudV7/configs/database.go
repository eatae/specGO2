package configs

import (
	"fmt"
	"gorm.io/gorm"
)

var DB *gorm.DB

// DBConfig ...
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLmode  string
}

// BuildDbConfig ...
func BuildDbConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     "localhost",
		Port:     "54320",
		User:     "postgres",
		Password: "2222",
		DBname:   "specgo2v7",
		SSLmode:  "disable",
	}
	return &dbConfig
}

// DbUri ...
func DbUri(dbConfig *DBConfig) string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBname, dbConfig.Port, dbConfig.SSLmode)
}
