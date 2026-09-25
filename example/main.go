package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"shadowns/shadow"
)

const (
	screenW = 1024
	screenH = 768
)

type ColoredBlock struct {
	block shadow.Block
	col   color.RGBA
	img   *ebiten.Image
}

type Game struct {
	blocks      []ColoredBlock
	light       *ebiten.Image
	lightX, lightY float64
}

func (g *Game) Update() error {
	speed := 6.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.lightY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.lightY += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.lightX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.lightX += speed
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{235, 230, 220, 255})

	for _, cb := range g.blocks {
		shadow.Draw(screen, cb.block, g.lightX, g.lightY, nil)
	}

	for _, cb := range g.blocks {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cb.block.X, cb.block.Y)
		screen.DrawImage(cb.img, op)
	}

	lop := &ebiten.DrawImageOptions{}
	lop.GeoM.Translate(g.lightX-7, g.lightY-7)
	screen.DrawImage(g.light, lop)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return screenW, screenH
}

func main() {
	raw := []struct {
		x, y, w, h float64
		r, g, b    uint8
	}{
		{200, 250, 80, 50, 200, 70, 70},
		{450, 180, 120, 60, 70, 150, 80},
		{350, 400, 100, 80, 60, 100, 200},
		{700, 300, 90, 55, 210, 160, 50},
		{550, 500, 110, 45, 170, 60, 150},
		{120, 500, 70, 70, 60, 180, 170},
		{760, 500, 85, 60, 180, 100, 50},
	}

	blocks := make([]ColoredBlock, len(raw))
	for i, r := range raw {
		img := ebiten.NewImage(int(r.w+0.5), int(r.h+0.5))
		img.Fill(color.RGBA{r.r, r.g, r.b, 255})
		blocks[i] = ColoredBlock{
			block: shadow.Block{X: r.x, Y: r.y, W: r.w, H: r.h},
			col:   color.RGBA{r.r, r.g, r.b, 255},
			img:   img,
		}
	}

	g := &Game{
		light:  ebiten.NewImage(14, 14),
		lightX: screenW / 2,
		lightY: screenH / 2,
		blocks: blocks,
	}
	g.light.Fill(color.RGBA{255, 240, 160, 255})

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("shadowns — point light")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}