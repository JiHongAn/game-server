# Game Server

게임 서버 애플리케이션입니다.

## 🚀 시작하기

### 설치 및 실행

1. **의존성 설치**
   ```bash
   go mod download
   ```

2. **개발 환경 서비스 시작 (MySQL, Redis)**
   ```bash
   make docker-up
   ```

3. **환경 변수 설정**
   ```bash
   # 필요시 configs/.env.development 파일 수정
   ```

4. **애플리케이션 실행**
   ```bash
   make run-dev
   ```

5. **개발 완료 후 서비스 정리**
   ```bash
   make docker-down
   ```

### 개발 모드

Hot reload를 지원하는 개발 모드로 실행:

```bash
make dev
```

## 🔧 사용 가능한 명령어

```bash
make build         # 애플리케이션 빌드
make run           # 애플리케이션 실행 (development)
make run-dev       # 개발 환경에서 실행
make run-test      # 테스트 환경에서 실행
make run-prod      # 프로덕션 환경에서 실행
make dev           # 개발 모드 (hot reload)
make test          # 테스트 실행
make test-coverage # 테스트 커버리지 리포트
make lint          # 코드 린팅
make fmt           # 코드 포맷팅
make clean         # 빌드 결과물 정리
make help          # 도움말
```