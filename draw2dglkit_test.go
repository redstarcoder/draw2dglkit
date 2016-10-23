package draw2dglkit

import (
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

func BenchmarkIsPointInShape(b *testing.B) {
	width, height = 800, 600
	
	glfw.WindowHint(glfw.Visible, glfw.False)
	offscreen, err := glfw.CreateWindow(width, height, "Offscreen (You shouldnt see this)", nil, nil)
	if err != nil {
		panic(err)
	}
	offscreen.MakeContextCurrent()
	reshape(width, height)
	
	glfw.WindowHint(glfw.Visible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "BenchmarkIsPointInShape", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(0)
	reshape(width, height)
	
	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, 0, 0, 10, 10)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsPointInShape(gc, offscreen, 0, 0, rect)
	}
	b.StopTimer()
	window.Destroy()
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
	
	glfw.WindowHint(glfw.Visible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "TestIsPointInShape", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	reshape(width, height)
	glfw.PollEvents()

	rect := &draw2d.Path{}
	draw2dkit.Rectangle(rect, 0, 0, 10, 10)
	
	
	if !IsPointInShape(gc, offscreen, 9, 9, rect) {
		t.Error("Point incorrectly not found in shape.")
	}
	if IsPointInShape(gc, offscreen, 10, 10, rect) {
		t.Error("Point incorrectly found in shape.")
	}
	
	window.Destroy()
	offscreen.Destroy()
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
