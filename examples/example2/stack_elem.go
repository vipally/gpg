package example2

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (this *Person) Show() string {
	return fmt.Sprintf("{Name:%s,Age:%d}", this.Name, this.Age)
}
