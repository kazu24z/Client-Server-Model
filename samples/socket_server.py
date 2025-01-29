import socket
import os

# ソケットドメイン, ソケットタイプを指定してソケットを作成
sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)

# UNIXソケットのパス（UNIXソケットはファイルを介して接続される）
# サーバが待ち受けるパスを指定する
server_address = '/tmp/socket_file'

try:
    os.unlink(server_address)
except FileNotFoundError:
    pass

print('Starting up on {}'.format(server_address))

# サーバアドレスにソケットをバインドする
sock.bind(server_address)

# ソケットが接続要求を待機する（これで待ち状態を作る）
sock.listen(0)
counter = 1
# クライアントの接続を待ち続ける
while True:
    
    print(f'{counter}回目')
    # クライアントからの接続を受け入れる
    ## conn: 新しいソケットオブジェクト(上のsockの別インスタンス)
    ## address: 接続先でソケットにバインドしているアドレス（トンネルの向こう側）
    connection, client_address = sock.accept() # accept()は接続, クライアントアドレスのペア
    
    try:
        print('connection from', client_address)

        while True:
            # recv(byte) １度に16byteまでを受け取る
            data = connection.recv(16)
            
            # dataはバイナリ形式なのでそれを文字列に変換する
            data_str = data.decode('utf-8')

            print('Received '+ data_str)
            
            if data:
                response = 'Processing ' + data_str
                
                # client側に送る
                connection.sendall(response.encode()) # encode()はバイナリにするって意味
                
            else:
                print('no data from', client_address)
                break
    finally:
        print('Closing current connection')
        connection.close()
    counter = counter + 1
