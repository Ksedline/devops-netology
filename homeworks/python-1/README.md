# Домашнее задание к занятию "4.2. Использование Python для решения типовых DevOps задач"

## Обязательные задания

1. Есть скрипт:
	```python
    #!/usr/bin/env python3
	a = 1
	b = '2'
	c = a + b
	```
	* Какое значение будет присвоено переменной c?
	* Как получить для переменной c значение 12?
	* Как получить для переменной c значение 3?

* a + b в этом случае не выполнится, ошибка `unsupported operand type(s) for +: 'int' and 'str'`

```python
	c = str(a) + b # 12
```

```python
	c = a + int(b) # 3
```

2. Мы устроились на работу в компанию, где раньше уже был DevOps Engineer. Он написал скрипт, позволяющий узнать, какие файлы модифицированы в репозитории, относительно локальных изменений. Этим скриптом недовольно начальство, потому что в его выводе есть не все изменённые файлы, а также непонятен полный путь к директории, где они находятся. Как можно доработать скрипт ниже, чтобы он исполнял требования вашего руководителя?

	```python
    #!/usr/bin/env python3

    import os

	bash_command = ["cd ~/netology/sysadm-homeworks", "git status"]
	result_os = os.popen(' && '.join(bash_command)).read()
    is_change = False
	for result in result_os.split('\n'):
        if result.find('modified') != -1:
            prepare_result = result.replace('\tmodified:   ', '')
            print(prepare_result)
            break

	```

  Доработка:

```python
  import os

  bash_command = ["cd ~/netology/sysadm-homeworks", "git status"]
  popen = os.popen(' && '.join(bash_command))
  result_os = popen.read()
  popen.close()
  for result in result_os.split('\n'):
    if result.find('modified:') != -1:
        prepare_result = result.replace('\tmodified:   ', '')
        print(f'{os.path.join(os.getcwd(), prepare_result)}')
```

3. Доработать скрипт выше так, чтобы он мог проверять не только локальный репозиторий в текущей директории, а также умел воспринимать путь к репозиторию, который мы передаём как входной параметр. Мы точно знаем, что начальство коварное и будет проверять работу этого скрипта в директориях, которые не являются локальными репозиториями.

```python
  #!/usr/bin/env python3
  import os
  import sys
  
  directory = os.getcwd()

  if len(sys.argv) >= 2: # если есть второй входной аргумент
    directory = sys.argv[1]

  bash_command = [f"cd {directory}", "cd ~/netology/sysadm-homeworks", "git status"]
  popen = os.popen(' && '.join(bash_command))
  result_os = popen.read()
  popen.close()
  for result in result_os.split('\n'):
    if result.find('modified:') != -1:
        prepare_result = result.replace('\tmodified:   ', '')
        print(f'{os.path.join(os.getcwd(), prepare_result)}')
```

4. Наша команда разрабатывает несколько веб-сервисов, доступных по http. Мы точно знаем, что на их стенде нет никакой балансировки, кластеризации, за DNS прячется конкретный IP сервера, где установлен сервис. Проблема в том, что отдел, занимающийся нашей инфраструктурой очень часто меняет нам сервера, поэтому IP меняются примерно раз в неделю, при этом сервисы сохраняют за собой DNS имена. Это бы совсем никого не беспокоило, если бы несколько раз сервера не уезжали в такой сегмент сети нашей компании, который недоступен для разработчиков. Мы хотим написать скрипт, который опрашивает веб-сервисы, получает их IP, выводит информацию в стандартный вывод в виде: <URL сервиса> - <его IP>. Также, должна быть реализована возможность проверки текущего IP сервиса c его IP из предыдущей проверки. Если проверка будет провалена - оповестить об этом в стандартный вывод сообщением: [ERROR] <URL сервиса> IP mismatch: <старый IP> <Новый IP>. Будем считать, что наша разработка реализовала сервисы: drive.google.com, mail.google.com, google.com.


```python
 #!/usr/bin/env python3
 
from datetime import datetime
import time
import socket

services = {
  "drive.google.com": "", 
  "mail.google.com": "", 
  "google.com": ""
}

while (True):
    for service, known_ip in services.items():
        try:
            new_ip = socket.gethostbyname(service)
            
            if (known_ip != "" and known_ip != new_ip):
                print(f'[ERROR] {service} IP mismatch: {known_ip} {new_ip}')
            else:
                print(f'{service} - {new_ip}')

            services[service] = new_ip

        except:
            print('Something goes wrong')

    print(datetime.now(), services)
    time.sleep(60)
```
