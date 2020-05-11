package skateboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kevinwylder/sbvision"
)

type Renderer struct {
	input  chan sbvision.Quaternion
	output chan Image
	x      *exec.Cmd

	program      uint32
	materialType int32
	rotation     int32
	translation  int32

	wheelVerts     uint32
	wheelElems     uint32
	wheelCount     int32
	deckVerts      uint32
	deckElems      uint32
	deckCount      int32
	perimeterVerts uint32
	perimeterCount int32
	capVerts       uint32
	capCount       int32
}

func NewRenderer() (*Renderer, error) {
	sb := &Renderer{
		output: make(chan Image),
		input:  make(chan sbvision.Quaternion),
	}
	errors := make(chan error)
	if _, exists := os.LookupEnv("DISPLAY"); !exists {
		// start a virtual display
		sb.x = exec.Command("Xvfb", ":99", "-screen", "0", "1024x768x16")
		sb.x.Start()
		go func() {
			err := sb.x.Wait()
			if !strings.Contains(err.Error(), "killed") {
				fmt.Println(err)
			}
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Println(string(exitErr.Stderr))
			}
		}()
		os.Setenv("DISPLAY", ":99.0")
		time.Sleep(time.Second)
	}

	go func() {
		runtime.LockOSThread()

		if err := glfw.Init(); err != nil {
			errors <- err
			return
		}
		window, err := glfw.CreateWindow(500, 500, "", nil, nil)
		if err != nil {
			errors <- err
			return
		}
		window.MakeContextCurrent()
		if err := gl.Init(); err != nil {
			errors <- err
			return
		}
		if err := sb.setup(); err != nil {
			errors <- err
			return
		}
		close(errors)
		sb.worker()
	}()
	err := <-errors
	if err != nil {
		fmt.Println(err)
		sb.Destroy()
		return nil, err
	}
	return sb, nil
}

func (sb *Renderer) Render(rotation sbvision.Quaternion) Image {
	sb.input <- rotation
	return <-sb.output
}

func (sb *Renderer) Destroy() {
	close(sb.input)
	glfw.Terminate()
	if sb.x != nil {
		sb.x.Process.Signal(os.Interrupt)
	}
}

func (sb *Renderer) setup() error {

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.DITHER)
	gl.ClearColor(1, 1, 1, 0)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRangef(0.0, 1.0)

	vertex, err := downloadShader("https://skateboardvision.net/skateboard/vertex.glsl", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragment, err := downloadShader("https://skateboardvision.net/skateboard/fragment.glsl", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)
	gl.UseProgram(program)

	sb.program = program
	sb.materialType = gl.GetUniformLocation(program, gl.Str("materialType\x00"))
	sb.rotation = gl.GetUniformLocation(program, gl.Str("rotation\x00"))
	sb.translation = gl.GetUniformLocation(program, gl.Str("translationAdditional\x00"))

	var vertBuffer, elementBuffer uint32
	var size int32
	vertBuffer, elementBuffer, size, err = downloadBuffer("https://skateboardvision.net/skateboard/wheel.json")
	if err != nil {
		return err
	}
	sb.wheelVerts = vertBuffer
	sb.wheelElems = elementBuffer
	sb.wheelCount = size
	vertBuffer, elementBuffer, size, err = downloadBuffer("https://skateboardvision.net/skateboard/deck.json")
	if err != nil {
		return err
	}
	sb.deckVerts = vertBuffer
	sb.deckElems = elementBuffer
	sb.deckCount = size
	vertBuffer, _, size, err = downloadBuffer("https://skateboardvision.net/skateboard/perimeter.json")
	if err != nil {
		return err
	}
	sb.perimeterVerts = vertBuffer
	sb.perimeterCount = size
	vertBuffer, _, size, err = downloadBuffer("https://skateboardvision.net/skateboard/cap.json")
	if err != nil {
		return err
	}
	sb.capVerts = vertBuffer
	sb.capCount = size
	return nil
}

func (sb *Renderer) worker() {
	var wheelX float32 = 0.52
	var wheelY float32 = 0.21
	var wheelZ float32 = -0.12
	var outer float32 = 0.25
	var inner float32 = 0.17

	for {
		rotation, more := <-sb.input
		if !more {
			return
		}

		gl.Viewport(0, 0, 500, 500)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		sb.bindBuffer(sb.deckVerts, sb.deckElems)

		gl.Uniform4f(sb.rotation, float32(rotation[0]), float32(rotation[1]), float32(rotation[2]), float32(rotation[3]))
		gl.Uniform1i(sb.materialType, 0) // DECK_GRAPHIC
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.deckCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.Uniform1i(sb.materialType, 1) // GRIP_TAPE
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.deckCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))

		gl.Uniform4f(sb.rotation, -float32(rotation[3]), float32(rotation[2]), -float32(rotation[1]), float32(rotation[0]))
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.deckCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.Uniform1i(sb.materialType, 0) // DECK_GRAPHIC
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.deckCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))

		sb.bindBuffer(sb.perimeterVerts, 0)
		gl.Uniform1i(sb.materialType, 2) // BOARD_RAIL
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, sb.perimeterCount)
		gl.Uniform4f(sb.rotation, float32(rotation[0]), float32(rotation[1]), float32(rotation[2]), float32(rotation[3]))
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, sb.perimeterCount)

		sb.bindBuffer(sb.wheelVerts, sb.wheelElems)
		gl.Uniform1i(sb.materialType, 3) // WHEEL
		gl.Uniform3f(sb.translation, -wheelX, -wheelY, wheelZ)
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.wheelCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.Uniform3f(sb.translation, -wheelX, wheelY, wheelZ)
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.wheelCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.Uniform3f(sb.translation, wheelX, -wheelY, wheelZ)
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.wheelCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
		gl.Uniform3f(sb.translation, wheelX, wheelY, wheelZ)
		gl.DrawElements(gl.TRIANGLE_STRIP, sb.wheelCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))

		sb.bindBuffer(sb.capVerts, 0)
		gl.Uniform3f(sb.translation, -wheelX, -outer, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, -wheelX, -inner, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, -wheelX, outer, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, -wheelX, inner, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, wheelX, -outer, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, wheelX, -inner, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, wheelX, outer, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)
		gl.Uniform3f(sb.translation, wheelX, inner, wheelZ)
		gl.DrawArrays(gl.TRIANGLE_FAN, sb.capCount/12, sb.capCount/12)

		data := make([]byte, 500*500*4)
		gl.ReadPixels(0, 0, 500, 500, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
		sb.output <- data
	}
}

func (sb *Renderer) bindBuffer(verts, elems uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, verts)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elems)
	pos := uint32(gl.GetAttribLocation(sb.program, gl.Str("vertPos\000")))
	gl.VertexAttribPointer(pos, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(pos)
	norm := uint32(gl.GetAttribLocation(sb.program, gl.Str("vertNorm\000")))
	gl.VertexAttribPointer(norm, 3, gl.FLOAT, true, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(norm)
}

func checkError(call string) error {
	if err := gl.GetError(); err != gl.NO_ERROR {
		switch err {
		case gl.INVALID_OPERATION:
			return fmt.Errorf("%s INVALID_OPERATION", call)
		default:
			return fmt.Errorf("%s error 0x%04X\n", call, err)
		}
	}
	return nil
}

func downloadShader(url string, ty uint32) (uint32, error) {
	source, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer source.Body.Close()
	data, err := ioutil.ReadAll(source.Body)
	if err != nil {
		return 0, err
	}
	data = append(data, 0)
	sourceCode := string(data)
	program, free := gl.Strs(sourceCode)
	size := int32(len(data))
	shader := gl.CreateShader(ty)
	if err := checkError("createShader"); err != nil {
		return 0, err
	}
	gl.ShaderSource(shader, 1, program, &size)
	if err := checkError("shaderSource"); err != nil {
		return 0, err
	}
	gl.CompileShader(shader)
	if err := checkError("compileShader"); err != nil {
		return 0, err
	}
	free()
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success != gl.TRUE {
		var logSize int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logSize)
		log := make([]uint8, logSize+1)
		gl.GetShaderInfoLog(shader, logSize, &logSize, &log[0])
		fmt.Println(string(log))
		return 0, fmt.Errorf("Error compiling shader")
	}
	return shader, nil
}

func downloadBuffer(url string) (uint32, uint32, int32, error) {
	var data struct {
		Points   []float32 `json:"points"`
		Elements []uint16  `json:"elements"`
	}
	res, err := http.Get(url)
	if err != nil {
		return 0, 0, 0, err
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return 0, 0, 0, err
	}
	var vertBuffer, elementBuffer uint32
	var size int32
	gl.GenBuffers(1, &vertBuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(data.Points), gl.Ptr(data.Points), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	size = int32(len(data.Points))

	if len(data.Elements) > 0 {
		gl.GenBuffers(1, &elementBuffer)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBuffer)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 2*len(data.Elements), gl.Ptr(data.Elements), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		size = int32(len(data.Elements))
	}
	return vertBuffer, elementBuffer, size, nil
}
