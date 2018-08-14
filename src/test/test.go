package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
)

type Person struct{
	Name string
	Age int
}

func main() {
	log.Println("")
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil { panic(err) }
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("liupeng").C("test")
	p := Person{}
	for i:=0;i<5000;i++{
		p.Name = "liupeng"+strconv.Itoa(i)
		p.Age = i
		err = c.Insert(&p)
		if(err != nil){
			log.Println(err)
		}
	}
}
