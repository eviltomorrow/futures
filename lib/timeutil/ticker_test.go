package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTicker(t *testing.T) {
	ticker := NewAlignedTicker(time.Now(), 5*time.Second, 0, 0)
	go func() {
		for {
			now := time.Now()
			t := <-ticker.Elapsed()
			fmt.Println(now.Format(time.DateTime), "--", t.Format(time.DateTime))
		}
	}()
	select {}
}
