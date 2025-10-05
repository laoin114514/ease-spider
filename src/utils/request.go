package utils

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Request[T any] struct {
	c *resty.Client
}

func NewRequest[T any]() *Request[T] {
	return &Request[T]{
		c: resty.New(),
	}
}
func (r *Request[T]) Get(url string, params map[string]string) (T, error) {
	var result T
	cookie := JsonDB.Get("Cookie").(string)
	resp, err := r.c.R().
		SetHeader("Cookie", cookie).
		SetQueryParams(params).
		Get(url)
	if err != nil {
		return result, err
	}
	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("请求失败  %d", resp.StatusCode())
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return result, err
	}
	// fmt.Println("[Debug]", string(resp.Body()))
	return result, nil
}
