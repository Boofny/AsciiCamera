package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net"
	"strings"
)

// golang perfomince changes
/*
Number one made a bufferd chanel for the image
Made a seperate goroutine just for the printing in order to have better printings even at large scale
Found out about the select design pattern that skips the unused frames from a channl very cool
*/
var picWidth = flag.Int("s", 80, "Size of image output") // global size so any thing can use this
var camType = flag.Int("cam", 1, "The type of Camera display, 1 for Color pound sign, 2 for ascii color spaces, 3 list of ascii chars, 4 gray ascii output")

type Mode int

const (
	ColorPound Mode = iota + 1 //default if no flag is given or 1 is passed
	ASCIIColor
	ColorASCIIChars
	GreyASCII
)

func main() {
	flag.Parse()

	switch *camType{
		case ColorPound:
	}
	imgCh := make(chan string, 1)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	fmt.Print("\033[H\033[2J") // Clear screen and move to top-left

	go func() { // single goroutine in order to have good timing with the img channel
		for img := range imgCh {
			fmt.Printf("\033[%dA\n\n%s", *picWidth, img)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		go handleConnection(conn, imgCh)
	}
}

func handleConnection(conn net.Conn, imgCh chan string) {
	defer conn.Close()

	for {
		// Step 1: Read 4-byte length
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthBytes)
		if err != nil {
			fmt.Println("Connection closed or error:", err)
			return
		}
		length := binary.BigEndian.Uint32(lengthBytes)

		// Step 2: Read the full JPEG frame
		jpegData := make([]byte, length)
		_, err = io.ReadFull(conn, jpegData)
		if err != nil {
			fmt.Println("Error reading frame:", err)
			return
		}

		// Step 3: Optionally decode or save
		img, err := jpeg.Decode(bytes.NewReader(jpegData))
		if err != nil {
			fmt.Println("Failed to decode JPEG:", err)
			continue
		}

		// fmt.Println("Received image:", img.At(100, 100))
		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		outputWidth := *picWidth
		outputHeight := int(float64(height) / float64(width) * float64(outputWidth) * 0.5) // Adjust aspect ratio
		f := ColorASCII(outputHeight, outputWidth, height, width, img)
		select { // works by only running/displying if the goroutine is ready if not just skip a frame
		case imgCh <- f:
		default:
		}
	}
}

func ColoredASCIIPound(outputHeight, outputWidth, height, width int, img image.Image) string {
	const asciiChar = '#'
	resetColor := "\033[0m"
	var sb strings.Builder
	sb.Grow(outputHeight * outputWidth * 25)

	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth-(x+1)
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
			x = outputWidth-(x+1)
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

func colorSpaces(outputHeight, outputWidth, height, width int, img image.Image) string {
	resetColor := "\033[0m"
	var sb strings.Builder
	sb.Grow(outputHeight * outputWidth * 25)
	for y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth-(x+1)

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

func GrayScaleImage(outputHeight, outputWidth, height, width int, img image.Image)string{
	const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?-_+~<>i!lI;:,"
	var resp strings.Builder
	resp.Grow(outputHeight * outputWidth * 25)
	for  y := range outputHeight {
		for x := range outputWidth {
			x = outputWidth-(x+1)
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
