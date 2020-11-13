package main

import (
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/aquilax/go-perlin"
	eif "github.com/oskanberg/eif-go/float64"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func plotAnomalousness(data [][]float64) *plot.Plot {
	// set max tree depth to a large number so we get good range of values
	// for plotting
	t := eif.NewForest(data, eif.WithMaxTreeDepth(50), eif.WithTrees(100))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	xys := make(plotter.XYZs, len(data))
	var min float64 = 1
	var max float64 = 0
	var score float64
	for i, v := range data {
		score = t.Score(v)
		max = math.Max(max, score)
		min = math.Min(min, score)
		xys[i] = plotter.XYZ{
			X: v[0],
			Y: v[1],
			Z: score,
		}
	}

	s, err := plotter.NewScatter(xys)
	colors := moreland.SmoothPurpleOrange()
	colors.SetMax(max)
	colors.SetMin(min)

	s.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		_, _, z := xys.XYZ(i)
		c, err := colors.At(z)
		if err != nil {
			log.Panic(err)
		}
		var glyph draw.GlyphDrawer = draw.CircleGlyph{}
		return draw.GlyphStyle{Color: c, Radius: vg.Points(3), Shape: glyph}
	}
	p.Add(s)
	return p
}
func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	plots := []*plot.Plot{}

	// 2d norm
	dp := make([][]float64, 2000)
	for i := range dp {
		dp[i] = []float64{
			rand.NormFloat64() / 16,
			rand.NormFloat64() / 16,
		}
	}
	plots = append(plots, plotAnomalousness(dp))

	// circular
	for i := range dp {
		fi := float64(i) / 400
		dp[i] = []float64{
			math.Sin(fi) + rand.NormFloat64()/16,
			math.Cos(fi) + rand.NormFloat64()/16,
		}
	}
	plots = append(plots, plotAnomalousness(dp))

	// perlin
	p := perlin.NewPerlin(2, 2, 3, rand.Int63n(100))
	for i := range dp {
		fi := float64(i) / 500
		dp[i] = []float64{
			fi,
			10*p.Noise1D(fi) + rand.NormFloat64()/2,
		}
	}
	plots = append(plots, plotAnomalousness(dp))

	t := draw.Tiles{
		Rows: 1,
		Cols: 3,
	}

	img := vgimg.New(vg.Points(1200), vg.Points(400))
	dc := draw.New(img)
	canvases := plot.Align([][]*plot.Plot{plots}, t, dc)
	for i, v := range plots {
		v.Draw(canvases[0][i])
	}

	w, err := os.Create("plots/combined.png")
	if err != nil {
		panic(err)
	}

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}
