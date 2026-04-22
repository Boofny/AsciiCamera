package main

//TODO: make the function use byte streams rather than using strings
import (
	"fmt"
	"image"
	"strings"
)

type CamFunc func(ImgParam) string

type ImgParam struct {
	OutPutHeight,
	OutPutWidth,
	Height,
	Width int
	Img image.Image
}

func ColoredASCIIPound(outputHeight, outputWidth, height, width int, img image.Image) string {
	const asciiChar = '#'
	resetColor := "\033[0m"
	var sb strings.Builder
	sb.Grow(outputHeight * outputWidth * 25)

	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth - (x + 1)
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()
			fmt.Fprintf(&sb, "\u001b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8) // could slow down for future reference
			sb.WriteByte(asciiChar)
			sb.WriteString(resetColor)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func ColorASCII(outputHeight, outputWidth, height, width int, img image.Image) string {
	const asciiChars = "#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?"
	var resp strings.Builder
	resp.Grow(outputHeight * outputWidth * 25)
	resetColor := "\033[0m"
	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth - (x + 1)
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

			fmt.Fprint(&resp, correctColor, string(asciiChars[charIndex]), resetColor)
		}
		resp.WriteByte('\n')
	}
	return resp.String()
}

func ColorSpaces(outputHeight, outputWidth, height, width int, img image.Image) string {
	resetColor := "\033[0m"
	var sb strings.Builder
	sb.Grow(outputHeight * outputWidth * 25)
	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth - (x + 1)

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
			sb.WriteString(s)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func GrayScaleImage(outputHeight, outputWidth, height, width int, img image.Image) string {
	const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?-_+~<>i!lI;:,"
	var resp strings.Builder
	resp.Grow(outputHeight * outputWidth * 25)
	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth - (x + 1)
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			gray := (r + g + b) / 3

			charIndex := int(float64(gray) / 65535.0 * float64(len(asciiChars)-1))
			s := asciiChars[charIndex]
			resp.WriteByte(s)
		}
		resp.WriteByte('\n')
	}
	return resp.String()
}
