1. ipvs. Если при запросе на VIP сделать подряд несколько запросов (например, for i in {1..50}; do curl -I -s 172.28.128.200>/dev/null; done ), ответы будут получены почти мгновенно. Тем не менее, в выводе ipvsadm -Ln еще некоторое время будут висеть активные InActConn. Почему так происходит?

Балансировщик ipvs не видит прямые ответы от backend серверов клиентам, т.к. они идут напрямую минуя балансировщик. 
Поэтому ipvs не знает - ответил ли backend клиенту на запрос или нет. 
Из-за этого TCP соединение может висеть еще какое-то время. 
Эти соединения исчезают с балансировщика после timeout

2. На лекции мы познакомились отдельно с ipvs и отдельно с keepalived. Воспользовавшись этими знаниями, совместите технологии вместе (VIP должен подниматься демоном keepalived). Приложите конфигурационные файлы, которые у вас получились, и продемонстрируйте работу получившейся конструкции. Используйте для директора отдельный хост, не совмещая его с риалом! Подобная схема возможна, но выходит за рамки рассмотренного на лекции.


balancer1 keepalived.conf

```

root@balancer1:/home/vagrant# cat /etc/keepalived/keepalived.conf 
vrrp_instance VI_1 {
state MASTER      
interface eth1    
virtual_router_id 33
priority 100
advert_int 1

  virtual_ipaddress {       
    172.28.128.200/24 dev eth1
  }
}
```

balanser2 keepalived.conf

```
root@balancer2:/home/vagrant# cat /etc/keepalived/keepalived.conf 
vrrp_instance VI_1 {
state BACKUP
interface eth1    
virtual_router_id 33
priority 50
advert_int 1

  virtual_ipaddress {       
    172.28.128.200/24 dev eth1
  }    
}
```

отключим eth1 на balancer1 для проверки подключения:

```
root@balancer1:/home/vagrant# ip link set eth1 down
root@balancer1:/home/vagrant# cat /var/log/syslog | grep VI_1
Aug 8 20:20:18 vagrant Keepalived_vrrp[54462]: (VI_1) Entering BACKUP STATE (init)
Aug 8 20:20:22 vagrant Keepalived_vrrp[54462]: (VI_1) Entering MASTER STATE
Aug 8 20:21:45 vagrant Keepalived_vrrp[54462]: (VI_1) Entering FAULT STATE
Aug 8 20:21:45 vagrant Keepalived_vrrp[54462]: (VI_1) sent 0 priority
---
root@balancer2:/home/vagrant# cat /var/log/syslog | grep VI_1
Aug 8 20:31:43 vagrant Keepalived_vrrp[54389]: (VI_1) Entering BACKUP STATE (init)
Aug 8 20:31:48 vagrant Keepalived_vrrp[54389]: (VI_1) Entering MASTER STATE
---
vagrant@client:~$ for i in {1..50}; do curl -I -s 172.28.128.200>/dev/null; done
---
root@balancer2:/home/vagrant# ipvsadm -Ln --stats
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port               Conns   InPkts  OutPkts  InBytes OutBytes
  -> RemoteAddress:Port
TCP  172.28.128.200:80                  50      300        0    19950        0
  -> 172.28.128.30:80                   25      150        0     9975        0
  -> 172.28.128.40:80                   25      150        0     9975        0
---
root@balancer1:/home/vagrant# ipvsadm -Ln --stats
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port               Conns   InPkts  OutPkts  InBytes OutBytes
  -> RemoteAddress:Port
TCP  172.28.128.200:80                   0        0        0        0        0
  -> 172.28.128.30:80                    0        0        0        0        0
  -> 172.28.128.40:80                    0        0        0        0        0
---
```

Запустим eth1 на balancer1:

```
root@balancer1:/home/vagrant# ip link set eth1 up
root@balancer1:/home/vagrant# cat /var/log/syslog | grep VI_1
...
Aug 8 21:02:52 vagrant Keepalived_vrrp[54462]: (VI_1) Entering MASTER STATE
---
root@balancer2:/home/vagrant# cat /var/log/syslog | grep VI_1
Aug 8 21:31:43 vagrant Keepalived_vrrp[54389]: (VI_1) Entering BACKUP STATE (init)
Aug 8 21:40:48 vagrant Keepalived_vrrp[54389]: (VI_1) Entering MASTER STATE
root@balancer2:/home/vagrant# cat /var/log/syslog | grep VI_1
Aug 8 21:31:43 vagrant Keepalived_vrrp[54389]: (VI_1) Entering BACKUP STATE (init)
Aug 8 21:40:48 vagrant Keepalived_vrrp[54389]: (VI_1) Entering MASTER STATE
Aug 8 21:45:52 vagrant Keepalived_vrrp[54389]: (VI_1) Master received advert from 172.28.128.10 with higher priority 100, ours 50
Aug 8 21:45:52 vagrant Keepalived_vrrp[54389]: (VI_1) Entering BACKUP STATE
---
vagrant@client:~$ for i in {1..50}; do curl -I -s 172.28.128.200>/dev/null; done
---
root@balancer1:/home/vagrant# ipvsadm -Ln --stats
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port               Conns   InPkts  OutPkts  InBytes OutBytes
  -> RemoteAddress:Port
TCP  172.28.128.200:80                  50      300        0    19950        0
  -> 172.28.128.30:80                   25      150        0     9975        0
  -> 172.28.128.40:80                   25      150        0     9975        0
```

Оба balancer работают!


3. В лекции мы использовали только 1 VIP адрес для балансировки. У такого подхода несколько отрицательных моментов, один из которых – невозможность активного использования нескольких хостов (1 адрес может только переехать с master на standby). Подумайте, сколько адресов оптимально использовать, если мы хотим без какой-либо деградации выдерживать потерю 1 из 3 хостов при входящем трафике 1.5 Гбит/с и физических линках хостов в 1 Гбит/с? Предполагается, что мы хотим задействовать 3 балансировщика в активном режиме (то есть не 2 адреса на 3 хоста, один из которых в обычное время простаивает).

3 балансировщика в активном режиме должны использавать минимум 3 адреса. При входящем трафике 1,5 Гбит/с отказ одного из балансировщиков не приведет к замедлению обслуживания запросов. Таким образом, ответ: 3.
