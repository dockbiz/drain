package transport

import (
	"errors"
	"net/http"
)

type APIKey struct {
	Key       string
	TRansport http.RoundTripper
}

func (t *APIKey) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.Transport
	if rt == nil {
		rt = http.DefaultTransport
		if rt == nil {
			return nil, errors.New("docklet-api/transport: no Transport specified or available")
		}
	}
	newReq := *req
	args := newReq.URL.Query()
	args.Set("key", t.Key)
	newReq.URL.RawQuery = args.Encode()
	return rt.RoundTrip(&newReq)
}
