FROM scratch
MAINTAINER Nomous

EXPOSE 80

COPY gowebapp /gowebapp
COPY ca-certificates.crt /etc/ssl/certs/

WORKDIR /data
ENTRYPOINT ["/gowebapp"]