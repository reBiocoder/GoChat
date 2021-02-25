//Package
/*
 * @Author: zhoupeng
 * @Date: 2021-02-24 16:12:53
 * @LastEditTime: 2021-02-25 14:43:19
 * @LastEditors: zhoupeng
 * @Description: API结构定义
 * @FilePath: /GoChat/server/datamodel.go
 */
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"GoChat/server/store/types"
)

// 客户端发送消息结构定义
type ClientComMessage struct {
	Hi    *MsgClientHi    `json:"hi"`
	Acc   *MsgClientAcc   `json:"acc"`
	Login *MsgClientLogin `json:"login"`
	Sub   *MsgClientSub   `json:"sub"`
	Leave *MsgClientLeave `json:"leave"`
	Pub   *MsgClientPub   `json:"pub"`
	Get   *MsgClientGet   `json:"get"`
	Set   *MsgClientSet   `json:"set"`
	Del   *MsgClientDel   `json:"del"`
	Note  *MsgClientNote  `json:"note"`

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
	Timestamp time.Time `json:"-"`
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

const (
	constMsgMetaDesc = 1 << iota
	constMsgMetaSub
	constMsgMetaData
	constMsgMetaTags
	constMsgMetaDel
	constMsgMetaCred
)
const (
	constMsgDelTopic = iota + 1
	constMsgDelMsg
	constMsgDelSub
	constMsgDelUser
	constMsgDelCred
)

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

type MsgClientLeave struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	Unsub bool   `json:"unsub,omitempty"`
}

type MsgClientPub struct {
	Id      string                 `json:"id,omitempty"`
	Topic   string                 `json:"topic"`
	NoEcho  bool                   `json:"noecho,omitempty"`
	Head    map[string]interface{} `json:"head,omitempty"`
	Content interface{}            `json:"content"`
}
type MsgClientGet struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	MsgGetQuery
}

type MsgClientSet struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	MsgSetQuery
}
type MsgClientDel struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic,omitempty"`
	// What to delete:
	// * "msg" to delete messages (default)
	// * "topic" to delete the topic
	// * "sub" to delete a subscription to topic.
	// * "user" to delete or disable user.
	// * "cred" to delete credential (email or phone)
	What string `json:"what"`
	// Delete messages with these IDs (either one by one or a set of ranges)
	DelSeq []MsgDelRange `json:"delseq,omitempty"`
	// User ID of the user or subscription to delete
	User string `json:"user,omitempty"`
	// Credential to delete
	Cred *MsgCredClient `json:"cred,omitempty"`
	// Request to hard-delete objects (i.e. delete messages for all users), if such option is available.
	Hard bool `json:"hard,omitempty"`
}
type MsgDelRange struct {
	LowId int `json:"low,omitempty"`
	HiId  int `json:"hi,omitempty"`
}
type MsgClientNote struct {
	// There is no Id -- server will not akn {ping} packets, they are "fire and forget"
	Topic string `json:"topic"`
	// what is being reported: "recv" - message received, "read" - message read, "kp" - typing notification
	What string `json:"what"`
	// Server-issued message ID being reported
	SeqId int `json:"seq,omitempty"`
	// Client's count of unread messages to report back to the server. Used in push notifications on iOS.
	Unread int `json:"unread,omitempty"`
}

// server到客户端的数据结构

//user最后一次在线的ua和时间
type MsgLastSeenInfo struct {
	When      time.Time `json:"when,omitempty"`
	UserAgent string    `json:"ua,omitempty"`
}

func (src *MsgLastSeenInfo) describe() string {
	return "'" + src.UserAgent + "' @ " + src.When.String()
}

// MsgAccessMode 用来得到user具有的权限
type MsgAccessMode struct {
	Want  string `json:"want,omitempty"`
	Given string `json:"given,omitempty"`
	Mode  string `json:"mode,omitempty"`
}

func (src *MsgAccessMode) describe() string {
	var s string
	if src.Want != "" {
		s = "w=" + src.Want
	}
	if src.Given != "" {
		s += " g=" + src.Given
	}
	if src.Mode != "" {
		s += " m=" + src.Mode
	}
	return strings.TrimSpace(s)
}

type MsgTopicDesc struct {
	CreatedAt *time.Time `json:"created,omitempty"`
	UpdatedAt *time.Time `json:"updated,omitempty"`
	// Timestamp of the last message
	TouchedAt *time.Time `json:"touched,omitempty"`

	// Account state, 'me' topic only.
	State string `json:"state,omitempty"`

	// If the group topic is online.
	Online bool `json:"online,omitempty"`

	DefaultAcs *MsgDefaultAcsMode `json:"defacs,omitempty"`
	// Actual access mode
	Acs *MsgAccessMode `json:"acs,omitempty"`
	// Max message ID
	SeqId     int `json:"seq,omitempty"`
	ReadSeqId int `json:"read,omitempty"`
	RecvSeqId int `json:"recv,omitempty"`
	// Id of the last delete operation as seen by the requesting user
	DelId  int         `json:"clear,omitempty"`
	Public interface{} `json:"public,omitempty"`
	// Per-subscription private data
	Private interface{} `json:"private,omitempty"`
}

func (src *MsgTopicDesc) describe() string {
	var s string
	if src.State != "" {
		s = " state=" + src.State
	}
	s += " online=" + strconv.FormatBool(src.Online)
	if src.Acs != nil {
		s += " acs={" + src.Acs.describe() + "}"
	}
	if src.SeqId != 0 {
		s += " seq=" + strconv.Itoa(src.SeqId)
	}
	if src.ReadSeqId != 0 {
		s += " read=" + strconv.Itoa(src.ReadSeqId)
	}
	if src.RecvSeqId != 0 {
		s += " recv=" + strconv.Itoa(src.RecvSeqId)
	}
	if src.DelId != 0 {
		s += " clear=" + strconv.Itoa(src.DelId)
	}
	if src.Public != nil {
		s += " pub='...'"
	}
	if src.Private != nil {
		s += " priv='...'"
	}
	return s
}

// MsgTopicSub is topic subscription details, sent in Meta message.
type MsgTopicSub struct {
	// Fields common to all subscriptions

	// Timestamp when the subscription was last updated
	UpdatedAt *time.Time `json:"updated,omitempty"`
	// Timestamp when the subscription was deleted
	DeletedAt *time.Time `json:"deleted,omitempty"`

	// If the subscriber/topic is online
	Online bool `json:"online,omitempty"`

	// Access mode. Topic admins receive the full info, non-admins receive just the cumulative mode
	// Acs.Mode = want & given. The field is not a pointer because at least one value is always assigned.
	Acs MsgAccessMode `json:"acs,omitempty"`
	// ID of the message reported by the given user as read
	ReadSeqId int `json:"read,omitempty"`
	// ID of the message reported by the given user as received
	RecvSeqId int `json:"recv,omitempty"`
	// Topic's public data
	Public interface{} `json:"public,omitempty"`
	// User's own private data per topic
	Private interface{} `json:"private,omitempty"`

	// Response to non-'me' topic

	// Uid of the subscribed user
	User string `json:"user,omitempty"`

	// The following sections makes sense only in context of getting
	// user's own subscriptions ('me' topic response)

	// Topic name of this subscription
	Topic string `json:"topic,omitempty"`
	// Timestamp of the last message in the topic.
	TouchedAt *time.Time `json:"touched,omitempty"`
	// ID of the last {data} message in a topic
	SeqId int `json:"seq,omitempty"`
	// Id of the latest Delete operation
	DelId int `json:"clear,omitempty"`

	// P2P topics only:

	// Other user's last online timestamp & user agent
	LastSeen *MsgLastSeenInfo `json:"seen,omitempty"`
}

func (src *MsgTopicSub) describe() string {
	s := src.Topic + ":" + src.User + " online=" + strconv.FormatBool(src.Online) + " acs=" + src.Acs.describe()

	if src.SeqId != 0 {
		s += " seq=" + strconv.Itoa(src.SeqId)
	}
	if src.ReadSeqId != 0 {
		s += " read=" + strconv.Itoa(src.ReadSeqId)
	}
	if src.RecvSeqId != 0 {
		s += " recv=" + strconv.Itoa(src.RecvSeqId)
	}
	if src.DelId != 0 {
		s += " clear=" + strconv.Itoa(src.DelId)
	}
	if src.Public != nil {
		s += " pub='...'"
	}
	if src.Private != nil {
		s += " priv='...'"
	}
	if src.LastSeen != nil {
		s += " seen={" + src.LastSeen.describe() + "}"
	}
	return s
}

// MsgDelValues describes request to delete messages.
type MsgDelValues struct {
	DelId  int           `json:"clear,omitempty"`
	DelSeq []MsgDelRange `json:"delseq,omitempty"`
}

// MsgServerCtrl is a server control message {ctrl}.
type MsgServerCtrl struct {
	Id     string      `json:"id,omitempty"`
	Topic  string      `json:"topic,omitempty"`
	Params interface{} `json:"params,omitempty"`

	Code      int       `json:"code"`
	Text      string    `json:"text,omitempty"`
	Timestamp time.Time `json:"ts"`
}

// Deep-shallow copy.
func (src *MsgServerCtrl) copy() *MsgServerCtrl {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}
func (src *MsgServerCtrl) describe() string {
	return src.Topic + " id=" + src.Id + " code=" + strconv.Itoa(src.Code) + " txt=" + src.Text
}

// MsgServerData is a server {data} message.
type MsgServerData struct {
	Topic string `json:"topic"`
	// ID of the user who originated the message as {pub}, could be empty if sent by the system
	From      string                 `json:"from,omitempty"`
	Timestamp time.Time              `json:"ts"`
	DeletedAt *time.Time             `json:"deleted,omitempty"`
	SeqId     int                    `json:"seq"`
	Head      map[string]interface{} `json:"head,omitempty"`
	Content   interface{}            `json:"content"`
}

// Deep-shallow copy.
func (src *MsgServerData) copy() *MsgServerData {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func (src *MsgServerData) describe() string {
	s := src.Topic + " from=" + src.From + " seq=" + strconv.Itoa(src.SeqId)
	if src.DeletedAt != nil {
		s += " deleted"
	} else {
		if src.Head != nil {
			s += " head=..."
		}
		s += " content='...'"
	}
	return s
}

// MsgServerPres is presence notification {pres} (authoritative update).
type MsgServerPres struct {
	Topic     string        `json:"topic"`
	Src       string        `json:"src,omitempty"`
	What      string        `json:"what"`
	UserAgent string        `json:"ua,omitempty"`
	SeqId     int           `json:"seq,omitempty"`
	DelId     int           `json:"clear,omitempty"`
	DelSeq    []MsgDelRange `json:"delseq,omitempty"`
	AcsTarget string        `json:"tgt,omitempty"`
	AcsActor  string        `json:"act,omitempty"`
	// Acs or a delta Acs. Need to marshal it to json under a name different than 'acs'
	// to allow different handling on the client
	Acs *MsgAccessMode `json:"dacs,omitempty"`

	// UNroutable params. All marked with `json:"-"` to exclude from json marshalling.
	// They are still serialized for intra-cluster communication.

	// Flag to break the reply loop
	WantReply bool `json:"-"`

	// Additional access mode filters when sending to topic's online members. Both filter conditions must be true.
	// send only to those who have this access mode.
	FilterIn int `json:"-"`
	// skip those who have this access mode.
	FilterOut int `json:"-"`

	// When sending to 'me', skip sessions subscribed to this topic.
	SkipTopic string `json:"-"`

	// Send to sessions of a single user only.
	SingleUser string `json:"-"`

	// Exclude sessions of a single user.
	ExcludeUser string `json:"-"`
}

// Deep-shallow copy.
func (src *MsgServerPres) copy() *MsgServerPres {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func (src *MsgServerPres) describe() string {
	s := src.Topic
	if src.Src != "" {
		s += " src=" + src.Src
	}
	if src.What != "" {
		s += " what=" + src.What
	}
	if src.UserAgent != "" {
		s += " ua=" + src.UserAgent
	}
	if src.SeqId != 0 {
		s += " seq=" + strconv.Itoa(src.SeqId)
	}
	if src.DelId != 0 {
		s += " clear=" + strconv.Itoa(src.DelId)
	}
	if src.DelSeq != nil {
		s += " delseq"
	}
	if src.AcsTarget != "" {
		s += " tgt=" + src.AcsTarget
	}
	if src.AcsActor != "" {
		s += " actor=" + src.AcsActor
	}
	if src.Acs != nil {
		s += " dacs=" + src.Acs.describe()
	}

	return s
}

// MsgServerMeta is a topic metadata {meta} update.
type MsgServerMeta struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`

	Timestamp *time.Time `json:"ts,omitempty"`

	// Topic description
	Desc *MsgTopicDesc `json:"desc,omitempty"`
	// Subscriptions as an array of objects
	Sub []MsgTopicSub `json:"sub,omitempty"`
	// Delete ID and the ranges of IDs of deleted messages
	Del *MsgDelValues `json:"del,omitempty"`
	// User discovery tags
	Tags []string `json:"tags,omitempty"`
	// Account credentials, 'me' only.
	Cred []*MsgCredServer `json:"cred,omitempty"`
}

// Deep-shallow copy of meta message. Deep copy of Id and Topic fields, shallow copy of payload.
func (src *MsgServerMeta) copy() *MsgServerMeta {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

func (src *MsgServerMeta) describe() string {
	s := src.Topic + " id=" + src.Id

	if src.Desc != nil {
		s += " desc={" + src.Desc.describe() + "}"
	}
	if src.Sub != nil {
		var x []string
		for _, sub := range src.Sub {
			x = append(x, sub.describe())
		}
		s += " sub=[{" + strings.Join(x, "},{") + "}]"
	}
	if src.Del != nil {
		x, _ := json.Marshal(src.Del)
		s += " del={" + string(x) + "}"
	}
	if src.Tags != nil {
		s += " tags=[" + strings.Join(src.Tags, ",") + "]"
	}
	if src.Cred != nil {
		x, _ := json.Marshal(src.Cred)
		s += " cred=[" + string(x) + "]"
	}
	return s
}

// MsgServerInfo is the server-side copy of MsgClientNote with From and optionally Src added (non-authoritative).
type MsgServerInfo struct {
	// Topic to send event to.
	Topic string `json:"topic"`
	// Topic where the even has occured (set only when Topic='me').
	Src string `json:"src,omitempty"`
	// ID of the user who originated the message.
	From string `json:"from"`
	// The event being reported: "rcpt" - message received, "read" - message read, "kp" - typing notification.
	What string `json:"what"`
	// Server-issued message ID being reported.
	SeqId int `json:"seq,omitempty"`

	// UNroutable params. All marked with `json:"-"` to exclude from json marshalling.
	// They are still serialized for intra-cluster communication.

	// When sending to 'me', skip sessions subscribed to this topic.
	SkipTopic string `json:"-"`
}

// Deep copy
func (src *MsgServerInfo) copy() *MsgServerInfo {
	if src == nil {
		return nil
	}
	dst := *src
	return &dst
}

// Basic description
func (src *MsgServerInfo) describe() string {
	s := src.Topic
	if src.Src != "" {
		s += " src=" + src.Src
	}
	s += " what=" + src.What + " from=" + src.From
	if src.SeqId > 0 {
		s += " seq=" + strconv.Itoa(src.SeqId)
	}
	return s
}

type ServerComMessage struct {
	Ctrl *MsgServerCtrl `json:"ctrl,omitempty"`
	Data *MsgServerData `json:"data,omitempty"`
	Meta *MsgServerMeta `json:"meta,omitempty"`
	Pres *MsgServerPres `json:"pres,omitempty"`
	Info *MsgServerInfo `json:"info,omitempty"`

	// Internal fields.

	// MsgServerData has no Id field, copying it here for use in {ctrl} aknowledgements
	Id string `json:"-"`
	// Routable (expanded) name of the topic.
	RcptTo string `json:"-"`
	// User ID of the sender of the original message.
	AsUser string `json:"-"`
	// Timestamp for consistency of timestamps in {ctrl} messages
	// (corresponds to originating client message receipt timestamp).
	Timestamp time.Time `json:"-"`
	sess      *Session
	// Originating session to send an aknowledgement to. Could be nil.
	// Session ID to skip when sendng packet to sessions. Used to skip sending to original session.
	// Could be either empty.
	SkipSid string `json:"-"`
	uid     types.Uid
}

func (src *ServerComMessage) copy() *ServerComMessage {
	if src == nil {
		return nil
	}
	dst := &ServerComMessage{
		Id:        src.Id,
		RcptTo:    src.RcptTo,
		AsUser:    src.AsUser,
		Timestamp: src.Timestamp,
		sess:      src.sess,
		SkipSid:   src.SkipSid,
		uid:       src.uid,
	}

	dst.Ctrl = src.Ctrl.copy()
	dst.Data = src.Data.copy()
	dst.Meta = src.Meta.copy()
	dst.Pres = src.Pres.copy()
	dst.Info = src.Info.copy()

	return dst
}

func (src *ServerComMessage) describe() string {
	if src == nil {
		return "-"
	}

	switch {
	case src.Ctrl != nil:
		return "{ctrl " + src.Ctrl.describe() + "}"
	case src.Data != nil:
		return "{data " + src.Data.describe() + "}"
	case src.Meta != nil:
		return "{meta " + src.Meta.describe() + "}"
	case src.Pres != nil:
		return "{pres " + src.Pres.describe() + "}"
	case src.Info != nil:
		return "{info " + src.Info.describe() + "}"
	default:
		return "{nil}"
	}
}

// Generators of server-side error messages {ctrl}.

// NoErr indicates successful completion (200)
func NoErr(id, topic string, ts time.Time) *ServerComMessage {
	return NoErrParams(id, topic, ts, nil)
}

// NoErrExplicitTs indicates successful completion with explicit server and incoming request timestamps (200)
func NoErrExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return NoErrParamsExplicitTs(id, topic, serverTs, incomingReqTs, nil)
}

// NoErrReply indicates successful completion as a reply to a client message (200).
func NoErrReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return NoErrExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// NoErrParams indicates successful completion with additional parameters (200)
func NoErrParams(id, topic string, ts time.Time, params interface{}) *ServerComMessage {
	return NoErrParamsExplicitTs(id, topic, ts, ts, params)
}

// NoErrParamsExplicitTs indicates successful completion with additional parameters
// and explicit server and incoming request timestamps (200)
func NoErrParamsExplicitTs(id, topic string, serverTs, incomingReqTs time.Time, params interface{}) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusOK, // 200
		Text:      "ok",
		Topic:     topic,
		Params:    params,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// NoErrParamsReply indicates successful completion with additional parameters
// and explicit server and incoming request timestamps (200)
func NoErrParamsReply(msg *ClientComMessage, ts time.Time, params interface{}) *ServerComMessage {
	return NoErrParamsExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp, params)
}

// NoErrCreated indicated successful creation of an object (201).
func NoErrCreated(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusCreated, // 201
		Text:      "created",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// NoErrAccepted indicates request was accepted but not perocessed yet (202).
func NoErrAccepted(id, topic string, ts time.Time) *ServerComMessage {
	return NoErrAcceptedExplicitTs(id, topic, ts, ts)
}

// NoErrAcceptedExplicitTs indicates request was accepted but not perocessed yet
// with explicit server and incoming request timestamps (202).
func NoErrAcceptedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusAccepted, // 202
		Text:      "accepted",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// NoContentParams indicates request was processed but resulted in no content (204).
func NoContentParams(id, topic string, serverTs, incomingReqTs time.Time, params interface{}) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNoContent, // 204
		Text:      "no content",
		Topic:     topic,
		Params:    params,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// NoContentParamsReply indicates request was processed but resulted in no content
// in response to a client request (204).
func NoContentParamsReply(msg *ClientComMessage, ts time.Time, params interface{}) *ServerComMessage {
	return NoContentParams(msg.Id, msg.Original, ts, msg.Timestamp, params)
}

// NoErrEvicted indicates that the user was disconnected from topic for no fault of the user (205).
func NoErrEvicted(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusResetContent, // 205
		Text:      "evicted",
		Topic:     topic,
		Timestamp: ts}, Id: id}
}

// NoErrShutdown means user was disconnected from topic because system shutdown is in progress (205).
func NoErrShutdown(ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Code:      http.StatusResetContent, // 205
		Text:      "server shutdown",
		Timestamp: ts}}
}

// NoErrDeliveredParams means requested content has been delivered (208).
func NoErrDeliveredParams(id, topic string, ts time.Time, params interface{}) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusAlreadyReported, // 208
		Text:      "delivered",
		Topic:     topic,
		Params:    params,
		Timestamp: ts}, Id: id}
}

// 3xx

// InfoValidateCredentials requires user to confirm credentials before going forward (300).
func InfoValidateCredentials(id string, ts time.Time) *ServerComMessage {
	return InfoValidateCredentialsExplicitTs(id, ts, ts)
}

// InfoValidateCredentialsExplicitTs requires user to confirm credentials before going forward
// with explicit server and incoming request timestamps (300).
func InfoValidateCredentialsExplicitTs(id string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusMultipleChoices, // 300
		Text:      "validate credentials",
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// InfoChallenge requires user to respond to presented challenge before login can be completed (300).
func InfoChallenge(id string, ts time.Time, challenge []byte) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusMultipleChoices, // 300
		Text:      "challenge",
		Params:    map[string]interface{}{"challenge": challenge},
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// InfoAuthReset is sent in response to request to reset authentication when it was completed but login was not performed (301).
func InfoAuthReset(id string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusMovedPermanently, // 301
		Text:      "auth reset",
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// InfoUseOther is a response to a subscription request redirecting client to another topic (303).
func InfoUseOther(id, topic, other string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusSeeOther, // 303
		Text:      "use other",
		Topic:     topic,
		Params:    map[string]string{"topic": other},
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// InfoUseOtherReply is a response to a subscription request redirecting client to another topic (303).
func InfoUseOtherReply(msg *ClientComMessage, other string, ts time.Time) *ServerComMessage {
	return InfoUseOther(msg.Id, msg.Original, other, ts, msg.Timestamp)
}

// InfoAlreadySubscribed response means request to subscribe was ignored because user is already subscribed (304).
func InfoAlreadySubscribed(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotModified, // 304
		Text:      "already subscribed",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// InfoNotJoined response means request to leave was ignored because user was not subscribed (304).
func InfoNotJoined(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotModified, // 304
		Text:      "not joined",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// InfoNoAction response means request was ignored because the object was already in the desired state
// with explicit server and incoming request timestamps (304).
func InfoNoAction(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotModified, // 304
		Text:      "no action",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// InfoNoActionReply response means request was ignored because the object was already in the desired state
// in response to a client request (304).
func InfoNoActionReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return InfoNoAction(msg.Id, msg.Original, ts, msg.Timestamp)
}

// InfoNotModified response means update request was a noop (304).
func InfoNotModified(id, topic string, ts time.Time) *ServerComMessage {
	return InfoNotModifiedExplicitTs(id, topic, ts, ts)
}

// InfoNotModifiedReply response means update request was a noop
// in response to a client request (304).
func InfoNotModifiedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return InfoNotModifiedExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// InfoNotModifiedExplicitTs response means update request was a noop
// with explicit server and incoming request timestamps (304).
func InfoNotModifiedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotModified, // 304
		Text:      "not modified",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// InfoFound redirects to a new resource (307).
func InfoFound(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusTemporaryRedirect, // 307
		Text:      "found",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// 4xx Errors

// ErrMalformed request malformed (400).
func ErrMalformed(id, topic string, ts time.Time) *ServerComMessage {
	return ErrMalformedExplicitTs(id, topic, ts, ts)
}

// ErrMalformedReply request malformed
// in response to a client request (400).
func ErrMalformedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrMalformedExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrMalformedExplicitTs request malformed with explicit server and incoming request timestamps (400).
func ErrMalformedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusBadRequest, // 400
		Text:      "malformed",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrAuthRequired authentication required  - user must authenticate first (401).
func ErrAuthRequired(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusUnauthorized, // 401
		Text:      "authentication required",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrAuthRequiredReply authentication required  - user must authenticate first
// in response to a client request (401).
func ErrAuthRequiredReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrAuthRequired(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrAuthFailed authentication failed
// with explicit server and incoming request timestamps (400).
func ErrAuthFailed(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusUnauthorized, // 401
		Text:      "authentication failed",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrAuthUnknownScheme authentication scheme is unrecognized or invalid (401).
func ErrAuthUnknownScheme(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusUnauthorized, // 401
		Text:      "unknown authentication scheme",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrPermissionDenied user is authenticated but operation is not permitted (403).
func ErrPermissionDenied(id, topic string, ts time.Time) *ServerComMessage {
	return ErrPermissionDeniedExplicitTs(id, topic, ts, ts)
}

// ErrPermissionDeniedExplicitTs user is authenticated but operation is not permitted
// with explicit server and incoming request timestamps (403).
func ErrPermissionDeniedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusForbidden, // 403
		Text:      "permission denied",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrPermissionDeniedReply user is authenticated but operation is not permitted
// with explicit server and incoming request timestamps in response to a client request (403).
func ErrPermissionDeniedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrPermissionDeniedExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrAPIKeyRequired  valid API key is required (403).
func ErrAPIKeyRequired(ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Code:      http.StatusForbidden,
		Text:      "valid API key required",
		Timestamp: ts}}
}

// ErrSessionNotFound  valid API key is required (403).
func ErrSessionNotFound(ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Code:      http.StatusForbidden,
		Text:      "invalid or expired session",
		Timestamp: ts}}
}

// ErrTopicNotFound topic is not found
// with explicit server and incoming request timestamps (404).
func ErrTopicNotFound(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotFound,
		Text:      "topic not found", // 404
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrTopicNotFoundReply topic is not found
// with explicit server and incoming request timestamps
// in response to a client request (404).
func ErrTopicNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrTopicNotFound(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrUserNotFound user is not found
// with explicit server and incoming request timestamps (404).
func ErrUserNotFound(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotFound, // 404
		Text:      "user not found",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrUserNotFoundReply user is not found
// with explicit server and incoming request timestamps in response to a client request (404).
func ErrUserNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrUserNotFound(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrNotFound is an error for missing objects other than user or topic
// with explicit server and incoming request timestamps (404).
func ErrNotFound(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotFound, // 404
		Text:      "not found",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrNotFoundReply is an error for missing objects other than user or topic
// with explicit server and incoming request timestamps in response to a client request (404).
func ErrNotFoundReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrNotFound(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrOperationNotAllowed a valid operation is not permitted in this context (405).
func ErrOperationNotAllowed(id, topic string, ts time.Time) *ServerComMessage {
	return ErrOperationNotAllowedExplicitTs(id, topic, ts, ts)
}

// ErrOperationNotAllowedExplicitTs a valid operation is not permitted in this context
// with explicit server and incoming request timestamps (405).
func ErrOperationNotAllowedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusMethodNotAllowed, // 405
		Text:      "operation or method not allowed",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrOperationNotAllowedReply a valid operation is not permitted in this context
// with explicit server and incoming request timestamps (405).
func ErrOperationNotAllowedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrOperationNotAllowedExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrInvalidResponse indicates that the client's response in invalid
// with explicit server and incoming request timestamps (406).
func ErrInvalidResponse(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotAcceptable, // 406
		Text:      "invalid response",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrAlreadyAuthenticated invalid attempt to authenticate an already authenticated session
// Switching users is not supported (409).
func ErrAlreadyAuthenticated(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusConflict, // 409
		Text:      "already authenticated",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrDuplicateCredential attempt to create a duplicate credential
// with explicit server and incoming request timestamps (409).
func ErrDuplicateCredential(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusConflict, // 409
		Text:      "duplicate credential",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrAttachFirst must attach to topic first in response to a client message (409).
func ErrAttachFirst(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        msg.Id,
		Code:      http.StatusConflict, // 409
		Text:      "must attach first",
		Topic:     msg.Original,
		Timestamp: ts}, Id: msg.Id, Timestamp: msg.Timestamp}
}

// ErrAlreadyExists the object already exists (409).
func ErrAlreadyExists(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusConflict, // 409
		Text:      "already exists",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrCommandOutOfSequence invalid sequence of comments, i.e. attempt to {sub} before {hi} (409).
func ErrCommandOutOfSequence(id, unused string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusConflict, // 409
		Text:      "command out of sequence",
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrGone topic deleted or user banned (410).
func ErrGone(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusGone, // 410
		Text:      "gone",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrTooLarge packet or request size exceeded the limit (413).
func ErrTooLarge(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusRequestEntityTooLarge, // 413
		Text:      "too large",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrPolicy request violates a policy (e.g. password is too weak or too many subscribers) (422).
func ErrPolicy(id, topic string, ts time.Time) *ServerComMessage {
	return ErrPolicyExplicitTs(id, topic, ts, ts)
}

// ErrPolicyExplicitTs request violates a policy (e.g. password is too weak or too many subscribers)
// with explicit server and incoming request timestamps (422).
func ErrPolicyExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusUnprocessableEntity, // 422
		Text:      "policy violation",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrPolicyReply request violates a policy (e.g. password is too weak or too many subscribers)
// with explicit server and incoming request timestamps in response to a client request (422).
func ErrPolicyReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrPolicyExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrUnknown database or other server error (500).
func ErrUnknown(id, topic string, ts time.Time) *ServerComMessage {
	return ErrUnknownExplicitTs(id, topic, ts, ts)
}

// ErrUnknownExplicitTs database or other server error with explicit server and incoming request timestamps (500).
func ErrUnknownExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusInternalServerError, // 500
		Text:      "internal error",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrUnknownReply database or other server error in response to a client request (500).
func ErrUnknownReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrUnknownExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrNotImplemented feature not implemented with explicit server and incoming request timestamps (501).
// TODO: consider changing status code to 4XX.
func ErrNotImplemented(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusNotImplemented, // 501
		Text:      "not implemented",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrClusterUnreachable in-cluster communication has failed (502).
func ErrClusterUnreachable(id, topic string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusBadGateway, // 502
		Text:      "cluster unreachable",
		Topic:     topic,
		Timestamp: ts}, Id: id, Timestamp: ts}
}

// ErrServiceUnavailableReply server error in response to a client request (503).
func ErrServiceUnavailableReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrServiceUnavailableExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrServiceUnavailableExplicitTs server error (503).
func ErrServiceUnavailableExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusServiceUnavailable, // 503
		Text:      "service unavailable",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrLocked operation rejected because the topic is being deleted (503).
func ErrLocked(id, topic string, ts time.Time) *ServerComMessage {
	return ErrLockedExplicitTs(id, topic, ts, ts)
}

// ErrLockedReply operation rejected because the topic is being deleted with explicit server and
// incoming request timestamps in response to a client request (503).
func ErrLockedReply(msg *ClientComMessage, ts time.Time) *ServerComMessage {
	return ErrLockedExplicitTs(msg.Id, msg.Original, ts, msg.Timestamp)
}

// ErrLockedExplicitTs operation rejected because the topic is being deleted
// with explicit server and incoming request timestamps (503).
func ErrLockedExplicitTs(id, topic string, serverTs, incomingReqTs time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusServiceUnavailable, // 503
		Text:      "locked",
		Topic:     topic,
		Timestamp: serverTs}, Id: id, Timestamp: incomingReqTs}
}

// ErrVersionNotSupported invalid (too low) protocol version (505).
func ErrVersionNotSupported(id string, ts time.Time) *ServerComMessage {
	return &ServerComMessage{Ctrl: &MsgServerCtrl{
		Id:        id,
		Code:      http.StatusHTTPVersionNotSupported, // 505
		Text:      "version not supported",
		Timestamp: ts}, Id: id, Timestamp: ts}
}
