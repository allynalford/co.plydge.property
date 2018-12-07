FROM golang:1.10.5
ARG COMMIT_SHA=""

COPY . $GOPATH/src/co.plydge.property
WORKDIR $GOPATH/src/co.plydge.property

# get dependencies
RUN go get github.com/boltdb/bolt
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jmoiron/sqlx
RUN go get gopkg.in/mgo.v2
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/haisum/recaptcha
RUN go get github.com/gorilla/sessions
RUN go get github.com/josephspurrier/csrfbanana
RUN go get github.com/julienschmidt/httprouter
RUN go get github.com/justinas/alice
RUN go get github.com/chromedp/chromedp
RUN go get github.com/PuerkitoBio/goquery

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"'

COPY ca-certificates.crt /etc/ssl/certs/

EXPOSE 80

ENV COMMIT_SHA=${COMMIT_SHA}
ENTRYPOINT ["./co.plydge.property"]
