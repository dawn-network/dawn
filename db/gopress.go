package db

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	//"strconv"
)

type User struct {
	ID			string
	Username 		string
	Pubkey  		string
	UserRegistered 		string
	DisplayName 		string
}

type Post struct {
	ID                  int64
	PostAuthor          string
	PostDate            string
	PostContent         string
	PostTitle           string
	PostModified        string
	Thumb 		    string
}

//type PostMeta struct {
//	MetaId    int64
//	PostId    int64
//	MetaKey   string
//	MetaValue string
//}

//type Term struct {
//	TermId    int64
//	Name      string
//	Slug      string
//	TermGroup int64
//}
//
//type TermTaxonomy struct {
//	TermTaxonomyId int64
//	TermId         int64
//	Taxonomy       string
//	Description    string
//	Parent         int64
//	Count          int64
//}

type Category struct {
	ID    		int64
	Name     	string
	Count 		int64
}

type TermRelationship struct {
	ObjectId       int64
	TermTaxonomyId int64
	TermOrder      int
}

// ---------------------------------------------------------------------------------------------------------------------

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@/glogchain")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return nil, err
	}

	return db, nil
}

func Query (sql string) (*sql.Rows, error) {
	db, err := GetDB()
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}

	return rows, nil
}

// ---------------------------------------------------------------------------------------------------------------------

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
			log.Fatal(err)
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
	defer db.Close()


	_, err = db.Exec("INSERT INTO wp_users(ID, user_username, user_pubkey, user_registered, display_name) " +
		"VALUES(?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Pubkey, user.UserRegistered, user.DisplayName)




	//if err != nil {
	//	return err
	//}

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

func GetPost(postID int64) (Post, error)  {
	var post Post

	sql := fmt.Sprintf("SELECT * FROM wp_posts WHERE ID=%d", postID)

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
			&post.Thumb)
		if err != nil {
			log.Fatal(err)
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

// term_taxonomy_id is category
func GetPostsByCategory(term_taxonomy_id int64, page_no int64, records_per_page int64) ([]Post, error)  {
	var record_offset int64 = records_per_page * page_no

	sql := fmt.Sprintf("SELECT wp_posts.* FROM wp_posts LEFT JOIN wp_term_relationships ON (wp_posts.ID = wp_term_relationships.object_id) WHERE 1=1 AND ( wp_term_relationships.term_taxonomy_id IN (%d) ) GROUP BY wp_posts.ID ORDER BY wp_posts.post_date DESC LIMIT %d, %d", term_taxonomy_id, record_offset, records_per_page)
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
			&post.Thumb)
		if err != nil {
			log.Fatal(err)
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

func GetCategoryOfPost(postID int64) ([]Category, error)  {
	sql := fmt.Sprintf("SELECT * FROM tbl_cat WHERE ID IN (SELECT tbl_cat.ID FROM wp_term_relationships, tbl_cat WHERE wp_term_relationships.term_taxonomy_id=tbl_cat.ID AND object_id=%d )", postID)
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
			&item.Name,
			&item.Count)
		if err != nil {
			log.Fatal(err)
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

//func GetPostMetas(postID int64) ([]PostMeta, error)  {
//	sql := fmt.Sprintf("SELECT * FROM wp_postmeta WHERE post_id=%d", postID)
//	rows, err := Query (sql)
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//		return nil, err
//	}
//	defer rows.Close()
//
//	items := []PostMeta{}
//	for rows.Next() {
//		var item PostMeta
//		err := rows.Scan(
//			&item.MetaId,
//			&item.PostId,
//			&item.MetaKey,
//			&item.MetaValue)
//		if err != nil {
//			log.Fatal(err)
//			return nil, err
//		}
//		items = append(items, item)
//	}
//
//	err = rows.Err()
//	if err != nil {
//		log.Fatal(err)
//		return nil, err
//	}
//
//	return items, nil
//}

//func GetPostThumbnail(postID int64) string  {
//	postmetas, err := GetPostMetas(postID)
//	if err != nil {
//		return ""
//	}
//
//	for _, v := range postmetas {
//		if (v.MetaKey == "_thumbnail_id") {
//			thumbnail_id, err := strconv.ParseInt(v.MetaValue, 10, 64)
//			if (err == nil) {
//				thumbnail, err := GetPost(thumbnail_id)
//				if (err == nil) {
//					//log.Println("GetPostThumbnail", thumbnail.Guid)
//					return thumbnail.Guid
//				}
//			}
//
//			break
//		}
//	}
//
//	return "/static/img/slider-featured-image.png" // default
//}

// ---------------------------------------------------------------------------------------------------------------------

