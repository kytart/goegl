package main

import (
	"flag"
	"fmt"
	"github.com/remogatto/application"
	"github.com/remogatto/egl"
	"github.com/remogatto/egl/platform"
	gl "github.com/remogatto/opengles2"
	"log"
	"runtime"
	"time"
)

const FRAMES_PER_SECOND = 24

var (
	verticesArrayBuffer, colorsArrayBuffer uint32
	attrPos, attrColor                     uint32
	currWidth, currHeight                  int

	vertices = [12]float32{
		-1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
	}
	colors = [12]float32{
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
	}
)

// emulatorLoop sends a cmdRenderFrame command to the rendering backend
// (displayLoop) each 1/50 second.
type renderLoop struct {
	ticker           *time.Ticker
	pause, terminate chan int
}

// newRenderLoop returns a new renderLoop instance. It takes the
// number of frame-per-second as argument.
func newRenderLoop(fps int) *renderLoop {
	renderLoop := &renderLoop{
		ticker:    time.NewTicker(time.Duration(1e9 / fps)),
		pause:     make(chan int),
		terminate: make(chan int),
	}
	return renderLoop
}

// Pause returns the pause channel of the loop.
// If a value is sent to this channel, the loop will be paused.
func (l *renderLoop) Pause() chan int {
	return l.pause
}

// Terminate returns the terminate channel of the loop.
// If a value is sent to this channel, the loop will be terminated.
func (l *renderLoop) Terminate() chan int {
	return l.terminate
}

// Run runs renderLoop.
// The loop renders a frame and swaps the buffer for each tick
// received.
func (l *renderLoop) Run() {
	runtime.LockOSThread()
	initialize()
	gl.Viewport(0, 0, INITIAL_WINDOW_WIDTH, INITIAL_WINDOW_HEIGHT)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	initShaders()
	for {
		select {
		case <-l.pause:
			l.ticker.Stop()
			l.pause <- 0
		case <-l.terminate:
			cleanup()
			l.terminate <- 0
		case <-l.ticker.C:
			draw(currWidth, currHeight)
			egl.SwapBuffers(platform.Display, platform.Surface)
		}
	}
}

func check() {
	error := gl.GetError()
	if error != 0 {
		panic(fmt.Sprintf("An error occurred! Code: 0x%x", error))
	}
}

func initShaders() {
	program := Program(FragmentShader(fsh), VertexShader(vsh))
	gl.UseProgram(program)
	attrPos = uint32(gl.GetAttribLocation(program, "pos"))
	attrColor = uint32(gl.GetAttribLocation(program, "color"))
	gl.GenBuffers(1, &verticesArrayBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, verticesArrayBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.SizeiPtr(len(vertices))*4, gl.Void(&vertices[0]), gl.STATIC_DRAW)
	gl.GenBuffers(1, &colorsArrayBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorsArrayBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.SizeiPtr(len(colors))*4, gl.Void(&colors[0]), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(attrPos)
	gl.EnableVertexAttribArray(attrColor)
}

func draw(width, height int) {
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindBuffer(gl.ARRAY_BUFFER, verticesArrayBuffer)
	gl.VertexAttribPointer(attrPos, 4, gl.FLOAT, false, 0, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorsArrayBuffer)
	gl.VertexAttribPointer(attrColor, 4, gl.FLOAT, false, 0, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.Flush()
	gl.Finish()
}

func cleanup() {
	egl.DestroySurface(platform.Display, platform.Surface)
	egl.DestroyContext(platform.Display, platform.Context)
	egl.Terminate(platform.Display)
}

func reshape(width, height int) {
	currWidth, currHeight = width, height
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
}

func printInfo() {
	log.Printf("GL_RENDERER   = %s\n", gl.GetString(gl.RENDERER))
	log.Printf("GL_VERSION    = %s\n", gl.GetString(gl.VERSION))
	log.Printf("GL_VENDOR     = %s\n", gl.GetString(gl.VENDOR))
	log.Printf("GL_EXTENSIONS = %s\n", gl.GetString(gl.EXTENSIONS))
}

func main() {
	info := flag.Bool("info", false, "display OpenGL renderer info")
	flag.Parse()
	if *info {
		printInfo()
	}
	application.Register("render loop", newRenderLoop(FRAMES_PER_SECOND))
	exitCh := make(chan bool, 1)
	application.Run(exitCh)
	<-exitCh
}
