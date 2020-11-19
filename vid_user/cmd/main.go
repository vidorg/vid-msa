package main

import (
	"fmt"
	"github.com/douyu/jupiter"
	"log"
	"vid_user/internal/app/engine"
	"vid_user/internal/app/model"
)

func main() {
	eng := engine.NewEngine()
	eng.RegisterHooks(jupiter.StageAfterStop, func() error {
		fmt.Println("exit user service ...")
		return nil
	})

	model.Init()
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}
