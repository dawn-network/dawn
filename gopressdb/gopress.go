package gopressdb

import (
	"log"
	"github.com/asdine/storm"
	"time"

)
var Gdb *storm.DB

type User struct {
	ID			int `"storm:id,index"`
	Username 		string `"storm:index"`
	Email			string `"storm:unique"`
	lat			[]int
	long			[]int
	GpsInterval		int
	Pubkey  		string
	UserRegistered 		string
	DisplayName 		string
	CreatedAt		time.Time `"storm:index"`

}

type Post struct {
	ID                  	int `"storm:id,index"`
	PostAuthor          	[]string `"storm:index"`
	PostAuthorSplit		[]int
	CreatedAt            	time.Time `"storm:index"`
	PostContent         	string
	Hits			[]byte `"storm:index"`
	HitLat			[]float64  `"storm:index"`
	HitLong			[]float64 `"storm:index"`
	PostTitle		string `"storm:index"`
	ParentPostID		int `"storm:index"`
	ForkOf			bool
	PostType		string `"storm:index"`
	Parent			string	`"storm:index"`
	PostModified        	string	`"storm:index"`
	Thumb 		    	string  `"storm:index"`
	Cat 			[]string 	`"storm:index"`	// category
	Lat			int
	Long			int
}



// ---------------------------------------------------------------------------------------------------------------------

func GetDB() (*storm.DB) {
	ldb, err := storm.Open("my.db", storm.Batch(),  )
	Gdb = ldb
	defer Gdb.Close()
	if err != nil {
		log.Println(err)
	}
	return Gdb
}


// ---------------------------------------------------------------------------------------------------------------------
// User

//GetUser allows
func GetUser(ID string) (User, error)  {
	var user User
	err := Gdb.Find("ID", ID, &user)
	return user, err
}

//CreateUser allows non-users to become users
func CreateUser(user User) error {
	err := Gdb.Save(&user)
	return err
}

// ---------------------------------------------------------------------------------------------------------------------
// Post

//CreatePost allows users to create new posts
func CreatePost(post Post) error {
	err := Gdb.Save(&post)
	return err
}

//EditPost allows users to edit their posts.
func EditPost(post Post) error {
	err := Gdb.Update(&post)
	return err
}

//GetPost will return posts with a certain ID (ID=IPFS Hash)
func GetPost(postID int) (post Post)  {
	Gdb.Find("ID", postID, &post)
	return post
}
//GetPostsByCategory will return posts of N category
func GetPostsByCategory(category string) (posts []Post)  {
	Gdb.Find("cat", category, &posts )
	return posts
}
