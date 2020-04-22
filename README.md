# zkutils
This is zookeeper(zk) utility module to make "go-zookeeper" easier to use. 

Import it in your program as:
```go
      import "github.com/pochard/zkutils"
```

## API
### pakcage func
	func NewKeepWatcher(conn *zk.Conn) *KeepWatcher

### type KeepWatcher (Keep watching a specified zk node in an infinite loop)
	//listener func will be called whenever children of the path chanaged
	WatchChildren(path string, listener func(children []string, err error)) 
	//listener func will be called whenever data of the path chanaged
	WatchData(path string, listener func(data []byte, err error))

## KeepWatcher Example
```go
package main

import (
	"fmt"
	"github.com/pochard/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"net/http"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	var hosts = []string{"localhost:2181"}
	conn, _, _ := zk.Connect(hosts, time.Second*5)
	defer conn.Close()

	keepWatcher := zkutils.NewKeepWatcher(conn)
	path := "/hades/services/main/cms"
	go keepWatcher.WatchChildren(path, func(children []string, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", strings.Join(children, ","))
	})

	path = "/hades/configs/main/cms/config"
	go keepWatcher.WatchData(path, func(data []byte, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s\n", string(data))
	})

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}


```