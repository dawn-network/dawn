package db

import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
	"encoding/json"
)

type User struct {
	ID			string
	Username 		string
	Pubkey  		string
	UserRegistered 		string
	DisplayName 		string
}

type Post struct {
	ID                  	string
	PostAuthor          	string
	PostDate            	string
	PostContent         	string
	PostTitle           	string
	PostModified        	string
	Thumb 		    	string
	Cat 			string 		// category
}

type Comment struct {
	ID                  	string
	PostID                 	string 		// which post belong to
	Parent 			string 		// parrent comment, empty string means no parent
	Author          	string		// ID/Address of Account/User
	Date            	string 		// create datetime
	Content         	string
	Modified        	string 		// Last Modified datetime
}

type Category struct {
	ID    			string
	Count 			int64
}

type TermRelationship struct {
	ObjectId       		int64
	TermTaxonomyId 		int64
	TermOrder      		int
}

// ---------------------------------------------------------------------------------------------------------------------

var __db *sql.DB = nil

func GetDB() (*sql.DB, error) {
	if (__db == nil) {
		db, err := sql.Open("sqlite3", "db.db")

		if err != nil {
			panic(err)
		}

		if db == nil {
			panic("db nil")
		}

		////////////
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tbl_cat
		(
		    ID VARCHAR(200) NOT NULL PRIMARY KEY,
		    count INTEGER NOT NULL DEFAULT '0'
		);`)
		if err != nil {
			panic(err)
		}

		//////
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wp_users
		(
		    ID VARCHAR(40) NOT NULL PRIMARY KEY,
		    user_username VARCHAR(60) NOT NULL DEFAULT '',
		    user_pubkey VARCHAR(255) NOT NULL DEFAULT '',
		    user_registered DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
		    display_name VARCHAR(250) NOT NULL DEFAULT ''
		);`)
		// balance INTEGER NOT NULL DEFAULT '0'
		if err != nil {
			panic(err)
		}

		//////
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wp_posts
		(
		    ID VARCHAR(40) NOT NULL PRIMARY KEY,
		    post_author VARCHAR(40) NOT NULL DEFAULT '0',
		    post_date DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
		    post_content LONGTEXT NOT NULL,
		    post_title TEXT NOT NULL,
		    post_modified DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
		    thumb VARCHAR(255) NOT NULL DEFAULT '',
		    cat VARCHAR(255) NOT NULL DEFAULT ''
		);`)
		if err != nil {
			panic(err)
		}

		//////
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tbl_comments
		(
		    ID VARCHAR(40) NOT NULL PRIMARY KEY,
		    postID VARCHAR(40) NOT NULL DEFAULT '0',
		    cm_parent VARCHAR(40) NOT NULL DEFAULT '',
		    cm_author VARCHAR(40) NOT NULL DEFAULT '0',
		    cm_date DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
		    cm_content LONGTEXT NOT NULL,
		    cm_modified DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00'
		);`)
		if err != nil {
			panic(err)
		}

		__db = db
	}

	return __db, nil
}

func Query (sql string) (*sql.Rows, error) {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return nil, err
	}
	//defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}

	return rows, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// User

func GetUser(ID string) (User, error)  {
	var item User

	sql := fmt.Sprintf("SELECT * FROM wp_users WHERE ID=\"%s\"", ID)

	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return item, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&item.ID,
			&item.Username,
			&item.Pubkey,
			&item.UserRegistered,
			&item.DisplayName)
		if err != nil {
			log.Println("GetUser", err)
			return item, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return item, err
	}

	return item, nil
}

func CreateUser(user User) error {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return err
	}
	//defer db.Close()


	_, err = db.Exec("INSERT INTO wp_users(ID, user_username, user_pubkey, user_registered, display_name) " +
		"VALUES(?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Pubkey, user.UserRegistered, user.DisplayName)

	//if err != nil {
	//	return err
	//}

	return err
}

// ---------------------------------------------------------------------------------------------------------------------
// Post

func CreatePost(post Post) error {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return err
	}
	//defer db.Close()


	_, err = db.Exec("INSERT INTO wp_posts(ID, post_author, post_date, post_content, post_title, post_modified, thumb, cat) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		post.ID, post.PostAuthor, post.PostDate, post.PostContent, post.PostTitle, post.PostModified, post.Thumb, post.Cat)

	return err
}

func EditPost(post Post) error {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return err
	}
	//defer db.Close()


	_, err = db.Exec("UPDATE wp_posts SET post_content=?, post_title=?, post_modified=?, thumb=?, cat=? WHERE ID=?",
		post.PostContent, post.PostTitle, post.PostModified, post.Thumb, post.Cat, post.ID)

	return err
}

func GetPost(postID string) (Post, error)  {
	var post Post

	sql := fmt.Sprintf(`SELECT * FROM wp_posts WHERE ID="%s"`, postID)

	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return post, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&post.ID,
			&post.PostAuthor,
			&post.PostDate,
			&post.PostContent,
			&post.PostTitle,
			&post.PostModified,
			&post.Thumb,
			&post.Cat)
		if err != nil {
			log.Println("GetPost", err)
			return post, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return post, err
	}

	return post, nil
}

/**
 * if category is empty string then no filtering
 */
func GetPostsByCategory(category string, page_no int64, records_per_page int64) ([]Post, error)  {
	var record_offset int64 = records_per_page * page_no

	var sql string

	if (category != "") {
		sql = fmt.Sprintf(`SELECT * FROM wp_posts WHERE wp_posts.cat LIKE '%%"%s"%%'
			ORDER BY post_date
			DESC LIMIT %d, %d`, category, record_offset, records_per_page)
	} else {
		sql = fmt.Sprintf(`SELECT * FROM wp_posts
			ORDER BY post_date
			DESC LIMIT %d, %d`, record_offset, records_per_page)
	}


	//log.Println("GetPostsByCategory: sql=", sql)
	rows, err := Query (sql)
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
			&post.PostContent,
			&post.PostTitle,
			&post.PostModified,
			&post.Thumb,
			&post.Cat)
		if err != nil {
			log.Println("GetPostsByCategory", err)
			return nil, err
		}

		items = append(items, post)
	}


	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}

func GetCategoryOfPost(json_arr_str string) ([]Category, error)  {
	cats_string := []string{}
	json.Unmarshal([]byte(json_arr_str), &cats_string)

	items := []Category{}
	for _, item := range cats_string {
		cat := Category{ item, 0 }
		items = append(items, cat)
	}

	return items, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Comment

func CreateComment(cm Comment) error {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return err
	}
	//defer db.Close()

	_, err = db.Exec("INSERT INTO tbl_comments(ID, postID, cm_parent, cm_author, cm_date, cm_content, cm_modified) " +
		"VALUES(?, ?, ?, ?, ?, ?)",
		cm.ID, cm.PostID, cm.Parent, cm.Author, cm.Date, cm.Content, cm.Modified)

	return err
}

func EditComment(cm Comment) error {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return err
	}
	//defer db.Close()


	_, err = db.Exec("UPDATE tbl_comments SET cm_content=?, cm_modified=? WHERE ID=?", cm.Content, cm.Modified, cm.ID)

	return err
}

func GetComment(ID string) (Comment, error)  {
	var cm Comment

	sql := fmt.Sprintf(`SELECT * FROM tbl_comments WHERE ID="%s"`, ID)

	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return cm, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&cm.ID,
			&cm.PostID,
			&cm.Parent,
			&cm.Author,
			&cm.Date,
			&cm.Content,
			&cm.Modified,
		)
		if err != nil {
			log.Println("GetPost", err)
			return cm, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return cm, err
	}

	return cm, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Category

func GetTopCategories(max_records int64) ([]Category, error)  {
	sql := fmt.Sprintf(`SELECT * FROM tbl_cat ORDER BY count DESC LIMIT 0, %d`, max_records)

	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	defer rows.Close()

	items := []Category{}
	for rows.Next() {
		var item Category
		err := rows.Scan(
			&item.ID,
			&item.Count)
		if err != nil {
			log.Println("GetTopCategories", err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}