package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestHub(t *testing.T){
	data,e:= json.Marshal(&Message{Timestamp: time.Now(), Nick: "scott"})
	fmt.Println(string(data),e)

	b:=[]int8{1,2,3,4,5}
	b = append(b[:1],b[2:len(b)]...)
	fmt.Println(b)

	 x:= map[string] bool{
		 "abc":true,
		 "123":false,
	 }
	 n,m:=x["dd"]
	 fmt.Println( n,m )
}
