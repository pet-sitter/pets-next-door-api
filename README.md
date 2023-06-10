# pets-next-door-api

이웃집멍냥 백엔드 API 서버입니다.

## How to run

개발 환경에서 실행하는 방법입니다.

```bash
$ go run cmd/server/main.go
```

배포 환경에서 빌드 및 실행하는 방법입니다.

```bash
$ go build -o main cmd/server/main.go
$ ./main
```

## How to test

```bash
$ go test ./...
```
