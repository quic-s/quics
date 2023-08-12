package config

import (
	"os"
)

func GetServerAddress() string {
	// FIXME: 서버의 주소는 Option으로 입력받도록 하는 게 좋다.
	return os.Getenv("BASE_URL")
}

func GetServerPort() string {
	return os.Getenv("PORT")
}
