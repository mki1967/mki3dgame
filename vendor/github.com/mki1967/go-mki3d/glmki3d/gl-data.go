package glmki3d

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/mki1967/go-mki3d/mki3d"
)

// references to the objects defining the shape and parameters of mki3d object

// GLBufTr contains references to GL triangle buffers for triangle shader's input attributes
type GLBufTr struct {
	// buffer objects in GL
	// triangles:
	VertexCount int32 // the last argument for gl.DrawArrays
	PositionBuf uint32
	NormalBuf   uint32
	ColorBuf    uint32
}

// GLBufSeg contains references to GL segment buffers for segment shader's input attributes
type GLBufSeg struct {
	// buffer objects in GL
	// segments:
	VertexCount int32 // the last argument for gl.DrawArrays
	PositionBuf uint32
	ColorBuf    uint32
}

// GLBuf contains references to GL buffers for shaders' input attributes
type GLBuf struct {
	// buffer objects in GL
	// triangles:
	TrPtr *GLBufTr
	// segments:
	SegPtr *GLBufSeg
}

// Delete the buffers in GL, when they are not needed any more
func (glBuf *GLBufTr) Delete() {
	vbo := []uint32{glBuf.PositionBuf, glBuf.NormalBuf, glBuf.ColorBuf}
	gl.DeleteBuffers(3, &vbo[0])
}

// Delete the buffers in GL, when they are not needed any more
func (glBuf *GLBufSeg) Delete() {
	vbo := []uint32{glBuf.PositionBuf, glBuf.ColorBuf}
	gl.DeleteBuffers(3, &vbo[0])
}

// Delete the buffers in GL, when they are not needed any more
func (glBuf *GLBuf) Delete() {
	glBuf.TrPtr.Delete()
	glBuf.SegPtr.Delete()
}

// LoadTriangleBufs loads data from mki3dData to the GL buffers referenced by glBuf (and fills glBuf.NormalBuf with computed normals)
func (glBuf *GLBufTr) LoadTriangleBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.VertexCount = int32(3 * len(mki3dData.Model.Triangles))
	if glBuf.VertexCount == 0 {
		return // do not create empty buffers
	}
	dataPos := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataCol := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	dataNor := make([]float32, 0, 9*len(mki3dData.Model.Triangles)) // each triangle has 3*3 coordinates
	i := 0
	for _, triangle := range mki3dData.Model.Triangles {
		// compute normal
		a := mgl32.Vec3(triangle[0].Position)
		b := mgl32.Vec3(triangle[1].Position)
		c := mgl32.Vec3(triangle[2].Position)
		normal := (b.Sub(a)).Cross(c.Sub(a))
		if normal.Dot(normal) > 0 {
			normal = normal.Normalize()
		}
		// append to buffers
		for j := 0; j < 3; j++ {
			dataPos = append(dataPos, triangle[j].Position[0:3]...)
			dataCol = append(dataCol, triangle[j].Color[0:3]...)
			dataNor = append(dataNor, normal[0:3]...)
			i = i + 3
		}
	}

	/* transfer data to the GL memory */
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.PositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.ColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.NormalBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataNor)*4 /* 4 bytes per flat32 */, gl.Ptr(dataNor), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

// LoadSegmentBufs loads data from mki3dData to the GL buffers referenced by glBuf
func (glBuf *GLBufSeg) LoadSegmentBufs(mki3dData *mki3d.Mki3dType) {
	glBuf.VertexCount = int32(2 * len(mki3dData.Model.Segments))
	if glBuf.VertexCount == 0 {
		return // do not create empty buffers
	}
	dataPos := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	dataCol := make([]float32, 0, 6*len(mki3dData.Model.Segments)) // each segment has 2*3 coordinates
	i := 0
	for _, segment := range mki3dData.Model.Segments {
		for j := 0; j < 2; j++ {
			dataPos = append(dataPos, segment[j].Position[0:3]...)
			dataCol = append(dataCol, segment[j].Color[0:3]...)
			i = i + 2
		}
	}

	/* transfer data to the GL memory */
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.PositionBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataPos)*4 /* 4 bytes per flat32 */, gl.Ptr(dataPos), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, glBuf.ColorBuf)
	gl.BufferData(gl.ARRAY_BUFFER, len(dataCol)*4 /* 4 bytes per flat32 */, gl.Ptr(dataCol), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // unbind
}

// MakeGLBufTr either returns pointer to a new GLBufTr or an error
func MakeGLBufTr(mki3dData *mki3d.Mki3dType) (glBufPtr *GLBufTr, err error) {
	var glBuf GLBufTr
	var vbo [3]uint32 // 5 is the number of buffers
	gl.GenBuffers(3, &vbo[0])
	// TO DO: test for error ...

	// assign buffer ids from vbo array
	glBuf.PositionBuf = vbo[0]
	glBuf.NormalBuf = vbo[1]
	glBuf.ColorBuf = vbo[2]

	// load data from mki3dData
	glBuf.LoadTriangleBufs(mki3dData)

	return &glBuf, nil
}

// MakeGLBufSeg either returns pointer to a new GLBufSeg or an error
func MakeGLBufSeg(mki3dData *mki3d.Mki3dType) (glBufPtr *GLBufSeg, err error) {
	var glBuf GLBufSeg
	var vbo [2]uint32 // 2 is the number of buffers
	gl.GenBuffers(2, &vbo[0])
	// TO DO: test for error ...

	// assign buffer ids from vbo array
	glBuf.PositionBuf = vbo[0]
	glBuf.ColorBuf = vbo[1]

	// load data from mki3dData
	glBuf.LoadSegmentBufs(mki3dData)

	return &glBuf, nil
}

// MakeGLBuf either returns pointer to a new GLBuf or an error
func MakeGLBuf(mki3dData *mki3d.Mki3dType) (glBufPtr *GLBuf, err error) {

	glSegBufPtr, err := MakeGLBufSeg(mki3dData)
	if err != nil {
		return nil, err
	}
	glTrBufPtr, err := MakeGLBufTr(mki3dData)
	if err != nil {
		return nil, err
	}

	glBuf := GLBuf{TrPtr: glTrBufPtr, SegPtr: glSegBufPtr}
	return &glBuf, nil
}
