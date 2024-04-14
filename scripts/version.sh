#!/bin/bash
# example: ./scripts/version.sh
# cmd/server/main.go에 있는 주석을 분석(// @version Major.Minor.Patch) 하여 버전을 올립니다.

# 주석을 분석할 파일 경로
FILE_PATH="./cmd/server/main.go"

# 주석을 분석할 패턴
PATTERN="// @version"

# 파일이 존재하는지 확인합니다.
if [ ! -f $FILE_PATH ]; then
  echo "파일이 존재하지 않습니다. ($FILE_PATH)"
  return
fi

# 주석을 분석하여 버전을 가져옵니다.
VERSION=$(grep "$PATTERN" $FILE_PATH | awk '{print $3}')
echo "기존 버전: $VERSION"


# Major을 올릴건지, Minor를 올릴건지, Patch를 올릴건지 여부를 확인합니다.
echo "Major.Minor.Patch 중 어떤 버전을 올릴지 선택해주세요:"
echo "1. Major"
echo "2. Minor"
echo "3. Patch"
read -p "버전 선택: " VERSION_TYPE

# Major, Minor, Patch를 설정합니다.
case $VERSION_TYPE in
1)
  MAJOR=1
  MINOR=0
  PATCH=0
  ;;
2)
  MAJOR=0
  MINOR=1
  PATCH=0
  ;;
3)
  MAJOR=0
  MINOR=0
  PATCH=1
  ;;
*)
  echo "잘못된 버전을 선택하였습니다."
  return
  ;;
esac

# Major, Minor, Patch를 분리합니다.
# shellcheck disable=SC2207
VERSION_ARRAY=($(echo "$VERSION" | tr '.' '\n'))
MAJOR_VERSION=${VERSION_ARRAY[0]}
MINOR_VERSION=${VERSION_ARRAY[1]}
PATCH_VERSION=${VERSION_ARRAY[2]}

# Major, Minor, Patch를 올립니다.
if [ $MAJOR -eq 1 ]; then
  MAJOR_VERSION=$((MAJOR_VERSION + 1))
  MINOR_VERSION=0
  PATCH_VERSION=0
elif [ $MINOR -eq 1 ]; then
  MINOR_VERSION=$((MINOR_VERSION + 1))
  PATCH_VERSION=0
elif [ $PATCH -eq 1 ]; then
  PATCH_VERSION=$((PATCH_VERSION + 1))
fi

# 새로운 버전을 만듭니다.
NEW_VERSION="$MAJOR_VERSION.$MINOR_VERSION.$PATCH_VERSION"
echo "새로운 버전: $NEW_VERSION"

# 주석을 변경합니다.
sed "s#$PATTERN $VERSION#$PATTERN $NEW_VERSION#g" "$FILE_PATH" > "$FILE_PATH.tmp"
mv "$FILE_PATH.tmp" "$FILE_PATH"

# 성공 여부를 확인합니다.
if [ $? -eq 0 ]; then
  echo "$VERSION -> $NEW_VERSION 버전으로 변경되었습니다. ($FILE_PATH)"

  # 버전 변경 후, 문서를 생성합니다.
  make docs
else
  echo "버전 변경에 실패하였습니다."
fi
