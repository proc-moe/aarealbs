package e

const (
	DB_ERR         = 102
	DATA_NOT_FOUND = 100
	PARAM_ERR      = 101
	PARTIAL_FAIL   = 200
	TOKEN_FAILED   = 601
	UNAUTHORIZED   = 602
)

var Str map[int]string

func init() {
	Str = map[int]string{
		100: "DATA_NOT_FOUND",
		101: "PARAM_ERR",
		102: "DB_ERR",
		601: "TOKEN EXPIRED",
		602: "UNAUTHORIZED",
		200: "PARTIAL_FAIL",
	}

}
