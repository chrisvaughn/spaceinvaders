package spaceinvaders

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	minEnemyShotInterval = time.Second * 3
	enemyVelocity        = 1
)

var (
	globalEnemyImage *ebiten.Image
)

func init() {
	var err error
	globalEnemyImage, _, err = ebitenutil.NewImageFromFile("assets/alien.png")
	if err != nil {
		log.Fatal(err)
	}

}

type Enemy struct {
	x, y           float64
	lastShotFired  time.Time
	enemyDirection float64
	image          *ebiten.Image
}

func NewEnemy(x, y float64) *Enemy {
	return &Enemy{
		x:              x,
		y:              y,
		enemyDirection: 1,
		lastShotFired:  time.Now(),
		image:          globalEnemyImage,
	}
}

func (e *Enemy) Update() {
	if e.x+float64(e.image.Bounds().Dx()) > ScreenWidth {
		e.enemyDirection = -1
	} else if e.x < 0 {
		e.enemyDirection = 1
	}
	e.x += enemyVelocity * e.enemyDirection
}

func (e *Enemy) FireProjectile() *Projectile {
	if time.Since(e.lastShotFired) < minEnemyShotInterval {
		return nil
	}
	e.lastShotFired = time.Now()
	return NewProjectile(
		e.x+float64(e.image.Bounds().Dx()/2),
		e.y+float64(e.image.Bounds().Dy()),
		projectileTypeEnemy,
	)
}

func (e *Enemy) OnScreenRect() image.Rectangle {
	bounds := e.image.Bounds()
	return image.Rectangle{
		Min: image.Pt(int(e.x), int(e.y)),
		Max: image.Pt(int(e.x)+bounds.Dx(), int(e.y)+bounds.Dy()),
	}
}

func (e *Enemy) IsHit(projectileRect image.Rectangle) bool {
	return projectileRect.In(e.OnScreenRect())
}
