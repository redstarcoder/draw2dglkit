// Package draw2dglkit offers useful tools for using draw2d with OpenGL.
package draw2dglkit

import (
	"image/color"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/golang/freetype/raster"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dbase"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dimg"
)

// IsPointInShape uses offscreen as a backbuffer. It checks to see if x, y are inside of poly by drawing a
// small red line at x, y, filling poly as green, and finally running glReadPixels and returning true if x,
// y is green. It returns true if point is on the inside, it may return true or false if the point is on the
// edge of poly.
func IsPointInShape(gc draw2d.GraphicContext, offscreen *glfw.Window, x, y float64, poly *draw2d.Path) bool {
	gc.Save()
	window := glfw.GetCurrentContext()
	offscreen.MakeContextCurrent()

	// 1 added to solved ReadPixels bug regarding y 0
	gc.SetMatrixTransform(draw2d.NewTranslationMatrix(-x+1, -y+1))

	gc.BeginPath()
	gc.SetStrokeColor(color.RGBA{255, 0, 0, 0xff})
	gc.MoveTo(x, y)
	gc.LineTo(x+1, y+1)
	gc.Stroke()

	green := color.RGBA{0, 255, 0, 0xff}
	gc.BeginPath()
	gc.SetFillColor(green)
	gc.Fill(poly)

	_, height := offscreen.GetFramebufferSize()

	gl.ReadBuffer(gl.BACK)
	data := make([]byte, 4)
	// gl.ReadPixels is upside-down
	gl.ReadPixels(1, int32(height)-1, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	window.MakeContextCurrent()
	gc.Restore()
	return color.RGBA{data[0], data[1], data[2], data[3]} == green
}

// FillWithin does a gc.Fill, only painting within the specified boundaries.
// BUG(x) FillWithin uses a *draw2dgl.GraphicContext struct instead of a draw2d.GraphicContext interface.
func FillWithin(gc *draw2dgl.GraphicContext, x, y float64, width, height int, paths ...*draw2d.Path) {
	paths = append(paths, gc.Current.Path)
	rasterizer := &raster.Rasterizer{UseNonZeroWinding: gc.Current.FillRule == draw2d.FillRuleWinding}
	rasterizer.SetBounds(width, height)

	tr := gc.GetMatrixTransform()
	tr.Translate(-x, -y)
	flattener := draw2dbase.Transformer{Tr: tr, Flattener: draw2dimg.FtLineBuilder{Adder: rasterizer}}
	for _, p := range paths {
		draw2dbase.Flatten(p, flattener, tr.GetScale())
	}

	p := draw2dgl.NewPainter()
	gl.PushMatrix()
	gl.Translated(x, y, 0)

	//paint
	p.SetColor(gc.Current.FillColor)
	rasterizer.Rasterize(p)
	rasterizer.Clear()
	p.Flush()
	gc.Current.Path.Clear()

	gl.PopMatrix()
}
