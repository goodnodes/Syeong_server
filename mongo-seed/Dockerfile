# FROM golang
# :alpine AS builder

# ENV GO111MOduLE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux

# WORKDIR /build

# COPY . ./

# RUN go mod tidy

# ENTRYPOINT [ "go", "run", "main.go" ]

# RUN go build -o main .

# WORKDIR /dist

# RUN cp /build/main .

# FROM scratch

# COPY --from=builder /dist/main .

# COPY config/config.toml config/

# ENTRYPOINT ["./main"]





FROM mongo

COPY ./pools.json ./

# RUN mongoimport -d syeong-db -c pools --jsonArray --file="/data/pools.json"

COPY ./user.json ./

CMD mongoimport --uri "mongodb://mongodb:27017" -d syeong-db -c pools --jsonArray --file="./pools.json"


# 이미지 생성, 컨테이너 생성, 백그라운드 실행 - 처음 컨테이너 빌드 할 때
# docker run --name syeong-db -v $(pwd)/docker/mongo/data/db:/data/db -d -p 27017:27017 shyun96/syeong-db

# 도커 쉘 접속
# docker exec -it syeong-db bash

# 아래 두 녀석은 처음에 한 번씩만 실행해줌
# pools 데이터 추가하기
# mongoimport -d syeong-db -c pools --jsonArray --file="/data/pools.json"
# user 데이터 추가하기
# mongoimport -d syeong-db -c users --file="/data/user.json"

# 빠져나오기
# exit

# 컨테이너 종료시
# docker stop syeong-db

# 컨테이너 삭제시
# docker rm syeong-db

# 빌드되어있는 컨테이너 이름으로 재시작
# docker restart syeong-db



# RUN mongoimport -d syeong-db -c pools --jsonArray --file="/data/pools.json"

# RUN mongoimport --uri="mongodb://localhost:27017/syeongDB" --collection=pools --jsonArray --file="/data/db/pools.json"

# RUN mongod

# ENTRYPOINT [ "mongosh" ]

# CMD mongoimport -d syeong-db -c pools --jsonArray --file="/data/db/pools.json"


# 이미지 생성
# docker build -t syeong-mongo .

