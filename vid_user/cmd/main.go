package main

import (
	"fmt"
	"github.com/douyu/jupiter"
	"log"
	"vid_user/internal/app/engine"
	jwt2 "vid_user/internal/app/jwt"
	"vid_user/internal/app/model"
)

func main() {
	eng := engine.NewEngine()
	eng.RegisterHooks(jupiter.StageAfterStop, func() error {
		fmt.Println("exit user service ...")
		return nil
	})

	model.Init()
	jwt2.Init()
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}
