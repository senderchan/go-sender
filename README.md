# Sender Agent

这是一个信达的信令API的实现。这个程序将自动拉取信令，并根据配置文件执行对应的cmd命令。

## 获取Agent

运行前请确认已经安装`wget`和`curl`。

```shell
sudo sh -c "$(curl -sSL https://raw.githubusercontent.com/github.com/senderchan/go-sender/master/install.sh)"
```

第一次执行`sender`命令将在工作路径创建一个配置文件。在这个配置文件中设置Agent name和Access
Key。
```yaml
name: test
accesskey: 758e2e934f**************9ebafcb

action:
  hello:
    cmd: echo "hello"
  ls:
    cmd: ls -la
    dir: ~/
  aaa:
    script: ./run.sh
    dir: ~/
  send:
    cmd: echo "ok"
    forward: http://localhost:9000
  default:
    cmd: docker version
```

再次运行`sender`后Agent即开始运行。

使用`sender service`将在`/etc/systemd/system/sender.service`创建service配置文件。

```shell
sudo systemctl daemon-reload
sudo systemctl enable sender
sudo systemctl start sender
```

## SDK
```go
package main

import (
	"github.com/senderchan/go-sender/agent"
)

func main() {
	agent.SendMessage("1f7549f3edaf597fa770cca2d26ecbad", "hello", "", "", "")

	a := agent.Agent{
		AccessKey: "****************",
		Name:      "test",
	}

	a.Run(func(s []agent.Signaling) {
		for _, i := range s {
			i.Result = "abababababab"
			a.SubmitSignalingResult(i)
		}
	})
}

```
