FROM golang:1.10.5
ARG COMMIT_SHA=""

# Install Chrome
RUN apt-get update && apt-get install -y \
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg \
	hicolor-icon-theme \
	libcanberra-gtk* \
	libgl1-mesa-dri \
	libgl1-mesa-glx \
	libpango1.0-0 \
	libpulse0 \
	libv4l-0 \
	fonts-symbola \
	--no-install-recommends \
	&& curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
	&& echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list \
	&& apt-get update && apt-get install -y \
	google-chrome-stable \
	--no-install-recommends \
	&& apt-get purge --auto-remove -y curl \
	&& rm -rf /var/lib/apt/lists/*# Install Chrome
RUN apt-get update && apt-get install -y \
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg \
	hicolor-icon-theme \
	libcanberra-gtk* \
	libgl1-mesa-dri \
	libgl1-mesa-glx \
	libpango1.0-0 \
	libpulse0 \
	libv4l-0 \
	fonts-symbola \
	--no-install-recommends \
	&& curl -sSL https://dl.google.com/linux/linux_signing_key.pub | apt-key add - \
	&& echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google.list \
	&& apt-get update && apt-get install -y \
	google-chrome-stable \
	--no-install-recommends \
	&& apt-get purge --auto-remove -y curl \
	&& rm -rf /var/lib/apt/lists/*

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
RUN go get github.com/ivpusic/grpool

# build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"'

COPY ca-certificates.crt /etc/ssl/certs/

EXPOSE 80 443

ENV COMMIT_SHA=${COMMIT_SHA}
ENTRYPOINT ["./co.plydge.property"]
