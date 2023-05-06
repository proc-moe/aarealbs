# 如何debug？
0. 记得删掉自己数据库中的已有相关数据
1. git clone 
2. cd server/
2. go get -u
3.  vscode 全局搜索 dsn,然后在model/model.go替换上自己的数据库连接
```
dsn := "${数据库用户名}:${数据库密码}@tcp(localhost:3306)/${数据库名称}?parseTime=true"
```
4. go run . 启动服务
5. 部署后调用接口添加测试帐号 id: 114514, user_name: yjsnpi, token:555

URL:
```
localhost:9987/api/token/add
```
body:
```
{"token":"test"}
```
