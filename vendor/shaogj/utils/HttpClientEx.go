
package utils

import (
	"crypto/tls"
	"bytes"
	"io"
	"io/ioutil"
	"errors"
	"net"
	"time"
	"net/http"
	"encoding/json"
	"github.com/mkideal/log"

)

type CHttpClientEx struct {
	Headers  map[string]string
	BSkipVerify bool
}
func  (self* CHttpClientEx)Init()      {

	self.Headers=map[string]string{}
}
func  NewHttpClient(bSkipVerify bool,iovertime  int)  *http.Client     {
	var hClient * http.Client=nil
	if iovertime<=0 {
		if false==bSkipVerify{
			return  http.DefaultClient
		}else {
			return &http.Client{
				Transport:  &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}
	}
	if  iovertime>0 {
		if false==bSkipVerify{
			return &http.Client{
				Transport: &http.Transport{
					Dial: func(netw, addr string) (net.Conn, error) {
						deadline := time.Now().Add(time.Duration(iovertime)* time.Millisecond)
						c, err := net.DialTimeout(netw, addr, time.Millisecond*time.Duration(iovertime) )
						if err != nil {
							return nil, err
						}
						c.SetDeadline(deadline)
						return  c,nil
					},
				},
			}
		}else {
			return &http.Client{
				Transport: &http.Transport{
					Dial: func(netw, addr string) (net.Conn, error) {
						deadline := time.Now().Add(time.Duration(iovertime)* time.Millisecond)
						c, err := net.DialTimeout(netw, addr, time.Millisecond*time.Duration(iovertime) )
						if err != nil {
							return nil, err
						}
						c.SetDeadline(deadline)
						return  c,nil
					},
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}

	}
	return   hClient
}

//https 取消证书认证
func  (self* CHttpClientEx)TLSSkipVerify(bSkeip bool) {
	self.BSkipVerify=bSkeip
}
/*
"Content-Type", "text/json"
"Connection", "Keep-Alive"
"Content-Type", "application/x-www-form-urlencoded"

//设置文本格式
HeaderSet("Content-Type", "text/json")
HeaderSet("Content-Type", "application/json;charset=utf-8")
*/
func  (self* CHttpClientEx)HeaderSet(strKey,strVal string) {
	self.Headers[strKey]=strVal
}
//返回：  返回数据  网络状态, 错误码  错误信息
func (self* CHttpClientEx)RequestResponseByte(req []byte,url string,iover int) ([]byte, int, int,error) {
	return  self.RequestResponse(url,bytes.NewReader(req),iover)
}
//1102 add ,for Signauth

func (self* CHttpClientEx)RequestJsonSignResponseJson(url string,iover int,vreq interface{},signData string,vres interface{}) ([]byte,int, int, error)  {
	byReq,err:=json.Marshal(vreq)
	if nil!=err {
		return  []byte{},0,7017,err
	}
	return  self.RequestResponseJson(url,bytes.NewReader(byReq),iover,vres)
}
//end add

func (self* CHttpClientEx)RequestJsonResponseJson(url string,iover int,vreq interface{},vres interface{}) ([]byte,int, int, error)  {
	byReq,err:=json.Marshal(vreq)
	if nil!=err {
		return  []byte{},0,7017,err
	}
	return  self.RequestResponseJson(url,bytes.NewReader(byReq),iover,vres)
}
func (self* CHttpClientEx)RequestJsonResponse(url string,iover int,vreq interface{}) ([]byte,int, int, error)  {
	byReq,err:=json.Marshal(vreq)
	if nil!=err {
		return  []byte{},0,7017,err
	}
	return  self.RequestResponse(url,bytes.NewReader(byReq),iover)
}
//处理
func (self* CHttpClientEx)RequestResponseJson(url string,body io.Reader,iover int,vres interface{}) ([]byte,int, int, error)  {

	byResp,statusCode,errorcode,err:=self.RequestResponse(url,body,iover)
	if 200!=statusCode||  0!=errorcode  {
		stE:=""
		if nil!=err{
			stE=err.Error()
		}
		return   []byte{},statusCode,6012,errors.New(stE)
	}
	err=json.Unmarshal(byResp,vres)
	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(byResp),url,err.Error())
		return byResp,statusCode,6016,err
	}
	return byResp,statusCode,0,nil
}
//处理
func (self* CHttpClientEx)RequestResponseJsonByte(url string,body io.Reader,iover int,vres interface{}) ([]byte, int, int, error)  {
	byResp,statusCode,errorcode,err:=self.RequestResponse(url,body,iover)
	if 200!=statusCode||  0!=errorcode  {
		stE:=""
		if nil!=err{
			stE=err.Error()
		}
		return  byResp,statusCode,6012,errors.New(stE)
	}
	err=json.Unmarshal(byResp,vres)
	if  nil!=err  {
		log.Error("resp=%s,url=%s,err=%v",string(byResp),url,err.Error())
		return byResp,statusCode,6016,err
	}
	return byResp,statusCode,0,nil
}
//返回：  返回数据  网络状态, 错误码  错误信息
func (self* CHttpClientEx)RequestResponse(url string,body io.Reader,iover int) ([]byte, int, int, error) {
	var byNull []byte=[]byte("")
	var  hServer *http.Client
	hServer=NewHttpClient(self.BSkipVerify,iover)
	if nil==hServer {
		return  byNull,0,5012,errors.New("Can;t create http Client")
	}

	var request *http.Request
	var err error=nil
	if nil==body {
		request, err = http.NewRequest("GET", url, nil)
	}else {
		request, err = http.NewRequest("POST", url, body)
	}
	if err != nil {
		log.Error("http.NewRequest,[err=%s][url=%s]", err, url)
		return byNull,0,5012,err
	}
	for k,v:=range self.Headers {
		request.Header.Set(k, v)
	}
	var resp *http.Response
	resp, err = hServer.Do(request)
	if err != nil {
		log.Error("http.Do failed,[err=%s][url=%s]", err.Error(), url)
		return byNull,0,5012,err
	}
	defer resp.Body.Close()
	if 200!=resp.StatusCode {
		return byNull,resp.StatusCode,5013,nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("http.Do failed,[err=%s][url=%s][res=%s]", err, url,string(b))
		return byNull,resp.StatusCode,5014,err
	}
	return b,resp.StatusCode,0,nil
}