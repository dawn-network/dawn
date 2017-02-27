package db

import (
	//"testing"
	//"fmt"
	//"log"
	//"encoding/json"
)
//
//func Test_query_mysql(t *testing.T) {
//	//sql := fmt.Sprintf(`SELECT * FROM tbl_cat`)
//	//
//	//rows, err := Query (sql)
//	//if err != nil {
//	//	panic(err.Error()) // proper error handling instead of panic in your app
//	//	return
//	//}
//	//defer rows.Close()
//	//
//	//items := []Category{}
//	//for rows.Next() {
//	//	var item Category
//	//	err := rows.Scan(
//	//		&item.ID,
//	//		&item.Name,
//	//		&item.Count)
//	//	if err != nil {
//	//		log.Println("Test_query_mysql", err)
//	//		return
//	//	}
//	//	items = append(items, item)
//	//
//	//	log.Println(item.ID, item.Name, item.Count)
//	//}
//	//
//	//err = rows.Err()
//	//if err != nil {
//	//	log.Fatal(err)
//	//	return
//	//}
//
//	sql := fmt.Sprintf(`SELECT * FROM wp_posts`)
//
//	rows, err := Query (sql)
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var post Post
//
//		var scan_cat *string
//
//		err := rows.Scan(
//			&post.ID,
//			&post.PostAuthor,
//			&post.PostDate,
//			&post.PostContent,
//			&post.PostTitle,
//			&post.PostModified,
//			&post.Thumb,
//			&scan_cat)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//
//		log.Println("post.ID: ", post.ID, scan_cat)
//
//		if (scan_cat != nil) {
//			post.Cat = *scan_cat
//		}
//
//		//var string_categories string = ""
//		cats, err := GetCategoryOfPost(post.ID)
//
//		cats_string := []string{}
//		for _, cat := range cats {
//			//fmt.Printf("cat [%d] is [%s]\n", index, cat.Name)
//			//string_categories += cat.Name
//			//
//			//if (index < len(cats) - 1) {
//			//	string_categories += ", "
//			//}
//
//			cats_string = append(cats_string, cat.Name)
//		}
//
//		var buf []byte
//		buf, err = json.Marshal(cats_string)
//		if err != nil {
//			log.Fatal("Cannot encode to JSON ", err)
//		}
//
//		log.Println("string_categories=" + string(buf[:]))
//
//		UpdatePost(post.ID, string(buf[:]))
//	}
//
//	err = rows.Err()
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//}
//
//func UpdatePost(postId int64, json_cat string) error {
//	db, err := GetDB()
//	if err != nil {
//		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
//		return err
//	}
//	defer db.Close()
//
//
//	_, err = db.Exec("UPDATE wp_posts SET cat=? WHERE ID=?",
//		json_cat, postId)
//
//	//if err != nil {
//	//	return err
//	//}
//
//	return err
//}