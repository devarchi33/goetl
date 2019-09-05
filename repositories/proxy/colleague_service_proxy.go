package proxy

import (
	"clearance-adapter/config"
	"errors"
	"fmt"
	"net/http"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
)

// ColleagueServiceProxy colleague服务的代理
type ColleagueServiceProxy struct{}

// GetTokenByUsernameAndPassword 根据用户名密码获取token
func (ColleagueServiceProxy) GetTokenByUsernameAndPassword(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("用户名或密码不能为空")
	}

	var resp P2BrandAPIResponse
	url := config.GetP2BrandColleagueAPIRoot() + "/sso/login-user-name?appCode=CloudPortal&tenantCode=pangpang"
	data := make(map[string]interface{}, 0)
	data["username"] = username
	data["password"] = password
	statusCode, err := httpreq.New(http.MethodPost, url, data).
		WithBehaviorLogContext(behaviorlog.FromCtx(nil)).
		Call(&resp)

	if err != nil {
		return "", err
	}

	if statusCode >= 200 && statusCode <= 299 {
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("GetTokenByUsernameAndPassword convert Result to map[string]interface{} failed response: %v", resp)
		}
		token, ok := result["token"]
		if !ok {
			return "", errors.New("there is no key's name is token")
		}
		tokenValue, ok := token.(string)
		if !ok {
			return "", fmt.Errorf("GetTokenByUsernameAndPassword convert token to string failed token: %v", token)
		}
		return tokenValue, nil
	}

	return "", fmt.Errorf("GetTokenByUsernameAndPassword call colleague API failed status code: %v, response: %v", statusCode, resp)
}
