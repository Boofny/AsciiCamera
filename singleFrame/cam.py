import cv2 as cv
import socket
import struct

HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 8080 # The port used by the server

cam = cv.VideoCapture(0)

ret, frame = cam.read()

frame = cv.resize(frame, (640, 480))

if not ret:
    print("Error")
    exit()

success, encode = cv.imencode(
  ".jpg",
  frame,
  [int(cv.IMWRITE_JPEG_QUALITY), 70]
)

if not success:
    print("fail")
    exit()

jpeg_bytes = encode.tobytes()

print("jpeg size", len(jpeg_bytes))

cv.imshow('cam', frame)

if cv.waitKey(1) == 113:
    exit()


try:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        print(f"Connected to {HOST}:{PORT}")

        # Send data
        # s.sendall(MESSAGE)
        # print(f"Sent: {MESSAGE.decode()}")

        s.sendall(struct.pack(">I", len(jpeg_bytes)))  # 4-byte length
        s.sendall(jpeg_bytes)                          # actual JPEG

except ConnectionRefusedError:
    print(f"Connection failed. Make sure a server is running on {HOST}:{PORT}")
except Exception as e:
    print(f"An error occurred: {e}")

