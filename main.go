package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/nfnt/resize"
)

func main() {

	outputPtr := flag.String("o", "", "Output File Name")
	widthPtr := flag.Int("w", 40, "Width of the output")
	colorPtr := flag.Bool("c", false, "Colorful Output")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <image_file> [-o output_file, -w width]")
		os.Exit(1)
	}

	filename := args[0]

	asciiChars := " .,-=+*%@"
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	newWidth := uint(*widthPtr)

	ratio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	newHeight := uint(float64(newWidth) * ratio * 0.5)

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos2)
	bounds := resizedImage.Bounds()

	w, h := bounds.Max.X, bounds.Max.Y

	var writer io.Writer = os.Stdout

	if *outputPtr != "" {
		f, err := os.Create(*outputPtr)
		if err != nil {
			fmt.Println("Error creating output file: ", err)
			return
		}
		defer f.Close()
		writer = f

		if *colorPtr {
			fmt.Println("Warning: Color disabled for file output (files don't support ANSI colors well).")
			*colorPtr = false
		}
	}

	bufferedWriter := bufio.NewWriter(writer)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := resizedImage.At(x, y)

			gray := color.GrayModel.Convert(c).(color.Gray)

			i := int(gray.Y) * (len(asciiChars) - 1) / 255

			if *colorPtr {
				r, g, b, _ := c.RGBA()

				r8 := uint(r >> 8)
				g8 := uint(g >> 8)
				b8 := uint(b >> 8)

				fmt.Fprintf(bufferedWriter, "\x1b[38;2;%d;%d;%dm%c\x1b[0m", r8, g8, b8, asciiChars[i])
			} else {
				fmt.Fprintf(bufferedWriter, "%c", asciiChars[i])
			}
		}
		fmt.Fprintln(bufferedWriter)
	}

	bufferedWriter.Flush()

	if *outputPtr != "" {
		fmt.Println("\nASCII image written to: ", *outputPtr)
	}
}
