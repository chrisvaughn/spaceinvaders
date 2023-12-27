package spaceinvaders

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1024
	ScreenHeight = 768
)

var (
	background        *ebiten.Image
	enemies           []*Enemy
	player            *Player
	playerProjectiles []*Projectile
	enemyProjectiles  []*Projectile
	input             *Input
	globalHeartImage  *ebiten.Image
)

type Game struct{}

func NewGame() (*Game, error) {
	g := &Game{}
	err := g.init()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) init() error {
	var err error
	background, _, err = ebitenutil.NewImageFromFile("assets/space.png")
	if err != nil {
		return err
	}
	globalHeartImage, _, err = ebitenutil.NewImageFromFile("assets/heart.png")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		enemies = append(enemies, NewEnemy(float64(10+(i*100)), 20))
	}
	player = NewPlayer(ScreenWidth/2, ScreenHeight-120)
	input = NewInput()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	for _, enemy := range enemies {
		enemy.Update()
		if rand.Intn(100) < 2 {
			if projectile := enemy.FireProjectile(); projectile != nil {
				enemyProjectiles = append(enemyProjectiles, projectile)
			}
		}
	}

	switch input.MoveShip() {
	case MoveShipLeft:
		player.x -= playerVelocity
		if player.x < 0 {
			player.x = 0
		}
	case MoveShipRight:
		player.x += playerVelocity
		if player.x > float64(ScreenWidth-player.image.Bounds().Dx()) {
			player.x = float64(ScreenWidth - player.image.Bounds().Dx())
		}
	case NoMovement:
		break
	}

	if input.Fire() {
		if projectile := player.FireProjectile(); projectile != nil {
			playerProjectiles = append(playerProjectiles, projectile)
		}
	}

	playerProjectiles = updateProjectiles(playerProjectiles)
	enemyProjectiles = updateProjectiles(enemyProjectiles)

	for i, playerProjectile := range playerProjectiles {
		for j, enemy := range enemies {
			if enemy.IsHit(playerProjectile.OnScreenRect()) {
				playerProjectiles = append(playerProjectiles[:i], playerProjectiles[i+1:]...)
				enemies = append(enemies[:j], enemies[j+1:]...)
				break
			}
		}
	}

	for i, enemyProjectile := range enemyProjectiles {
		if player.IsHit(enemyProjectile.OnScreenRect()) {
			enemyProjectiles = append(enemyProjectiles[:i], enemyProjectiles[i+1:]...)
			player.lives -= 1
			break
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	// draw lives remaining
	for i := 0; i < player.lives; i++ {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(10+(i*40)), ScreenHeight-30)
		screen.DrawImage(globalHeartImage, op)
	}

	for _, enemy := range enemies {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(enemy.x, enemy.y)
		screen.DrawImage(enemy.image, op)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(player.x, player.y)
	screen.DrawImage(player.image, op)

	for _, projectile := range append(playerProjectiles, enemyProjectiles...) {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(projectile.x, projectile.y)
		screen.DrawImage(projectile.image, op)
	}
}

func updateProjectiles(projectiles []*Projectile) []*Projectile {
	offScreenIndex := -1
	for i := range projectiles {
		if projectiles[i].Update() {
			offScreenIndex = i
		}
	}
	if offScreenIndex > -1 {
		projectiles = projectiles[offScreenIndex+1:]
	}
	return projectiles
}
