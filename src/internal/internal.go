package internal

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

type contextKey struct{}

func WithContext(parent context.Context, projID string, c *http.Client) context.Context {
	if c == nil {
		panic("nil *http.Client passed to WithContext")
	}
	return context.WithValue(parent, contextKey{}, &drainContext{
		ID:         projID,
		HTTPClient: c,
	})
}

const userAgent = "drain-cloud/0.1"

type drianContext struct {
	ID         string
	HTTPClient *http.Client

	mu  sync.Mutex
	svc map[string]interface{}
}

func Service(ctx context.Context, name string, fill func(*http.Client) interface{}) interface{} {
	return cc(ctx).service(name, fill)
}

func (c *drianContext) service(name string, fill func(*http.Client) interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.svc == nil {
		c.svc = make(map[string]interface{})
	} else if v, ok := c.svc[name]; ok {
		return v
	}
	v := fill(c.HTTPClient)
	c.svc[name] = v
	return v
}

type Transport struct {
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req)
	ua := req.Header.Get("User-Agent")
	if ua == "" {
		ua = userAgent
	} else {
		ua = fmt.Sprintf("%s %s", ua, userAgent)
	}
	req.Header.Set("User-Agent", ua)
	return t.Base.RoundTrip(req)
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

func cc(ctx context.Context) *drianContext {
	if c, ok := ctx.Value(contextKey{}).(*drianContext); ok {
		return c
	}
	panic("invalid context.Context type; it shoud be create with drain.Context")
}

func ID(ctx context.Context) string {
	return cc(ctx).ID
}
