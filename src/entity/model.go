package entity

type Model struct {
	Name 	string 	    		`json:"name"` //bson:"name" 表示mongodb数据库中对应的字段名称
	Age 	int 	       		`json:"age"`
}