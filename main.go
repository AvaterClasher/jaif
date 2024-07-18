package main

import (
	"fmt"
	"os"
)

const TEMP_RESULT_PATH = "temp.png"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command> <path>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "convert":
		if len(os.Args) < 3 {
			fmt.Println("Second argument ('path') is required. Example: go run main.go convert <path>")
			os.Exit(1)
		}

		path := os.Args[2]
		err := PngToJaif(path)
		if err != nil {
			fmt.Println("Failed to convert PNG to JAIF:", err)
			os.Exit(1)
		} else {
			fmt.Println("Converted PNG to JAIF:", path)
		}

	case "show":
		if len(os.Args) < 3 {
			fmt.Println("File path not provided. Example: `go run main.go show <path>`")
			os.Exit(1)
		}

		filePath := os.Args[2]
		width, height, err := JaifToPng(filePath)
		if err != nil {
			fmt.Println("Failed to convert JAIF to PNG:", err)
			os.Exit(1)
		} else {
			fmt.Println("Successfully converted JAIF to PNG")
			ShowImage(width, height)
		}

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands: convert, show")
		os.Exit(1)
	}
}
