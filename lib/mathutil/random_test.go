package mathutil

import (
	"fmt"
	"testing"
)

func TestGenRandomInt(t *testing.T) {
	result := GenRandomInt(0, 0)
	fmt.Println(result)
}
