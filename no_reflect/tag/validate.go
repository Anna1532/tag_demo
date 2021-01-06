package tag

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const ValidTagName="valid"

type Validation struct {
	Data   interface{} //接收各种结构体数据，空接口可以接收任何数据类型
	Errors []error     //错误列表
}

func NewValidation(structData interface{})*Validation{
	vd:=new(Validation)
	vd.Data=structData
	return vd
}
func (v *Validation)Validate(){
	//编写验证逻辑
	//初始化错误列表
	errs :=[]error{}
	//获取runtime数据
	runData:=reflect.ValueOf(v.Data)
	//校验数据是否是结构体
	if runData.Type().Kind().String()!="struct"{
		panic("数据必须是结构体")
	}
	//循环遍历结构体字段，获取标签内容，进行校验
	for i:=0;i<runData.NumField();i++{
		//获取tag内容,Get传递的是标签的key值
		tag:=runData.Type().Field(i).Tag.Get(ValidTagName)
		//根据tag内容获取当前字段的验证器
		curValidator:=v.GetValidator(tag)
		//fmt.Println("get value",curValidator)
		//调用验证器的验证方法，记录错误,把字段的值传递进去
		//runData.Field(i).Interface()把字段类型转化成接口类型
		valid,err:=curValidator.Validate(runData.Field(i).Interface())
		//把验证失败的错误添加到错误列表
		if !valid&&err!=nil{
			errs=append(errs,err)
		}
	}
	//验证完毕之后把错误列表赋值给验证器类
	v.Errors=errs
}

//通过tag内容返回验证器OK
var MaxRe=regexp.MustCompile(`\Amaxsize\(\d+\)\z`)
var MinRe=regexp.MustCompile(`\Aminsize\(\d+\)\z`)
func (v *Validation)GetValidator(tag string)Validator {
	//解析tag内容
	tagArgs:=strings.Split(tag,";")
	//根据标志生成对应验证器
	switch tagArgs[0] {
	case "string":
		strValidator:=&StringValidator{}
		for _,str:=range tagArgs[1:]{
			switch  {
			case MaxRe.MatchString(str):
				fmt.Sscanf(str,"maxsize(%d)",&strValidator.Max)
			case MinRe.MatchString(str):
				fmt.Sscanf(str,"minsize(%d)",&strValidator.Min)
			}
		}
		if strValidator.Min>strValidator.Max{
			panic("max value must be bigger than min")
		}
		return strValidator
	case "number":
		numValidator:=NumberValidator{}//{0,0}
		fmt.Sscanf(tagArgs[1],"range(%d,%d)",&numValidator.Min,&numValidator.Max)
		if numValidator.Min>numValidator.Max{
			panic("max value must be bigger than min")
		}
		return &numValidator
	case "email":
		return &EmailValidator{}
	}
	//若接口的方法绑定的是指针类型，转化接口的时候也需要使用指针或地址
	return &DefaultValidator{}
}

//多态：调用一个行为（方法），呈现不同结果，Go中通过接口实现多态
//定义验证器接口
type Validator interface {
	Validate(interface{})(bool,error)
}

//缺省的验证器类，不验证
type DefaultValidator struct {

}

func (v *DefaultValidator)Validate(data interface{})(bool,error){
	return true,nil
}

//string验证器类
type StringValidator struct {
	Min int
	Max int
}

func (v *StringValidator)Validate(data interface{})(bool,error){
	//把数据从接口类型转换成string类型
	str:=data.(string)
	//获取字符串长度，用于校验验证
	lenth:=len(str)
	//fmt.Println(v.Max)
	//比较长度进行验证
	if lenth<v.Min{
		//error没有格式化,这里直接构造错误
		return false,fmt.Errorf("姓名长度不能小于%d",v.Min)
	}
	if lenth>v.Max{
		return false,fmt.Errorf("姓名长度不能大于%d",v.Max)
	}
	return true,nil
}

type NumberValidator struct {
	Min int
	Max int
}
func (v *NumberValidator)Validate(data interface{})(bool,error){
	//把接口数据转化成int类型
	number:=data.(int)
	if number<v.Min{
		return false,fmt.Errorf("年龄不能小于%d",v.Min)
	}
	if number>v.Max{
		return false,fmt.Errorf("年龄不能大于%d",v.Max)
	}
	return true,nil
}

//定义email验证器,通过正则验证邮箱
var emailRe=regexp.MustCompile(`\A[\w+\d\-.]+@[a-z\d\-]+(\.[a-z\d\-]+)*\.[a-z]+\z`)
type EmailValidator struct{

}
func(v *EmailValidator)Validate(data interface{})(bool,error){
	//将接口类型转换成字符串
	email:=data.(string)
	if !emailRe.MatchString(email){
		return false,fmt.Errorf("邮箱地址格式错误")
	}
	return true,nil
}