
import math

# RPCで提供する関数を持ったクラス
class Service:
    method: str
    params: list
    param_types: list
    id: int
    
    def __init__(self,request):
        expected_keys = {"method", "params", "param_types", "id"}
        if set(request.keys()) == expected_keys:
            self.method =request["method"]
            self.params = request["params"]
            self.param_types = request["param_types"]
            self.id = request["id"]
        else:
            raise ValueError("JSON Key Error: property not collect")

    # floor(): float: x, 10 進数 x を最も近い整数に切り捨て、その結果を整数で返す。
    def floor(self)->int:
        # 入力のバリデーション
        if len(self.params) != 1:
            raise ValueError("Invalid number of parameters")
        
        if self.param_types[0] != "float":
            raise ValueError("Invalid parameter type")
        
        x = self.params[0]
        
        result = int(math.floor(x))
        return {"result": result, "result_type": type(result).__name__, "id": self.id}
    
    # nroot(): int n, int x,  方程式 r**n = x における、r の値を計算する。
    def nroot(self)->float:
        # 入力のバリデーション
        if len(self.params) != 2:
            raise ValueError("Invalid number of parameters")
        
        if self.param_types[0] != "int" or self.param_types[1] != "int":
            raise ValueError("Invalid parameter type")
        
        n = self.params[0]
        x = self.params[1]
        
        # 計算
        r = x ** (1/n)
        
        return {"result": r, "result_type": type(r).__name__, "id": self.id}
        
    # reverse(): str s, 文字列 s を入力として受け取り、入力文字列の逆である新しい文字列を返す。
    def reverse(self)->str:
        # 入力のバリデーション
        if len(self.params) != 1:
            raise ValueError("Invalid number of parameters")
        
        if self.param_types[0] != "str":
            raise ValueError("Invalid parameter type")
        
        s = self.params[0]
        
        new = s[::-1]
        
        return {"result": new, "result_type": type(new).__name__, "id": self.id}
    
    # 2つの文字列を入力として受け取り，2 つの入力文字列が互いにアナグラムであるかどうかを示すブール値を返す。
    def validAnagram(self)->bool:
        # 入力のバリデーション
        if len(self.params) != 2:
            raise ValueError("Invalid number of parameters")
        
        if self.param_types[0] != "str" or self.param_types[1] != "str":
            raise ValueError("Invalid parameter type")
        
        s1 = self.params[0]
        s2 = self.params[1]
        
        s1_sorted = sorted(s1)
        s2_sorted = sorted(s2)
        
        return {"result": s1_sorted == s2_sorted, "result_type": type(s1_sorted == s2_sorted).__name__, "id": self.id}
        
    # 文字列の配列を入力として受け取り、その配列をソートして、ソート後の文字列の配列を返す。
    def sort(self)->list[str]:
        # 入力のバリデーション
        if self.param_types[0] != "list":
            raise ValueError("Invalid parameter type")
        
        arr = self.params[0]
        
        return {"result": sorted(arr), "result_type": type(sorted(arr)).__name__, "id": self.id}
