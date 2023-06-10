# pets-next-door-api

이웃집멍냥 백엔드 API 서버입니다.

## How to run

개발 환경에서는 docker-compose를 활용한 Postgres DB를 제공합니다.

```bash
$ docker-compose up -d # Postgres 컨테이너 실행
$ docker-compose down -v # Postgres 컨테이너 중지 및 볼륨 삭제
```

개발 환경에서 실행하는 방법입니다.

```bash
$ go run cmd/server/main.go
```

배포 환경에서 빌드 및 실행하는 방법입니다.

```bash
$ go build -o main cmd/server/main.go
$ ./main
```

Production 환경을 위한 Dockerfile도 제공합니다.

```bash
$ docker build -t pets-next-door-api .
$ docker run -p 8080:8080 pets-next-door-api
```

## How to test

```bash
$ go test ./...
```
