package db
//
//import (
//	"database/sql"
//	_ "github.com/mattn/go-sqlite3"
//)

//// TestItem is test only and should be removed soon
//type TestItem struct {
//	Id	string
//	Name	string
//	Phone	string
//}
//
//var EmbeddedDB *sql.DB // global var
//
//func InitDB() {
//	EmbeddedDB, err := sql.Open("sqlite3", "db.db")
//
//	if err != nil {
//		panic(err)
//	}
//
//	if EmbeddedDB == nil {
//		panic("db nil")
//	}
//
//	//return db
//
//
//	////////////
//	_, err = EmbeddedDB.Exec(`CREATE TABLE IF NOT EXISTS tbl_cat
//		(
//		    ID VARCHAR(200) NOT NULL PRIMARY KEY,
//		    count INTEGER NOT NULL DEFAULT '0'
//		);`)
//	if err != nil {
//		panic(err)
//	}
//
//	//////
//	_, err = EmbeddedDB.Exec(`CREATE TABLE wp_users
//		(
//		    ID VARCHAR(40) NOT NULL PRIMARY KEY,
//		    user_username VARCHAR(60) NOT NULL DEFAULT '',
//		    user_pubkey VARCHAR(255) NOT NULL DEFAULT '',
//		    user_registered DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
//		    display_name VARCHAR(250) NOT NULL DEFAULT ''
//		);`)
//	if err != nil {
//		panic(err)
//	}
//
//	//////
//	_, err = EmbeddedDB.Exec(`CREATE TABLE wp_posts
//		(
//		    ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
//		    post_author VARCHAR(40) NOT NULL DEFAULT '0',
//		    post_date DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
//		    post_content LONGTEXT NOT NULL,
//		    post_title TEXT NOT NULL,
//		    post_modified DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00',
//		    thumb VARCHAR(255) NOT NULL DEFAULT '',
//		    cat VARCHAR(255) NOT NULL DEFAULT ''
//		);`)
//	if err != nil {
//		panic(err)
//	}
//}



//func StoreItem(db *sql.DB, items []TestItem) {
//	sql_additem := `
//	INSERT OR REPLACE INTO items(
//		Id,
//		Name,
//		Phone,
//		InsertedDatetime
//	) values(?, ?, ?, CURRENT_TIMESTAMP)
//	`
//
//	stmt, err := db.Prepare(sql_additem)
//	if err != nil { panic(err) }
//	defer stmt.Close()
//
//	for _, item := range items {
//		_, err2 := stmt.Exec(item.Id, item.Name, item.Phone)
//		if err2 != nil { panic(err2) }
//	}
//}
//
//func ReadItem(db *sql.DB) []TestItem {
//	sql_readall := `
//	SELECT Id, Name, Phone FROM items
//	ORDER BY datetime(InsertedDatetime) DESC
//	`
//
//	rows, err := db.Query(sql_readall)
//	if err != nil { panic(err) }
//	defer rows.Close()
//
//	var result []TestItem
//	for rows.Next() {
//		item := TestItem{}
//		err2 := rows.Scan(&item.Id, &item.Name, &item.Phone)
//		if err2 != nil { panic(err2) }
//		result = append(result, item)
//	}
//	return result
//}