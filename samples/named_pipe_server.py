# 名前付きパイプ
# プロセス間通信の一つで、プロセス間でデータをやり取りするための仕組み
import os
import json
import time

config = json.load(open('config.json'))

if os.path.exists(config['file_path']):
    os.remove(config['file_path'])
    
# mkfifo()で指定したパスに名前付きパイプを作成する
# 名前付きパイプはファイルシステム上に作成される
# 0o600はパーミッションを指定している(chmod 600と同じ)
os.mkfifo(config['file_path'], 0o600)

print("FIFO named '% s' is created successfully." % config['file_path'])
print("Type in what you would like to send to clients.")

flag = True

# 5秒ごとに実行する
while flag:
    inputstr = input()
    
    if(inputstr == 'exit'):
        flag = False
    else:
        with open(config['file_path'], 'w') as fifo:
            fifo.write(inputstr)
            print("Sent: ", inputstr)
            
os.remove(config['file_path'])
