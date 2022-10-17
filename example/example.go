package example

// 1. 模版代码
// 2. 服务提供的功能，如何在客户端直接调用，而不是重新实现一次

import (
	"context"
	_ "fmt"

	"github.com/wanglihui/restyless/types"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserProvider
// host=http://www.baidu.com
// params named header will put into http.header, params named query will put into querystring
// body will put into body
//
//go:generate go run github.com/wanglihui/restyless UserProvider
type UserProvider interface {
	//GetUser
	//host=http://www.baidu.com,url=/id
	GetUser(ctx context.Context, userID types.QueryParam) (User, error)
	//PostUser
	//host=http://www.dixincaigang.cn,url=/user/
	PostUser(ctx context.Context, user User) (User, error)
	//DeleteUser
	//url=/user/{userID}
	DeleteUser(ctx context.Context, userID types.PathParam) error
	//PutUser
	//url=/user/{userID}
	PutUser(ctx context.Context, userID types.PathParam, user User) (User, error)
	//PostUser2
	//url=/user2
	PostUser2(ctx context.Context, uid types.HeaderParam, token types.HeaderParam, user User) (User, error)
	//url=/users
	GetUsers(ctx context.Context) ([]User, error)
	//url=/user/{userId}/points
	GetUserPoint(ctx context.Context) (*User, error)
	//url=/user/{userId}/age
	GetUserAge(ctx context.Context, userId types.PathParam) (int, error)
	//url=/user/map
	PostUserUseMap(ctx context.Context, user map[string]string) (User, error)
}
