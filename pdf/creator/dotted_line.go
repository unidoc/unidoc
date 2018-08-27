package creator

import (
	"math"

	"github.com/unidoc/unidoc/pdf/model"
	"fmt"
	"github.com/unidoc/unidoc/pdf/contentstream/draw"
)

// Line defines a line between point 1 (X1,Y1) and point 2 (X2,Y2).  The line ending styles can be none (regular line),
// or arrows at either end.  The line also has a specified width, color and opacity.
// Implements the Drawable interface and can be drawn on PDF using the Creator.
type DottedLine struct {
	x1        float64
	y1        float64
	x2        float64
	y2        float64
	lineColor *model.PdfColorDeviceRGB
	lineWidth float64
}

// NewLine creates a new Line with default parameters between (x1,y1) to (x2,y2).
func NewDottedLine(x1, y1, x2, y2 float64) *DottedLine {
	l := &DottedLine{}

	l.x1 = x1
	l.y1 = y1
	l.x2 = x2
	l.y2 = y2

	l.lineColor = model.NewPdfColorDeviceRGB(0, 0, 0)
	l.lineWidth = 1.0

	return l
}

// GetCoords returns the (x1, y1), (x2, y2) points defining the Line.
func (l *DottedLine) GetCoords() (float64, float64, float64, float64) {
	return l.x1, l.y1, l.x2, l.y2
}

// SetLineWidth sets the line width.
func (l *DottedLine) SetLineWidth(lw float64) {
	l.lineWidth = lw
}

// SetColor sets the line color.
// Use ColorRGBFromHex, ColorRGBFrom8bit or ColorRGBFromArithmetic to make the color object.
func (l *DottedLine) SetColor(col Color) {
	l.lineColor = model.NewPdfColorDeviceRGB(col.ToRGB())
}

// Length calculates and returns the line length.
func (l *DottedLine) Length() float64 {
	return math.Sqrt(math.Pow(l.x2-l.x1, 2.0) + math.Pow(l.y2-l.y1, 2.0))
}

// GeneratePageBlocks draws the line on a new block representing the page. Implements the Drawable interface.
func (l *DottedLine) GeneratePageBlocks(ctx DrawContext) ([]*Block, DrawContext, error) {
	block := NewBlock(ctx.PageWidth, ctx.PageHeight)

	distanceY := l.y2 - l.y1
	distance := math.Sqrt(math.Pow(l.x2-l.x1, 2) + math.Pow(l.y2-l.y1, 2))

	isHorizontal := true
	if distanceY != 0 {
		isHorizontal = false
	}

	fmt.Println("isHorizontal : ", isHorizontal)
	fmt.Println("distance : ", distance)
	fmt.Println("width : ", l.lineWidth)

	for i := 0.0; i <= distance; i += 2 {
		drawline := draw.Line{
			LineWidth: l.lineWidth,
			Opacity:   1.0,
			LineColor: l.lineColor,
		}

		if isHorizontal {
			drawline.X1 = l.x1 + i
			drawline.Y1 = ctx.PageHeight - l.y1
			drawline.X2 = l.x1 + i + 1
			drawline.Y2 = ctx.PageHeight - l.y2
		} else {
			drawline.X1 = l.x1
			drawline.Y1 = (ctx.PageHeight - l.y1) - i
			drawline.X2 = l.x1
			drawline.Y2 = (ctx.PageHeight - l.y1) - (i - 1)
		}

		contents, _, err := drawline.Draw("")
		if err != nil {
			return nil, ctx, err
		}

		err = block.addContentsByString(string(contents))
		if err != nil {
			return nil, ctx, err
		}
	}

	return []*Block{block}, ctx, nil
}
