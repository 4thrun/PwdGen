# build 
FROM golang:alpine AS build-env 
WORKDIR /
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"
COPY . .
RUN go build -o pwdgen .

# run
FROM scratch
COPY --from=build-env /pwdgen /
ENTRYPOINT [ "/pwdgen" ]