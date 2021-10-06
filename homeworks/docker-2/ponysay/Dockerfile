FROM archlinux:latest

RUN pacman -Sy && \
pacman -S ponysay --noconfirm
 
ENTRYPOINT ["/usr/bin/ponysay"]
CMD ["Hey, netology”]
