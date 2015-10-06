package api

import (
	"bytes"
	"errors"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"sync"
	"time"
)

const (
	Version     = "0.1"
	UserAgent   = "docklet-api-client/" + Version
	uploadPause = 1 * time.Second
)

var (
	rangeRE         = regexp.MustCompile(``)
	chunkSize int64 = 1 << 18
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Body    string
	Errors  []ErrorItem
}

type ErrorItem struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	if len(e.Errors) == 0 && e.Message == "" {
		return
	}
}
func Progress() int64 {

}
