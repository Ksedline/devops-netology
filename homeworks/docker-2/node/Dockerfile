FROM node:latest

RUN wget -P /usr/app https://github.com/simplicitesoftware/nodejs-demo/archive/refs/heads/master.zip && cd /usr/app && unzip master.zip

WORKDIR /usr/app/nodejs-demo-master

RUN npm install

EXPOSE 3000
ENTRYPOINT ["npm"]

CMD ["start", "0.0.0.0", "3000"]
