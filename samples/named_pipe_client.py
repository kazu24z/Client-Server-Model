import os
import json

config = json.load(open('config.json'))

f = open(config['file_path'], 'r')

flag = True

while flag:
    if not os.path.exists(config['file_path']):
        flag = False
    
    data = f.read()
    
    if len(data) != 0:
        print(f'Data received from pip: {data}')

f.close()
