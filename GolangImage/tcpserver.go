package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net"
)

// golang perfomince changes
/*
Number one made a bufferd chanel for the image
Made a seperate goroutine just for the printing in order to have better printings even at large scale
Found out about the select design pattern that skips the unused frames from a channl very cool
*/
var picWidth = flag.Int("s", 80, "Size of image output") // global size so any thing can use this

var camType = flag.Int("mode", 1,
	"The type of Camera display, 1 for Color pound sign, 2 for ascii color spaces, 3 list of ascii chars, 4 gray ascii output")

type Mode uint8

const (
	ColorPound Mode = iota + 1 //default if no flag is given or 1 is passed
	ASCIIColor
	ColorASCIIChars
	GreyASCII
)

// PickMode just picks the function that the flag chooses and the default for the flag is always 1 so no need for default 
func PickMode(mode Mode)CamFunc{
	var resp CamFunc
	fmt.Println(mode)
	switch mode{
	case ColorPound:
		resp = ColoredASCIIPound
	case ASCIIColor: 
		resp = ColorSpaces
	case ColorASCIIChars:
		resp = ColorASCII
	case GreyASCII:
		resp = GrayScaleImage
	default:
		log.Fatal("Not a mode select 1-4")
	}
	return resp
}

var camMode CamFunc// may work?

func main() {
	flag.Parse()
	camMode = PickMode(Mode(*camType)) // may work?

	// camMode := PickMode(Mode(*camType))
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
		f := camMode(outputHeight, outputWidth, height, width, img)
		select { // works by only running/displying if the goroutine is ready if not just skip a frame
		case imgCh <- f:
		default:
		}
	}
}
