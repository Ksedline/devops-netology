FROM centos:7

ADD https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.13.4-linux-x86_64.tar.gz /
ADD https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.13.4-linux-x86_64.tar.gz.sha512 /

RUN yum update -y && \
  yum install perl-Digest-SHA -y && \
  shasum -a 512 -c elasticsearch-7.13.4-linux-x86_64.tar.gz.sha512 && \
  tar -xzf elasticsearch-7.13.4-linux-x86_64.tar.gz && \
  cd elasticsearch-7.13.4/ && \
  useradd elasticuser && \
  chown -R elasticuser:elasticuser /elasticsearch-7.13.4/ && \
  rm -rf /elasticsearch-7.13.4-linux-x86_64.tar.gz.sha512 /elasticsearch-7.13.4-linux-x86_64.tar.gz

RUN mkdir /var/lib/{data,logs} && \
  chown -R elasticuser:elasticuser /var/lib/data && \
  chown -R elasticuser:elasticuser /var/lib/logs

WORKDIR /elasticsearch-7.13.4

RUN mkdir snapshots && \
  chown -R elasticuser:elasticuser snapshots

COPY elasticsearch.yml /elasticsearch-7.13.4/config/

RUN chown -R elasticuser:elasticuser /elasticsearch-7.13.4/config

USER elasticuser

EXPOSE 9200 9300

CMD ["./bin/elasticsearch", "-Ecluster.name=netology_cluster", "-Enode.name=netology_test"]
