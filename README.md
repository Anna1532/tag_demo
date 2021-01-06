# tag_demo
# 自定义结构体标签

### 使用场景：

1.ORM结构体映射数据库

2.数据验证：gin中的form表单数据

3.序列化与反序列化：JSON模块-结构体转JSON字符串，JSON字符串转结构体

![image-20210105140124131](C:\Users\23013\Desktop\Go学习\笔记\form表单.png)



### Beego form表单数据校验

1.自定义结构体标签，使用valid标签作为标签

2.动态解析不同的字段调用对应字段的验证方法，并返回结果

![image-20210105141030025](C:\Users\23013\Desktop\Go学习\笔记\标签示例代码.png)
