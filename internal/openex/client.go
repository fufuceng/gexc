package openex

import (
	"encoding/json"
	"fmt"
	rsp "github.com/fufuceng/gexc/response"
	"github.com/fufuceng/gexc/time"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	Latest(params LatestParams) (*rsp.SingleDate, error)
	SingleDate(params SingleDateParams) (*rsp.SingleDate, error)
	History(params HistoryParams) (*rsp.History, error)
}

type httpGetter func(url string) (*http.Response, error)

var defaultHttpGetter = http.Get

type client struct {
	config     Config
	httpGetter httpGetter
}

type response struct {
	Body []byte
	Code int
}

func (c client) toUrl(path string, qp string) url.URL {
	return url.URL{
		Scheme:   c.config.Protocol,
		Host:     c.config.BaseUrl,
		Path:     path,
		RawQuery: qp,
	}
}

func (c client) doRequest(method string, url url.URL) (*response, error) {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		resp, err := c.httpGetter(url.String())
		if err != nil {
			return nil, err
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &response{
			Body: body,
			Code: resp.StatusCode,
		}, nil

	default:
		panic("unsupported method: " + method)
	}
}

func (c client) bindJson(from []byte, to interface{}) error {
	return json.Unmarshal(from, to)
}

func (c client) parseResp(resp *response, successResp interface{}) error {
	if resp.Code != http.StatusOK {
		var errResp errorResponse
		if err := c.bindJson(resp.Body, &errResp); err != nil {
			return fmt.Errorf("error while binding response: %v", err)
		}

		return fmt.Errorf("request failed: %v", errResp.Error)
	}

	if err := c.bindJson(resp.Body, &successResp); err != nil {
		return fmt.Errorf("error while binding response: %v", err)
	}

	return nil
}

func (c client) Latest(params LatestParams) (*rsp.SingleDate, error) {
	qp, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(http.MethodGet, c.toUrl("/latest", qp.Encode()))
	if err != nil {
		return nil, err
	}

	var singleDateResponse rsp.SingleDate
	if err := c.parseResp(resp, &singleDateResponse); err != nil {
		return nil, err
	}

	return &singleDateResponse, nil
}

func (c client) SingleDate(params SingleDateParams) (*rsp.SingleDate, error) {
	qp, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(http.MethodGet, c.toUrl("/"+params.Date.Format(time.GexcLayout), qp.Encode()))
	if err != nil {
		return nil, err
	}

	var singleDateResponse rsp.SingleDate
	if err := c.parseResp(resp, &singleDateResponse); err != nil {
		return nil, err
	}

	return &singleDateResponse, nil
}

func (c client) History(params HistoryParams) (*rsp.History, error) {
	qp, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	// add start_at and end_at parameters explicitly
	// go-querystring does not expose any interfaces to implement
	// custom types
	qp.Set("start_at", params.StartAt.String())
	qp.Set("end_at", params.EndAt.String())

	resp, err := c.doRequest(http.MethodGet, c.toUrl("/history", qp.Encode()))
	if err != nil {
		return nil, err
	}

	var histResponse rsp.History
	if err := c.parseResp(resp, &histResponse); err != nil {
		return nil, err
	}

	return &histResponse, nil
}

func NewDefaultClient() Client {
	return &client{
		config:     defaultConfig,
		httpGetter: defaultHttpGetter,
	}
}
