package weather

import "github.com/Rocksus/pogo/configs"

var def Repository

type Repository interface {
}

func Init(config configs.WeatherConfig) {
	def = &weatherRepo{
		APIKey: config.APIKey,
	}
}
