package textractor

import (
	"fmt"
	"math"
	"strings"
)

type BoundingBox struct {
	height float64
	left   float64
	top    float64
	width  float64
}

// Bottom returns the bottom coordinate of the bounding box.
func (bb *BoundingBox) Bottom() float64 {
	return bb.Top() + bb.Height()
}

// HorizontalCenter returns the horizontal center coordinate of the bounding box.
func (bb *BoundingBox) HorizontalCenter() float64 {
	return bb.Left() + bb.Width()/2
}

func (bb *BoundingBox) Height() float64 {
	return bb.height
}

func (bb *BoundingBox) Left() float64 {
	return bb.left
}

func (bb *BoundingBox) Top() float64 {
	return bb.top
}

func (bb *BoundingBox) Width() float64 {
	return bb.width
}

// Right returns the right coordinate of the bounding box.
func (bb *BoundingBox) Right() float64 {
	return bb.Left() + bb.Width()
}

// VerticalCenter returns the vertical center coordinate of the bounding box.
func (bb *BoundingBox) VerticalCenter() float64 {
	return bb.Top() + bb.Height()/2
}

// Area calculates and returns the area of the bounding box.
// If either the width or height of the bounding box is less than zero,
// the area is considered zero to prevent negative area values.
func (bb *BoundingBox) Area() float64 {
	if bb.Width() < 0 || bb.Height() < 0 {
		return 0
	}

	return bb.Width() * bb.Height()
}

// Intersection returns a new bounding box that represents the intersection of two bounding boxes.
func (bb *BoundingBox) Intersection(other *BoundingBox) *BoundingBox {
	vtop := math.Max(bb.Top(), other.Top())
	vbottom := math.Min(bb.Bottom(), other.Bottom())
	visect := math.Max(0, vbottom-vtop)
	hleft := math.Max(bb.Left(), other.Left())
	hright := math.Min(bb.Right(), other.Right())
	hisect := math.Max(0, hright-hleft)

	if hisect > 0 && visect > 0 {
		return &BoundingBox{
			height: vbottom - vtop,
			left:   hleft,
			top:    vtop,
			width:  hright - hleft,
		}
	}

	return nil
}

// String returns a string representation of the bounding box.
func (bb *BoundingBox) String() string {
	return fmt.Sprintf("[width: %f, height: %f, left: %f, top: %f]", bb.Width(), bb.Height(), bb.Left(), bb.Top())
}

type BoundingBoxAccessor interface {
	BoundingBox() *BoundingBox
}

// NewEnclosingBoundingBox returns a new bounding box that represents the union of multiple bounding boxes.
func NewEnclosingBoundingBox[T BoundingBoxAccessor](accessors ...T) *BoundingBox {
	if len(accessors) == 0 {
		return nil
	}

	bboxes := make([]*BoundingBox, 0, len(accessors))
	for _, a := range accessors {
		bboxes = append(bboxes, a.BoundingBox())
	}

	left, top, right, bottom := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)

	for _, bb := range bboxes {
		if bb == nil {
			continue
		}

		left = math.Min(left, bb.Left())
		top = math.Min(top, bb.Top())
		right = math.Max(right, bb.Right())
		bottom = math.Max(bottom, bb.Bottom())
	}

	return &BoundingBox{
		height: bottom - top,
		left:   left,
		top:    top,
		width:  right - left,
	}
}

type Polygon []*Point

func (p Polygon) String() string {
	points := make([]string, len(p))
	for i, point := range p {
		points[i] = point.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(points, ", "))
}

// Point represents a 2D point.
type Point struct {
	x, y float64
}

// X returns the X coordinate of the point.
func (p *Point) X() float64 {
	return p.x
}

// Y returns the Y coordinate of the point.
func (p *Point) Y() float64 {
	return p.y
}

// String returns a string representation of the Point, including its X and Y coordinates.
func (p *Point) String() string {
	return fmt.Sprintf("(x: %f, y: %f)", p.x, p.y)
}

// Orientation represents the orientation of a geometric element.
type Orientation struct {
	point0 *Point
	point1 *Point
}

// Radians returns the orientation in radians.
func (o *Orientation) Radians() float64 {
	return math.Atan2(o.point1.Y()-o.point0.Y(), o.point1.X()-o.point0.X())
}

// Degrees returns the orientation in degrees.
func (o *Orientation) Degrees() float64 {
	return (o.Radians() * 180) / math.Pi
}
