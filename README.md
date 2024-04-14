# pets-next-door-api

이웃집멍냥 백엔드 API 서버입니다.

## Setup

`.env.template` 파일을 참고하여 `.env` 파일을 루트 디렉토리에 생성합니다.

`.env`는 로컬 환경에서만 사용되며, `.gitignore`에 등록되어 있습니다.

```bash
$ cp .env.template .env
```

필요에 따라 `swagger-gen.sh` 파일에 exec 권한을 부여합니다.

```bash
$ chmod +x ./scripts/swagger-gen.sh
```

추가로, 파이어베이스 프로젝트의 서비스 계정 키를 다운로드하여 `firebase-credentials.json` 파일로 루트 디렉토리에 저장합니다.

## How to run

개발 환경에서는 docker-compose를 활용한 Postgres DB를 제공합니다.

```bash
$ make db:up # Postgres 컨테이너 실행
$ make db:down # Postgres 컨테이너 중지
$ make db:destroy # Postgres 컨테이너 삭제
```

개발 환경에서 실행하는 방법입니다.

```bash
$ make run
```

배포 환경에서 빌드 및 실행하는 방법입니다.

```bash
$ make build
$ ./bin/server
```

Production 환경을 위한 Dockerfile도 제공합니다.

```bash
$ docker build -t pets-next-door-api .
$ docker run -p 8080:8080 pets-next-door-api
```

## How to test

```bash
$ make test
```

## How to upgrade API docs version

`version.sh` 파일에 exec 권한을 부여합니다.

```bash
$ chmod +x ./scripts/version.sh
```

버전을 업그레이드하려면 다음 명령을 실행합니다.

```bash
$ make version
```

## Set Up

### Firebase

Firebase Auth를 사용하기 위해, Firebase 프로젝트를 생성하고, 서비스 계정 키를 발급받아야 합니다.

기본값으로 `firebase-credentials.json` 파일을 사용하며, 환경변수 `FIREBASE_CREDENTIALS_PATH`, 또는 `FIREBASE_CREDENTIALS_PATH` 비우고 `FIREBASE_CREDENTIALS_*`를 통해 파일 경로를 지정할 수 있습니다.

### Database Migration

Postgres DB를 사용하며, 마이그레이션은 [golang-migrate/migrate](https://github.com/golang-migrate/migrate)를 사용합니다.


마이그레이션 생성은 다음과 같이 합니다.

```bash
$ make migrate:create name=MIGRATION_NAME
```

마이그레이션 실행은 다음과 같이 합니다.

```bash
$ make migrate:up # 또는 go run ./cmd/migrate
```

테스트 환경은 `docker-compose-test.yml`를 통해 자동으로 마이그레이션됩니다.

배포 환경에서는 다음과 같은 스크립트를 제공합니다.

```bash
$ go build -o migrate ./cmd/migrate
$ ./migrate
```

## API Docs

API 문서는 [swaggo](https://github.com/swaggo/swag)를 사용하여 자동으로 생성됩니다.

```bash
$ make docs
```

`swagger/index.html` 또는 `/swagger/doc.json`을 열어 확인할 수 있습니다.

```bash
make docs:open
```
