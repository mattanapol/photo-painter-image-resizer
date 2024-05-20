package image_processing

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_imageResizer_Resize_LandscapeSuccess(t *testing.T) {
	imageResizerResizeSuccess(t,
		imageResizerTestInput{
			FilePath: "../assets/test/landscape.jpg",
			Width:    500,
			Height:   500,
		},
		imageResizerTestExpected{
			FilePath:       "../bin/test",
			Width:          500,
			Height:         500,
			OutputFileName: "landscape__500x500.jpg",
		},
	)
}

func Test_imageResizer_Resize_PortraitSuccess(t *testing.T) {
	imageResizerResizeSuccess(t,
		imageResizerTestInput{
			FilePath: "../assets/test/portrait.jpg",
			Width:    500,
			Height:   500,
		},
		imageResizerTestExpected{
			FilePath:       "../bin/test",
			Width:          500,
			Height:         500,
			OutputFileName: "portrait__500x500.jpg",
		},
	)
}

func Test_imageResizer_Resize_SquareSuccess(t *testing.T) {
	imageResizerResizeSuccess(t,
		imageResizerTestInput{
			FilePath: "../assets/test/square.jpg",
			Width:    500,
			Height:   500,
		},
		imageResizerTestExpected{
			FilePath:       "../bin/test",
			Width:          500,
			Height:         500,
			OutputFileName: "square__500x500.jpg",
		},
	)
}

func Test_imageResizer_Resize_PngSuccess(t *testing.T) {
	imageResizerResizeSuccess(t,
		imageResizerTestInput{
			FilePath: "../assets/test/png1.png",
			Width:    500,
			Height:   500,
		},
		imageResizerTestExpected{
			FilePath:       "../bin/test",
			Width:          500,
			Height:         500,
			OutputFileName: "png1__500x500.png",
		},
	)
}

func Test_imageResizer_Resize_BugSuccess(t *testing.T) {
	imageResizerResizeSuccess(t,
		imageResizerTestInput{
			FilePath: "../assets/test/pl.jpg",
			Width:    720,
			Height:   1200,
		},
		imageResizerTestExpected{
			FilePath:       "../bin/test",
			Width:          720,
			Height:         1200,
			OutputFileName: "pl__720x1200.jpg",
		},
	)
}

func Test_imageResizer_Resize_Corrupted(t *testing.T) {
	// Arrange
	generator := NewImageResizer()
	file, err := os.Open("../assets/test/corrupted.jpg")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return
	}

	// Act
	_, err = generator.ResizeImage(context.Background(),
		ImageResizerInput{
			OriginalFile: buffer.Bytes(),
			Height:       200,
			Width:        200,
			Quality:      80,
		},
	)

	// Assert
	assertions := assert.New(t)
	assertions.NotNil(err)
}

type imageResizerTestInput struct {
	FilePath string
	Height   int
	Width    int
}

type imageResizerTestExpected struct {
	FilePath       string
	Height         int
	Width          int
	OutputFileName string
}

func imageResizerResizeSuccess(t *testing.T,
	input imageResizerTestInput,
	expected imageResizerTestExpected,
) {
	// Arrange
	generator := NewImageResizer()
	file, err := os.Open(input.FilePath)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	buffer := &bytes.Buffer{}
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return
	}

	// Act
	startTime := time.Now()
	outputFile, err := generator.ResizeImage(context.Background(),
		ImageResizerInput{
			OriginalFile:     buffer.Bytes(),
			OriginalFileName: input.FilePath,
			Height:           input.Height,
			Width:            input.Width,
			Quality:          80,
		},
	)
	fmt.Printf("Time ues: %s", time.Since(startTime))

	// Assert
	assertions := assert.New(t)
	assertions.Nil(err)
	assertions.NotNil(outputFile)
	assertions.Equal(expected.Height, outputFile.Height)
	assertions.Equal(expected.Width, outputFile.Width)
	assertions.Equal(expected.OutputFileName, outputFile.FileName)

	err = os.MkdirAll(expected.FilePath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	out, err := os.Create(path.Join(expected.FilePath, outputFile.FileName))
	if err != nil {
		panic(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}(out)
	_, err = out.Write(outputFile.ImageContent)
	if err != nil {
		panic(err)
	}
}
