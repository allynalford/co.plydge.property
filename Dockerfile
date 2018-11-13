FROM scratch
MAINTAINER Nomous

EXPOSE 80

COPY co.plydge.property /co.plydge.property
COPY ca-certificates.crt /etc/ssl/certs/

WORKDIR /data
ENTRYPOINT ["/gowebapp"]