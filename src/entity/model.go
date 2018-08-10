package entity

type Model struct {
	//ID   string  		`bson:"id"`
	NAME string        	`bson:"name"` //bson:"name" 表示mongodb数据库中对应的字段名称
	AGE  int          `bson:"age"`
}

