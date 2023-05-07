package provider

type BaseRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var SuccessBaseRsp BaseRsp = BaseRsp{
	Code: 0,
	Msg:  "success",
}
