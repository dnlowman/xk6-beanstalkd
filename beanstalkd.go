package beanstalkd

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/beanstalkd", new(Beanstalkd))
}

type Beanstalkd struct{}

func (*Beanstalkd) NewClient(addr string) (*Client, error) {
	return NewClient(addr)
}
