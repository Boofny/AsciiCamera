package main

// what this package will contain is the full image to ascii logic and
// a backup random picture from an external api that will be shown as an example to the user

import (
	"flag"
	"fmt"
	"image"
	"net/http"
	"time"

	_ "image/jpeg" // Import for JPEG support
	_ "image/png"  // Import for PNG support
)

// asciiChars represents an ordered set of characters from dark to light const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "

const URL = "https://picsum.photos/200/300"
func main() {
	picSize := flag.Int("s", 80, "size of ascii image in terminal")

	for range 10{ // i guess this would work in the tcp server 
		resp, err := http.Get(URL)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			return
		}

		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		outputWidth := *picSize
		outputHeight := int(float64(height) / float64(width) * float64(outputWidth) * 0.5) // Adjust aspect ratio
		time.Sleep(time.Second / 2)
		f := GrayScaleImage(outputHeight, outputWidth, height, width, img)
		fmt.Println("\r" + f)
	}
}

func GrayScaleImage(outputHeight, outputWidth, height, width int, img image.Image)string{
	im := make([][]string, outputHeight)

	for y := range outputHeight{
		im[y] = make([]string, outputWidth)
	}

	const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?-_+~<>i!lI;:,"
	for y := range outputHeight {
		for x := range outputWidth {
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			gray := (r + g + b) / 3

			charIndex := int(float64(gray) / 65535.0 * float64(len(asciiChars)-1))
			im[y][x] = string( asciiChars[charIndex] )
		}
	}
	var r string

	for x := range im{
		for _, i := range im[x] {
			r += i
		}
		r += "\n"
	}
	return r
}


func ColorSpaces(outputHeight, outputWidth, height, width int, img image.Image)[][]string{
	im := make([][]string, outputHeight)
	// for y := 0; y < outputHeight; y++ {
	for y := range outputHeight{
		im[y] = make([]string, outputWidth)
	}
	resetColor := "\033[0m"
	for y := range outputHeight {
		for x := range outputWidth {
			// Get pixel from original image, scaled to output dimensions
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			red := r >> 8
			green := g >> 8
			blue := b >> 8

			correctColor := fmt.Sprintf("\x1b[48;2;%d;%d;%dm \x1B[0m", red, green, blue)

			s := fmt.Sprint(correctColor, resetColor)
			im[y][x] = s
		}
	}
	return im
}

func ColorANSI(outputHeight, outputWidth, height, width int, img image.Image, asciiChar string)[][]string{
	im := make([][]string, outputHeight)
	// for y := 0; y < outputHeight; y++ {
	for y := range outputHeight{
		im[y] = make([]string, outputWidth)
	}
	resetColor := "\033[0m"
	for y := range outputHeight {
		for x := range outputWidth {
			// Get pixel from original image, scaled to output dimensions
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			red := r >> 8
			green := g >> 8
			blue := b >> 8

			correctColor := fmt.Sprintf("\u001b[38;2;%d;%d;%dm", red, green, blue)

			s := fmt.Sprint(correctColor, asciiChar, resetColor)
			im[y][x] = s
		}
	}
	return im
}

func ColorASCII(outputHeight, outputWidth, height, width int, img image.Image)[][]string{
	const asciiChars = "#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?"
	im := make([][]string, outputHeight)
	for y := range outputHeight{
		im[y] = make([]string, outputWidth)
	}
	resetColor := "\033[0m"
	for y := range outputHeight {
		for x := range outputWidth {
			// Get pixel from original image, scaled to output dimensions
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			gray := (r + g + b) / 3

			// Map grayscale to ASCII character
			charIndex := int(float64(gray) / 65535.0 * float64(len(asciiChars)-1))

			red := r >> 8
			green := g >> 8
			blue := b >> 8

			correctColor := fmt.Sprintf("\u001b[38;2;%d;%d;%dm", red, green, blue)

			// fmt.Print(correctColor, string(asciiChars[charIndex]), resetColor)
			var s string
			if string(asciiChars[charIndex]) == "@"{
				s = " "
			}else{
				s = fmt.Sprint(correctColor, string(asciiChars[charIndex]), resetColor)
			}
			im[y][x] = s
		}
	}
	return im
}

