package spaceinvaders

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	playerVelocity        = 5
	minPlayerShotInterval = time.Millisecond * 500
	playerWidth           = 80
)

var (
	globalPlayerImage *ebiten.Image
)

func init() {
	var err error
	globalPlayerImage, _, err = ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Player struct {
	x, y          float64
	lastShotFired time.Time
	lives         int
	image         *ebiten.Image
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		x:             x,
		y:             y,
		lastShotFired: time.Now(),
		lives:         3,
		image:         globalPlayerImage,
	}
}

func (p *Player) FireProjectile() *Projectile {
	if time.Since(p.lastShotFired) < minPlayerShotInterval {
		return nil
	}
	p.lastShotFired = time.Now()
	return NewProjectile(
		p.x+float64(playerWidth/2),
		p.y,
		projectileTypePlayer,
	)
}

func (p *Player) OnScreenRect() image.Rectangle {
	bounds := p.image.Bounds()
	return image.Rectangle{
		Min: image.Pt(int(p.x), int(p.y)),
		Max: image.Pt(int(p.x)+bounds.Dx(), int(p.y)+bounds.Dy()),
	}
}

func (p *Player) IsHit(projectileRect image.Rectangle) bool {
	return projectileRect.In(p.OnScreenRect())
}
