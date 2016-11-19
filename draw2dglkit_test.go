package draw2dglkit

import (
	"image/color"
	"os"
	"runtime"
	"testing"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dgl"
	"github.com/llgcode/draw2d/draw2dkit"
)

var (
	width, height int
	gc            draw2d.GraphicContext
)

func BenchmarkFillWithin(b *testing.B) {
	width, height = 800, 600

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, 0, 0, 150, 150)

	gc.BeginPath()
	gc2 := gc.(*draw2dgl.GraphicContext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FillWithin(gc2, 50, 50, 50, 50, rect)
	}
	b.StopTimer()
	offscreen.Destroy()
}

func getcolor(x, y, height int32) color.RGBA {
	data := make([]byte, 4)
	// gl.ReadPixels is upside-down
	gl.ReadPixels(x, height-y, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	return color.RGBA{data[0], data[1], data[2], data[3]}
}

func TestFillWithin(t *testing.T) {
	width, height = 800, 600

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, 0, 0, 150, 150)

	gc.BeginPath()
	gc2 := gc.(*draw2dgl.GraphicContext)

	_, h := offscreen.GetFramebufferSize()
	height := int32(h)
	gl.ReadBuffer(gl.BACK)

	red := color.RGBA{255, 0, 0, 0xff}
	green := color.RGBA{0, 255, 0, 0xff}

	gc.SetFillColor(red)
	gc.Fill(rect)
	gc.SetFillColor(green)
	FillWithin(gc2, 50, 50, 50, 50, rect)

	if getcolor(1, 1, height) != red {
		t.Error("(1, 1)")
	}
	if getcolor(50, 50, height) != green {
		t.Error("(50, 50)")
	}
	if getcolor(99, 99, height) != green {
		t.Error("(99, 99)")
	}
	if getcolor(149, 149, height) != red {
		t.Error("(149, 149)")
	}

	offscreen.Destroy()
}

func BenchmarkIsPointInShape(b *testing.B) {
	width, height = 800, 600

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, 0, 0, 10, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsPointInShape(gc, offscreen, 9, 9, rect)
	}
	b.StopTimer()
	offscreen.Destroy()
}

func TestIsPointInShape(t *testing.T) {
	width, height = 800, 600

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, -1, -1, 1000, 1000)

	if IsPointInShape(gc, offscreen, -2, -2, rect) {
		t.Error("Point incorrectly found in shape at (-2, -2).")
	}
	if !IsPointInShape(gc, offscreen, -1, -1, rect) {
		t.Error("Point incorrectly not found in shape at (-1, -1).")
	}
	if !IsPointInShape(gc, offscreen, 0, 0, rect) {
		t.Error("Point incorrectly not found in shape at (0, 0).")
	}
	if !IsPointInShape(gc, offscreen, 900, 900, rect) {
		t.Error("Point incorrectly not found in shape at (900, 900).")
	}
	if !IsPointInShape(gc, offscreen, 999, 999, rect) {
		t.Error("Point incorrectly not found in shape at (999, 999).")
	}
	if IsPointInShape(gc, offscreen, 1000, 1000, rect) {
		t.Error("Point incorrectly found in shape at (1000, 1000).")
	}

	offscreen.Destroy()
}

// This test is to ensure IsPointInShape ignores the current transformation matrix
func TestIsPointInShapeTranslate(t *testing.T) {
	width, height = 800, 600

	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, -1, -1, 1000, 1000)

	gc.Translate(1000, 1000)

	if IsPointInShape(gc, offscreen, -2, -2, rect) {
		t.Error("Point incorrectly found in shape at (-2, -2).")
	}
	if !IsPointInShape(gc, offscreen, -1, -1, rect) {
		t.Error("Point incorrectly not found in shape at (-1, -1).")
	}
	if !IsPointInShape(gc, offscreen, 0, 0, rect) {
		t.Error("Point incorrectly not found in shape at (0, 0).")
	}
	if !IsPointInShape(gc, offscreen, 900, 900, rect) {
		t.Error("Point incorrectly not found in shape at (900, 900).")
	}
	if !IsPointInShape(gc, offscreen, 999, 999, rect) {
		t.Error("Point incorrectly not found in shape at (999, 999).")
	}
	if IsPointInShape(gc, offscreen, 1000, 1000, rect) {
		t.Error("Point incorrectly found in shape at (1000, 1000).")
	}

	offscreen.Destroy()
}

var offscreen *glfw.Window // Exists solely so the example can compile

func ExampleIsPointInShape() {
	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, -1, -1, 1000, 1000)

	IsPointInShape(gc, offscreen, -1, -1, rect)     // returns true
	IsPointInShape(gc, offscreen, 1000, 1000, rect) // returns false

}

func TestMain(m *testing.M) {
	r := m.Run()
	glfw.Terminate()
	os.Exit(r)
}

func reshape(w, h int) {
	gl.ClearColor(1, 1, 1, 1)
	/* Establish viewing area to cover entire window. */
	gl.Viewport(0, 0, int32(w), int32(h))
	/* PROJECTION Matrix mode. */
	gl.MatrixMode(gl.PROJECTION)
	/* Reset project matrix. */
	gl.LoadIdentity()
	/* Map abstract coords directly to window coords. */
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	/* Invert Y axis so increasing Y goes down. */
	gl.Scalef(1, -1, 1)
	/* Shift origin up to upper-left corner. */
	gl.Translatef(0, float32(-h), 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	width, height = w, h
	/* Recreate graphic context with new width & height. */
	gc = draw2dgl.NewGraphicContext(width, height)
}

func init() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	err = gl.Init()
	if err != nil {
		panic(err)
	}
}
