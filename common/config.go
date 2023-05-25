package common

import "os"

type Config struct {
	DriverName  string
	Port        string
	Url         string
	AccessToken string
}

func InitConfig() Config {
	return Config{
		DriverName:  os.Getenv("DRIVER_NAME"),
		Port:        os.Getenv("PORT"),
		Url:         os.Getenv("DATABASE_URL"),
		AccessToken: os.Getenv("ACCESS_TOKEN"),
	}

}
