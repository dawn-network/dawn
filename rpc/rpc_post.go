package rpc

import (
	"log"
	"net/http"
)

type RpcPost struct {

}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//// List
//// to get list of posts
//
//type RpcPost_List_Args struct {
//
//}
//
//type RpcPost_List_Result struct {
//	A	string
//}
//
//func (t *RpcPost) List(r *http.Request, args *RpcPost_List_Args, result *RpcPost_List_Result) error {
//	log.Println("List")
//	*result = RpcPost_List_Result { A: "AAAA" }
//	return nil
//}


/////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create
// to create a new post

type RpcPost_Create_Args struct {

}

type RpcPost_Create_Result struct {
	A	string
}

func (t *RpcPost) Create(r *http.Request, args *RpcPost_Create_Args, result *RpcPost_Create_Result) error {
	log.Println("Create")
	*result = RpcPost_Create_Result { A: "AAAA" }
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
// Get
// to get an existing post

type RpcPost_Get_Args struct {

}

type RpcPost_Get_Result struct {
	A	string
}

func (t *RpcPost) Get(r *http.Request, args *RpcPost_Get_Args, result *RpcPost_Get_Result) error {
	log.Println("Get")
	*result = RpcPost_Get_Result { A: "AAAA" }
	return nil
}

