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

	"golang.org/x/term"
)

// asciiChars represents an ordered set of characters from dark to light const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. "

const URLS = "https://picsum.photos/200/300"
func main() {
	arr := make([]string, 10)
	ch := make(chan string)
	picSize := flag.Int("s", 80, "size of ascii image in terminal")
	t := 10
	flag.Parse()

	// now := time.Now()
	for i := 0; i < t; i++ {
		go getImage(*picSize, ch)
	}

	// fmt.Println(time.Since(now))

	for i := 0; i < t; i++ {
		frame := <-ch
		arr[i] = frame
	}
	// sizer := <-asciiHeight

	w, h, err := term.GetSize(0)
	if err != nil {
		return 
	}
	fmt.Print(w, h)
	time.Sleep(time.Second)
	// var lines int 
	for i := 0; i < t; i++ {
	 	if i >= 0 {
		 	fmt.Printf("\033[%dA", *picSize)
	 	}
		fmt.Println("\n\n")
		fmt.Print(arr[i])
		time.Sleep(500 * time.Millisecond)
	}
  // fmt.Printf("\033[%dB", *picSize/3)
}

func getImage(picSize int, ch chan<- string) {
	resp, err := http.Get(URLS)
	if err != nil {
		fmt.Println("Error decoding image:", err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	outputWidth := picSize
	outputHeight := int(float64(height) / float64(width) * float64(outputWidth) * 0.5) // Adjust aspect ratio
	f := colorASCII(outputHeight, outputWidth, height, width, img)
	ch <- f
}

func colorASCII(outputHeight, outputWidth, height, width int, img image.Image)string{
	const asciiChars = "#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?"
	// const asciiChars = "#$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/()1{}[]?-_+~<>i!lI;:,"
	// const asciiChars = "#"
	// const asciiChars = "󰝤"
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
			resp += fmt.Sprint(correctColor, string(asciiChars[charIndex]), resetColor)
		}
		resp += "\n"
	}
	return resp
}

func rayScaleImage(outputHeight, outputWidth, height, width int, img image.Image)string{
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
