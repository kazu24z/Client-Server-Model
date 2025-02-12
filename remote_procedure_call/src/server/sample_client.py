import socket
import sys
import json

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

server_address = '/tmp/rpc_socket_file'
print('connecting to {}'.format(server_address))

try:
    sock.connect(server_address)
except socket.error as err:
    print(err)
    sys.exit(1)

try:
    # ソケットで送るときはバイナリ形式
    message = b'{"method": "sort", "params": [["hello", "thanks", "bomb"]], "param_types": ["list"]}'
    # message = b'{key: 12}'
    sock.sendall(message)
    
    # サーバからの応答待ち時間
    sock.settimeout(2)
    
    try:
        while True:
            # サーバから受け取るデータ
            data = sock.recv(4096).decode('utf-8') # バイトで受け取ってそれをstrにデコード
            
            if data :  
                # デコード
                # response = json.loads(data) # json文字列(str)をdictに変換
                print('Server response: ' + data)
            else:
                break
    # 2秒間サーバから応答がない場合
    except(TimeoutError):
        print('Socket timeout, ending listening for server messages')
finally:
    print('closing socket')
    sock.close()
