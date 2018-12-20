FROM ubuntu

EXPOSE 9090

WORKDIR /usr/src/

COPY ./serv/ /usr/src/serv

COPY ./client/ /usr/src/client

COPY ./docker-entrypoint.sh /usr/local/bin

ENTRYPOINT docker-entrypoint.sh


