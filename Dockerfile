FROM ubuntu

EXPOSE 9090

WORKDIR /usr/src/

COPY ./serv/ /usr/src/serv

COPY ./client/ /usr/src/client

#CMD ["./serv/serv", "&"]

#CMD ["echo", "starting client"]

#CMD ["./client/client"]


