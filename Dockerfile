FROM mongo

WORKDIR /data/db

COPY ./pools.json ../

# RUN mongoimport -d syeong-db -c pools --jsonArray --file="/data/pools.json"

# RUN mongoimport --uri="mongodb://localhost:27017/syeong-db" --collection=pools --jsonArray --file="/data/db/pools.json"

# RUN mongod

# ENTRYPOINT [ "mongosh" ]

# CMD mongoimport -d syeong-db -c pools --jsonArray --file="/data/db/pools.json"


# 이미지 생성
# docker build -t syeong-mongo .

# 컨테이너 생성, 백그라운드 실행 - 처음 컨테이너 빌드 할 때
# docker run --name syeong-mongo -v $(pwd)/docker/mongo/data/db:/data/db -d -p 27017:27017 syeong-mongo

# 컨테이너 종료시
# docker stop syeong-mongo

# 컨테이너 삭제시
# docker rm syeong-mongo

# 빌드되어있는 컨테이너 이름으로 재시작
# docker restart syeong-server