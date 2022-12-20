package main

import (
	"fmt"

	"github.com/red-rocket-software/reminder-go/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.HTTP.IP)

}
