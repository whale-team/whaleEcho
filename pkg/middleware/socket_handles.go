package middleware

import "github.com/whale-team/whaleEcho/pkg/wsserver"

func WsErrorHandle(c *wsserver.Context, err error) {

}

func WsConnBuildHandle(c *wsserver.Context) {

}

func WsConnCloseHandle(c *wsserver.Context) error {
	return nil
}
