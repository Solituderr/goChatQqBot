package service

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetJsonStr(filename string) map[string]interface{} {
	viper.SetConfigType("json")
	viper.SetConfigName(filename)
	viper.AddConfigPath("../touzi")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}
	m := viper.AllSettings()
	return m
}
