// Package shadow реализует направленные бесконечные тени для 2D сцены
// с точечным источником света. Тени строятся как теневые объёмы (shadow volumes):
// от каждого блока экструдируются вершины, образующие трапециевидный полигон.
//
// Использование:
//
//	import "shadowns/shadow"
//
//	b := shadow.Block{X: 100, Y: 100, W: 80, H: 50}
//	shadow.Draw(screen, b, lightX, lightY, &shadow.Options{Alpha: 0.5})
package shadow

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Block описывает прямоугольный объект, отбрасывающий тень.
type Block struct {
	X, Y, W, H float64
}

// Options управляет параметрами тени.
type Options struct {
	// Color — цвет тени (по умолчанию чёрный).
	Color color.RGBA
	// Alpha — непрозрачность тени [0, 1] (по умолчанию 0.43).
	Alpha float64
	// Extrude — длина экструзии (вытягивания) тени в пикселях (по умолчанию 3000).
	Extrude float64
	// AntiAlias — сглаживание краёв тени (по умолчанию true).
	AntiAlias bool
}

var defaultOpts = Options{
	Color:     color.RGBA{0, 0, 0, 255},
	Alpha:     0.43,
	Extrude:   3000,
	AntiAlias: true,
}

// Draw рисует тень от одного блока от источника света (lightX, lightY).
// opts может быть nil для значений по умолчанию.
func Draw(screen *ebiten.Image, b Block, lightX, lightY float64, opts *Options) {
	if opts == nil {
		opts = &defaultOpts
	}

	verts := [][2]float64{
		{b.X, b.Y},
		{b.X, b.Y + b.H},
		{b.X + b.W, b.Y + b.H},
		{b.X + b.W, b.Y},
	}

	extrude := func(p [2]float64) [2]float64 {
		dx := p[0] - lightX
		dy := p[1] - lightY
		d := math.Hypot(dx, dy)
		if d < 1 {
			d = 1
		}
		return [2]float64{p[0] + dx/d*opts.Extrude, p[1] + dy/d*opts.Extrude}
	}

	edgeFacesLight := make([]bool, 4)
	for i := 0; i < 4; i++ {
		p1 := verts[i]
		p2 := verts[(i+1)%4]
		dx := p2[0] - p1[0]
		dy := p2[1] - p1[1]
		nx := -dy
		ny := dx

		mx := (p1[0] + p2[0]) / 2
		my := (p1[1] + p2[1]) / 2
		ldx := lightX - mx
		ldy := lightY - my

		edgeFacesLight[i] = nx*ldx+ny*ldy > 0
	}

	silIdx := make([]int, 0, 2)
	for i := 0; i < 4; i++ {
		prev := (i + 3) % 4
		if edgeFacesLight[prev] != edgeFacesLight[i] {
			silIdx = append(silIdx, i)
		}
	}
	if len(silIdx) != 2 {
		return
	}

	var startIdx, endIdx int
	if !edgeFacesLight[silIdx[0]] {
		startIdx = silIdx[0]
		endIdx = silIdx[1]
	} else if !edgeFacesLight[(silIdx[0]+3)%4] {
		startIdx = silIdx[1]
		endIdx = silIdx[0]
	} else {
		return
	}

	farVerts := [][2]float64{}
	i := startIdx
	for {
		farVerts = append(farVerts, verts[i])
		if i == endIdx {
			break
		}
		i = (i + 1) % 4
	}

	n := len(farVerts)
	if n < 2 {
		return
	}

	eFirst := extrude(farVerts[0])
	eLast := extrude(farVerts[n-1])

	var p vector.Path
	p.MoveTo(float32(farVerts[0][0]), float32(farVerts[0][1]))
	for k := 1; k < n; k++ {
		p.LineTo(float32(farVerts[k][0]), float32(farVerts[k][1]))
	}
	p.LineTo(float32(eLast[0]), float32(eLast[1]))
	p.LineTo(float32(eFirst[0]), float32(eFirst[1]))
	p.Close()

	var cs ebiten.ColorScale
	cs.Set(
		float32(opts.Color.R)/255,
		float32(opts.Color.G)/255,
		float32(opts.Color.B)/255,
		float32(opts.Alpha),
	)

	vector.FillPath(screen, &p, &vector.FillOptions{}, &vector.DrawPathOptions{
		AntiAlias:  opts.AntiAlias,
		ColorScale: cs,
	})
}