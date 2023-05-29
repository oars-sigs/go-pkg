# flow

编写程序`main.go`

```go

package main

import (
	"fmt"
	"os"

	"pkg.oars.vip/go-pkg/flow"
)

type Sum struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
}

func (a *Sum) Do(conf *Config, params interface{}) (interface{}, error) {
	args := params.(Sum)
	return args.X + args.Y, nil
}
func (a *Sum) Params() interface{} {
	return Sum{}
}

func (a *Sum) Scheme() string {
	return ""
}

func main() {
	flow.AddCustomActions("sum", &Sum{})
	err := flow.Run(os.Args[1])
	fmt.Println(err)
}

```

编写脚本`test.yaml`

```yaml

tasks:
- sum:
    x: 1
    "y": 1
  output: $.values.sum
- print: $.values.sum

```

运行脚本

```
go run main.go test.yaml

```

