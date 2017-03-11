FROM golang:stretch
RUN apt update 
RUN apt upgrade -y
RUN go get .
RUN go get github.com/Masterminds/glide
RUN go get github.com/gorilla/rpc
RUN glide install
RUN go build .
RUN go install .
EXPOSE 8000
CMD ["glogchain"]
