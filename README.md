
# serializer 
Serializer helps you to select struct's fields to export depengind on your context.

Most basic common usage is an API that returns Users' informations. Depengind whether user that make the call requests his own informations or the the informations from another one, you probably do not want to send the same fields. 
# Install

```bash
go get github.com/BorisLeMeec/serializer
```

# Example

```go
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func main() {
	user := &User{} // Retrieve user from a database, for instance.

	userJSON, _ := json.Marshal(user)
	fmt.Println(userJSON) // will print fields Username and Email.
}
```

At this point golang does not give you the possibility to return only `Username` field or both `Username` and `Email` without creating another struct.

This is the same code with `serializer` :
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/BorisLeMeec/serializer"
)

type User struct {
	Username string `json:"username" serialize:"public,private"`
	Email    string `json:"email" serialize:"private"`
	Password string `json:"-"`
}

func main() {
	user := &User{} // Retrieve user from a database, for instance.

	userJSON, _ := json.Marshal(serializer.Serialize(user, "public"))
	fmt.Println(userJSON) // will only print the field Username.
}
```

# Usages

To be able to `Serialize` a struct, you must start by adding the `serializer` tag to its fields.

By default a field without a `serializer` tag will be considered as non-desired, except for struct and Anonymous field.

# Compatibility

Right now `serializer` support most of the types.


