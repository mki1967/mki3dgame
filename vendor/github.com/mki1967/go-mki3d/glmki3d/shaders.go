package glmki3d

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

// Vertex shader for drawing triangles
var vertexShaderT = `
#version 330

/* attributes */
layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec3 color;

/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
uniform vec3 light;
uniform float ambient; 
 
/* output to fragment shader */
out vec4 vColor;

void main() {
    /* compute shaded color */
    vec4 modelNormal=model*vec4(normal, 0.0);
    float shade= abs( dot( modelNormal.xyz, light ) ); 
    vColor= (ambient+(1.0-ambient)*shade)*vec4(color, 1.0);
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

// vertex shader for drawing segments
var vertexShaderS = `
#version 330

/* attributes */
layout (location = 0) in vec3 position;
layout (location = 2) in vec3 color;

/* uniforms */
uniform mat4 model; 
uniform mat4 view;
uniform mat4 projection;
 
/* output to fragment shader */
out vec4 vColor;

void main() {
    /* compute shaded color */
    vColor= vec4(color, 1.0);
    /* compute projected position */
    gl_Position = projection*view*model*vec4(position, 1.0);
}
` + "\x00"

// fragment shader - the same for segments and triangles
var fragmentShader = `
#version 330

/* input from vertex shader */
in vec4 vColor;

/* fragment color output */
out vec4 outputColor;

void main() {
    outputColor = vColor ;
}
` + "\x00"

// from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// from https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

// structure for mki3d shader for drawing triangles
// with references to attributes and uniform locations.
type ShaderTr struct {
	// program Id
	ProgramId uint32
	// locations of attributes
	PositionAttr uint32
	NormalAttr   uint32
	ColorAttr    uint32
	// locations of uniforms ( why int32 instead of uint32 ? )
	ProjectionUni int32
	ViewUni       int32
	ModelUni      int32
	LightUni      int32
	AmbientUni    int32
}

// MakeShaderTr compiles  mki3d shader and
// returns ShaderTr structure with reference to the program and its attributes and uniforms
// or error
func MakeShaderTr() (shaderPtr *ShaderTr, err error) {
	program, err := newProgram(vertexShaderT, fragmentShader)
	if err != nil {
		return nil, err
	}

	var shader ShaderTr

	// set ProgramId
	shader.ProgramId = program

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00")) // test

	// set attributes
	shader.PositionAttr = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	shader.NormalAttr = uint32(gl.GetAttribLocation(program, gl.Str("normal\x00")))
	shader.ColorAttr = uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))

	// set uniforms
	shader.ProjectionUni = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	shader.ViewUni = gl.GetUniformLocation(program, gl.Str("view\x00"))
	shader.ModelUni = gl.GetUniformLocation(program, gl.Str("model\x00"))
	shader.LightUni = gl.GetUniformLocation(program, gl.Str("light\x00"))
	shader.AmbientUni = gl.GetUniformLocation(program, gl.Str("ambient\x00"))
	return &shader, nil
}

// ShaderSeg is a structure for mki3d shader for drawing segments
// with references to attributes and uniform locations.
type ShaderSeg struct {
	// program Id
	ProgramId uint32
	// locations of attributes
	PositionAttr uint32
	ColorAttr    uint32
	// locations of uniforms ( why int32 instead of uint32 ? )
	ProjectionUni int32
	ViewUni       int32
	ModelUni      int32
}

// MakeShaderSeg compiles  mki3d shader and
// returns ShaderSeg structure with reference to the program and its attributes and uniforms
// or error.
func MakeShaderSeg() (shaderPtr *ShaderSeg, err error) {
	program, err := newProgram(vertexShaderS, fragmentShader)
	if err != nil {
		return nil, err
	}

	var shader ShaderSeg

	// set ProgramId
	shader.ProgramId = program

	// set attributes
	shader.PositionAttr = uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	shader.ColorAttr = uint32(gl.GetAttribLocation(program, gl.Str("color\x00")))

	// set uniforms
	shader.ProjectionUni = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	shader.ViewUni = gl.GetUniformLocation(program, gl.Str("view\x00"))
	shader.ModelUni = gl.GetUniformLocation(program, gl.Str("model\x00"))
	return &shader, nil
}

// Both shaders in one struct.
type Shader struct {
	SegPtr *ShaderSeg
	TrPtr  *ShaderTr
}

// Creates both shaders at once and returns in Shader structure
func MakeShader() (shaderPtr *Shader, err error) {
	shaderSeg, err := MakeShaderSeg()
	if err != nil {
		return nil, err
	}

	shaderTr, err := MakeShaderTr()
	if err != nil {
		return nil, err
	}

	return &Shader{SegPtr: shaderSeg, TrPtr: shaderTr}, err

}
