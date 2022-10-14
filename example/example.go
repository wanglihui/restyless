package example

// 1. 模版代码
// 2. 服务提供的功能，如何在客户端直接调用，而不是重新实现一次

import (
	"context"
	_ "fmt"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UserProvider
// host=http://www.baidu.com
// params named header will put into http.header, params named query will put into querystring
// body will put into body
//
//go:generate go run github.com/wanglihui/restyless
type UserProvider interface {
	//GetUser
	//host=http://www.baidu.com,url=/id
	GetUser(ctx context.Context, userID string) (User, error)
	//PostUser
	//host=http://www.dixincaigang.cn,url=/user/
	PostUser(ctx context.Context, body User) (User, error)
	//DeleteUser
	//url=/user/:id
	DeleteUser(ctx context.Context, userID string) error
	//PutUser
	//url=/user/:id
	PutUser(ctx context.Context, user User) (User, error)
	//PostUser2
	//url=/user2
	PostUser2(ctx context.Context, headers string, user User) (User, error)
}
