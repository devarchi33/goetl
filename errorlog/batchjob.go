package errorlog

import (
	"clearance-adapter/config"
	clrConst "clearance-adapter/domain/clr-constants"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/pangpanglabs/goutils/httpreq"
)

type ErrorLog struct {
	LogID int64
}

//CreateLog ...
func (ErrorLog) CreateLog(errType clrConst.ClrErrorType) (int64, error) {
	serverName := config.GetServerName()
	data := make(map[string]interface{})
	param := make(map[string]string)
	param["type"] = errType
	data["service"] = serverName
	data["param"] = param
	id, err := createErrData(data)
	if err != nil {
		log.Printf(err.Error())
		return id, err
	}
	return id, nil
}

//CreateErrData ...
func createErrData(data map[string]interface{}) (int64, error) {
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

//AppendLogs ...
func (ErrorLog) AppendLogs(id int64, message map[string]string) {
	url := config.GetBatchJobAPIRoot()
	logid := strconv.FormatInt(id, 10)
	data := make(map[string]string)
	data["level"] = "error"
	msg, _ := json.Marshal(message)
	data["message"] = string(msg)
	httpreq.New(http.MethodPost, url+"/"+logid+"/logs", data).Call(&map[string]string{})
}

//Finish ...
func (ErrorLog) Finish(id int64) {
	url := config.GetBatchJobAPIRoot()
	logid := strconv.FormatInt(id, 10)
	httpreq.New(http.MethodPost, url+"/"+logid+"/finish", "").Call(&map[string]string{})
}
