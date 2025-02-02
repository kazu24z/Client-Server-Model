import socket
import os # OSに依存する処理をまとめたライブラリ（ファイル作成,削除など）
from faker import Faker

def main():
    # ソケットインスタンス作成
    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    # サーバ起動
    server_address = '/tmp/socket_file'
    
    # 既存のソケットがあれば落とす
    try:
        os.unlink(server_address) #ファイルを削除する関数
    except FileNotFoundError:
        pass
    
    print('Starting up on {}'.format(server_address))

    # サーバアドレスにソケットをバインドする
    sock.bind(server_address)
    
    sock.listen(0)
    # ユーザーの入力待ち
    while True:
        # クライアントからの接続を受け入れ
        connection, client_address = sock.accept()
        
        print(f'connection success')
        
        try:
            while True:
                # 入力があったら処理
                data = connection.recv(16)
                if data:
                    data_str = data.decode('utf-8')
                    print(f'receive: {data_str}')
                    
                    response  = create_response(data_str)
                    
                    if response == 'invalid_request':
                        raise ValueError('invalid request')
                    
                    connection.sendall(response.encode())
                    print(f'response: {response}')
                else:
                    print('no data from', client_address)
                    break
        finally:
            print('Closing current connection')
            connection.close()

def create_response(request: str):
    fake = Faker()
    if request == 'name':
        return fake.name()
    if request == 'address':
        return fake.address()
    if request == 'email':
        return fake.email()
    
    return 'invalid_request'
    
    
if __name__ == '__main__':
    main()
