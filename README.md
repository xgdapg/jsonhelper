# jsonhelper
jsonhelper

```go
package main

import (
	"fmt"
	"github.com/xgdapg/jsonhelper"
)

func main() {
	s := []byte(`{"a":1,"b":[1,2],"c":true,"d":"asdf"}`)
	n, err := jsonhelper.Parse(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n.Key("b").Index(1).ToInt()) // 2 <nil>
	fmt.Println(n.Key("f").Index(1).ToInt()) // 0 Key `f` not exist
}
```