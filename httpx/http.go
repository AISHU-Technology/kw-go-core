package httpx

import (
	"bytes"
	"crypto/tls"
	"github.com/AISHU-Technology/kw-go-core/utils"

	"errors"
	"github.com/json-iterator/go"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"
)

const (
	ContentTypeJson       string = "application/json;charset=UTF-8"
	ContentTypeUrlEncoded string = "application/x-www-form-urlencoded"
	ContentTypeFromData   string = "multipart/form-data"
	ContentType           string = "Content-Type"
)

func HttpGetJson(url string) (data string, err error) {
	header := map[string]string{
		ContentType: ContentTypeJson,
	}
	return HttpGet(url, header)
}

func HttpGetParams(url string, paramMap map[string]interface{}) (data string, err error) {
	header := map[string]string{
		ContentType: ContentTypeUrlEncoded,
	}
	return HttpGet(url, header, paramMap)
}
func HttpGet(url string, params ...any) (data string, err error) {
	res, err := Httpx(http.MethodGet, url, params...)
	return res, err
}

func HttpPostJson(url string, bodyMap map[string]interface{}) (data string, err error) {
	bytesData, err := jsoniter.Marshal(bodyMap)
	if err != nil {
		return "", err
	}
	header := map[string]string{
		ContentType: ContentTypeJson,
	}
	return HttpPost(url, header, nil, bytesData)
}
func HttpPostFormData(url string, postData map[string]interface{}) (data string, err error) {
	header := map[string]string{
		ContentType: ContentTypeFromData,
	}
	return HttpPost(url, header, nil, postData)
}

func HttpPostUrlencoded(url string, postData map[string]interface{}) (data string, err error) {
	header := map[string]string{
		ContentType: ContentTypeUrlEncoded,
	}
	return HttpPost(url, header, nil, postData)
}

// HttpPost send post http request.
func HttpPost(url string, params ...any) (data string, err error) {
	res, err := Httpx(http.MethodPost, url, params...)
	return res, err
}

func Httpx(method string, url string, params ...any) (data string, err error) {
	res, err := doHttpRequest(method, url, params...)
	if err != nil {
		return "", err
	}
	// 最后关闭res.Body文件
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	// 使用ioutil.ReadAll将res.Body中的数据读取出来,并使用body接收
	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	// byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&body))
	return *str, nil
}

// HttpPut send put http request.
func HttpPut(url string, params ...any) (*http.Response, error) {
	return doHttpRequest(http.MethodPut, url, params...)
}

// HttpDelete send delete http request.
func HttpDelete(url string, params ...any) (*http.Response, error) {
	return doHttpRequest(http.MethodDelete, url, params...)
}

// HttpPatch send patch http request.
func HttpPatch(url string, params ...any) (*http.Response, error) {
	return doHttpRequest(http.MethodPatch, url, params...)
}

// HttpRequest struct is a composed http request
type HttpRequest struct {
	RawURL      string
	Method      string
	Headers     http.Header
	QueryParams url.Values
	FormData    url.Values
	File        *File
	Body        []byte
}

// HttpClientConfig contains some configurations for http client
type HttpClientConfig struct {
	SSLEnabled       bool
	TLSConfig        *tls.Config
	Compressed       bool
	HandshakeTimeout time.Duration
	ResponseTimeout  time.Duration
	Verbose          bool
}

// defaultHttpClientConfig defalut client config
var defaultHttpClientConfig = &HttpClientConfig{
	Compressed:       false,
	HandshakeTimeout: 20 * time.Second,
	ResponseTimeout:  40 * time.Second,
}

// HttpClient is used for sending http request
type HttpClient struct {
	*http.Client
	TLS     *tls.Config
	Request *http.Request
	Config  HttpClientConfig
}

// NewHttpClient make a HttpClient instance
func NewHttpClient() *HttpClient {
	client := &HttpClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   defaultHttpClientConfig.HandshakeTimeout,
				ResponseHeaderTimeout: defaultHttpClientConfig.ResponseTimeout,
				DisableCompression:    !defaultHttpClientConfig.Compressed,
			},
		},
		Config: *defaultHttpClientConfig,
	}

	return client
}

// NewHttpClientWithConfig make a HttpClient instance with pass config
func NewHttpClientWithConfig(config *HttpClientConfig) *HttpClient {
	if config == nil {
		config = defaultHttpClientConfig
	}

	client := &HttpClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   config.HandshakeTimeout,
				ResponseHeaderTimeout: config.ResponseTimeout,
				DisableCompression:    !config.Compressed,
			},
		},
		Config: *config,
	}

	if config.SSLEnabled {
		client.TLS = config.TLSConfig
	}

	return client
}

// SendRequest send http request.
// Play: https://go.dev/play/p/jUSgynekH7G
func (client *HttpClient) SendRequest(request *HttpRequest) (*http.Response, error) {
	err := validateRequest(request)
	if err != nil {
		return nil, err
	}

	rawUrl := request.RawURL

	req, err := http.NewRequest(request.Method, rawUrl, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}
	client.setTLS(rawUrl)
	client.setHeader(req, request.Headers)

	err = client.setQueryParam(req, rawUrl, request.QueryParams)
	if err != nil {
		return nil, err
	}

	if request.FormData != nil {
		if request.File != nil {
			err = client.setFormData(req, request.FormData, setFile(request.File))
		} else {
			err = client.setFormData(req, request.FormData, nil)
		}
	}

	client.Request = req

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DecodeResponse decode response into target object.
// Play: https://go.dev/play/p/jUSgynekH7G
func (client *HttpClient) DecodeResponse(resp *http.Response, target any) error {
	if resp == nil {
		return errors.New("invalid target param")
	}
	defer resp.Body.Close()
	return jsoniter.NewDecoder(resp.Body).Decode(target)
}

// setTLS set http client transport TLSClientConfig
func (client *HttpClient) setTLS(rawUrl string) {
	if strings.HasPrefix(rawUrl, "https") {
		if transport, ok := client.Client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = client.TLS
		}
	}
}

// setHeader set http request header
func (client *HttpClient) setHeader(req *http.Request, headers http.Header) {
	if headers == nil {
		headers = make(http.Header)
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = []string{"*/*"}
	}
	if _, ok := headers["Accept-Encoding"]; !ok && client.Config.Compressed {
		headers["Accept-Encoding"] = []string{"deflate, gzip"}
	}

	req.Header = headers
}

// setQueryParam set http request query string param
func (client *HttpClient) setQueryParam(req *http.Request, reqUrl string, queryParam url.Values) error {
	if queryParam != nil {
		if !strings.Contains(reqUrl, "?") {
			reqUrl = reqUrl + "?" + queryParam.Encode()
		} else {
			reqUrl = reqUrl + "&" + queryParam.Encode()
		}
		u, err := url.Parse(reqUrl)
		if err != nil {
			return err
		}
		req.URL = u
	}
	return nil
}

// setFormData set http request FormData param
func (client *HttpClient) setFormData(req *http.Request, values url.Values, setFile SetFileFunc) error {
	if setFile != nil {
		err := setFile(req, values)
		if err != nil {
			return err
		}
	} else {
		formData := []byte(values.Encode())
		req.Body = io.NopCloser(bytes.NewReader(formData))
		req.ContentLength = int64(len(formData))
	}
	return nil
}

type SetFileFunc func(req *http.Request, values url.Values) error

// File struct is a combination of file attributes
type File struct {
	Content   []byte
	Path      string
	FieldName string
	FileName  string
}

// setFile set parameters for http request formdata file upload
func setFile(f *File) SetFileFunc {
	return func(req *http.Request, values url.Values) error {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, vals := range values {
			for _, val := range vals {
				err := writer.WriteField(key, val)
				if err != nil {
					return err
				}
			}
		}

		if f.Content != nil {
			part, err := writer.CreateFormFile(f.FieldName, f.FileName)
			if err != nil {
				return err
			}
			part.Write(f.Content)
		} else if f.Path != "" {
			file, err := os.Open(f.Path)
			if err != nil {
				return err
			}
			defer file.Close()

			part, err := writer.CreateFormFile(f.FieldName, f.FileName)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return err
			}
		}

		err := writer.Close()
		if err != nil {
			return err
		}

		req.Body = io.NopCloser(body)
		req.Header.Set(ContentType, writer.FormDataContentType())
		req.ContentLength = int64(body.Len())
		return nil
	}
}

// validateRequest check if a request has url, and valid method.
func validateRequest(req *HttpRequest) error {
	if utils.IsBlank(req.RawURL) {
		return errors.New("invalid request url")
	}
	// common HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH",
		"HEAD", "CONNECT", "OPTIONS", "TRACE"}
	if !utils.IsArrayContain(strings.ToUpper(req.Method), methods) {
		return errors.New("invalid request method")
	}
	return nil
}
