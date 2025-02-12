import socket
import os
import sys

# ソケットの作成、


class UnixSocket:
    sock: socket.socket
    address: str
    
    def __init__(self, address, is_all_in_one):
        self.sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.address = address
        
        if is_all_in_one:
            try:
                self.bind()
                self.listen(1)
            except OSError as e:
                raise e
    
    def bind(self):
        try:
            try:
                os.unlink(self.address) 
            except FileNotFoundError:
                pass
            self.sock.bind(self.address)
        except OSError as e:
            print(f"Socket error occurred: {e.strerror} (errno={e.errno})")
            sys.exit(1)
        
    def listen(self, n: int=1):
        self.sock.listen(n)
    
    def accept(self):
        conn, client_address = self.sock.accept()
        return conn, client_address
    
