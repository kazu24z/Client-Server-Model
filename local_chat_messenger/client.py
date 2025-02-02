import socket
import sys

def main():
    # ユーザープロパティを入力し、それに基づくデータを返すアプリケーション
    # ex) 住所→address, 名前→name

    # TCPでやりとりする

    # socketインスタンス生成（UNIXドメイン, TCP）
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # アドレスファイルを指定（サーバ側の）
    # server_address = '/tmp/socket_file'
    server_address = 'localhost'
    server_port = 50007

    # サーバーへの接続を確立
    try:
        sock.connect((server_address, server_port))
    except socket.error as err:
        print(err)
        sys.exit(1)

    try:
        # ユーザーからの入力待ち
        user_input = input('input name, address or email...') # 文字列
        is_validated = validate_input(user_input)
        if not is_validated:
            raise ValueError('The argument is not valid.')

        binarized = user_input.encode('utf-8')
        
        # 受け取ったデータをサーバーに送信
        sock.sendall(binarized)
        
        sock.settimeout(3)
        # サーバからのレスポンスを受け取る
        try:
            while True:
                data = sock.recv(128).decode('utf-8') #最大32バイトを文字列に
                
                if data:
                    print(f'{user_input}: {data}')
                else:
                    break
        except ValueError as e:
            print(e)
        except TimeoutError:
            print('Socket timeout, ending listening for server messages')
    finally:
        print('closing socket')
        sock.close()

def validate_input(input: str):
    input_list = ['name', 'address', 'email']
    
    if(input in input_list):
        return True
    else:
        return False

if __name__ == '__main__':
    main()
