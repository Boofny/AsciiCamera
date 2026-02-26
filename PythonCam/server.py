import socket

HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 8080 # The port used by the server
MESSAGE = b"hello world!" # Data is sent as bytes

try:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        print(f"Connected to {HOST}:{PORT}")

        # Send data
        s.sendall(MESSAGE)
        print(f"Sent: {MESSAGE.decode()}")

        # Receive a response from the server (optional)
        data = s.recv(1024)
        print(f"Received from server: {data.decode()}")

except ConnectionRefusedError:
    print(f"Connection failed. Make sure a server is running on {HOST}:{PORT}")
except Exception as e:
    print(f"An error occurred: {e}")

