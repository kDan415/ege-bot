package main

import (
	"ege/app/ege"
	"ege/app/vk"
	"log"
)

func main() {

	log.Print("getting config")
	config, err := GetConfig()
	if err != nil {
		if initError, ok := err.(*InitializationError); ok {
			if initError.CompareCode(MissingConfigFile) {
				err = OverwriteConfig()
				if err != nil {
					log.Fatal(err)
				}
				log.Print("Создан новый конфиг файл. Введите параметры в файл Config.")
				return
			}
		}
		log.Fatal(err)
	}
	log.Print("config received")

	log.Print("initialization grabber")
	grabber, err := ege.New(config.Ege)
	if err != nil {
		log.Print(err)
	}
	log.Print("grabber initialized")

	log.Print("initialization vk")
	bot, err := vk.NewBot(config.VK, grabber)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("vk initialized")

	log.Print("starting longpoll")
	err = bot.RunLongPool()
	if err != nil {
		log.Fatal(err)
	}

}
