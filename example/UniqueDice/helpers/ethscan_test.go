package helpers_test

import (
	"UniqueDice/helpers"
	"fmt"
	"testing"
)

func Test_spider(t *testing.T) {
	tokens, e := helpers.GetAllERC20Tokens()
	if e != nil {
		t.Error(e)
	}
	fmt.Println(tokens)

	t.Error("tokens")
}
