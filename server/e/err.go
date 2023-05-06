package e

const (
	DATA_NOT_FOUND = 100
	PARAM_ERR      = 101
	UNAUTHORIZED   = 102
	PARTIAL_FAIL   = 200
)

var Str map[int]string

func init() {
	Str = map[int]string{
		100: "DATA_NOT_FOUND",
		101: "PARAM_ERR",
		102: "UNAUTHORIZED",
		200: "PARTIAL_FAIL",
	}

}
