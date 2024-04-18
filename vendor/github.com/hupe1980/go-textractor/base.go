package textractor

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// base represents the base information shared among different types of blocks.
type base struct {
	id          string          // Identifier for the block
	confidence  float64         // Confidence for the block
	blockType   types.BlockType // Type of the block
	boundingBox *BoundingBox    // Bounding box information
	polygon     Polygon         // Polygon information
	page        *Page           // Page information
	raw         types.Block     // Raw block data
}

// newBase creates a new base instance from the provided Textract block and page information.
func newBase(b types.Block, p *Page) base {
	polygon := make(Polygon, len(b.Geometry.Polygon))
	for i, p := range b.Geometry.Polygon {
		polygon[i] = &Point{
			x: float64(p.X),
			y: float64(p.Y),
		}
	}

	return base{
		id:         aws.ToString(b.Id),
		confidence: float64(aws.ToFloat32(b.Confidence)),
		blockType:  b.BlockType,
		boundingBox: &BoundingBox{
			height: float64(b.Geometry.BoundingBox.Height),
			left:   float64(b.Geometry.BoundingBox.Left),
			top:    float64(b.Geometry.BoundingBox.Top),
			width:  float64(b.Geometry.BoundingBox.Width),
		},
		polygon: polygon,
		page:    p,
		raw:     b,
	}
}

// ID returns the identifier of the block.
func (b *base) ID() string {
	return b.id
}

// Confidence returns the confidence of the block.
func (b *base) Confidence() float64 {
	return b.confidence
}

// BlockType returns the type of the block.
func (b *base) BlockType() types.BlockType {
	return b.blockType
}

// BoundingBox returns the bounding box information of the block.
func (b *base) BoundingBox() *BoundingBox {
	return b.boundingBox
}

// Polygon returns the polygon information of the block.
func (b *base) Polygon() Polygon {
	return b.polygon
}

// PageNumber returns the page number associated with the block.
func (b *base) PageNumber() int {
	return b.page.Number()
}

// Raw returns the raw block data.
func (b *base) Raw() types.Block {
	return b.raw
}

// OCRConfidence represents the confidence scores (mean, max, min) from OCR processing.
type OCRConfidence struct {
	mean, max, min float64
}

// Mean returns the mean (average) confidence score.
func (ocr *OCRConfidence) Mean() float64 {
	return ocr.mean
}

// Max returns the maximum confidence score.
func (ocr *OCRConfidence) Max() float64 {
	return ocr.max
}

// Min returns the minimum confidence score.
func (ocr *OCRConfidence) Min() float64 {
	return ocr.min
}
