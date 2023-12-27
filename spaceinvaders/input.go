package spaceinvaders

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ShipDirection int

const (
	NoMovement ShipDirection = iota
	MoveShipLeft
	MoveShipRight
)

type Input struct {
	leftKeyHeld  bool
	rightKeyHeld bool
}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) checkKeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		i.leftKeyHeld = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		i.leftKeyHeld = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		i.rightKeyHeld = true
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		i.rightKeyHeld = false
	}
}

func (i *Input) MoveShip() ShipDirection {
	i.checkKeys()
	if i.rightKeyHeld && i.leftKeyHeld {
		return NoMovement
	}
	if i.leftKeyHeld {
		return MoveShipLeft
	}
	if i.rightKeyHeld {
		return MoveShipRight
	}
	return NoMovement
}

func (i *Input) Fire() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}
