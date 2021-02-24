//Package
/*
 * @Author: zhoupeng
 * @Date: 2021-02-24 16:12:53
 * @LastEditTime: 2021-02-24 17:57:57
 * @LastEditors: zhoupeng
 * @Description: API结构定义
 * @FilePath: /GoChat/server/datamodel.go
 */
package main

import "time"

// 客户端发送消息结构定义
type ClientComMessage struct {

	// 内部字段

	//客户端消息ID
	Id string `json:"-"` //不进行序列化
	// 没有格式化的topic名称
	Original string `json:"-"`
	//格式化的topic名称
	RcptTo string `json:"-"`
	//该消息的发送者
	AsUser string `json:"-"`
	//消息发送者的权限水平
	AuthLvl int `json:"-"`
	//消息中的what元数据
	MetaWhat int `json:"-"`
	//消息被server接受时的时间戳
	TimeStamp time.Time `json:"-"`
}

// 客户端发送给服务器的消息

// MsgClientHi 是客户端给server发送的握手消息
type MsgClientHi struct {
	//消息ID
	Id string `json:"id,omitempty"`
	//用户代理标识
	UserAgent string `json:"ua,omitempty"`
	//消息的协议版本
	Version string `json:"ver,omitempty"`
	//客户端设备唯一id
	DeviceID string `json:"dev,omitempty"`
	//使用的语言,国际化
	Lang string `json:"lang,omitempty"`
	//平台代码: ios,android,web
	Platform string `json:"platf,omitempty"`
	//session的状态，后台还是前台
	Background bool `json:"bkg,omitempty"`
}

//MsgClientAcc 是客户端发起创建用户，或者更新用户状态的消息结构
type MsgClientAcc struct {
	//消息Id
	Id string `json:"id,omitempty"`
	// "newXYZ"  创建一个用户，或者传入一个UserId进行更新
	User string `json:"user,omitempty"`
	// 用户的状态: normal, suspended
	State string `json:"status,omitempty"`
	//当不是当前用户在进行更新操作时，需要认证权限:"","auth","anon",默认是"""
	AuthLevel string
	//当进行修改密码等操作时，需要token鉴权
	Token []byte `json:"token,omitempty"`
	//用户可以使用的认证方案
	Scheme string `json:"scheme,omitempty"`
	//密码
	Secret []byte `json:"secret,omitempty"`
	// 当在进行登录时，需要对token进行鉴权
	Login bool `json:"login,omitempty"`
	// 用户搜索时的标签
	Tags []string `json:"tags,omitempty"`
	//创建一个新用户时，一些初始化数据
	Desc *MsgSetDesc `json:"desc,omitempty"`
	//验证方式:email,phone,captcha
	Cred []MsgCredClient `json:"cred,omitempty"`
}

//MsgClientLogin 是客户端发起的登录请求
type MsgClientLogin struct {
	//消息ID
	Id string `json:"id,omitempty"`
	//认证方案
	Scheme string `json:"scheme,omitempty"`
	//密码
	Secret []byte `json:"secret"`
	//验证方式
	Cred []MsgCredClient `json:"cred,omitempty"`
}

//MsgClientSub是客户端发送的订阅消息
type MsgClientSub struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`

	Set *MsgSetQuery `json:"set,omitempty"`
	Get *MsgGetQuery `json:"get,omitempty"`

	//该subscription是否会新建一个topic
	Created bool `json:"-"`
	//是否是一个新的订阅
	NewSub bool `json:"-"`
}

//MsgCredClient 是账户的验证方式
type MsgCredClient struct {
	//验证的方式,比如"email","tel"
	Method string `json:"meth,omitempty"`
	//存储的值
	Value string `json:"val,omitempty"`
	//验证响应值
	Response string `json:"resp,omitempty"`
	//请求的相关参数
	Params map[string]interface{} `json:"params,omitempty"`
}

// MsgSetQuery 用来更新topic的meta data
type MsgSetQuery struct {
	Desc *MsgSetDesc    `json:"desc,omitempty"`
	Sub  *MsgSetSub     `json:"sub,omitempty"`
	Tags []string       `json:"tags,omitempty"`
	Cred *MsgCredClient `json:"cred,omitempty"`
}

// MsgGetQuery 用来查看topic的相关数据
type MsgGetQuery struct {
	What string `json:"what"`

	//Desc需要参数：IfModifiedSince
	Desc *MsgGetOpts `json:"desc,omitempty"`
	//sub需要参数: User,Topic, IfModifiedSince,Limit
	Sub *MsgGetOpts `json:"sub,omitempty"`
	// data参数：since,before,limit
	Data *MsgGetOpts `json:"data,omitempty"`
	// del参数:since,before,limit
	Del *MsgGetOpts `json:"del,omitempty"`
}

// MsgGetOpts 定义了客户端需要查询的参数
type MsgGetOpts struct {
	//用户ID
	User string `json:"user,omitempty"`
	//topic的名称
	Topic string `json:"topic,omitempty"`
	//返回从指定时间修改之后的结果
	IfModifiedSince *time.Time `json:"ims,omitempty"`
	//加载从此Id开始的消息
	SinceId int `json:"since,omitempty"`
	//加载小于此id的消息
	BeforeId int `json:"before,omitempty"`
	//限制显示消息的数量
	Limit int `json:"limit,omitempty"`
}

// MsgSetSub 用户更新topic的订阅
type MsgSetSub struct {
	User string `json:"user,omitempty"`
	Mode string `json:"mode,omitempty"`
}

// MsgSetDesc 是用户描述信息
type MsgSetDesc struct {
	DefaultAcs *MsgDefaultAcsMode `json:"defacs,omitempty"`
	Public     interface{}        `json:"public:omitempty"`
	Private    interface{}        `json:"private,omitempty"`
}

//MsgDefaultAcsMode 是一个topic默认的权限模式
type MsgDefaultAcsMode struct {
	Auth string `json:"auth,omitempty"`
	Anon string `json:"anon,omitempty"`
}
