package channel

import (
	"golang.org/x/net/context"
)

func Create(c context.Context, clientID string) (token string, err error) {

}

func Query(conns []Conn, query string) Result {
	ch := make(chan Result, len(conns))
	for _, conn := range conns {
		go func(c Conn) {
			ch <- c.DoQ(query)
		}(conn)
	}
}
