package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	mgl "github.com/go-gl/mathgl/mgl32"

	"github.com/nicholasblaskey/Conways-Game-Of-Life/glfwBoilerplate"
	"github.com/nicholasblaskey/go-learn-opengl/includes/shader"
)

func init() {
	runtime.LockOSThread()
}

func randFloat(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
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
	n := 10000
	window := glfwBoilerplate.InitGLFW("", imageSize, imageSize, false)
	defer glfw.Terminate()

	rand.Seed(time.Now().Unix())
	triangleVAO, triangleVBO := makeBuffers([]float32{
		0.0, 0.0,
		0.5, 0.5,
		-0.5, 0.5,
	})
	// Divide 0.5 by two to ensure both shapes have same starting area
	squareVAO, squareVBO := makeBuffers([]float32{
		-0.5 / 2.0, -0.5 / 2.0,
		0.5 / 2.0, -0.5 / 2.0,
		-0.5 / 2.0, 0.5 / 2.0,
		0.5 / 2.0, 0.5 / 2.0,
	})
	defer gl.DeleteVertexArrays(1, &triangleVAO)
	defer gl.DeleteVertexArrays(1, &triangleVBO)
	defer gl.DeleteVertexArrays(1, &squareVAO)
	defer gl.DeleteVertexArrays(1, &squareVBO)

	generateTriangle := true
	ourShader := shader.MakeShaders("genImages.vs", "genImages.fs")

	start := time.Now()
	i := 0
	for !window.ShouldClose() && i < n {
		// Reset
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if i%100 == 0 {
			fmt.Println(i)
		}

		// Apply a random rotation / slide / scale
		transform := mgl.Translate3D(
			randFloat(-0.45, 0.45), randFloat(-0.45, 0.45), 0)
		transform = transform.Mul4(
			mgl.HomogRotate3D(randFloat(-9, 9), mgl.Vec3{0.0, 0.0, 1.0}))
		transform = transform.Mul4(
			mgl.Scale3D(randFloat(0.5, 1.5), randFloat(0.5, 1.5), 0))
		ourShader.SetMat4("transform", transform)

		// Render to screen
		ourShader.Use()
		fileName := fmt.Sprintf("%d", i)
		if generateTriangle {
			gl.BindVertexArray(triangleVAO)
			gl.DrawArrays(gl.TRIANGLES, 0, 6)
			fileName = "./triangle/" + fileName
		} else {
			gl.BindVertexArray(squareVAO)
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 8)
			fileName = "./square/" + fileName
		}

		time.Sleep(time.Millisecond * 0)

		// Save screen
		screenToFile(imageSize, fileName)

		generateTriangle = !generateTriangle
		i += 1
		window.SwapBuffers()
		glfw.PollEvents()
	}

	fmt.Println(time.Now().Sub(start))
}
