package repositories

import (
	"clearance-adapter/config"
	"errors"
	"net/http"

	"github.com/pangpanglabs/goutils/httpreq"
)

//CreateErrData ...
func CreateErrData(data map[string]interface{}) (int64, error) {
	url := config.GetBatchJobAPIRoot()
	var resp struct {
		Result  int64       `json:"result"`
		Success bool        `json:"success"`
		Error   interface{} `json:"error"`
	}
	statusCode, err := httpreq.New(http.MethodPost, url, data).Call(&resp)
	if err != nil {
		return resp.Result, err
	}
	if statusCode >= 200 && statusCode <= 299 && resp.Success == true {
		return resp.Result, nil
	}
	return 0, errors.New(resp.Error.(string))
}
