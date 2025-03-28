package errors

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func InvalidInput() *AppError {
	return &AppError{
		Code:       "INVALID_INPUT",
		Message:    "필수 입력값이 누락되었거나 형식이 올바르지 않습니다",
		StatusCode: 400,
	}
}

func BadRequest() *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    "잘못된 요청입니다",
		StatusCode: 400,
	}
}

func NotFound() *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    "API 요청에 대한 리소스를 찾을 수 없습니다",
		StatusCode: 404,
	}
}

func Unauthorized() *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "유효하지 않은 토큰입니다",
		StatusCode: 401,
	}
}

func InternalServerError() *AppError {
	return &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "서버 오류가 발생했습니다",
		StatusCode: 500,
	}
}

func DBError() *AppError {
	return &AppError{
		Code:       "DB_ERROR",
		Message:    "데이터베이스 오류가 발생했습니다",
		StatusCode: 500,
	}
}
