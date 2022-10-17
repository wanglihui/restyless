
// Code generated by file example.go and line 23 DO NOT EDIT.
// For more detail see https://github.com/wanglihui/restyless
package example
import (
	"context"
	_ "fmt"

	"github.com/wanglihui/restyless/types"
	 "github.com/go-resty/resty/v2"
	 "github.com/wanglihui/httperror"
)
func NewUserProviderImpl(r *resty.Client) UserProvider {
	return &UserProvider{
		r : r,
	}
}

type UserProviderImpl	struct {
	r *resty.Client
}

func (it *UserProvider) GetUser (ctx context.Context,userID types.QueryParam) (User,error) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	
	
	r=r.SetQueryParam("userID", string(userID))
	
	
	
	
	var ret User
	r = r.SetResult(&ret)
	_, err := r.Get("http://www.baidu.com/id")
	return ret, err
	
}

func (it *UserProvider) PostUser (ctx context.Context,user User) (User,error) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	
	
	
	r=r.SetBody(user)
	
	var ret User
	r = r.SetResult(&ret)
	_, err := r.Post("http://www.dixincaigang.cn/user/")
	return ret, err
	
}

func (it *UserProvider) DeleteUser (ctx context.Context,userID types.PathParam) (error) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	
	
	
	r=r.SetPathParam("userID", string(userID))
	
	r=r.SetBody(ctx)
	
	_, err := r.Delete("http://www.baidu.com/user/{userID}")
	return err
	
}

func (it *UserProvider) PutUser (ctx context.Context,userID types.PathParam,user User) (User,error) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	
	
	
	r=r.SetPathParam("userID", string(userID))
	
	r=r.SetBody(user)
	
	var ret User
	r = r.SetResult(&ret)
	_, err := r.Put("http://www.baidu.com/user/{userID}")
	return ret, err
	
}

func (it *UserProvider) PostUser2 (ctx context.Context,uid types.HeaderParam,token types.HeaderParam,user User) (User,error) {
	var e httperror.HTTPError
	r := it.r.R().SetError(e)
	
	r = r.SetHeader("uid", string(uid))
	
	r = r.SetHeader("token", string(token))
	
	
	
	r=r.SetBody(user)
	
	var ret User
	r = r.SetResult(&ret)
	_, err := r.Post("http://www.baidu.com/user2")
	return ret, err
	
}

