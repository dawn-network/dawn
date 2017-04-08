package types

type Comment struct {
	ID                  	string
	PostID                 	string 		// which post belong to
	Parent 			string 		// parrent comment, empty string means no parent
	Author          	string		// ID/Address of Account/User
	Date            	string 		// create datetime
	Content         	string
	Modified        	string 		// Last Modified datetime
}


