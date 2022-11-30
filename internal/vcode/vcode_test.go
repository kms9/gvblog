package vcode

import (
	"fmt"
	"testing"

	"github.com/kms9/gvblog/utils"
)

func TestVcode(t *testing.T) {
	s := "12345"
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}
	fmt.Println(NewImage(utils.RandDigitStr(6)).Base64())
}
