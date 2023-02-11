FROM golang
# :alpine AS builder

# ENV GO111MOduLE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux

WORKDIR /build

COPY . ./

RUN go mod tidy

RUN ln -sf /usr/share/zoneinfo/Asia/Seoul /etc/localtime

# RUN go build -o main .

# WORKDIR /dist

# RUN cp /build/main .

# FROM scratch

# COPY --from=builder /dist/main .

# COPY config/config.toml config/

# ENTRYPOINT ["./main"]



