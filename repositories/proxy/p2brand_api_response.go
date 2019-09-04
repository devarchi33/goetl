package proxy

// P2BrandAPIResponse API 返回结果
type P2BrandAPIResponse struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details"`
	} `json:"error"`
}
