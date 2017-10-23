## Qick Start for Learining How OpenSDS Client Works

To better learn how opensds client works for connecting with OpenSDS service,
here is a three-step example showing how to use the client.

Before these three steps, you have to make sure client package has been imported
in your local repo.

### Step 1: Initialize Client object
It's a simple and easy step for user to create a new Client object like below:
```go
package main

import (
	"fmt"
	
	"github.com/opensds/opensds/client"
)

func CreateClient() {
	c1 := NewClient(&client.Config{})
	c2 := NewClient(&client.Config{
		Endpoint: ":8080",
	})
	
	fmt.Printf("c1 is %v, c2 is %v\n", c1, c2)
}
```
As you can see from code above, user has two ways to create ```Client``` object:
parsing ```Config``` object or fetching the endpoint from environment variable
(```os.Getenv("OPENSDS_ENDPOINT")```), you can choose one with your reference.

### Step 2: Update request content in Client object
It's also easy to understand by looking the example below:
```go
package main

import (
	"fmt"
	
	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
)

func main() {
	c := NewClient(&client.Config{
		Endpoint: ":8080",
	})
	
	c.UpdateRequestContent("volume", &modelVolumeSpec{Name: "test"})
	fmt.Printf("Updated c is %v\n", c)
}
```

### Step 3: Call method in Client object
In the last step, you can just call method in Client object which is created
in step 2 like this:
```go
vol, err := c.CreateVolume()
if err != nil {
	panic(err)
}
```