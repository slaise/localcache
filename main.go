package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localCache "github.com/slaise/localcache/localcache"
)

var c *localCache.Cache

type Param struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func PutKey(w http.ResponseWriter, r *http.Request) {
	var keyVal Param
	_ = json.NewDecoder(r.Body).Decode(&keyVal)
	fmt.Printf("put key: %v and val: %v \n", keyVal.Key, keyVal.Val)
	c.Put(keyVal.Key, keyVal.Val)
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	var keyVal Param
	_ = json.NewDecoder(r.Body).Decode(&keyVal)
	val, _ := c.Get(keyVal.Key)
	fmt.Printf("get val: %v from key: %v \n", val, keyVal.Key)
	res := []byte(val.(string))
	w.Write(res)
}

func main() {
	c, _ = localCache.NewCache(1000)
	http.HandleFunc("/put", PutKey)
	http.HandleFunc("/get", GetKey)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
