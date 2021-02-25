package requests

import (
	//"encoding/json"
	//"errors"
	//"io"
	//"io/ioutil"
	//"net/http"
	//"net/url"
	//"os"
	//"strings"
	//"wqk/pubilc"
	//"net/http"
)

//const (
//	ContentType               = "Content-Type"
//	ApplicationJSON           = "application/json"
//	ApplicationFormUrlencoded = "application/x-www-form-urlencoded"
//)

//// RequestInterceptor 请求拦截器
//// 返回不为nil，即有错误会终止后续执行
//type RequestInterceptor func(request *http.Request) error
//
//// requestInterceptorChain 请求拦截链
//type requestInterceptorChain struct {
//	mutex        *sync.RWMutex
//	interceptors []RequestInterceptor
//}
//
//// defaultRequestInterceptorChain 默认的请求拦截链实例
//var defaultRequestInterceptorChain = &requestInterceptorChain{
//	mutex:        new(sync.RWMutex),
//	interceptors: make([]RequestInterceptor, 0),
//}

// Client 封装了http的参数等信息
//type Client struct {
//	// 自定义Client
//	Client *http.Client
//
//	url    string
//	method string
//	header http.Header
//	params url.Values
//
//	form      url.Values
//	json      interface{}
//	multipart FileForm
//}
//
//// FileForm form参数和文件参数
//type FileForm struct {
//	Value url.Values
//	File  map[string]string
//}
//
//// Result http响应结果
//type Result struct {
//	Resp *http.Response
//	Err  error
//}



//// Post http `POST` 请求
//func Post(url string) *Client {
//	return newClient(url, http.MethodPost, nil)
//}
//
//// Put http `PUT` 请求
//func Put(url string) *Client {
//	return newClient(url, http.MethodPut, nil)
//}
//
//// Delete http `DELETE` 请求
//func Delete(url string) *Client {
//	return newClient(url, http.MethodDelete, nil)
//}
//
//// Request 用于自定义请求方式，比如`HEAD`、`PATCH`、`OPTIONS`、`TRACE`
//// client参数用于替换DefaultClient，如果为nil则会使用默认的
//func Request(url, method string, client *http.Client) *Client {
//	return newClient(url, method, client)
//}

//// Params http请求中url参数
//func (c *Client) params(params url.Values) *Client {
//	for k, v := range params {
//		c.Params[k] = v
//	}
//	return c
//}

//// Header http请求头
//func (c *Client) header(k, v string) *Client {
//	c.header.Set(k, v)
//	return c
//}
//
//// Headers http请求头
//func (c *Client) Headers(header http.Header) *Client {
//	for k, v := range header {
//		c.header[k] = v
//	}
//	return c
//}
//
//// Form 表单提交参数
//func (c *Client) form(form url.Values) *Client {
//	c.header.Set(pubilc.ContentType, pubilc.ApplicationFormUrlencoded)
//	c.form = form
//	return c
//}
//
//// Json json提交参数
//// 如果是string，则默认当作是json字符串；否则会序列化为json字节数组，再发送
//func (c *Client) json(json interface{}) *Client {
//	c.header.Set(pubilc.ContentType, pubilc.ApplicationJSON)
//	c.json = json
//	return c
//}
//
//// Multipart form-data提交参数
//func (c *Client) multipart(multipart FileForm) *Client {
//	c.multipart = multipart
//	return c
//}
//
//
//
//
//
//
//
//// createForm 创建application/x-www-form-urlencoded请求
//func (c *Client) createForm() *Result {
//	var result = new(Result)
//
//	form := c.form.Encode()
//
//	req, err := http.NewRequest(c.method, c.url, strings.NewReader(form))
//	if err != nil {
//		result.Err = err
//		return result
//	}
//
//	req.Header = c.header
//	c.doSend(req, result)
//	return result
//}

// createEmptyBody 没有内容的body
//func (c *Client) createEmptyBody() *Result {
//	var result = new(Result)
//
//	req, err := http.NewRequest(c.method, c.url, nil)
//	if err != nil {
//		result.Err = err
//		return result
//	}
//
//	req.Header = c.header
//	c.doSend(req, result)
//	return result
//}

//// doSend 发送请求
//func (c *Client) doSend(req *http.Request, result *Result) {
//	// 调用拦截器，遇到错误就退出
//	if err := c.beforeSend(req); err != nil {
//		result.Err = err
//		return
//	}
//
//	// 发送请求
//	result.Resp, result.Err = c.Client.Do(req)
//}
//
//// beforeSend 发送请求前，调用拦截器
//func (c *Client) beforeSend(req *http.Request) error {
//	mutex := defaultRequestInterceptorChain.mutex
//	mutex.RLock()
//	defer mutex.RUnlock()
//
//	// 遍历调用拦截器
//	for _, interceptor := range defaultRequestInterceptorChain.interceptors {
//		err := interceptor(req)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//// StatusOk 判断http响应码是否为200
//func (r *Result) StatusOk() *Result {
//	if r.Err != nil {
//		return r
//	}
//	if r.Resp.StatusCode != http.StatusOK {
//		r.Err = errors.New("status code is not 200")
//		return r
//	}
//
//	return r
//}
//
//// Status2xx 判断http响应码是否为2xx
//func (r *Result) Status2xx() *Result {
//	if r.Err != nil {
//		return r
//	}
//	if r.Resp.StatusCode < http.StatusOK || r.Resp.StatusCode >= http.StatusMultipleChoices {
//		r.Err = errors.New("status code is not match [200, 300)")
//		return r
//	}
//
//	return r
//}



// newClient 创建Client
//func newClient(u string, method string, client *http.Client) *Client {
//	// client为nil则使用默认的DefaultClient
//	if client == nil {
//		client = http.DefaultClient
//	}
//
//	return &Client{
//		Client: client,
//		url:    u,
//		method: method,
//		header: make(http.Header),
//		params: make(url.Values),
//		form:   make(url.Values),
//	}
//}

// AddRequestInterceptors 添加请求拦截器
func AddRequestInterceptors(interceptors ...RequestInterceptor) {
	mutex := defaultRequestInterceptorChain.mutex
	mutex.Lock()
	defer mutex.Unlock()

	// 添加到拦截器链
	for _, interceptor := range interceptors {
		defaultRequestInterceptorChain.interceptors = append(defaultRequestInterceptorChain.interceptors, interceptor)
	}
}