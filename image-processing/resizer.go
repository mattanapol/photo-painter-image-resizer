package image_processing

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"path"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"golang.org/x/image/draw"
)

type ImageResizer interface {
	ResizeImage(ctx context.Context,
		inputFile ImageResizerInput,
	) (ImageResizerOutput, error)
}

type imageResizer struct {
}

func NewImageResizer() imageResizer {
	return imageResizer{}
}

func (t imageResizer) ResizeImage(ctx context.Context,
	inputFile ImageResizerInput,
) (ImageResizerOutput, error) {

	log.Printf("Resize image for size: %dx%d", inputFile.Width, inputFile.Height)
	outputFile, err := t.processFile(ctx,
		inputFile.OriginalFile,
		inputFile.Height,
		inputFile.Width,
		inputFile.Quality,
	)
	if err != nil {
		log.Printf("Error processing file: %v", err)
		return ImageResizerOutput{}, err
	}

	return ImageResizerOutput{
		ImageContent: outputFile.OutputFile,
		Height:       int(outputFile.Height),
		Width:        int(outputFile.Width),
		FileName: getImageResizedFileName(inputFile.OriginalFileName,
			inputFile.Height, // right now we are use request height and width for resized file name
			inputFile.Width,
		),
	}, nil
}

func (t imageResizer) processFile(ctx context.Context,
	input []byte,
	height int,
	width int,
	jpegQuality int,
) (*processFileOutput, error) {
	kind, err := validateImageFileType(input)
	if err != nil {
		return nil, fmt.Errorf("error validating image file type: %s", err)
	}

	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("error decoding image config: %s", err)
	}

	originalWidth, originalHeight := img.Bounds().Dx(), img.Bounds().Dy()
	newWidth, newHeight := calculateWidthHeight(width, height,
		originalHeight, originalWidth,
	)
	resizedImg := image.NewRGBA(image.Rect(0, 0, int(newWidth), int(newHeight)))
	draw.BiLinear.Scale(resizedImg, resizedImg.Rect, img, img.Bounds(), draw.Over, nil)

	combinedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	offsetX := (width - int(newWidth)) / 2
	offsetY := (height - int(newHeight)) / 2
	// fill white background
	draw.Draw(combinedImage, combinedImage.Bounds(), image.Black, image.Point{}, draw.Src)

	// fill image
	draw.Draw(combinedImage,
		image.Rect(offsetX, offsetY, offsetX+int(newWidth), offsetY+int(newHeight)),
		resizedImg,
		resizedImg.Rect.Min,
		draw.Over)

	buf := &bytes.Buffer{}
	if strings.Contains(kind.MIME.Subtype, "png") {
		err = png.Encode(buf, combinedImage)
	} else { // default to jpeg encoding
		err = jpeg.Encode(buf, combinedImage, &jpeg.Options{Quality: jpegQuality})
	}

	if err != nil {
		return nil, err
	}

	return &processFileOutput{
		OutputFile: buf.Bytes(),
		Height:     uint(height),
		Width:      uint(width),
	}, nil
}

func calculateWidthHeight(width int, height int, originalHeight int, originalWidth int) (uint, uint) {
	// Calculate the aspect ratios
	canvasRatio := float64(width) / float64(height)
	originalRatio := float64(originalWidth) / float64(originalHeight)

	var newWidth, newHeight uint

	// Determine the new dimensions based on the aspect ratio comparison
	if originalRatio > canvasRatio {
		// The original image is wider compared to the canvas
		newWidth = uint(width)
		newHeight = uint(float64(width) / originalRatio)
	} else {
		// The original image is taller compared to the canvas
		newHeight = uint(height)
		newWidth = uint(float64(height) * originalRatio)
	}

	return newWidth, newHeight
}

func getImageResizedFileName(originalFileName string,
	height int,
	width int,
) string {
	fileName := path.Base(originalFileName)
	fileExt := path.Ext(fileName)
	fileNameWithoutExt := strings.TrimSuffix(fileName, fileExt)
	return fmt.Sprintf("%s__%dx%d%s", fileNameWithoutExt, width, height, fileExt)
}

func validateImageFileType(input []byte) (types.Type, error) {
	kind, _ := filetype.Match(input[:261])
	if kind == filetype.Unknown {
		return kind, fmt.Errorf("unknown file type")
	}

	if !strings.Contains(kind.MIME.Type, "image") {
		return kind, fmt.Errorf("not an image")
	}

	return kind, nil
}
