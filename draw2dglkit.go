// Package draw2dglkit offers useful tools for using draw2d with OpenGL.
package draw2dglkit

import (
	"image/color"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
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

	gc.SetStrokeColor(color.RGBA{255, 0, 0, 0xff})
	gc.MoveTo(x, y)
	gc.LineTo(x+1, y+1)
	gc.Stroke()

	green := color.RGBA{0, 255, 0, 0xff}
	gc.SetFillColor(green)
	gc.Fill(poly)

	_, height := offscreen.GetSize()

	gl.ReadBuffer(gl.BACK)
	data := make([]byte, 4)
	// gl.ReadPixels is upside-down
	gl.ReadPixels(1, int32(height)-1, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	window.MakeContextCurrent()
	gc.Restore()
	return color.RGBA{data[0], data[1], data[2], data[3]} == green
}
