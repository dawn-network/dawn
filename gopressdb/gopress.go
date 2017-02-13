package db

import (
	"log"
	"github.com/asdine/storm"
	"github.com/blevesearch/bleve"
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


// ---------------------------------------------------------------------------------------------------------------------

func GetDB() (*storm.DB){
	db, err := storm.Open("my.db", storm.Batch())
	defer db.Close()
	if err != nil {
		log.Println(err)
	}
	return db
}

func Bleve() (bleve.Index) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("bleve", mapping)
	if err != nil {
		panic(err)
	}
	return index
}


// ---------------------------------------------------------------------------------------------------------------------
// User

//GetUser allows users and non-users to search for users
func GetUser(ID string) (User, error)  {
	var user User
	db := GetDB()
	err := db.Find("ID", ID, &user)
	return user, err
}



//CreateUser allows non-users to become users
func CreateUser(user User) {
	db := GetDB()
	index := Bleve()
	index.Index(user.ID, user)
	db.Save(&user)
}

// ---------------------------------------------------------------------------------------------------------------------
// Post

//CreatePost allows users to create new posts
func CreatePost(post Post)  {
	db := GetDB()
	index := Bleve()
	db.Save(&post)
	index.Index(post.ID, post)
}

//EditPost allows users to edit their posts.
func EditPost(post Post)  {
	db := GetDB()
	index := Bleve()
	db.Update(post)
	
	index.Index(post.ID, post)
}

//GetPost will return posts with a certain ID (ID=IPFS Hash)
func GetPost(postID string) (post Post)  {
	db := GetDB()
	db.Find("ID", postID, post)
	return post
}
//GetPostsByCategory will return posts of N category
func GetPostsByCategory(category string) (posts []Post)  {
	db := GetDB()
	db.Find("cat", category, posts )
	return posts
}