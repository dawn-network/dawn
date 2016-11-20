package protocol

import (
	"testing"
	"fmt"
	"reflect"
)


func TestUnmarshal(t *testing.T) {
	const jsonstring = `
{
	"type": "PostOperation",
	"Operation": {
		"Title": "the Title",
		"Body": "the Body",
		"Author": "the Author"
	}
}
`
	obj, err := UnMarshal(jsonstring)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(reflect.TypeOf(obj))

	t.Log("TestUnmarshal Completed")

}