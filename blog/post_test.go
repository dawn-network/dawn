package blog

import (
	"testing"
)

func TestCreatePost(t *testing.T) {
	err := CreatePost(nil)

	if err != nil {
		t.Fatal(err)
	}
}