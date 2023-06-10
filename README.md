# pets-next-door-api

이웃집멍냥 백엔드 API 서버입니다.

## Setup

.env.template 파일을 참고하여 .env 파일을 루트 디렉토리에 생성합니다.

.env는 로컬 환경에서만 사용되며, .gitignore에 등록되어 있습니다.

```bash
$ cp .env.template .env
```

추가로, 파이어베이스 프로젝트의 서비스 계정 키를 다운로드하여 `firebase-credentials.json` 파일로 루트 디렉토리에 저장합니다.

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
