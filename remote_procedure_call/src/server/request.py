import json
import uuid
from typing import TypedDict, Literal, List, Any

class JSONRequestData(TypedDict):
    method: Literal["subtract", "floor", "nroot", "reverse", "validAnagram", "sort"]
    params: List[Any]  # 整数のリスト
    param_types: List[Any]
    id: str  # UUID型（strで表現）
    
# クライアントから送られて来たデータ（JSON）の管理
class JSONRequest:
    data: JSONRequestData
    def __init__(self, req_data: str):
        parsed_data = json.loads(req_data)
        
        # UUIDを生成する
        parsed_data["id"] = str(uuid.uuid4())  

        self.data = parsed_data 
