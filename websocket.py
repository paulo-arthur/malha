import socket

def server(host='localhost', port=3878):
    DATA_PAYLOAD = 2048

    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_address = (host, port)

    print(f"Starting server on port {port} | {server_address}")

    sock.bind(server_address)
    sock.listen(5)

    i = 0
    try:
        while True:
            print(f"Listening... | {i} connections were made.")
            client, address = sock.accept()
            data = client.recv(DATA_PAYLOAD)
            print(f"A request was received from {client} | {address}.")

            if data:
                print(data)

            res = f"Connection no. {i} received. Congratulations!"
            http_response = (
                "HTTP/1.1 200 OK\r\n"
                "Content-Type: text/plain; charset=utf-8\r\n"
                f"Content-Length: {len(res)}\r\n"
                "Connection: close\r\n"
                "\r\n"  # Linha em branco obrigatória separando cabeçalho e corpo
                f"{res}"
            )

            client.sendall(http_response.encode('utf-8'))
            client.close()
            i += 1
    except Exception as e:
        print("\nClosing server...")
        print(e)
        sock.close()

server()