# repreq

go请求回复模式
# 如何使用

服务器代码
~~~
package main
import (
	"github.com/lixxix/repreq"
)

func main() {
	serv := repreq.RepServer(nil)
	serv.Listen(":7788")
	serv.Run()
}
~~~

客户端代码
~~~
package main

import (
	"fmt"
	"time"

	"github.com/lixxix/repreq"
)

func main() {
	clie := repreq.CreateReqClient()
	clie.Connect("127.0.0.1:7788")
	clie.SetTimeout(time.Second * 3) //设置3秒没回复就自动跳过， 返回空的byte
	rec := clie.Request([]byte("whatever"))
	fmt.Println(rec)
}
~~~

# 服务器处理消息
通过接口`IRespone`对消息进行处理回复
~~~
package main

import (
	"github.com/lixxix/repreq"
)

type Respone struct {
}

func (r *Respone) OnRespone(buf []byte) []byte {
	nbuf := "my respone :" + string(buf)
	return []byte(nbuf)
}

func main() {
	serv := repreq.RepServer(&Respone{})
	serv.Listen(":7788")
	serv.Run()
}
~~~
`serv := repreq.RepServer(&Respone{})`设置了处理消息的回调方法客户端的运行结果将改变

`my respone :whatever`

服务器自动重连功能！