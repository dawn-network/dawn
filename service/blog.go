package service

import (
	"math/rand"
)

//func categories_normalize(jsonstr string) (string, error) {
//	cats_string := []string{}
//	json.Unmarshal([]byte(jsonstr), &cats_string)
//
//	//items := []Category{}
//	for i, item := range cats_string {
//	//	cat := Category{ item, 0 }
//	//	items = append(items, cat)
//
//		item = strings.ToLower(item)
//		item = strings.TrimSpace(item)
//		cats_string[i] = item
//	}
//
//	return items, nil
//
//
//}


var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

/**
 * Generate random string
 * http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
 */
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}