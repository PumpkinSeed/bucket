package bucket

import (
	"fmt"
	"testing"
)

func TestHandler_GetSubDocuments(t *testing.T) {
	var ws = generate()

	resultset := th.getSubDocuments("webshop", ws)
	for k, v := range resultset {
		fmt.Printf("k: %s, v: %+v\n", k, v)
	}
}
