/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package extractor

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/contentstream"
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/internal/transform"
	"github.com/unidoc/unidoc/pdf/model"
)

// ExtractPageImages returns the image contents of the page extractor, including data
// and position, size information for each image.
func (e *Extractor) ExtractPageImages() (*PageImages, error) {
	ctx := &imageExtractContext{}

	err := ctx.extractContentStreamImages(e.contents, e.resources)
	if err != nil {
		return nil, err
	}

	return &PageImages{
		Images: ctx.extractedImages,
	}, nil
}

// PageImages represents extracted images on a PDF page with spatial information:
// display position and size.
type PageImages struct {
	Images []ImageMark
}

// ImageMark represents an image drawn on a page and its position in device coordinates.
// All coordinates are in device coordinates.
type ImageMark struct {
	Image *model.Image

	// Dimensions of the image as displayed in the PDF.
	Width  float64
	Height float64

	// Position of the image in PDF coordinates (lower left corner).
	X float64
	Y float64

	// Angle in degrees, if rotated.
	Angle float64

	CTM transform.Matrix
}

// String returns a string describing `mark`.
func (mark ImageMark) String() string {
	img := mark.Image
	imgStr := fmt.Sprintf("%dx%d cpts=%d bpp=%d",
		img.Width, img.Height, img.ColorComponents, img.BitsPerComponent)
	return fmt.Sprintf("%.1fx%.1f (%.1f,%.1f) ϴ=%.1f img=[%s]",
		mark.Width, mark.Height, mark.X, mark.Y, mark.Angle, imgStr)
}

// Clip returns `mark`.Image clipped to `box`.
// TODO(peterwilliams): Return image in orginal colorspace. The github.com/disintegration/imaging
// library we are using converts all images to image.NRGBA.
// This function can be used to clip extracted images the same way they are clipped in the PDF they
// are extracted from to give the same image the user sees in the enclosing PDF
func (mark ImageMark) Clip(box model.PdfRectangle) (*image.NRGBA, error) {
	inv, hasInverse := mark.CTM.Inverse()
	if !hasInverse {
		return nil, errors.New("CTM has no inverse")
	}
	clp := model.PdfRectangle{}
	clp.Llx, clp.Lly = inv.Transform(box.Llx, box.Lly)
	clp.Urx, clp.Ury = inv.Transform(box.Urx, box.Ury)
	clp.Llx, clp.Lly = maxFloat(0, clp.Llx), maxFloat(0, clp.Lly)
	clp.Urx, clp.Ury = minFloat(1, clp.Urx), minFloat(1, clp.Ury)

	img, err := mark.Image.ToGoImage()
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	w := float64(b.Max.X - b.Min.X)
	h := float64(b.Max.Y - b.Min.Y)

	rect := image.Rectangle{
		Min: image.Point{
			X: round(w * clp.Llx),
			Y: round(h * clp.Lly),
		},
		Max: image.Point{
			X: round(w * clp.Urx),
			Y: round(h * clp.Ury),
		},
	}

	imgRgb := imaging.Crop(img, rect)
	return imgRgb, nil
}

// PageView returns `mark`.Image transformed to appear as it appears the PDF page it was extracted
// from.
//    `bbox` is a clipping rectangle. It should be the clipping path in effect when the image was
//          rendered. TODO(peterwilliams97) support non-rectangular clipping paths.
//    If `doScale` is true the image is scaled as it is on the PDF page. `doScale` will typically
//          only be set false for debugging to check it the scaling is correct.
func (mark ImageMark) PageView(bbox model.PdfRectangle, doScale bool) (*image.NRGBA, error) {
	img, err := mark.Clip(bbox)
	if err != nil {
		return nil, err
	}
	bgColor := color.White
	img = imaging.Rotate(img, -mark.Angle, bgColor)

	if doScale {
		W, H := int(mark.Image.Width), int(mark.Image.Height)
		wf, hf := float64(W), float64(H)
		w, h := mark.Width, mark.Height
		fmt.Printf("W,H = %d,%d (%.2f) w,h=%g,%g (%.2f)\n", W, H, hf/wf, w, h, h/w)
		if w*hf != wf*h {
			if w*hf > wf*h {
				W0 := W
				W = round(hf * (w / h))
				fmt.Printf("W %d->%d (%.2f)\n", W0, W, float64(H)/float64(W))
			} else {
				H0 := H
				H = round(wf * (h / w))
				fmt.Printf("H %d->%d (%.2f)\n", H0, H, float64(H)/float64(W))
			}
			img = imaging.Resize(img, W, H, imaging.CatmullRom)
		}
	}

	return img, nil
}

// round returns `x` rounded the nearest int.
func round(x float64) int {
	return int(math.Round(x))
}

// round64 returns `x` rounded the nearest int64.
func round64(x float64) int64 {
	return int64(math.Round(x))
}

// Provide context for image extraction content stream processing.
type imageExtractContext struct {
	extractedImages []ImageMark
	inlineImages    int
	xObjectImages   int
	xObjectForms    int

	// Cache to avoid processing same image many times.
	cacheXObjectImages map[*core.PdfObjectStream]*cachedImage
}

type cachedImage struct {
	image *model.Image
	cs    model.PdfColorspace
}

func (ctx *imageExtractContext) extractContentStreamImages(contents string,
	resources *model.PdfPageResources) error {
	cstreamParser := contentstream.NewContentStreamParser(contents)
	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	if ctx.cacheXObjectImages == nil {
		ctx.cacheXObjectImages = map[*core.PdfObjectStream]*cachedImage{}
	}

	processor := contentstream.NewContentStreamProcessor(*operations)
	processor.AddHandler(contentstream.HandlerConditionEnumAllOperands, "",
		func(op *contentstream.ContentStreamOperation, gs contentstream.GraphicsState,
			resources *model.PdfPageResources) error {
			return ctx.processOperand(op, gs, resources)
		})

	return processor.Process(resources)
}

// Process individual content stream operands for image extraction.
func (ctx *imageExtractContext) processOperand(op *contentstream.ContentStreamOperation,
	gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
	if op.Operand == "BI" && len(op.Params) == 1 {
		// BI: Inline image.
		iimg, ok := op.Params[0].(*contentstream.ContentStreamInlineImage)
		if !ok {
			return nil
		}

		return ctx.extractInlineImage(iimg, gs, resources)
	} else if op.Operand == "Do" && len(op.Params) == 1 {
		// Do: XObject.
		name, ok := core.GetName(op.Params[0])
		if !ok {
			return errTypeCheck
		}

		_, xtype := resources.GetXObjectByName(*name)
		switch xtype {
		case model.XObjectTypeImage:
			return ctx.extractXObjectImage(name, gs, resources)
		case model.XObjectTypeForm:
			return ctx.extractFormImages(name, gs, resources)
		}
	}
	return nil
}

func (ctx *imageExtractContext) extractInlineImage(iimg *contentstream.ContentStreamInlineImage,
	gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
	img, err := iimg.ToImage(resources)
	if err != nil {
		return err
	}

	cs, err := iimg.GetColorSpace(resources)
	if err != nil {
		return err
	}
	if cs == nil {
		// Default if not specified?
		cs = model.NewPdfColorspaceDeviceGray()
	}

	rgbImg, err := cs.ImageToRGB(*img)
	if err != nil {
		return err
	}

	imgMark := ImageMark{
		Image:  &rgbImg,
		CTM:    gs.CTM,
		Width:  gs.CTM.ScalingFactorX(),
		Height: gs.CTM.ScalingFactorY(),
		Angle:  gs.CTM.Angle(),
	}
	imgMark.X, imgMark.Y = gs.CTM.Translation()

	ctx.extractedImages = append(ctx.extractedImages, imgMark)
	ctx.inlineImages++
	return nil
}

func (ctx *imageExtractContext) extractXObjectImage(name *core.PdfObjectName,
	gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
	common.Log.Debug("extractXObjectImage: name=%#q", name)
	stream, _ := resources.GetXObjectByName(*name)
	if stream == nil {
		return nil
	}

	// Cache on stream pointer so can ensure that it is the same object (better than using name).
	cimg, cached := ctx.cacheXObjectImages[stream]
	if !cached {
		ximg, err := resources.GetXObjectImageByName(*name)
		if err != nil {
			return err
		}
		if ximg == nil {
			return nil
		}

		img, err := ximg.ToImage()
		if err != nil {
			return err
		}

		cimg = &cachedImage{
			image: img,
			cs:    ximg.ColorSpace,
		}
		ctx.cacheXObjectImages[stream] = cimg
	}
	img := cimg.image
	cs := cimg.cs

	rgbImg, err := cs.ImageToRGB(*img)
	if err != nil {
		return err
	}

	common.Log.Debug("@Do CTM: %s", gs.CTM.String())
	imgMark := ImageMark{
		Image:  &rgbImg,
		CTM:    gs.CTM,
		Width:  gs.CTM.ScalingFactorX(),
		Height: gs.CTM.ScalingFactorY(),
		Angle:  gs.CTM.Angle(),
	}
	imgMark.X, imgMark.Y = gs.CTM.Translation()

	ctx.extractedImages = append(ctx.extractedImages, imgMark)
	ctx.xObjectImages++
	return nil
}

// Go through the XObject Form content stream (recursive processing).
func (ctx *imageExtractContext) extractFormImages(name *core.PdfObjectName,
	gs contentstream.GraphicsState, resources *model.PdfPageResources) error {
	xform, err := resources.GetXObjectFormByName(*name)
	if err != nil {
		return err
	}
	if xform == nil {
		return nil
	}

	formContent, err := xform.GetContentStream()
	if err != nil {
		return err
	}

	// Process the content stream in the Form object too:
	formResources := xform.Resources
	if formResources == nil {
		formResources = resources
	}

	// Process the content stream in the Form object too:
	err = ctx.extractContentStreamImages(string(formContent), formResources)
	if err != nil {
		return err
	}
	ctx.xObjectForms++
	return nil
}
