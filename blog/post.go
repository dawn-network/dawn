package blog

import (
	"fmt"
	"os"
	"glogchain/protocol"
	"github.com/extemporalgenome/slug"
	"glogchain/config"
)

// we use hugo blog
// to create a post, just create a file in the content/post folder of the hugo site
// the content of the file must be in hugo format

//const HUGO_CONTENT_POST_PATH = "/Users/tuanpa/Projects/glogchain/hugocontent/bookshelf/content/"

func CreatePost(post *protocol.PostOperation) error {
	// use lib https://github.com/extemporalgenome/slug
	slugTitle := slug.Slug(post.Title)

	f := createFile(config.GlogchainConfigGlobal.HugoPostPath + slugTitle + ".md")

	defer f.Close()

	fmt.Fprintln(f, "+++")
	fmt.Fprintln(f, "title = \"" + post.Title + "\"")
	fmt.Fprintln(f, "draft = true")
	fmt.Fprintln(f, "date = \"2016-10-22T12:50:18+07:00\"")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "+++")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, post.Body)

	return nil
}

func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

