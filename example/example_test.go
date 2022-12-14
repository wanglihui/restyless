package example

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/wanglihui/httperror"
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

func TestUserPoint(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	responser, err := httpmock.NewJsonResponder(200, User{
		ID:   "1",
		Name: "test",
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
	if user, err := userProvider.GetUserPoint(context.Background()); err != nil {
		t.Error(err)
	} else if user.ID != "1" {
		t.Fail()
	}
}

func TestUserAge(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	responser, err := httpmock.NewJsonResponder(200, 1)
	if err != nil {
		t.Error(err)
	}
	reg, err := regexp.Compile(`http://www.baidu.com/*`)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterRegexpResponder("GET", reg, responser)
	userProvider := NewUserProviderImpl(r)
	if age, err := userProvider.GetUserAge(context.Background(), "1"); err != nil {
		t.Error(err)
	} else if age != 1 {
		t.Fail()
	}
}

func TestHttpError(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	responser, err := httpmock.NewJsonResponder(400, httperror.BadRequest("test", 4000))
	if err != nil {
		t.Error(err)
	}
	reg, err := regexp.Compile(`http://www.baidu.com/*`)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterRegexpResponder("GET", reg, responser)
	userProvider := NewUserProviderImpl(r)
	if _, err := userProvider.GetUserAge(context.Background(), "1"); err == nil {
		// fmt.Println("err=>", err)
		t.Fail()
	} else {
		// fmt.Println(err)
		if err.(*httperror.HTTPError).Code != 4000 {
			t.Fail()
		}
	}
}

func TestOtherError(t *testing.T) {
	httpmock.ActivateNonDefault(r.GetClient())
	m := map[string]string{
		"code1": "10010",
		"msg":   "test error",
	}
	responser, err := httpmock.NewJsonResponder(400, m)
	if err != nil {
		t.Error(err)
	}
	reg, err := regexp.Compile(`http://www.baidu.com/*`)
	if err != nil {
		t.Error(err)
	}
	httpmock.RegisterRegexpResponder("GET", reg, responser)
	userProvider := NewUserProviderImpl(r)
	if _, err := userProvider.GetUserAge(context.Background(), "1"); err == nil {
		// fmt.Println("err=>", err)
		t.Fail()
	} else {
		// fmt.Println(err)
		if err.(*httperror.HTTPError).StatusCode != 400 {
			t.Fail()
		}
	}
}
