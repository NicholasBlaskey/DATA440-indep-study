package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"
)

func init() {
	runtime.LockOSThread()
}

func randFloats(n int, min, max float32) []float32 {
	floats := make([]float32, n)
	for i := 0; i < n; i++ {
		floats[i] = min + rand.Float32()*(max-min)
	}
	return floats
}

// Vertices is formatted like this
// []float32{x0, y0, x1, y1, ...}
func makeBuffers(vertices []float32) (uint32, uint32) {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
		gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	return VAO, VBO
}

func screenToFile(imageSize int, fileName string) {
	pixels := make([]byte, imageSize*imageSize*4)
	gl.ReadPixels(0, 0, int32(imageSize), int32(imageSize),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	out := ""
	for i := 0; i < len(pixels); i += 4 {
		if pixels[i] == 0 {
			out += "0"
		} else {
			out += "1"
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(out)
	if err != nil {
		panic(err)
	}
}

func main() {
	imageSize := 32
	n := 200
	window := glfwBoilerplate.InitGLFW("", imageSize, imageSize, false)
	defer glfw.Terminate()

	rand.Seed(time.Now().Unix())
	//vertices := randFloats(8, -1, 1) //[]float32{0.2, 0.2, 0.3, 0.5, 0.9, 0.7, 0.3, 0.2}
	VAO, VBO := makeBuffers([]float32{-1})
	defer gl.DeleteVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VBO)

	generateTriangle := true

	ourShader := shader.MakeShaders("genImages.vs", "genImages.fs")
	i := 0
	for !window.ShouldClose() && i < n {
		// Reset
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update VBO (aka update vertices)
		numVerts := 8
		if generateTriangle {
			numVerts = 6
		}
		vertices := randFloats(numVerts, -1, 1)
		gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
		gl.BufferData(gl.ARRAY_BUFFER, numVerts*4,
			gl.Ptr(vertices), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Render to screen
		fileName := fmt.Sprintf("%d", i)
		ourShader.Use()
		gl.BindVertexArray(VAO)
		if generateTriangle {
			gl.DrawArrays(gl.TRIANGLES, 0, int32(numVerts))
			fileName = "./triangle/" + fileName
		} else {
			fileName = "./quad/" + fileName
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(numVerts))
		}
		gl.BindVertexArray(0)

		time.Sleep(time.Millisecond * 0)

		// Save screen (default framebuffer to file)
		screenToFile(imageSize, fileName)

		generateTriangle = !generateTriangle
		i += 1
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
