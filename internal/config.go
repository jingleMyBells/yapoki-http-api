package internal

import (
	"encoding/json"
	"log"
	"os"
)


type Config struct {
	Port int `json:"port"`
	DBHost string `json:"db_host"`
	DBPort int `json:"db_port"`
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	DBName string `json:"db_name"`
	DBConnSchema string `json:"db_conn_schema"`
	CookieName string `json:"cookie_name"`
	CookieLength int `json:"cookie_lenght"`
	TestUserName string `json:"test_username"`
	TestUserPass string `json:"test_userpass"`
}

var DevConfig *Config


func NewConfig(filename string) error {

	filebytes, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Создаю конфиг по умолчанию. Ошибка чтения конфига: %v", err)
		DevConfig = &Config{
			Port: 8080,
			DBHost: "localhost",
			DBPort: 5432,
			DBUser: "postgres",
			DBPass: "12345",
			DBName: "postgres",
			DBConnSchema: "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			CookieName: "YapokiAuth",
			CookieLength: 42,
			TestUserName: "admin",
			TestUserPass: "12345",
		}
		return nil
	}

	var config Config
	err = json.Unmarshal(filebytes, &config)
	DevConfig = &config
	return nil
}


func GetAppConfig() Config {
	return *DevConfig
}