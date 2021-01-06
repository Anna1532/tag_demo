package main

import (
	"fmt"
	"tag_demo/no_reflect/tag"
)

type User struct {
	ID    int    `-`
	Name  string `valid:"string;maxsize(6);minsize(3)"`
	Age   int    `valid:"number;range(0,150)"`
	Email string `valid:"email"`
}

func main() {
	user:=User{
		ID:1,
		Name:"Mrsdv",
		Age:199,
		Email:"xxxx@qq.com.com",
	}
	//创建一个验证器对象
	vd:=tag.NewValidation(user)
	//进行user验证
	vd.Validate()
	for _,err:=range(vd.Errors){
		fmt.Println(err.Error())
	}
	//TODO:啊啊啊啊报错啦，获取不到标签中的规定长度
}
