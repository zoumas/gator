package main

import (
	"fmt"
	"log"

	"github.com/zoumas/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	err = cfg.SetUser("zoumas")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cfg.DBURL, cfg.CurrentUserName)
}
