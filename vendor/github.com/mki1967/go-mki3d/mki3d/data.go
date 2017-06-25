package mki3d

/* data structures */

// 3D vector in MKI3D - represents coordinates and RGB colors
type Vector3dType [3]float32

// 3x3 Matrix in MKI3D - represents linear transformation
type Matrix3dType [3]Vector3dType

// Endpoint contains Position, Color, and Set index
type EndpointType struct {
	Position Vector3dType `json:"position"`
	Color    Vector3dType `json:"color"`
	Set      int          `json:"set"`
}

// Segment consists of two endpoints
type SegmentType [2]EndpointType

// SegmentsType is a sequence of segments
type SegmentsType []SegmentType

// Triangle consists of three endpoints
type TriangleType [3]EndpointType

// TrianglesType is a sequence of triangles
type TrianglesType []TriangleType

// Model consists of a sequence of segments and a sequence of triangles
type ModelType struct {
	Segments  SegmentsType  `json:"segments"`
	Triangles TrianglesType `json:"triangles"`
}

// ViewType contains view parameters from MKI3D editor
type ViewType struct {
	FocusPoint     Vector3dType `json:"focusPoint"`
	RotationMatrix Matrix3dType `json:"rotationMatrix"`
	Scale          float32      `json:"scale"`
	ScreenShift    Vector3dType `json:"screenShift"`
	// more fields
}

// Projection contains camera parametres from MKI3D editor
type ProjectionType struct {
	ZNear float32 `json:"zNear"`
	ZFar  float32 `json:"zFar"`
	ZoomY float32 `json:"zoomY"`
}

// CursorType is a state of cursor
type CursorType struct {
	Position Vector3dType  `json:"position"`
	Marker1  *EndpointType `json:"marker1"`
	Marker2  *EndpointType `json:"marker2"`
	Color    Vector3dType  `json:"color"`
	Step     float32       `json:"step"`
}

// Light is described by:
//    Vector - direction of diffuse light, and
//    AbmientFraction - the fraction of light that is ambient
type LightType struct {
	Vector          Vector3dType `json:"vector"`
	AmbientFraction float32      `json:"ambientFraction"`
}

// Set - the current set index
type SetType struct {
	Current int `json:"current"`
}

// The type of MKI3D data in Go.
type Mki3dType struct {
	Model           ModelType      `json:"model"`
	View            ViewType       `json:"view"`
	Projection      ProjectionType `json:"projection`
	BackgroundColor Vector3dType   `json:"backgroundColor"`
	Cursor          CursorType     `json:"cursor"`
	Light           LightType      `json:"light"`
	ClipMaxVector   Vector3dType   `json:"clipMaxVector"`
	ClipMinVector   Vector3dType   `json:"clipMinVector"`
	Set             SetType        `json:"set"`
}
