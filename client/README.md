## Quick Start for Learning How OpenSDS Client Works

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

func main() {
	c1, _ := client.NewClient(&client.Config{})
	c2, _ := client.NewClient(&client.Config{
		Endpoint: ":50040",
	})
	
	fmt.Printf("c1 is %v, c2 is %v\n", c1, c2)
}
```
As you can see from code above, user has two ways to create ```Client``` object:
parsing ```Config``` object or fetching the endpoint from environment variable
(```os.Getenv("OPENSDS_ENDPOINT")```), you can choose one with your reference.

### Step 2: Call method in Client object
In the second step, you can just call method in Client object which is created
in step 1 like this:
```go
package main

import (
	"fmt"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
)

func main() {
	c, _ := client.NewClient(&client.Config{
		Endpoint:    ":50040",
		AuthOptions: client.LoadNoAuthOptionsFromEnv(),
	})

	vol, err := c.CreateVolume(&model.VolumeSpec{
		Name:        "test",
		Description: "This is a volume for test",
		Size:        int64(1),
	})
	if err != nil {
		fmt.Println(err)
	}

	result, err := c.GetVolume(vol.Id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Volume created, get result:", result)

	if err = c.DeleteVolume(vol.Id, nil); err != nil {
		fmt.Println(err)
	}
}
```

### Step 3: Destory Client object
If you want to reset the Client object, just run ```c.Reset()``` and it will
clear all data in it and return a empty object.
