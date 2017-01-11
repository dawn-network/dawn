package db

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strconv"
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
	PostExcerpt         string
	PostStatus          string
	CommentStatus       string
	PingStatus          string
	PostPassword        string
	PostName            string
	ToPing              string
	Pinged              string
	PostModified        string
	PostContentFiltered string
	PostParent          int64
	Guid                string
	MenuOrder           int
	PostType            string
	PostMimeType        string
	CommentCount        int64
}

type PostMeta struct {
	MetaId    int64
	PostId    int64
	MetaKey   string
	MetaValue string
}

type Term struct {
	TermId    int64
	Name      string
	Slug      string
	TermGroup int64
}

type TermTaxonomy struct {
	TermTaxonomyId int64
	TermId         int64
	Taxonomy       string
	Description    string
	Parent         int64
	Count          int64
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
			&post.PostExcerpt,
			&post.PostStatus,
			&post.CommentStatus,
			&post.PingStatus,
			&post.PostPassword,
			&post.PostName,
			&post.ToPing,
			&post.Pinged,
			&post.PostModified,
			&post.PostContentFiltered,
			&post.PostParent,
			&post.Guid,
			&post.MenuOrder,
			&post.PostType,
			&post.PostMimeType,
			&post.CommentCount )
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

	sql := fmt.Sprintf("SELECT wp_posts.* FROM wp_posts LEFT JOIN wp_term_relationships ON (wp_posts.ID = wp_term_relationships.object_id) WHERE 1=1 AND ( wp_term_relationships.term_taxonomy_id IN (%d) ) AND wp_posts.post_type = 'post' AND (wp_posts.post_status = 'publish' OR wp_posts.post_status = 'private') GROUP BY wp_posts.ID ORDER BY wp_posts.post_date DESC LIMIT %d, %d", term_taxonomy_id, record_offset, records_per_page)
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
			&post.PostExcerpt,
			&post.PostStatus,
			&post.CommentStatus,
			&post.PingStatus,
			&post.PostPassword,
			&post.PostName,
			&post.ToPing,
			&post.Pinged,
			&post.PostModified,
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

		items = append(items, post)
	}


	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}

func GetCategoryOfPost(postID int64) ([]Term, error)  {
	sql := fmt.Sprintf("SELECT * FROM wp_terms WHERE term_id IN (SELECT wp_term_taxonomy.term_id FROM wp_term_relationships, wp_term_taxonomy WHERE wp_term_relationships.term_taxonomy_id=wp_term_taxonomy.term_taxonomy_id AND wp_term_taxonomy.taxonomy='category' AND object_id=%d )", postID)
	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	defer rows.Close()

	items := []Term{}
	for rows.Next() {
		var item Term
		err := rows.Scan(
			&item.TermId,
			&item.Name,
			&item.Slug,
			&item.TermGroup)
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

func GetPostMetas(postID int64) ([]PostMeta, error)  {
	sql := fmt.Sprintf("SELECT * FROM wp_postmeta WHERE post_id=%d", postID)
	rows, err := Query (sql)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
		return nil, err
	}
	defer rows.Close()

	items := []PostMeta{}
	for rows.Next() {
		var item PostMeta
		err := rows.Scan(
			&item.MetaId,
			&item.PostId,
			&item.MetaKey,
			&item.MetaValue)
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

func GetPostThumbnail(postID int64) string  {
	postmetas, err := GetPostMetas(postID)
	if err != nil {
		return ""
	}

	for _, v := range postmetas {
		if (v.MetaKey == "_thumbnail_id") {
			thumbnail_id, err := strconv.ParseInt(v.MetaValue, 10, 64)
			if (err == nil) {
				thumbnail, err := GetPost(thumbnail_id)
				if (err == nil) {
					//log.Println("GetPostThumbnail", thumbnail.Guid)
					return thumbnail.Guid
				}
			}

			break
		}
	}

	return "/static/img/slider-featured-image.png" // default
}

// ---------------------------------------------------------------------------------------------------------------------

