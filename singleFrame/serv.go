package main

import (
	"image/jpeg"
	"image"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// so this works with netcat but not with the currect example of ./clientExample.go

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		go handleConnection(conn)
	}
}


func handleConnection(conn net.Conn) {
 	defer conn.Close()

	// Read 4-byte length
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBytes)
	if err != nil {
		fmt.Println("Connection closed or error:", err)
		return
	}
	length := binary.BigEndian.Uint32(lengthBytes)

	// Read the full JPEG frame
	jpegData := make([]byte, length)
	_, err = io.ReadFull(conn, jpegData)
	if err != nil {
		fmt.Println("Error reading frame:", err)
		return
	}

	// Optionally decode or save
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		fmt.Println("Failed to decode JPEG:", err)
	}

	fmt.Println("Received image:", img.At(100, 100))
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	outputWidth := 80
	outputHeight := int(float64(height) / float64(width) * float64(outputWidth) * 0.5) // Adjust aspect ratio
	f := ColorASCII(outputHeight, outputWidth, height, width, img)
	fmt.Println(f)
}

func GrayScaleImage(outputHeight, outputWidth, height, width int, img image.Image)string{
	const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?-_+~<>i!lI;:,"
	var resp string
	for  y := range outputHeight {
		for x := range outputWidth {
			originalX := int(float64(x) / float64(outputWidth) * float64(width))
			originalY := int(float64(y) / float64(outputHeight) * float64(height))
			pixel := img.At(originalX, originalY)
			r, g, b, _ := pixel.RGBA()

			gray := (r + g + b) / 3

			charIndex := int(float64(gray) / 65535.0 * float64(len(asciiChars)-1))
			resp += string(asciiChars[charIndex])
		}
		resp += "\n"
	}
	return resp
}

func ColorSpaces(outputHeight, outputWidth, height, width int, img image.Image)string{
	var resp string
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
			resp += s
		}
		resp += "\n"
	}
	return resp
}

func ColorASCII(outputHeight, outputWidth, height, width int, img image.Image)string{
	// const asciiChars = "#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?"
	const asciiChars = "#"
	var resp string
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
			resp += s
		}
		resp += "\n"
	}
	return resp
}

