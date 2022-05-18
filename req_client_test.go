package repreq

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := CreateReqClient()
	if client != nil {
		err := client.Connect("127.0.0.1:7777")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		rec := client.Request([]byte("hello"))
		panic(string(rec))
	}
}
