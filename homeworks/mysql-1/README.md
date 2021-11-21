# Домашнее задание к занятию "6.3. MySQL"

## Введение

Перед выполнением задания вы можете ознакомиться с 
[дополнительными материалами](https://github.com/netology-code/virt-homeworks/tree/master/additional/README.md).

## Задача 1

- **Используя docker поднимите инстанс MySQL (версию 8). Данные БД сохраните в volume.**

```
version: "3.1"

services:
  mysqldb:
    container_name: mysql-1
    image: mysql:8
    restart: always
    environment:
      - MYSQL_ROOT_PASSWORD=drivemecrazy
    volumes:
      - "/home/dockeruser/docker/mysql-1/data:/var/lib/mysql/"
      - "/home/dockeruser/docker/mysql-1/backups:/var/lib/backups/"
    ports:
      - "3306:3306"
    command: mysqld --default-authentication-plugin=mysql_native_password --character-set-server=utf8 --collation-server=utf8_general_ci

```

- **Изучите [бэкап БД](https://github.com/netology-code/virt-homeworks/tree/master/06-db-03-mysql/test_data) и 
восстановитесь из него.**

```
docker exec -it mysql-1 mysql -u root -p
mysql> create database test_db;
Query OK, 1 row affected (0.01 sec)

docker exec -it mysql-1 bash
mysql -u root -p test_db < /var/lib/backups/test_dump.sql
```

- **Перейдите в управляющую консоль `mysql` внутри контейнера.**

```
docker exec -it mysql-1 mysql -u root -p
```

- Используя команду `\h` получите список управляющих команд. Найдите команду для выдачи статуса БД и **приведите в ответе** из ее вывода версию сервера БД.

```
mysql> \s
--------------
mysql  Ver 8.0.25 for Linux on x86_64 (MySQL Community Server - GPL)
...
Server version:         8.0.25 MySQL Community Server - GPL
...
```

- Подключитесь к восстановленной БД и получите список таблиц из этой БД.

```
mysql> use test_db;

mysql> show tables;
+-------------------+
| Tables_in_test_db |
+-------------------+
| orders            |
+-------------------+
1 row in set (0.00 sec)
```

- **Приведите в ответе** количество записей с `price` > 300.

```
mysql> SELECT COUNT(*) `cnt` FROM `orders` WHERE `price` > 300;
+-----+
| cnt |
+-----+
|   1 |
+-----+
1 row in set (0.00 sec)
```

## Задача 2

Создайте пользователя test в БД c паролем test-pass, используя:
- плагин авторизации mysql_native_password
- срок истечения пароля - 180 дней 
- количество попыток авторизации - 3 
- максимальное количество запросов в час - 100
- аттрибуты пользователя:
    - Фамилия "Pretty"
    - Имя "James"

```
CREATE USER IF NOT EXISTS 'test'@'localhost' IDENTIFIED WITH mysql_native_password BY 'test-pass'
WITH MAX_CONNECTIONS_PER_HOUR 100
PASSWORD EXPIRE INTERVAL 180 DAY
FAILED_LOGIN_ATTEMPTS 3
PASSWORD_LOCK_TIME 1
ATTRIBUTE '{"lastname":"Pretty", "name":"James"}';
```

- Предоставьте привелегии пользователю `test` на операции SELECT базы `test_db`.

`GRANT SELECT ON test_db.* TO 'test'@'localhost';`
    
- Используя таблицу INFORMATION_SCHEMA.USER_ATTRIBUTES получите данные по пользователю `test` и **приведите в ответе к задаче**.

```
mysql> SELECT * FROM INFORMATION_SCHEMA.USER_ATTRIBUTES WHERE USER = 'test';
+------+-----------+-----------------------------------------+
| USER | HOST      | ATTRIBUTE                               |
+------+-----------+-----------------------------------------+
| test | localhost | {"name": "James", "lastname": "Pretty"} |
+------+-----------+-----------------------------------------+
1 row in set (0.00 sec)
```

## Задача 3

Установите профилирование `SET profiling = 1`.
Изучите вывод профилирования команд `SHOW PROFILES;`.

Исследуйте, какой `engine` используется в таблице БД `test_db` и **приведите в ответе**.

Измените `engine` и **приведите время выполнения и запрос на изменения из профайлера в ответе**:
- на `MyISAM`
- на `InnoDB`

```
mysql> SET profiling = 1;
Query OK, 0 rows affected, 1 warning (0.00 sec)

mysql> SHOW PROFILES;
Empty set, 1 warning (0.00 sec)

mysql> SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES where TABLE_SCHEMA = 'test_db';
+------------+--------+
| TABLE_NAME | ENGINE |
+------------+--------+
| orders     | InnoDB |
+------------+--------+
1 row in set (0.00 sec)

mysql> ALTER TABLE orders ENGINE = MyISAM;
Query OK, 5 rows affected (0.02 sec)
Records: 5  Duplicates: 0  Warnings: 0

mysql> ALTER TABLE orders ENGINE = InnoDB;
Query OK, 5 rows affected (0.03 sec)
Records: 5  Duplicates: 0  Warnings: 0

mysql> SHOW PROFILES;
+----------+------------+-----------------------------------------------------------------------------------------+
| Query_ID | Duration   | Query                                                                                   |
+----------+------------+-----------------------------------------------------------------------------------------+
|        1 | 0.00893050 | SELECT * FROM information_schema.TABLES where TABLE_SCHEMA = 'test_db'                  |
|        2 | 0.00151925 | SELECT TABLE_NAME, ENGINE FROM information_schema.TABLES where TABLE_SCHEMA = 'test_db' |
|        3 | 0.02233201 | ALTER TABLE orders ENGINE = MyISAM                                                      |
|        4 | 0.02500800 | ALTER TABLE orders ENGINE = InnoDB                                                      |
+----------+------------+-----------------------------------------------------------------------------------------+
4 rows in set, 1 warning (0.00 sec)
```

## Задача 4 

Изучите файл `my.cnf` в директории /etc/mysql.

Измените его согласно ТЗ (движок InnoDB):
- Скорость IO важнее сохранности данных
- Нужна компрессия таблиц для экономии места на диске
- Размер буффера с незакомиченными транзакциями 1 Мб
- Буффер кеширования 30% от ОЗУ
- Размер файла логов операций 100 Мб

Приведите в ответе измененный файл `my.cnf`.

```
[mysqld]
pid-file        = /var/run/mysqld/mysqld.pid
socket          = /var/run/mysqld/mysqld.sock
datadir         = /var/lib/mysql
secure-file-priv= NULL

# Custom config should go here
!includedir /etc/mysql/conf.d/

innodb_flush_log_at_trx_commit = 0 # Скорость IO важнее сохранности данных
innodb_file_per_table = ON # Нужна компрессия таблиц для экономии места на диске
innodb_log_buffer_size = 1M # Размер буффера с незакомиченными транзакциями 1 Мб
innodb_buffer_pool_size = 1G # Буффер кеширования 30% от ОЗУ
innodb_log_file_size = 100M # Размер файла логов операций 100 Мб
```