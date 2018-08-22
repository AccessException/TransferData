package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"strconv"
	"github.com/modern-go/reflect2"
	"reflect"

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
	for i:=0;i<0;i++{
		p.Name = "liupeng"+strconv.Itoa(i)
		p.Age = i
		err = c.Insert(&p)
		if(err != nil){
			log.Println(err)
		}
	}
	//var m entity.Model
	model := reflect2.TypeByName("main.Person")
	d := model.Type1()
	log.Println(d.Field(1))
	log.Println(model)
	log.Println("&&&&&&&",reflect.TypeOf(model))

	//typByPkg := reflect2.TypeByPackageName(
	//	"github.com/modern-go/reflect2-tests",
	//	"MyStruct")
	//log.Println(typByPkg)

}
