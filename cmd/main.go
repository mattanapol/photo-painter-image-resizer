package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"slices"

	"github.com/google/uuid"
	image_processing "github.com/mattanapol/photo-painter-image-resizer/image-processing"
)

var (
	ResizeToWidth   = 720
	ResizeToHeight  = 1200
	ResizeToQuality = 100
	RandomName      = false
)

// go run ./cmd/main.go <image_input_path> -width <width> -height <height> -quality <quality> -random-name
func main() {
	ctx := context.Background()
	// Get the directory path from the first argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a directory path as an argument.")
		return
	}
	flag.IntVar(&ResizeToWidth, "width", 720, "Resize image to width")
	flag.IntVar(&ResizeToHeight, "height", 1200, "Resize image to height")
	flag.IntVar(&ResizeToQuality, "quality", 100, "Resize image to quality")
	flag.BoolVar(&RandomName, "random-name", false, "Use random name for output file")

	dirPath := os.Args[1]

	// Read directory contents
	files, err := getAllFilePaths(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Loop through files and print names
	for _, file := range files {
		fmt.Println(file)
	}

	imageResizer := image_processing.NewImageResizer()
	for _, file := range files {
		// read file into byte array
		fileContent, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %s", file, err)
			continue
		}

		output, err := imageResizer.ResizeImage(ctx, image_processing.ImageResizerInput{
			OriginalFile:     fileContent,
			OriginalFileName: file,
			Height:           ResizeToHeight,
			Width:            ResizeToWidth,
			Quality:          ResizeToQuality,
		})
		if err != nil {
			log.Printf("Error resizing image %s: %s", file, err)
			continue
		}

		if true {
			output.FileName = getRandomFileName(path.Ext(file))
		}

		err = image_processing.WriteFile("output", output.FileName, output.ImageContent)
		if err != nil {
			log.Printf("Error writing file %s: %s", output.FileName, err)
			continue
		}
	}

}

var (
	ignoreFileList = []string{".DS_Store"}
)

func getAllFilePaths(dirPath string) ([]string, error) {
	var filePaths []string
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() && !slices.Contains(ignoreFileList, file.Name()) {
			filePath := path.Join(dirPath, file.Name())
			filePaths = append(filePaths, filePath)
		}
	}
	return filePaths, nil
}

func getRandomFileName(ext string) string {
	return fmt.Sprintf("%s%s", uuid.New(), ext)
}
