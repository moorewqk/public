package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)


// 自定义Client
type Client struct {
	Client    *http.Client   //原请求客户端
	Url       string         //请求地址
	Method    string         //请求方法
	Timeout    time.Duration         //请求方法
	Header    http.Header    //headers
	Params    url.Values     //入参
	Form      url.Values     //form入参
	Json      interface{}    //json入参
	Multipart FileForm       //多文件入参
	Resp      *http.Response //返回请求结果
	Err       error          //err信息
}

// FileForm form参数和文件参数
type FileForm struct {
	Value url.Values
	File  map[string]string
}

// Result http响应结果
type Result struct {
	Resp *http.Response   	//原response数据
	StatusCode int		 	//外层http Code,从response中引出
	Err  error				//结果错误对象
	Message string			//结果错误字符串信息,用于人识别
	Json  interface{}		//response body转换为数据对象
	Text   string			//response body转换为字符串类型
}

//初始化请求客户端
func newClient(u string, method string,timeout int, client *http.Client) *Client {
	var (
		timet time.Duration
	)


	// client为nil则使用默认的DefaultClient
	if client == nil {
		client = http.DefaultClient
	}

	//超时处理
	if timeout == 0{
		timet = time.Duration(HTTP_REQUEST_TIMEOUT)*time.Second
	}else {
		timet = time.Duration(timeout)*time.Second
	}
	////长链接处理
	//transport := &http.Transport{
	//	TLSHandshakeTimeout: 5 * time.Second,
	//	TLSClientConfig:     nil,
	//	Dial: (&net.Dialer{
	//		Timeout:   500 * time.Second,
	//		KeepAlive: 30 * time.Second,
	//	}).Dial,
	//	ResponseHeaderTimeout: 1 * time.Second,
	//}

	client.Timeout = timet
	return &Client{
		Client:    client,
		Url:       u,
		Method:    method,
		Header:    make(http.Header),
		Params:    make(url.Values),
		Form:      make(url.Values),
		Multipart: FileForm{},
	}
}

// RequestInterceptor 请求拦截器
// 返回不为nil，即有错误会终止后续执行
type RequestInterceptor func(request *http.Request) error

// requestInterceptorChain 请求拦截链
type requestInterceptorChain struct {
	mutex        *sync.RWMutex
	interceptors []RequestInterceptor
}

// defaultRequestInterceptorChain 默认的请求拦截链实例
var defaultRequestInterceptorChain = &requestInterceptorChain{
	mutex:        new(sync.RWMutex),
	interceptors: make([]RequestInterceptor, 0),
}


// Get http `GET` 请求
func Get(url string,params url.Values,timeout int) (result *Result) {
	client :=newClient(url, http.MethodGet, timeout,nil)
	client.addParams(params)
	result = client.Send()
	result.Response()
	return result
}

func Post(url string,data interface{},timeout int) (result *Result) {
	client :=newClient(url, http.MethodPost, timeout,nil)
	client.addJson(data)
	result = client.Send()
	result.Response()
	return result
}


//结构化response结果数据
func (r *Result) Response() (err error) {
	var (
		v interface{}
	)

	b, err := r.Raw()
	if err != nil {
		r.Err = err
		r.Message = fmt.Sprintf("%s:%s", HttpMessageFail,err.Error())
		return err
	}
	r.Text = string(b)
	r.Message = fmt.Sprintf("%s", HttpMessageSuccess)
	err = json.Unmarshal(b, &v)
	if err!=nil{
		r.Err = err
		r.Message = fmt.Sprintf("%s:%s",HttpMessageFail,err.Error())
		return err
	}
	r.Json = v
	return nil
}

// Params http请求中url参数
func (c *Client) addParams(params url.Values) {
	for k, v := range params {
		c.Params[k] = v
	}
	return
}

// Json json提交参数
// 如果是string，则默认当作是json字符串；否则会序列化为json字节数组，再发送
func (c *Client) addJson(json interface{}) *Client {
	c.Header.Set(ContentType, ApplicationJSON)
	c.Json = json
	return c
}


//Send 发送http请求
func (c *Client) Send() *Result {
	var result *Result
	// 处理query string
	if c.Params != nil && len(c.Params) != 0 {
		// 如果url中已经有query string参数，则只需要&拼接剩下的即可
		encoded := c.Params.Encode()
		if strings.Index(c.Url, "?") == -1 {
			c.Url += "?" + encoded
		} else {
			c.Url += "&" + encoded
		}
	}
	// 根据不同的Content-Type设置不同的http body
	contentType := c.Header.Get(ContentType)
	if c.Multipart.Value != nil || c.Multipart.File != nil {
		result = c.multiFormReq()
	} else if strings.HasPrefix(contentType, ApplicationJSON) {
		result = c.jsonReq()
	} else if strings.HasPrefix(contentType,ApplicationFormUrlencoded) {
		result = c.formReq()
	} else {
		// 不是以上类型，就不设置http body
		result = c.emptyBodyReq()
	}
	return result
}

// Raw 获取http响应内容，返回字节数组
func (r *Result) Raw() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	b, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		r.Err = err
		return nil, r.Err
	}
	defer r.Resp.Body.Close()
	return b, r.Err
}


// Text 获取http响应内容，返回字符串
func (r *Result) text() (bodystr string,code int, err error) {
	b, err := r.Raw()
	if err != nil {
		r.Err = err
		return "",401, r.Err
	}

	if r.Resp.StatusCode != http.StatusOK {
		r.Err = errors.New("status code is not 200")
		return "",r.Resp.StatusCode,r.Err
	}
	return string(b),200,nil
}

// Json 获取http响应内容，返回json
func (r *Result) json() error {
	var (
		v  interface{}
	)
	b, err := r.Raw()
	if err != nil {
		r.Err = err
		return r.Err
	}
	return json.Unmarshal(b, v)
}

// Save 获取http响应内容，保存为文件
func (r *Result) Save(name string) error {
	if r.Err != nil {
		return r.Err
	}

	f, err := os.Create(name)
	if err != nil {
		r.Err = err
		return r.Err
	}
	defer f.Close()

	_, err = io.Copy(f, r.Resp.Body)
	if err != nil {
		r.Err = err
		return r.Err
	}
	defer r.Resp.Body.Close()

	return nil
}

// doSend 发送请求
func (c *Client) doSend(req *http.Request, result *Result) {
	// 调用拦截器，遇到错误就退出
	if err := c.beforeSend(req); err != nil {
		result.Err = err
		return
	}
	// TODO 发送请求
	result.Resp, result.Err = c.Client.Do(req)
	result.StatusCode = result.Resp.StatusCode
}

// beforeSend 发送请求前，调用拦截器,查询是否有请求错误,如果有就返回
func (c *Client) beforeSend(req *http.Request) error {
	mutex := defaultRequestInterceptorChain.mutex
	mutex.RLock()
	defer mutex.RUnlock()
	// 遍历调用拦截器
	for _, interceptor := range defaultRequestInterceptorChain.interceptors {
		err := interceptor(req)
		if err != nil {
			return err
		}
	}
	return nil
}

// createForm 创建application/json请求
func (c *Client) jsonReq() *Result {
	var result = new(Result)
	b, err := json.Marshal(c.Json)
	if err != nil {
		result.Err = err
		return result
	}
	req, err := http.NewRequest(c.Method, c.Url, bytes.NewReader(b))
	if err != nil {
		result.Err = err
		return result
	}
	req.Header = c.Header
	c.doSend(req, result)
	return result
}


// createEmptyBody 处理没有内容的body
func (c *Client) emptyBodyReq() *Result {
	var result = new(Result)
	//请求对象
	req, err := http.NewRequest(c.Method, c.Url, nil)
	if err != nil {
		result.Err = err
		return result
	}
	req.Header = c.Header
	c.doSend(req, result)
	return result
}

// createForm 创建application/x-www-form-urlencoded请求
func (c *Client) formReq() *Result {
	var result = new(Result)

	form := c.Form.Encode()

	req, err := http.NewRequest(c.Method, c.Url, strings.NewReader(form))
	if err != nil {
		result.Err = err
		return result
	}

	req.Header = c.Header
	c.doSend(req, result)
	return result
}

// createMultipartForm 创建form-data的请求
func (c *Client) multiFormReq() *Result {
	var result = new(Result)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 设置文件字节
	for name, filename := range c.Multipart.File {
		file, err := os.Open(filename)
		if err != nil {
			c.Err = err
			return result
		}
		part, err := writer.CreateFormFile(name, filename)
		if err != nil {
			result.Err = err
			return result
		}
		// todo 这里的io.Copy实现，会把file文件都读取到内存里面，然后当做一个buffer传给NewRequest。对于大文件来说会占用很多内存
		_, err = io.Copy(part, file)
		if err != nil {
			result.Err = err
			return result
		}

		err = file.Close()
		if err != nil {
			result.Err = err
			return result
		}
	}

	// 设置field
	for name, values := range c.Multipart.Value {
		for _, value := range values {
			_ = writer.WriteField(name, value)
		}
	}
	err := writer.Close()
	if err != nil {
		result.Err = err
		return result
	}

	req, err := http.NewRequest(c.Method, c.Url, body)
	req.Header = c.Header
	req.Header.Set(ContentType, writer.FormDataContentType())
	c.doSend(req, result)
	return result
}





////////////////////

//func ReqDo(method,url string, timeout int, header http.Header ) (*http.Response, []byte, error) {
//	req, err := http.NewRequest(method,url,nil)
//	if err != nil {
//		return nil, nil, err
//	}
//	if len(header) > 0 {
//		req.Header = header
//	}
//	dt := makeTime(timeout)
//	client := http.Client{Timeout: dt,}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, nil, err
//	}
//	return resp, body, nil
//
//}
//
//func HttpGet(url string, params url.Values, timeout int, header http.Header) (*http.Response, []byte, error) {
//	url = makeUrlWithData(url, params)
//	return ReqDo("GET",url,timeout,header)
//}
//
//func HttpPOST(url string, params url.Values, timeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
//	startTime := time.Now().UnixNano()
//	//client := http.Client{
//	//	Timeout: time.Duration(msTimeout) * time.Millisecond,
//	//}
//	if contextType == "" {
//		contextType = "application/x-www-form-urlencoded"
//	}
//	req, err := http.NewRequest("POST", urlString, strings.NewReader(urlParams.Encode()))
//	if len(header) > 0 {
//		req.Header = header
//	}
//	req = addTrace2Header(req, trace)
//	req.Header.Set("Content-Type", contextType)
//	resp, err := client.Do(req)
//	if err != nil {
//		Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
//			"url":       urlString,
//			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//			"method":    "POST",
//			"args":      urlParams,
//			"err":       err.Error(),
//		})
//		return nil, nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
//			"url":       urlString,
//			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//			"method":    "POST",
//			"args":      urlParams,
//			"result":    string(body),
//			"err":       err.Error(),
//		})
//		return nil, nil, err
//	}
//	Log.TagInfo(trace, DLTagHTTPSuccess, map[string]interface{}{
//		"url":       urlString,
//		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//		"method":    "POST",
//		"args":      urlParams,
//		"result":    string(body),
//	})
//	return resp, body, nil
//}
//
//func HttpJSON(trace *TraceContext, urlString string, jsonContent string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
//	startTime := time.Now().UnixNano()
//	client := http.Client{
//		Timeout: time.Duration(msTimeout) * time.Millisecond,
//	}
//	req, err := http.NewRequest("POST", urlString, strings.NewReader(jsonContent))
//	if len(header) > 0 {
//		req.Header = header
//	}
//	req = addTrace2Header(req, trace)
//	req.Header.Set("Content-Type", "application/json")
//	resp, err := client.Do(req)
//	if err != nil {
//		Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
//			"url":       urlString,
//			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//			"method":    "POST",
//			"args":      jsonContent,
//			"err":       err.Error(),
//		})
//		return nil, nil, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		Log.TagWarn(trace, DLTagHTTPFailed, map[string]interface{}{
//			"url":       urlString,
//			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//			"method":    "POST",
//			"args":      jsonContent,
//			"result":    string(body),
//			"err":       err.Error(),
//		})
//		return nil, nil, err
//	}
//	Log.TagInfo(trace, DLTagHTTPSuccess, map[string]interface{}{
//		"url":       urlString,
//		"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
//		"method":    "POST",
//		"args":      jsonContent,
//		"result":    string(body),
//	})
//	return resp, body, nil
//}
//
//func makeUrlWithData(url string, data url.Values) string {
//	if strings.Contains(url, "?") {
//		url = url + "&"
//	} else {
//		url = url + "?"
//	}
//	return fmt.Sprintf("%s%s",url, data.Encode())
//}