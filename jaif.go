package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Converts a PNG image to a JAIF file
func PngToJaif(path string) error {
	// Open the PNG file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	var str strings.Builder
	bounds := img.Bounds()

	// Iterate through each pixel and write its RGB values in hexadecimal format
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y > bounds.Min.Y {
			str.WriteString("\n")
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			hexColor := fmt.Sprintf("%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
			str.WriteString(hexColor)
		}
	}

	// Get image dimensions
	height := uint32(bounds.Max.Y - bounds.Min.Y)
	width := uint32(bounds.Max.X - bounds.Min.X)

	// Generate JAIF file path
	jaifPath := strings.Replace(filepath.Base(path), ".png", ".jaif", 1)
	err = os.Remove(jaifPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create the JAIF file
	outFile, err := os.OpenFile(jaifPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write the image dimensions and pixel data to the JAIF file
	err = binary.Write(outFile, binary.LittleEndian, width)
	if err != nil {
		return err
	}
	err = binary.Write(outFile, binary.LittleEndian, height)
	if err != nil {
		return err
	}
	_, err = outFile.WriteString(str.String())
	if err != nil {
		return err
	}

	return nil
}

// Converts a JAIF file back to a PNG image
func JaifToPng(path string) (uint32, uint32, error) {
	// Open the JAIF file
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read the image dimensions
	var width, height uint32
	err = binary.Read(reader, binary.LittleEndian, &width)
	if err != nil {
		return 0, 0, err
	}
	err = binary.Read(reader, binary.LittleEndian, &height)
	if err != nil {
		return 0, 0, err
	}

	// Read the pixel data
	content, err := reader.ReadString(0)
	if err != nil && err.Error() != "EOF" {
		return 0, 0, err
	}
	sanitizedContent := strings.ReplaceAll(content, "\n", "")

	// Create a new RGBA image with the read dimensions
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Set the pixel data from the JAIF file
	for i := 0; i < len(sanitizedContent); i += 6 {
		colorStr := sanitizedContent[i : i+6]
		r, g, b := hexToRgb(colorStr)
		x := (i / 6) % int(width)
		y := (i / 6) / int(width)
		img.Set(x, y, color.RGBA{r, g, b, 255})
	}

	// Save the decoded image to a temporary PNG file
	outFile, err := os.Create(TEMP_RESULT_PATH)
	if err != nil {
		return 0, 0, err
	}
	defer outFile.Close()

	err = png.Encode(outFile, img)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

// Converts a hexadecimal color string to RGB values
func hexToRgb(hexStr string) (uint8, uint8, uint8) {
	var r, g, b uint8
	fmt.Sscanf(hexStr, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// Displays an image using the Fyne GUI framework
func ShowImage(width, height uint32) {
	a := app.New()
	w := a.NewWindow("Image preview")

	// Open the temporary PNG file
	imgFile, err := os.Open(TEMP_RESULT_PATH)
	if err != nil {
		fmt.Println("Failed to open image:", err)
		return
	}
	defer imgFile.Close()

	// Decode the image
	img, err := png.Decode(imgFile)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return
	}

	// Create a canvas image and set its properties
	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillContain
	canvasImg.SetMinSize(fyne.NewSize(float32(width), float32(height)))

	// Set the window content and display the image
	w.SetContent(container.NewCenter(canvasImg))
	w.ShowAndRun()
}
