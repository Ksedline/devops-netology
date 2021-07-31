1. Мы выгрузили JSON, который получили через API запрос к нашему сервису:
```{ "info" : "Sample JSON output from our service\t",
    "elements" :[
        { "name" : "first",
        "type" : "server",
        "ip" : 7175 
        },
        { "name" : "second",
        "type" : "proxy",
        "ip : 71.78.22.43
        }
    ]
}
```

Нужно найти и исправить все ошибки, которые допускает наш сервис

Исправленный JSON:

```json
{
  "info": "Sample JSON output from our service \\t",
  "elements": [
    {
      "name": "first",
      "type": "server",
      "ip": "71.75.22.43"
    },
    {
      "name": "second",
      "type": "proxy",
      "ip": "71.78.22.43"
    }
  ]
}
```
```\t``` - управляющий символ, который нужно дополнительно экранировать

elements - это массив, логичнее предположить, что объекты с ключами name, type, ip должн быть одинаковыми. ключ ip - должен быть строкой из-за наличия точек в ip адресе. Во втором элементе массива сломался ключ для ip

2. В прошлый рабочий день мы создавали скрипт, позволяющий опрашивать веб-сервисы и получать их IP. К уже реализованному функционалу нам нужно добавить возможность записи JSON и YAML файлов, описывающих наши сервисы. Формат записи JSON по одному сервису: { "имя сервиса" : "его IP"}. Формат записи YAML по одному сервису: - имя сервиса: его IP. Если в момент исполнения скрипта меняется IP у сервиса - он должен так же поменяться в yml и json файле.

```python python
#!/usr/bin/env python3

from datetime import datetime
from os import path
import time
import socket
import yaml
import json

services = {
  "drive.google.com": "", 
  "mail.google.com": "", 
  "google.com": ""
}

def write_files():
    with open('services.yaml', 'w') as yaml_file:
        yaml.dump(services, yaml_file, default_flow_style=False)
    with open('services.json', 'w') as json_file:
        json.dump(services, json_file, indent=2)
        

while (True):
    for service, known_ip in services.items():
        try:
            new_ip = socket.gethostbyname(service)
            
            if (not path.exists('services.yaml') or not path.exists('services.json')):
                write_files()
            
            if (known_ip != "" and known_ip != new_ip):
                print(f'[ERROR] {service} IP mismatch: {known_ip} {new_ip}')
                write_files()
                
            else:
                print(f'{service} - {new_ip}')

            services[service] = new_ip
                
        except:
            print('Something goes wrong')
    
    time.sleep(10)
```