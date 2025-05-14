package casbinz

//  "github.com/iris-contrib/middleware/casbin"
// 这里对casbin middleware 改造成我的方式
// 每个用户本身已经缓存了pas roles ，所以不用每次都来查表了

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/casbinz/enum"
	"github.com/aarioai/golib/sdk/casbinz/models"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"log"
	"regexp"
	"strconv"
)

var (
	readHttpMethods  = map[string]struct{}{"GET": {}, "HEAD": {}, "OPTIONS": {}, "TRACE": {}, "CONNECT": {}}
	writeHttpMethods = map[string]struct{}{"PUT": {}, "DELETE": {}, "POST": {}, "PATCH": {}}
)

func HttpMethodToAct(method string) enum.Act {
	if _, ok := readHttpMethods[method]; !ok {
		return enum.Read
	}
	if _, ok := writeHttpMethods[method]; !ok {
		return enum.Write
	}
	return enum.UnknownAct
}

func HandleHttpPath(path string) string {
	re, _ := regexp.Compile("^.*/v1/(.+)")
	re.FindString(path)
	return path
}

func init() {
	context.SetHandlerName("github.com/iris-contrib/middleware/casbin.*", "iris-contrib.casbin")
}

// Middleware is the auth service which contains the casbin enforcer.
type Middleware struct {
	enforcer *casbin.Enforcer
	// RolesExtractor is used to extract the
	// current request's subject for the casbin role enforcer.
	// Defaults to the `Roles` package-level function which
	// extracts the subject from a prior registered authorization middleware's
	// username (e.g. basicauth or JWT).
	RolesExtractor func(iris.Context) []string

	// UnauthorizedHandler sets a custom handler to be executed
	// when the role checks fail.
	// Defaults to a handler which sends a status forbidden (403) status code
	UnauthorizedHandler iris.Handler
}

// NewHttpMiddleware returns the middleware based on the given casbin.Enforcer instance.
// The authorization determines a request based on `{subject, object, action}`.
// Please refer to: https://github.com/casbin/casbin to understand how it works first.
//
// The object is the current request's path and the action is the current request's method.
// The subject that casbin requires is extracted by:
//   - RolesExtractor
//   - casbin.Roles
//     | set with casbin.SetSubject
//   - Context.VUser().GetUsername()
//     | by a prior auth middleware through Context.SetUser.
func NewHttpMiddleware(e *casbin.Enforcer) *Middleware {
	return &Middleware{
		enforcer: e,
		RolesExtractor: func(ictx iris.Context) []string {
			return Roles(ictx)
		},
		UnauthorizedHandler: func(ictx iris.Context) {
			_, resp, _ := httpsvr.New(ictx)
			resp.WriteE(ae.ErrorForbidden)
		},
	}
}

// ServeHTTP is the iris compatible casbin handler which should be passed to specific routes or parties.
// Responds with Status Forbidden on unauthorized clients.
// Usage:
// - app.Use(authMiddleware)
// - app.Use(casbinMiddleware.ServeHTTP) OR
// - app.UseRouter(casbinMiddleware.ServeHTTP) OR per route:
// - app.Get("/dataset1/resource1", casbinMiddleware.ServeHTTP, myHandler)
func (c *Middleware) ServeHTTP(ictx iris.Context) {
	if !c.Check(ictx) {
		c.UnauthorizedHandler(ictx)
		return
	}

	ictx.Next()
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
//
// It's an Iris Filter.
// Usage:
// - inside a handler
// - using the iris.NewConditionalHandler
func (c *Middleware) Check(ictx iris.Context) bool {
	roles := c.RolesExtractor(ictx)
	var ok bool
	var err error
	for _, role := range roles {
		obj := HandleHttpPath(ictx.Path())
		act := HttpMethodToAct(ictx.Method())
		ok, err = c.enforcer.Enforce(role, obj, act)
		if ok {
			return ok
		}
		if err != nil {
			log.Println(err.Error())
		}
	}
	return false
}

func Roles(ictx iris.Context) []string {
	uid, _ := ictx.Values().GetUint64(enumz.IctxParamUid)
	if uid == 0 {
		return nil
	}
	return []string{strconv.FormatUint(uid, 10)}
}

func NewDefaultHttpMiddleware(a persist.Adapter) (*Middleware, error) {
	m, err := model.NewModelFromString(models.DefaultModel)
	if err != nil {
		return nil, err
	}
	enf, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	// 可以自定义 match 函数
	//  enf.AddFunction("actionMatch", matcher.ActionMatch)
	//  enf.AddFunction("objMatch", matcher.ObjectMatch)
	return NewHttpMiddleware(enf), nil
}
