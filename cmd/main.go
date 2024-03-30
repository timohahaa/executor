package main

import "github.com/timohahaa/executor/internal/app"

const configFilePath = "./config/config.yaml"

func main() {
	app.Run(configFilePath)
}
