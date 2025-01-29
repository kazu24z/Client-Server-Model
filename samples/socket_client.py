import socket
import sys

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

server_address = '/tmp/socket_file'
print('connecting to {}'.format(server_address))

try:
    sock.connect(server_address)
except socket.error as err:
    print(err)
    sys.exit(1)

try:
    # ソケットで送るときはバイナリ形式
    message = b'Sending a message to the server side'
    sock.sendall(message)
    
    # サーバからの応答待ち時間
    sock.settimeout(2)
    
    try:
        while True:
            # サーバから受け取るデータ
            # 最大サイズ32バイト
            data = str(sock.recv(32))
            
            if data :  
                print('Server response: ' + data)
            else:
                break
    # 2秒間サーバから応答がない場合
    except(TimeoutError):
        print('Socket timeout, ending listening for server messages')
finally:
    print('closing socket')
    sock.close()
