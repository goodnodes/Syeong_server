FROM golang
# :alpine AS builder

# ENV GO111MOduLE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux

WORKDIR /build

COPY . ./

RUN go mod tidy

# ENTRYPOINT [ "go", "run", "main.go" ]

# RUN go build -o main .

# WORKDIR /dist

# RUN cp /build/main .

# FROM scratch

# COPY --from=builder /dist/main .

# COPY config/config.toml config/

# ENTRYPOINT ["./main"]



