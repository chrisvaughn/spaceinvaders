package spaceinvaders

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Mode int

const (
	ScreenWidth       = 1024
	ScreenHeight      = 768
	ModeTitle    Mode = iota
	ModeGame
	ModeGameOver
)

var (
	psNormalFont font.Face
)

func init() {
	ttfBytes, err := os.ReadFile("assets/pressstart2P.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(ttfBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	psNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	mode              Mode
	background        *ebiten.Image
	enemies           [][]*Enemy
	player            *Player
	playerProjectiles []*Projectile
	enemyProjectiles  []*Projectile
	input             *Input
	heartImage        *ebiten.Image
}

func NewGame() (*Game, error) {
	g := &Game{
		mode: ModeTitle,
	}
	g.Reset()

	var err error
	g.background, _, err = ebitenutil.NewImageFromFile("assets/space.png")
	if err != nil {
		return nil, err
	}
	g.heartImage, _, err = ebitenutil.NewImageFromFile("assets/heart.png")
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) Reset() {
	const startRows = 5
	g.player = NewPlayer(ScreenWidth/2-playerWidth/2, ScreenHeight-120)
	g.enemies = make([][]*Enemy, startRows)
	for j := 0; j < startRows; j++ {
		g.enemies[j] = make([]*Enemy, 0)
		for i := 0; i < 10; i++ {
			g.enemies[j] = append(g.enemies[j], NewEnemy(float64(10+(i*100)), float64(20+j*80)))
		}
	}
	g.enemyProjectiles = []*Projectile{}
	g.playerProjectiles = []*Projectile{}
	g.input = NewInput()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) isKeyJustPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if g.isKeyJustPressed() {
			g.mode = ModeGame
		}
	case ModeGame:
		for _, enemyRow := range g.enemies {
			for _, enemy := range enemyRow {
				enemy.Update()
				if rand.Intn(100) < 2 {
					if projectile := enemy.FireProjectile(); projectile != nil {
						g.enemyProjectiles = append(g.enemyProjectiles, projectile)
					}
				}
			}
		}
		switch g.input.MoveShip() {
		case MoveShipLeft:
			g.player.x -= playerVelocity
			if g.player.x < 0 {
				g.player.x = 0
			}
		case MoveShipRight:
			g.player.x += playerVelocity
			if g.player.x > float64(ScreenWidth-g.player.image.Bounds().Dx()) {
				g.player.x = float64(ScreenWidth - g.player.image.Bounds().Dx())
			}
		case NoMovement:
			break
		}

		if g.input.Fire() {
			if projectile := g.player.FireProjectile(); projectile != nil {
				g.playerProjectiles = append(g.playerProjectiles, projectile)
			}
		}

		g.playerProjectiles = updateProjectiles(g.playerProjectiles)
		g.enemyProjectiles = updateProjectiles(g.enemyProjectiles)

		for i, playerProjectile := range g.playerProjectiles {
			for ei, enemyRow := range g.enemies {
				for j, enemy := range enemyRow {
					if enemy.IsHit(playerProjectile.OnScreenRect()) {
						g.playerProjectiles = append(g.playerProjectiles[:i], g.playerProjectiles[i+1:]...)
						g.enemies[ei] = append(g.enemies[ei][:j], g.enemies[ei][j+1:]...)
						break
					}
				}
			}
		}
		for i, enemyProjectile := range g.enemyProjectiles {
			if g.player.IsHit(enemyProjectile.OnScreenRect()) {
				g.enemyProjectiles = append(g.enemyProjectiles[:i], g.enemyProjectiles[i+1:]...)
				g.player.lives -= 1
				break
			}
		}
		if g.player.lives < 1 {
			g.mode = ModeGameOver
		}
	case ModeGameOver:
		if g.isKeyJustPressed() {
			g.mode = ModeTitle
			g.Reset()
		}
	default:
		panic("unhandled default case")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(g.background, op)

	switch g.mode {
	case ModeTitle:
		msg := "SPACE INVADERS"
		x, _ := GetRenderedStringLengthInPixels(msg, psNormalFont)
		text.Draw(screen, msg, psNormalFont, ScreenWidth/2-x/2, 250, color.White)

		msg = "PRESS SPACE TO START"
		x, _ = GetRenderedStringLengthInPixels(msg, psNormalFont)
		text.Draw(screen, msg, psNormalFont, ScreenWidth/2-x/2, 350, color.White)
	case ModeGame:
		// draw lives remaining
		for i := 0; i < g.player.lives; i++ {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(10+(i*40)), ScreenHeight-30)
			screen.DrawImage(g.heartImage, op)
		}

		for _, enemyRow := range g.enemies {
			for _, enemy := range enemyRow {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(enemy.x, enemy.y)
				screen.DrawImage(enemy.image, op)
			}
		}
		for _, projectile := range append(g.playerProjectiles, g.enemyProjectiles...) {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(projectile.x, projectile.y)
			screen.DrawImage(projectile.image, op)
		}
	case ModeGameOver:
		msg := "Game Over"
		x, _ := GetRenderedStringLengthInPixels(msg, psNormalFont)
		text.Draw(screen, msg, psNormalFont, ScreenWidth/2-x/2, 350, color.White)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.x, g.player.y)
	screen.DrawImage(g.player.image, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
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

func GetRenderedStringLengthInPixels(str string, fnt font.Face) (int, int) {
	// Measure the size of the rendered string
	bounds, _ := font.BoundString(fnt, str)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()

	return width, height
}
