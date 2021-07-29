1. Какой системный вызов делает команда `cd`? В прошлом ДЗ мы выяснили, что `cd` не является самостоятельной программой, 
   это `shell builtin`, поэтому запустить `strace` непосредственно на `cd` не получится. Тем не менее вы можете запустить `strace` на `/bin/bash -c 'cd /tmp'`. 
   В этом случае вы увидите полный список системных вызовов, которые делает сам `bash` при старте. Вам нужно найти тот единственный, который относится именно к `cd`.
   

`cd` использует системный вызов `chdir("/tmp")`

```bash
vagrant@vagrant:~$ strace /bin/bash -c 'cd /tmp' 2>&1 | grep tmp
execve("/bin/bash", ["/bin/bash", "-c", "cd /tmp"], 0x7ffd1ea44130 /* 25 vars */) = 0
stat("/tmp", {st_mode=S_IFDIR|S_ISVTX|0777, st_size=4096, ...}) = 0
stat("/tmp", {st_mode=S_IFDIR|S_ISVTX|0777, st_size=4096, ...}) = 0
stat("/tmp", {st_mode=S_IFDIR|S_ISVTX|0777, st_size=4096, ...}) = 0
chdir("/tmp")                           = 0
```

2. Попробуйте использовать команду `file` на объекты разных типов на файловой системе. Например:
    ```bash
    vagrant@netology1:~$ file /dev/tty
    /dev/tty: character special (5/0)
    vagrant@netology1:~$ file /dev/sda
    /dev/sda: block special (8/0)
    vagrant@netology1:~$ file /bin/bash
    /bin/bash: ELF 64-bit LSB shared object, x86-64
    ```
   Используя `strace` выясните, где находится база данных `file` на основании которой она делает свои догадки.
   

База данных располагается по адресу `/usr/share/misc/magic.mgc`. 
До этого проверяются пользовательские файлы `/home/vagrant/.magic.mgc`, но из не существуют.

```bash
strace -e trace=file /bin/bash -c 'file /dev/tty' 2>&1 | tail -17
---
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libmagic.so.1", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libc.so.6", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/liblzma.so.5", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libbz2.so.1.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libz.so.1", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libpthread.so.0", O_RDONLY|O_CLOEXEC) = 3
openat(AT_FDCWD, "/usr/lib/locale/locale-archive", O_RDONLY|O_CLOEXEC) = 3
stat("/home/vagrant/.magic.mgc", 0x7ffce9405b40) = -1 ENOENT (No such file or directory)
stat("/home/vagrant/.magic", 0x7ffce9405b40) = -1 ENOENT (No such file or directory)
openat(AT_FDCWD, "/etc/magic.mgc", O_RDONLY) = -1 ENOENT (No such file or directory)
stat("/etc/magic", {st_mode=S_IFREG|0644, st_size=111, ...}) = 0
openat(AT_FDCWD, "/etc/magic", O_RDONLY) = 3
openat(AT_FDCWD, "/usr/share/misc/magic.mgc", O_RDONLY) = 3
openat(AT_FDCWD, "/usr/lib/x86_64-linux-gnu/gconv/gconv-modules.cache", O_RDONLY) = 3
lstat("/dev/tty", {st_mode=S_IFCHR|0666, st_rdev=makedev(0x5, 0), ...}) = 0
/dev/tty: character special (5/0)
+++ exited with 0 +++
```

3. Предположим, приложение пишет лог в текстовый файл. Этот файл оказался удален (deleted в lsof), однако возможности 
   сигналом сказать приложению переоткрыть файлы или просто перезапустить приложение нет. Так как приложение продолжает 
   писать в удаленный файл, место на диске постепенно заканчивается. Основываясь на знаниях о перенаправлении потоков, 
   предложите способ обнуления открытого удаленного файла (чтобы освободить место на файловой системе).

Нужно получить файловый дескриптор для pid где происходит запись в файл, а потом в дескриптор перенаправить пустую строку

```
ps au | grep nano # нашли что pid 10000

ll /proc/10000/fd # нашли номер файлового дескриптора

/proc/10000/fd$ echo '' > /proc/10000/fd/5
```

4. Занимают ли зомби-процессы какие-то ресурсы в ОС (CPU, RAM, IO)?

Зомби-процессы - не потребляют ресурсов, но блокируют записи в таблице процессов, размер которой ограничен для каждого пользователя и системы в целом

5. В iovisor BCC есть утилита `opensnoop`:
    ```bash
    root@vagrant:~# dpkg -L bpfcc-tools | grep sbin/opensnoop
    /usr/sbin/opensnoop-bpfcc
    ```
   На какие файлы вы увидели вызовы группы `open` за первую секунду работы утилиты? Воспользуйтесь пакетом `bpfcc-tools` 
   для Ubuntu 20.04. Дополнительные [сведения по установке](https://github.com/iovisor/bcc/blob/master/INSTALL.md).

Установка:

```sudo apt-get install bpfcc-tools linux-headers-$(uname -r)```

Результат:

```
vagrant@vagrant:~$ dpkg -L bpfcc-tools | grep sbin/opensnoop
/usr/sbin/opensnoop-bpfcc
vagrant@vagrant:~$ sudo /usr/sbin/opensnoop-bpfcc
PID    COMM               FD ERR PATH
1      systemd            12   0 /proc/403/cgroup
773    vminfo              4   0 /var/run/utmp
593    dbus-daemon        -1   2 /usr/local/share/dbus-1/system-services
593    dbus-daemon        18   0 /usr/share/dbus-1/system-services
593    dbus-daemon        -1   2 /lib/dbus-1/system-services
593    dbus-daemon        18   0 /var/lib/snapd/dbus-1/system-services/
773    vminfo              4   0 /var/run/utmp
593    dbus-daemon        -1   2 /usr/local/share/dbus-1/system-services
593    dbus-daemon        18   0 /usr/share/dbus-1/system-services
593    dbus-daemon        -1   2 /lib/dbus-1/system-services
593    dbus-daemon        18   0 /var/lib/snapd/dbus-1/system-services/
604    irqbalance          6   0 /proc/interrupts
604    irqbalance          6   0 /proc/stat
604    irqbalance          6   0 /proc/irq/20/smp_affinity
604    irqbalance          6   0 /proc/irq/0/smp_affinity
604    irqbalance          6   0 /proc/irq/1/smp_affinity
604    irqbalance          6   0 /proc/irq/8/smp_affinity
604    irqbalance          6   0 /proc/irq/12/smp_affinity
604    irqbalance          6   0 /proc/irq/14/smp_affinity
604    irqbalance          6   0 /proc/irq/15/smp_affinity
```

6. Какой системный вызов использует `uname -a`? Приведите цитату из man по этому системному вызову, где описывается 
   альтернативное местоположение в `/proc`, где можно узнать версию ядра и релиз ОС.

```bash
vagrant@vagrant:~$ strace uname -a

uname({sysname="Linux", nodename="vagrant", ...}) = 0
fstat(1, {st_mode=S_IFCHR|0620, st_rdev=makedev(0x88, 0), ...}) = 0
uname({sysname="Linux", nodename="vagrant", ...}) = 0
uname({sysname="Linux", nodename="vagrant", ...}) = 0

```  

```man 2 uname | grep -C 10 proc```

Цитата: ```Part  of the utsname information is also accessible via /proc/sys/kernel/{ostype, hostname, osre‐
lease, version, domainname}.```

7. Чем отличается последовательность команд через `;` и через `&&` в bash? Например:
    ```bash
    root@netology1:~# test -d /tmp/some_dir; echo Hi
    Hi
    root@netology1:~# test -d /tmp/some_dir && echo Hi
    root@netology1:~#
    ```
   Есть ли смысл использовать в bash `&&`, если применить `set -e`?


При `;` команды выполняются последовательно, независимо от результата работы предыдущей.

При `&&` следующая команда выполнится только, если предшествующая завершится успешно.  


shell при `set -e` не завершается в том случае, если выражение, вернуло ненулевой статус. Нет смысла использовать `&&` после вызова `set -e`

8. Из каких опций состоит режим bash `set -euxo pipefail` и почему его хорошо было бы использовать в сценариях?

```man bash | grep pipefail```

`-e` - завершить сразу, если pipeline, список команд или составная команда завершится с ненулевым статусом, если эта команда выполняется в `while` или `until`, проверяется в условиях `if` или `elif`, является частью команд в списке с `&&` или `||` (кроме последней), её код завершения инвертируется. 

`-u` - считать неустановленные переменные и параметры, кроме @ и астерикса *, ошибкой

`-x` - после каждой команды отражать расширенную информацию с аргументами и трейсом

`-o` - возвращаемое значение pipeline. Если ненулевой код завершения самой последней команды или нулевой а также, если все команды завершились успешно.  

Опцию pipefail хорошо использовать в сценариях, из-за дополнительного логирования, а также завершения сценария при наличии ошибок. Бонусом неустановленные переменные считает ошибками.

9. Используя `-o stat` для `ps`, определите, какой наиболее часто встречающийся статус у процессов в системе. 
   В `man ps` ознакомьтесь (`/PROCESS STATE CODES`), что значат дополнительные к основной заглавной буквы статуса процессов. 
   Его можно не учитывать при расчете (считать S, Ss или Ssl равнозначными).
   
`ps a -o stat`

В основном встречается статус S. S это спящие, ожидающие, процессы, которые можно прервать. Дополнительные буквы статуса - это приоритет, многопоточность
