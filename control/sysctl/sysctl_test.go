package sysctl

import (
	"log"
	"net/url"
	"os"
	"testing"
)

// TestURLEncode url编码
func TestURLEncode(t *testing.T) {
	defer os.RemoveAll("logs")
	u := "https://blog.kms9.com/auth/qq.html"
	log.Println(url.QueryEscape(u))
}
