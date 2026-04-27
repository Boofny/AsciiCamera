// Package camera holds the code for running the tcpserver and the camera mode methods
package camera

import (
	"bytes"
	"context"
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

func RunCam(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	flag.Parse()
	camMode = PickMode(Mode(*camType)) // may work?

	// camMode := PickMode(Mode(*camType))
	imgCh := make(chan string, 1)
	errCh := make(chan error, 1)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Print("\033[H\033[2J") // Clear screen and move to top-left

	go func() { // single goroutine in order to have good timing with the img channel
		for img := range imgCh {
			fmt.Printf("\033[%dA\n\n%s", *picWidth, img)
		}
	}()

	go func(){
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Error accepting conn:", err)
				continue
			}
			go handleConnection(conn, imgCh, errCh, cancel)
		}
	}()

	// could use this somehow to implement the errCh aswell as the ctx channel
	for range 2{
		select{
			case <-ctx.Done:
				listener.Close();
				return ctx.Err()
			case err := <-errCh:
				return err
		}
	}
	// <-ctx.Done()
	// listener.Close()
	return ctx.Err()
}

// NOTE: error Chanel is not being used eithor find a way to use it or remove later 
func handleConnection(conn net.Conn, imgCh chan string, errCh chan<- error, cancel context.CancelFunc){
	defer conn.Close()
	for {
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthBytes)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("Connection closed cleanly")
				cancel() // or however you signal shutdown
				return
			}
			fmt.Println("Connection error:", err)
			cancel()
			return
		}
		length := binary.BigEndian.Uint32(lengthBytes)
		fmt.Printf("received length: 0x%X\n", length) // add this
		if length == 0xFFFFFFFF {
			fmt.Println("Shutdown signal received")
			cancel()
			return
		}

		jpegData := make([]byte, length)
		_, err = io.ReadFull(conn, jpegData)
		if err != nil {
			fmt.Println("Error reading frame:", err)
			errCh <- err
		}

		img, err := jpeg.Decode(bytes.NewReader(jpegData))
		if err != nil {
			fmt.Println("Failed to decode JPEG:", err)
			errCh <- err
		}

		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y
		outputWidth := *picWidth
		outputHeight := int(float64(height) / float64(width) * float64(outputWidth) * 0.5)
		f := camMode(outputHeight, outputWidth, height, width, img)

		select {
		case imgCh <- f:
		default:
		}
	}
}
