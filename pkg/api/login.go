package api

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/shiningrush/droplet"
	"github.com/shiningrush/droplet/data"
	"github.com/shiningrush/droplet/wrapper"
	"github.com/shiningrush/droplet/wrapper/gorestful"
	"github.com/shiningrush/goreq"
	"github.com/shiningrush/hspider/pkg/utils"
	"net/http"
	"net/url"
	"reflect"
)

const (
	UrlMakeQRCode = "https://user.91160.com/makeQrcode.html"
	UrlPull       = "https://user.91160.com/poll.html"
	UrlPullLogin  = "https://user.91160.com/login.html"
)

func NewLogin() *Login {
	svc := &Login{}
	return svc
}

type Login struct {
}

func (l *Login) ConfigRoutes(container *restful.Container) {
	ws := &restful.WebService{}
	ws.Path("/login")
	ws.Route(ws.GET("/getQRCode").To(gorestful.Wraps(l.GetLoginQRCode)))
	ws.Route(ws.GET("/cacheInfo").To(gorestful.Wraps(l.CacheLoginInfo,
		wrapper.InputType(reflect.TypeOf(&CacheLoginInfoInput{})))))

	container.Add(ws)
}

var QRCodeBase = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>QRCode</title>
</head>
<body>
	<a href="./cacheInfo?sceneId=%s">cache for %s</a>
    <img src="%s" />
</body>
</html>
`

type GetLoginQRCodeInput struct {
}

func (l *Login) GetLoginQRCode(ctx droplet.Context) (interface{}, error) {
	type makeQRCodeResp struct {
		SceneID string `json:"scene_id"`
		State   int    `json:"state"`
		Url     string `json:"url"`
	}

	var resp makeQRCodeResp
	if err := goreq.Get(UrlMakeQRCode, goreq.SetHeader(fakeClientHeader()), goreq.JsonResp(&resp)).Do(); err != nil {
		return nil, fmt.Errorf("get qrcode failed: %w", err)
	}

	return &data.RawResponse{
		StatusCode: http.StatusOK,
		Body:       []byte(fmt.Sprintf(QRCodeBase, resp.SceneID, resp.SceneID, resp.Url)),
	}, nil
}

type CacheLoginInfoInput struct {
	SceneID string `auto_read:"sceneId,query"`
}

func (l *Login) CacheLoginInfo(ctx droplet.Context) (interface{}, error) {
	input := ctx.Input().(*CacheLoginInfoInput)

	_, err := getUserName(input.SceneID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getUserName(sceneId string) (string, error) {
	type pullResp struct {
		Code     int    `json:"code"`
		Type     string `json:"type"`
		UserName string `json:"username"`
	}

	var resp pullResp
	err := goreq.Post(UrlPull,
		goreq.SetHeader(fakeClientHeader()),
		goreq.FormReq(url.Values{
			"scene_id": []string{sceneId},
		}),
		goreq.JsonResp(&resp)).Do()
	if err != nil {
		return "", err
	}

	if resp.Code != 1 {
		bs, _ := json.Marshal(resp)
		return "", fmt.Errorf("pull code is not correct: %v ", bs)
	}

	return resp.UserName, nil
}

func getLoginUrl(userName string) (string, error) {
	encryptedUserName, err := utils.Encrypt(userName)
	if err != nil {
		return "", err
	}

	var loginResp http.Response
	var loginRespBody []byte
	err = goreq.Post(UrlPullLogin,
		goreq.SetHeader(fakeClientHeader()),
		goreq.FormReq(url.Values{
			"username": []string{encryptedUserName},
			"password": []string{},
			"type":     []string{"x"},
			"target":   []string{"https://www.91160.com"},
		}),
		goreq.RawResp(&loginResp, &loginRespBody)).Do()
	if err != nil {
		return "", err
	}

	if loginResp.Header.Get("Location") == "" {
		return "", fmt.Errorf("get location failed, code: %d, body: %s", loginResp.StatusCode, loginRespBody)
	}

	return loginResp.Header.Get("Location"), nil
}

func fakeClientHeader() http.Header {
	return http.Header{
		"Referer":    []string{"https://user.91160.com/login.html"},
		"Host":       []string{"https://user.91160.com"},
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36"},
	}
}
