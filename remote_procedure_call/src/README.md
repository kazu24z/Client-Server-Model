# Local Chat Messenger

## 使い方
1. main.pyでサーバを起動します。
2.sample_client.pyにrequestデータを持たせ、実行します。

```
# サーバを起動
python3 main.py
```

```
# クライアントを起動
python3 sample_client.py
```

## RPCで呼び出せる関数
```
- floor(): 10 進数 x を最も近い整数に切り捨て、その結果を整数で返す。
- nroot(): 方程式 rn = x における、r の値を計算する。
- reverse(): 文字列 s を入力として受け取り、入力文字列の逆である新しい文字列を返す。
- validAnagram(): 2 つの文字列を入力として受け取り，2 つの入力文字列が互いにアナグラムであるかどうかを示すブール値を返す。
- sort(): 文字列の配列を入力として受け取り、その配列をソートして、ソート後の文字列の配列を返す
```

## リクエスト
各関数に対応するリクエストの例
```
{"method": "floor", "params": [12.34], "param_types": ["double"]}

{"method": "nroot", "params": [10, 3], "param_types": ["int", "int"]}

{"method": "reverse", "params": ["hello"], "param_types": ["str"]}

{"method": "validAnagram", "params": ["hello", "loleh"], "param_types": ["str", "str"]}

{"method": "sort", "params": [["hello", "thanks", "bomb"]], "param_types": ["str", "str", "str"]}
```
