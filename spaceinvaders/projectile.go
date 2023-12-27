package spaceinvaders

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type projectileType int

const (
	projectileTypePlayer projectileType = iota
	projectileTypeEnemy
	projectileTypePlayerVelocity = 5
	projectileTypeEnemyVelocity  = 4
)

var (
	globalProjectileImages map[projectileType]*ebiten.Image
)

func init() {
	globalProjectileImages = make(map[projectileType]*ebiten.Image)
	laserImage, _, err := ebitenutil.NewImageFromFile("assets/laser.png")
	if err != nil {
		log.Fatal(err)
	}
	globalProjectileImages[projectileTypePlayer] = laserImage

	bombImage, _, err := ebitenutil.NewImageFromFile("assets/alien_bomb.png")
	if err != nil {
		log.Fatal(err)
	}
	globalProjectileImages[projectileTypeEnemy] = bombImage
}

type Projectile struct {
	x, y           float64
	projectileType projectileType
	image          *ebiten.Image
}

func NewProjectile(x, y float64, bt projectileType) *Projectile {
	b := &Projectile{
		x:              x,
		y:              y,
		projectileType: bt,
		image:          globalProjectileImages[bt],
	}
	return b
}

func (p *Projectile) Update() bool {
	offScreen := false
	switch p.projectileType {
	case projectileTypePlayer:
		p.y -= projectileTypePlayerVelocity
		// if the full bullet is offscreen
		offScreen = p.y < 0-float64(p.image.Bounds().Dy())
	case projectileTypeEnemy:
		p.y += projectileTypeEnemyVelocity
		offScreen = p.y > ScreenHeight
	default:
		panic("unhandled default case")
	}
	return offScreen
}

func (p *Projectile) OnScreenRect() image.Rectangle {
	bounds := p.image.Bounds()
	return image.Rectangle{
		Min: image.Pt(int(p.x), int(p.y)),
		Max: image.Pt(int(p.x)+bounds.Dx(), int(p.y)+bounds.Dy()),
	}
}
