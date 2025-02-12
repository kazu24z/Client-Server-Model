import json
import sys
from sockets import UnixSocket
from request import JSONRequest
from services import Service

# from ./request import 

def main():
    server_address = '/tmp/rpc_socket_file'
    sock = UnixSocket(address = server_address, is_all_in_one = True)
    
    # clientから送られてくるのを待機するために無限ループ
    while True:
        connection, client_address = sock.accept()
        print(f'client connected')
        
        try:
            while True:
                data = connection.recv(4096)
                if data:
                    decoded_data = data.decode('utf-8')
                    print(f"Received data: {decoded_data}")  # デバッグ用
                    
                    # 受け取ったデータをもとにRequestオブジェクトを生成する
                    json_request = JSONRequest(decoded_data)
                    
                    service = Service(json_request.data)
                    method_to_call = getattr(service, service.method)
                    
                    # メソッドを動的に実行する
                    try:
                        result = method_to_call()
                    except Exception as e:
                        result = {
                            "error": str(e),
                            "error_type": e.__class__.__name__,
                            "id": json_request.data["id"]
                        }
                    # dict を json文字列(str)に変換
                    result_json = json.dumps(result)

                    # サービスの実行結果をクライアントに送り返す
                    connection.send(result_json.encode('utf-8')) # json文字列(str)をバイトに変換
                else:
                    print('No Data received')
                    break
        except Exception as e:
            print(f"Exception occurred: {e.__class__.__name__} - {e}")
        finally:
            print('closing current connection')
    # リクエストの処理
    ## 入力のバリデーション
    ## 対象関数呼び出し
    ## レスポンス作成
    ## 返却
    

if __name__ == '__main__':
    main()
