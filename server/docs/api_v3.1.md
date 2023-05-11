v3:
1. 单词列表增加: in_queue bool,是否在背诵队列内
2. 增加：获取所有人背诵列表中的词条（只是url路由不一样，返回体一致）
3. 增加：4.3.1 查看所有用户的背诵历史 （只是url路由不一样，返回体一致）
4. 增加：背诵队列（4.5.1，4.5.2）
```
 -------------
|           |                       -------- 
| 词条       | ====导入===>        | 进入队列|   ====到期====> 复习
|           |                       --------
 --------------
```

BaseURL: localhost:9987


# 1. 登陆
URL: localhost:9987/api/login

METHOD: POST 

请求头

无

请求体
```
名称      类型        必填        描述
token   string     yes      朝telegram机器人索要的有效期为60min的token
```
实例：
```
{
    "token": 114514
}
```
响应体
```
code int 错误码，非0表示失败
msg string 错误描述
user_id int 用户的tgid（也是这里的唯一id）
name string 用户名 
token string token
```

响应体实例
``` json
{
    "code": 0,
    "msg": "success",
    "user_id": 53453,
    "name": "235345",
    "token": "2350493fjx"
}
```
# 2. 词条

## 2.1 获取用户上传的词条
```
URL:localhost:9987/api/record/upload/:user_id?offset=${int}&limit=${int}
URL:localhost:9987/api/record/upload/893893?offset=23&limit=2
```
method: GET

请求头
```
authorization - string 
```
请求体
```
无
```
响应体
```
    code - int
    total - int
    record - array
     - id int
     - created_at string
     - updated_at string
     - deleted_at string
     - batch_info_id int
     - msg string
     - resp_msg string
     - user_info_id int 
     - in_queue bool 是否在背诵队列内 [v3修改]
    msg - string
```

``` json
{
    "code": 0,
    "total": 124,
    "record": [
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "batch_info_id": 1234, // 所属单词集id 不属于任何单词集默认为空
            "msg": "苹果的英语是什么？", // msg
            "resp_msg": "apple",    // 
            "user_info_id": 893893, // 用户id
            "in_queue": true // 是否在队列中
        },
        {
            "id": 2,
            "created_at": "2012-02-24 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "batch_info_id": 324, // 所属单词集id 不属于任何单词集默认为空
            "msg": "苹果的英语是什么？", // msg
            "resp_msg": "banana",
            "user_info_id": 893893,
            "in_queue": false,
        }
    ],
    "msg": "success",
}
```


## 2.2 修改词条
```
URL:localhost:9987/api/record/edit
```
method: POST

请求头 HEADERS
```
authorization - string 
```

请求体
```
id - int
msg - string
resp_msg - string
pattern_id - string
```

``` json
{
    "id": 5,
    "msg": "苹果的英语",
    "resp_msg": "apple",
    "pattern_id": 0,
}
```

响应体
```
    code - int
    msg - int
```

``` json
{
    "code": 0,
    "msg": "success",
}
```



## 2.3 删除词条
```
URL:localhost:9987/api/record/delete
```
method: POST

请求头 HEADERS
```
authorization - string 
```

请求体
```
id - int
```

``` json
{
    "id": 5
}
```

响应体
```
    code - int
    msg - int
```

``` json
{
    "code": 0,
    "msg": "success",
}
```


``` json
{
    "code": 223,
    "msg": "database error",
}
```

## 2.4 批量添加词条


```
URL: localhost:9987/api/record/batch/add
```

method 
```
POST
```

请求头
```
authorization - string 
```

请求体 

```
一个数组，里面的每个对象包括
batch_name - string
msg - string
resp_msg - string
user_id - int
```

``` json
{
    [
        {
            "batch_name": "a batch file",
            "msg": "苹果",
            "resp_msg": "apple",
            "user_id": 1,
        },
        {
            "batch_name": "a batch file",
            "msg": "火车",
            "resp_msg": "banana",
            "user_id": 1,
        },
        {
            "batch_name": "a batch file",
            "msg": "哭了",
            "resp_msg": "pear",
            "user_id": 1,
        },
    ]
}

```


响应体：
```
code - int  -  0
msg - string - "success"
count - int - 计数
```

## 2.5 添加词条
```
URL:localhost:9987/api/record/add
```
method: POST

请求头 HEADERS
```
authorization - string 
```

请求体
```
msg - string                词条
resp_msg - string           解释
user_id - int               作者id
pattern_id - int         pattern id 默认是 0
```

``` json
{
    "msg": "苹果的英语",
    "resp_msg": "apple",
    "user_id": 114514,
    "pattern_id": 0
}
```

响应体
```
    code - int
    msg - int
```

``` json
{
    "code": 0,
    "msg": "success",
}
```


# 3. 用户信息
## 3.1 获取所有用户信息
获取所有用户信息，管理员视角下调用用户信息可以使用这个
```
URL: localhost:9987/api/user/all?offset=${offset}&limit=${limit}
URL: localhost:9987/api/user/all?offset=123&limit=44
```

method: GET

请求头 HEADERS
```
authorization - string 
```

响应体
```
    code - int
    msg - int
    total - int             用户的总共数量
    user_infos: array
        id: int             数据库内的id， 不具备使用意义，只是个主键，可以忽略
        user_id: int        userid,值等价于tgid
        user_name: string   user_name: 用户名称
        status: int         账号状态
```

``` json
{
    "code": 0,
    "message": "success",
    "user_infos": [
        {
            "id": 1,
            "user_id": 2234234,
            "user_name": "好朋友",
            "status": 1,
        },
        {
            "id": 2,
            "user_id": 2234234,
            "user_name": "坏朋友",
            "status": 2,
        }
    ]
}
```



## 获取单用户信息
比如说获取用户本人的信息的时候可以用这个
URL: localhost:9987/api/user/:id

method: GET

请求头 HEADERS

```
authorization - string 
```

响应体体 
``` json
{ 
    "code": 0,
    "user":{
        "id": 1,
        "created_at": "2012-02-24 12:13:14", // 创建时间
        "updated_at": "2012-02-03 12:13:14", // 更新时间
        "deleted_at": "2012-02-03 12:13:14", // 删除时间
        "user_id": 234235, // tgID
        "user_name": "proc-moe", // tgName
        "status": 1, // 1 ok  2 管理员 3 banned
    },
    "msg": "success",
}
```


## 3.2 修改用户信息
- 只有管理员能调这个接口，后台会做校验，失败响应体401

URL: localhost:9987/api/user/edit

method: POST

请求头 HEADERS

```
authorization - string 
```

请求体
``` json
{ 
    "user_id": 234235,
    "status": 2, // 1 ok  2 管理员 3 banned
}
```



# 4. 背诵
## 4.1 “开始记忆～“ 把一个词条导入背诵队列，也就是说几分钟后就会需要复习这个词条了
- 只有自己才能开始一个词条的背诵

```
URL:localhost:9987/api/queue/add
```

method POST

请求头 HEADERS
```
authorization - string 
```

请求体
```
user_info_id: int
record_info_id: int
```

``` json
{
    "user_info_id": 1324,
    "record_info_id": 3425, 
}
```

## 4.2 获取用户(user_id)背诵列表中的词条 
```
URL:localhost:9987/api/queue/user/:user_id?offset=${int}&limit=${int}
URL:localhost:9987/api/queue/user/893893?offset=23&limit=2
```
method: GET

请求头 HEADERS
```
authorization - string 
```

请求体

无

响应体
```
    code - int
    record - array
     - id int
     - created_at string
     - updated_at string
     - deleted_at string
     - user_id id
     - remind_time_unix int 下一次提醒的时间
     - round 已经背诵的轮数
     - round_max 最多几轮背诵
     - status 状态
    msg - string
```

``` json
{
    "code": 0,
    "record": [
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178725, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178388, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
    ],
    "msg": "success",
}
```

## 4.2.1 获取所有人背诵列表中的词条 v3修改
```
URL:localhost:9987/api/queue/user/all?offset=${int}&limit=${int}
URL:localhost:9987/api/queue/user/all?offset=23&limit=2
```

method: GET

请求头 HEADERS
```
authorization - string 
```

请求体

无

响应体
```
    code - int
    record - array
     - id int
     - created_at string
     - updated_at string
     - deleted_at string
     - user_id int
     - record_id int
     - remind_time_unix int 下一次提醒的时间
     - round 已经背诵的轮数
     - round_max 最多几轮背诵
     - status 状态
    msg - string
```

``` json
{
    "code": 0,
    "record": [
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178725, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178388, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
    ],
    "msg": "success",
}
```


## 4.3 查看用户的背诵历史
```
URL:localhost:9987/api/record_history/:user_id?offset=${offset}&limit=${limit}
URL:localhost:9987/api/record_history/893893?offset=23&limit=2
```
method: GET

请求头 HEADERS
```
authorization - string 
```

请求体

无

响应体
```
    code - int
    total         int 总数量
    history
        created_at	    string
        updated_at	    string
        deleted_at	    string
        id              int     id
        user_info_id	int     用户id
        record_info_id	int     记录id
        time_gap	    int     时间间隔
        time_gap_est	int     预计时间间隔
        result	int 背诵结果 - 0 记忆成功 1 失败
    msg - string
```

``` json
{
    "code": 0,
    "total": 123, // 这样应该方便翻页
    "history": [
        {
            "id": 54,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 历史记录属于哪个用户
            "record_info_id": 1234, // 用户的记忆词条
            "time_gap": 1683178725, // unix 时间戳，预计需要背诵的时间
            "time_gap_est": 0, // 已经进行了几轮背诵
            "result": 0, // 背诵结果 - 0 - 1 - 2
        },
        {
            "id": 23,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 历史记录属于哪个用户
            "record_info_id": 1234, // 用户的记忆词条
            "time_gap": 1683178725, // unix 时间戳，预计需要背诵的时间
            "time_gap_est": 0, // 已经进行了几轮背诵
            "result": 0, // 背诵结果 - 0 - 1 - 2
        },
    ],
    "msg": "success",
}
```

## 4.3.1 查看所有用户的背诵历史 v3修改 

注：与4.3的区别，只是把user_id变成了all
```
URL:localhost:9987/api/record_history/all?offset=${offset}&limit=${limit}
URL:localhost:9987/api/record_history/893893?offset=23&limit=2
```
method: GET

请求头 HEADERS
```
authorization - string 
```

请求体

无

响应体
```
    code - int
    total         int 总数量
    history
        created_at	    string
        updated_at	    string
        deleted_at	    string
        id              int     id
        user_info_id	int     用户id
        record_info_id	int     记录id
        time_gap	    int     时间间隔
        time_gap_est	int     预计时间间隔
        result	int 背诵结果 - 0 记忆成功 1 失败
    msg - string
```

``` json
{
    "code": 0,
    "total": 123, // 这样应该方便翻页
    "history": [
        {
            "id": 54,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 历史记录属于哪个用户
            "record_info_id": 1234, // 用户的记忆词条
            "time_gap": 1683178725, // unix 时间戳，预计需要背诵的时间
            "time_gap_est": 0, // 已经进行了几轮背诵
            "result": 0, // 背诵结果 - 0 - 1 - 2
        },
        {
            "id": 23,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 历史记录属于哪个用户
            "record_info_id": 1234, // 用户的记忆词条
            "time_gap": 1683178725, // unix 时间戳，预计需要背诵的时间
            "time_gap_est": 0, // 已经进行了几轮背诵
            "result": 0, // 背诵结果 - 0 - 1 - 2
        },
    ],
    "msg": "success",
}
```






## 4.4 背诵操作
```
URL:localhost:9987/api/recite/:queue_id
URL:localhost:9987/api/recite/2342
```
method: POST

请求头 HEADERS
```
authorization - string 
```

请求体

```
    recite_id: int 背诵队列的id
    result: int // 0 成功 1 忘记了/失败
```

``` json
{
    "recite_id": 324,
    "result": 1,
}

```
响应体
```
    code - int
    msg - string
```

``` json
{
    "code": 0,
    "msg": "success",
}
```



## 4.5 该用户的背诵队列（待复习队列的是到了复习时间的背诵队列） v3修改

```
URL:localhost:9987/api/timeup_queue/user/:user_id?offset=${offset}&limit=${limit}
URL:localhost:9987/api/timeup_queue/user/:user_id?offset=${offset}&limit=${limit}
```

method: GET

请求头 HEADERS
```
authorization - string 
```

响应体
```
    code - int
    record - array
     - id int
     - created_at string
     - updated_at string
     - deleted_at string
     - user_id int
     - record_id int
     - remind_time_unix int 下一次提醒的时间
     - round 已经背诵的轮数
     - round_max 最多几轮背诵
     - status 状态
    msg - string
```

``` json
{
    "code": 0,
    "record": [
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178725, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178388, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
    ],
    "msg": "success",
}
```

## 4.5.1 所有用户的背诵队列（待复习队列的是到了复习时间的背诵队列） v3修改

注：与4.5的区别，只是把:user_id变成了all
```
URL:localhost:9987/api/timeup_queue/user/all?offset=${offset}&limit=${limit}
URL:localhost:9987/api/timeup_queue/user/all?offset=${offset}&limit=${limit}
```

method: GET

请求头 HEADERS
```
authorization - string 
```

响应体
```
    code - int
    record - array
     - id int
     - created_at string
     - updated_at string
     - deleted_at string
     - user_id int
     - record_id int
     - remind_time_unix int 下一次提醒的时间
     - round 已经背诵的轮数
     - round_max 最多几轮背诵
     - status 状态
    msg - string
```

``` json
{
    "code": 0,
    "record": [
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178725, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
        {
            "id": 1,
            "created_at": "2012-02-03 12:13:14", // 创建时间
            "updated_at": "2012-02-03 12:13:14", // 更新时间
            "deleted_at": "2012-02-03 12:13:14", // 删除时间
            "user_info_id": 1234, // 这个单词是谁应该记忆的
            "record_info_id": 1234, // 这个单词的具体record是什么
            "remind_time_unix": 1683178388, // unix 时间戳，预计需要背诵的时间
            "round": 0, // 已经进行了几轮背诵
            "round_max": 15, // 最多几轮背诵
            "status": 0, // 状态
        },
    ],
    "msg": "success",
}
```




# 5 系统监控
## 获取一条监控历史
```
URL:localhost:9987/api/monitor/single/:id
```

method GET


请求头 HEADERS
```
authorization - string 
```

响应
```
code    int
msg     string
monitor 
    - cpu_load	int
    - mem_load	int
    - record_count	int 
    - user_count	int 
```

## 获取监控历史条目数
```
URL: localhost:9987/api/monitor/count
```

method GET

响应体
```
code    int
msg     string
count   int
```

# 6 用户效率统计
## 获取一个用户的效率
```
URL:localhost:9987/api/effiency/single/:id
```

method GET

响应体
```
code    int
msg     string
effiency  
    - user_info_id	int
	- average_remember_time	string // 平均记忆时间
	- forget_rate	string // 平均遗忘速度
	- response_time	string // 平均响应时间

```

## 获取效率统计总条目数
```
URL: localhost:9987/api/effiency/count
```

method GET

响应体
```
code    int
msg     string
count   int // 总条目数量
```


# 7. 复习模型(pattern)
复习模型是指第i轮记忆后，下次的记忆时间是什么时候。比如说没有复习过的话，要5分钟间隔， 复习过一次的话30分钟复习，复习过两次的话 下一次3小时复习。那么就有关系

**!!!一个pattern_id最多15轮!!!**
| id | pattern_id  |   round （复习过i轮） | time_gap （下一次的时间） | 
| :---: | :---: | --- | --- |
| 1 | 1  | 0 |  5 * 60 |
| 2 | 1  | 1 | 30 * 60 |
| 3 | 1   | 2 |  3 * 60 * 60 |




## 通过pattern_id获取复习模型(不是id哦)
```
URL: localhost:9987/api/pattern/get/:pattern_id
```

method GET

响应体：
```
code: int
msg: string
patterns:  arr
    - pattern_id int
    - round int
    - time_gap int
```

``` json
{
    "code": "0",
    "msg": "success",
    "patterns": [
        {
            "pattern_id": 1,
            "round": 0,
            "time_gap": 300
        },
        {
            "pattern_id": 2,
            "round": 1,
            "time_gap": 1800 
        }
    ] //  pattern length <= 15
}
```
## 新增复习模型

```
URL: localhost:9987/api/pattern/add
```

没时间做批量了，如果遇到批量的话，可以直接多次调这个接口么？qaq


method POST

请求体：
```
pattern_id int
round int
time_gap int
```

``` json
{    
    "pattern_id": 1,
    "round": 0,
    "time_gap": 300
}
```

## 修改复习模型
请求体：
```
URL: localhost:9987/api/pattern/edit/:id
```

没时间做批量了，如果遇到批量的话，可以直接多次调这个接口么？qaq


method POST

请求体：
```
pattern_id int
round int
time_gap int
```

``` json
{    
    "pattern_id": 1,
    "round": 0,
    "time_gap": 300
}
```