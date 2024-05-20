package main

import (
	"couplebot/utils"
	"fmt"
)

func main() {
	fmt.Println("Hello, world.")

	telegramBotToken, _ := utils.GetEnvironmentVariable("TELEGRAM_BOT_TOKEN", "")

	fmt.Println(telegramBotToken)
}
