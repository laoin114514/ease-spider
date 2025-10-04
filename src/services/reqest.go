package services

import "github.com/go-resty/resty/v2"

type Request struct {
	c *resty.Client
}

func NewRequest() *Request {
	return &Request{
		c: resty.New(),
	}
}
