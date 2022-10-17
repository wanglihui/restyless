package example

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

var r = resty.New()

// // 1. 模版代码
// // 2. 服务提供的功能，如何在客户端直接调用，而不是重新实现一次
// func (it *UserProviderImp) GetUser(ctx context.Context, userID string) (User, error) {
// }
func TestSome(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	responser, err := httpmock.NewJsonResponder(200, map[string]string{
		"id":   "1",
		"name": "test",
	})
	if err != nil {
		t.Error(err)
	}
	reg, err := regexp.Compile(`http://www.baidu.com/*`)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterRegexpResponder("GET", reg, responser)
	userProvider := NewUserProviderImpl(r)
	if _, err := userProvider.GetUser(context.Background(), "1"); err != nil {
		t.Error(err)
	}
}

func TestUserList(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	responser, err := httpmock.NewJsonResponder(200, []User{{
		ID:   "1",
		Name: "test",
	},
	})
	if err != nil {
		t.Error(err)
	}
	reg, err := regexp.Compile(`http://www.baidu.com/*`)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterRegexpResponder("GET", reg, responser)
	userProvider := NewUserProviderImpl(r)
	if users, err := userProvider.GetUsers(context.Background()); err != nil {
		t.Error(err)
	} else if len(users) != 1 {
		t.Fail()
	}
}
