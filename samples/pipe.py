import os

def main():
    # 読み込み用、書き込み用のパイプを作成
    r, w = os.pipe()
    # プロセスをフォーク
    pid = os.fork()

    # 親プロセルの処理
    if pid > 0:
        # 親→子にデータを送るだけなので、読み込み用のパイプは使わない
        os.close(r)
        message = 'Message from parent with pid {}'.format(os.getpid())
        print('Parent: {}'.format(message, os.getpid()))
        # メッセージをパイプに書き込む
        os.write(w, message.encode('utf-8'))

    else:
        os.close(w)
        print('Fork is 0, this is child process PID:', os.getpid())
        # パイプからファイル記述子のデータを読み込む
        f = os.fdopen(r)
        
        print("Incoming string:", f.read())

if __name__ == '__main__':
    main()
