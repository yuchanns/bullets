# bullets
there is no silver bullet but sometimes may be helpful,
## usage
* gin middlewares: print request and handle panic
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/yuchanns/bullet/common/middlewares"
)

func main() {
    engine := gin.New()
    engine.Use(
        middlewares.NewDefaultRequestInterceptor(),
        middlewares.NewDefaultRequestInterceptor(),
    )
}
```
* log
```go
package main

import "github.com/yuchanns/bullet/common"

func main() {
    common.Logger.
        Fields(map[string]interface{}{"foo": "bar"}).
    	DebugInfo("debug")
}
```
* return json response and log regular error
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
    "github.com/yuchanns/bullet/common"
    "os"
)

func openNotExistFile() error {
	if _, err := os.Open("NOT_EXIST_FILE"); err != nil {
        // must wrap error with pkg/errors
		return errors.Wrap(err, "open error")
	}
	return nil
}

func OpenNotExistFileHandler(ctx *gin.Context) {
	if err := openNotExistFile(); err != nil {
		common.JsonFailWithStack(ctx, err, nil)
		return
	}
	common.JsonSuccess(ctx, "success", nil)
}
```
