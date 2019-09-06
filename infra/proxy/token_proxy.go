package proxy

import (
	"clearance-adapter/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/httpreq"
)

type userClaim struct {
	Aud           string
	ColleagueName string
	Exp           int64
	ID            int64
	Iss           string
	Nbf           int64
	TenantCode    string
	TenantID      int64
	Username      string
}

// TokenProxy 从colleague获取token的代理
type TokenProxy struct {
	Username   string
	Password   string
	Token      string
	Expiration int64
}

var tokenProxyInstance *TokenProxy
var tokenProxyOnce sync.Once

// GetInstance 创建 TokenProxy 对象实例
func (TokenProxy) GetInstance() *TokenProxy {
	tokenProxyOnce.Do(func() {
		username, password := config.GetColleagueClearanceUsernameAndPassword()
		tokenProxyInstance = &TokenProxy{
			Username: username,
			Password: password,
		}
	})
	return tokenProxyInstance
}

func (TokenProxy) parseToken(token string) (userClaim, error) {
	var claim userClaim

	si := strings.Index(token, ".")
	li := strings.LastIndex(token, ".")
	if si == -1 || li == -1 || si == li {
		return claim, errors.New("Invalid token")
	}

	payload := token[si+1 : li]
	if payload == "" {
		return claim, errors.New("Invalid token")
	}

	decodeSegment := func(seg string) ([]byte, error) {
		if l := len(seg) % 4; l > 0 {
			seg += strings.Repeat("=", 4-l)
		}

		return base64.URLEncoding.DecodeString(seg)
	}

	payloadBytes, err := decodeSegment(payload)
	if err != nil {
		return claim, err
	}

	if err := json.Unmarshal(payloadBytes, &claim); err != nil {
		return claim, err
	}

	return claim, nil
}

// GetToken 根据用户名密码获取token
func (proxy *TokenProxy) GetToken() (string, error) {
	if proxy.Expiration > time.Now().Unix() {
		return proxy.Token, nil
	}

	token, err := proxy.getToken()
	if err != nil {
		return "", err
	}

	claim, err := TokenProxy{}.parseToken(token)
	if err != nil {
		return "", err
	}

	proxy.Token = token
	proxy.Expiration = claim.Exp

	return token, nil
}

func (proxy *TokenProxy) getToken() (string, error) {
	if proxy.Username == "" || proxy.Password == "" {
		return "", errors.New("用户名或密码不能为空")
	}

	var resp P2BrandAPIResponse
	url := config.GetP2BrandColleagueAPIRoot() + "/sso/login-user-name?appCode=CloudPortal&tenantCode=pangpang"
	data := make(map[string]interface{}, 0)
	data["username"] = proxy.Username
	data["password"] = proxy.Password
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
