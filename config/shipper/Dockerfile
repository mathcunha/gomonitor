############################################################
# Dockerfile to install Go Monitor
# Based on docker.io/golang Image
############################################################

FROM docker.io/logstash:latest

RUN echo "#!/bin/bash \n sed 's;ELASTIC_PORT_9200_TCP_ADDR;'\"\$ELASTIC_PORT_9200_TCP_ADDR\"';g' /config-dir/logstash.template > /config-dir/logstash.conf" > /entrypoint.sh
RUN echo "\n sed -i 's;REDIS_PORT_6379_TCP_ADDR;'\"\$REDIS_PORT_6379_TCP_ADDR\"';g' /config-dir/logstash.conf" >> /entrypoint.sh
RUN echo "logstash -f /config-dir/logstash.conf" >> /entrypoint.sh

RUN chmod 777 /entrypoint.sh

VOLUME /config-dir

ENTRYPOINT ["/entrypoint.sh"]

MAINTAINER Matheus Cunha <mathcunha@gmail.com>
