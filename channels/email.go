package channels

import (
	"fmt"

	"github.com/benmanns/goworker"
)

func init() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
