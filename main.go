package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/chrisvaughn/spaceinvaders/spaceinvaders"
)

func main() {
	game, err := spaceinvaders.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(spaceinvaders.ScreenWidth, spaceinvaders.ScreenHeight)
	ebiten.SetWindowTitle("Space Invaders")
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
