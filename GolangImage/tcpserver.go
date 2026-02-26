package main

import (
	"image/jpeg"
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

		fmt.Println("Received image:", img.At(100, 100))
	}
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
//
// 	reader := bufio.NewReader(conn)
//
// 	for {
// 		msg, err := reader.ReadByte()
// 		if err != nil {
// 			return
// 		}
//
// 		fmt.Print(string(msg))
// 	}
// }
