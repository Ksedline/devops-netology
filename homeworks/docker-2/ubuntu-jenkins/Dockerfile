FROM ubuntu:latest

RUN apt-get update && apt-get install -y openjdk-11-jdk

ADD https://get.jenkins.io/war-stable/latest/jenkins.war /root/
WORKDIR /root
EXPOSE 8080
ENTRYPOINT ["java"]
CMD ["-jar","jenkins.war"]
