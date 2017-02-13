package db

import (
	"log"
	"github.com/asdine/storm"
	"time"
)

type User struct {
	ID			string `"storm:id,index"`
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
	ID                  	string `"storm:id,index"`
	PostAuthor          	string `"storm:index"`
	Author			string
	Date			string
	PostDate		string
	PostAuthorSplit		[]int
	PostID			string
	CreatedAt            	time.Time `"storm:index"`
	PostContent         	string
	Hits			[]byte `"storm:index"`
	HitLat			[]float64  `"storm:index"`
	HitLong			[]float64 `"storm:index"`
	PostTitle		string `"storm:index"`
	ParentPostID		string `"storm:index"`
	ForkOf			string
	PostType		string `"storm:index"`
	Parent			string	`"storm:index"`
	PostModified        	string	`"storm:index"`
	Thumb 		    	string  `"storm:index"`
	Cat 			string 	`"storm:index"`	// category
	Lat			int
	Long			int
}

var Gdb *storm.DB
var err error

// ---------------------------------------------------------------------------------------------------------------------

func GetDB() {
	Gdb, err = storm.Open("my.db", storm.Batch())
	defer Gdb.Close()
	if err != nil {
		log.Println(err)
	}
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
func CreateUser(user User) {
	Gdb.Save(&user)
}

// ---------------------------------------------------------------------------------------------------------------------
// Post

//CreatePost allows users to create new posts
func CreatePost(post Post)  {
	Gdb.Save(&post)
}

//EditPost allows users to edit their posts.
func EditPost(post Post)  {
	Gdb.Update(&post)
}

//GetPost will return posts with a certain ID (ID=IPFS Hash)
func GetPost(postID string) (post Post)  {
	Gdb.Find("ID", postID, &post)
	return post
}
//GetPostsByCategory will return posts of N category
func GetPostsByCategory(category string) (posts []Post)  {
	Gdb.Find("cat", category, &posts )
	return posts
}