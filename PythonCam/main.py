# import cv2 as cv
# import socket
# import struct
#
# HOST = '127.0.0.1'  # The server's hostname or IP address
# PORT = 8080 # The port used by the server
#
# cam = cv.VideoCapture(0)
#
# while True:
#     ret, frame = cam.read()
#
#     frame = cv.resize(frame, (640, 480))
#
#     if not ret:
#         print("Error")
#         break
#
#     success, encode = cv.imencode(
#       ".jpg",
#       frame,
#       [int(cv.IMWRITE_JPEG_QUALITY), 70]
#     )
#
#     if not success:
#         print("fail")
#         break
#
#     jpeg_bytes = encode.tobytes()
#
#     print("jpeg size", len(jpeg_bytes))
#
#     cv.imshow('cam', frame)
#
#     if cv.waitKey(1) == 113:
#         break
#
#     try:
#         with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
#             s.connect((HOST, PORT))
#             print(f"Connected to {HOST}:{PORT}")
#
#             # Send data
#             # s.sendall(MESSAGE)
#             # print(f"Sent: {MESSAGE.decode()}")
#
#             s.sendall(struct.pack(">I", len(jpeg_bytes)))  # 4-byte length
#             s.sendall(jpeg_bytes)                          # actual JPEG
#
#     except ConnectionRefusedError:
#         print(f"Connection failed. Make sure a server is running on {HOST}:{PORT}")
#     except Exception as e:
#         print(f"An error occurred: {e}")
#
#

import cv2 as cv
import socket
import struct

HOST = '127.0.0.1'
PORT = 8080

cam = cv.VideoCapture(0)

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
try:
    s.connect((HOST, PORT))
    print(f"Connected to {HOST}:{PORT}")
except ConnectionRefusedError:
    print(f"Connection failed. Make sure a server is running on {HOST}:{PORT}")
    exit(1)

try:
    while True:
        ret, frame = cam.read()
        if not ret:
            print("Camera read failed")
            break

        frame = cv.resize(frame, (640, 480))

        # Encode JPEG
        success, encode = cv.imencode(".jpg", frame, [int(cv.IMWRITE_JPEG_QUALITY), 70])
        if not success:
            print("JPEG encode failed")
            break

        jpeg_bytes = encode.tobytes()
        # print("jpeg size", len(jpeg_bytes))

        # Send length + JPEG
        s.sendall(struct.pack(">I", len(jpeg_bytes)))  # 4-byte length
        s.sendall(jpeg_bytes)

        # Show frame locally
        cv.imshow('cam', frame)
        if cv.waitKey(1) == 113:  # 'q' key
            break
finally:
    s.close()
    cam.release()
    cv.destroyAllWindows()
