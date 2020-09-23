package main

import (
	"fmt"
	"math"
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

func makePolygon(n int, radius float32) []float32 {
	angleIncrement := mgl.DegToRad(360.0 / float32(n))
	verts := []float32{}
	angle := 0.0
	for i := 0; i < n; i++ {
		verts = append(verts,
			radius*float32(math.Cos(angle)),
			radius*float32(math.Sin(angle)))
		angle += float64(angleIncrement)
	}
	return verts
}

// Vertices is formatted like this
// []float32{x0, y0, x1, y1, ...}
func makeBuffers(vertices []float32) uint32 {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4,
		gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	return VAO
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
	n := 1000
	window := glfwBoilerplate.InitGLFW("", imageSize, imageSize, false)
	defer glfw.Terminate()

	radius := float32(0.5)
	shapes := []struct {
		method uint32
		n      int32
		VAO    uint32
	}{
		{gl.TRIANGLE_FAN, 3, makeBuffers(makePolygon(3, radius))},
		{gl.TRIANGLE_FAN, 4, makeBuffers(makePolygon(4, radius))},
		{gl.TRIANGLE_FAN, 5, makeBuffers(makePolygon(5, radius))},
		{gl.TRIANGLE_FAN, 6, makeBuffers(makePolygon(6, radius))},
	}
	ourShader := shader.MakeShaders("genImages.vs", "genImages.fs")

	shapeIndex := 0
	i := 0

	start := time.Now()
	rand.Seed(start.Unix())
	for !window.ShouldClose() && shapeIndex != len(shapes) {
		// Reset
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Apply a random rotation / slide / scale
		transform := mgl.Translate3D(
			randFloat(-0.45, 0.45), randFloat(-0.45, 0.45), 0)
		transform = transform.Mul4(
			mgl.HomogRotate3D(randFloat(-9, 9), mgl.Vec3{0.0, 0.0, 1.0}))
		transform = transform.Mul4(
			mgl.Scale3D(randFloat(0.5, 1.5), randFloat(0.5, 1.5), 0))
		ourShader.SetMat4("transform", transform)

		// Render and save to file
		ourShader.Use()
		gl.BindVertexArray(shapes[shapeIndex].VAO)
		gl.DrawArrays(shapes[shapeIndex].method, 0, shapes[shapeIndex].n*2)
		screenToFile(imageSize, fmt.Sprintf("./%d/%d", shapes[shapeIndex].n, i))

		window.SwapBuffers()
		glfw.PollEvents()

		i += 1
		if i%100 == 0 {
			fmt.Println(i)
		}

		if i == n {
			fmt.Println("Finished class")
			i = 0
			shapeIndex += 1
		}
	}

	fmt.Println(time.Now().Sub(start))
}
