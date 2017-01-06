package db

import (
	//"time"
	//"math/rand"
	"log"
	//. "github.com/jasonknight/gopress"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//// DateTime A simple struct to represent DateTime fields
//type DateTime struct {
//	// The day as an int
//	Day int
//	// the month, as an int
//	Month int
//	// The year, as an int
//	Year int
//	// the hours, in 24 hour format
//	Hours int
//	// the minutes
//	Minutes int
//	// the seconds
//	Seconds  int
//}
type Post struct {
	ID                  int64
	PostAuthor          int64
	PostDate            string
	PostDateGmt         string
	PostContent         string
	PostTitle           string
	PostExcerpt         string
	PostStatus          string
	CommentStatus       string
	PingStatus          string
	PostPassword        string
	PostName            string
	ToPing              string
	Pinged              string
	PostModified        string
	PostModifiedGmt     string
	PostContentFiltered string
	PostParent          int64
	Guid                string
	MenuOrder           int
	PostType            string
	PostMimeType        string
	CommentCount        int64
}

//var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetFeaturedPosts() ([]Post, error)  {
	db, err := sql.Open("mysql", "root:123456@/wordpress")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return nil, err
	}
	defer db.Close()

	// Query
	rows, err := db.Query("SELECT wp_posts.* FROM wp_posts LEFT JOIN wp_term_relationships ON (wp_posts.ID = wp_term_relationships.object_id) WHERE 1=1 AND ( wp_term_relationships.term_taxonomy_id IN (20) ) AND wp_posts.post_type = 'post' AND (wp_posts.post_status = 'publish' OR wp_posts.post_status = 'private') GROUP BY wp_posts.ID ORDER BY wp_posts.post_date DESC LIMIT 0, 2")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	defer rows.Close()

	items := []Post{}
	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.ID,
			&post.PostAuthor,
			&post.PostDate,
			&post.PostDateGmt,
			&post.PostContent,
			&post.PostTitle,
			&post.PostExcerpt,
			&post.PostStatus,
			&post.CommentStatus,
			&post.PingStatus,
			&post.PostPassword,
			&post.PostName,
			&post.ToPing,
			&post.Pinged,
			&post.PostModified,
			&post.PostModifiedGmt,
			&post.PostContentFiltered,
			&post.PostParent,
			&post.Guid,
			&post.MenuOrder,
			&post.PostType,
			&post.PostMimeType,
			&post.CommentCount )
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		//log.Println("record: ", post.ID,
		//	post.PostAuthor,
		//	post.PostDate,
		//	post.PostDateGmt,
		//	post.PostContent,
		//	post.PostTitle,
		//	post.PostExcerpt,
		//	post.PostStatus,
		//	post.CommentStatus,
		//	post.PingStatus,
		//	post.PostPassword,
		//	post.PostName,
		//	post.ToPing,
		//	post.Pinged,
		//	post.PostModified,
		//	post.PostModifiedGmt,
		//	post.PostContentFiltered,
		//	post.PostParent,
		//	post.Guid,
		//	post.MenuOrder,
		//	post.PostType,
		//	post.PostMimeType,
		//	post.CommentCount)
		items = append(items, post)
	}


	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}

//func (o *Post) Find(_findByID int64) (bool, error) {
//
//	var _modelSlice []*Post
//	q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = '%d'", o._table, "ID", _findByID)
//	results, err := o._adapter.Query(q)
//	if err != nil {
//		return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
//	}
//
//	for _, result := range results {
//		ro := NewPost(o._adapter)
//		err = ro.FromDBValueMap(result)
//		if err != nil {
//			return false, o._adapter.Oops(fmt.Sprintf(`%s`, err))
//		}
//		_modelSlice = append(_modelSlice, ro)
//	}
//
//	if len(_modelSlice) == 0 {
//		// there was an error!
//		return false, o._adapter.Oops(`not found`)
//	}
//	o.FromPost(_modelSlice[0])
//	return true, nil
//
//}

//func randomInteger() int {
//	rand.Seed(time.Now().UnixNano())
//	x := rand.Intn(10000) + 100
//	if x == 0 {
//		return randomInteger()
//	}
//	return x + 100
//}

//func randomDateTime(a Adapter) *DateTime {
//	rand.Seed(time.Now().UnixNano())
//	d := NewDateTime(a)
//	d.Year = rand.Intn(2017)
//	d.Month = rand.Intn(11)
//	d.Day = rand.Intn(28)
//	d.Hours = rand.Intn(23)
//	d.Minutes = rand.Intn(59)
//	d.Seconds = rand.Intn(56)
//	if d.Year < 1000 {
//		d.Year = d.Year + 1000
//	}
//	return d
//}
//
//func randomString(n int) string {
//	rand.Seed(time.Now().UnixNano())
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}