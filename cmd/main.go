package main

import (
	"github.com/AnuragProg/compose-select/internal/ui"
)

func main(){
	screen := ui.NewUI()
	if err := screen.Run(); err != nil{
		panic(err.Error())
	}
}
