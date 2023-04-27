package fay

type MsgType uint32

// 无 = 0,
// 弹幕消息 = 1,
// 点赞消息 = 2,
// 进直播间 = 3,
// 关注消息 = 4,
// 礼物消息 = 5,
// 直播间统计 = 6,
// 粉丝团消息 = 7,
// 直播间分享 = 8,
// 下播 = 9
const (
	MsgType_None     MsgType = iota // 无
	MsgType_DanMu                   // 弹幕信息
	MsgType_Dianzan                 // 点赞信息
	MsgType_JoinRoom                // 进直播间
	MsgType_Star                    // 关注信息
	MsgType_Gift                    // 礼物信息
	MsgType_Repo                    // 直播间统计
	MsgType_FansMsg                 // 粉丝信息
	MsgType_Share                   // 直播间分享
	MsgType_Offline                 // 下播
)

type FansclubType uint32

//无 = 0,
//粉丝团升级 = 1,
//加入粉丝团 = 2

const (
	FansclubType_None     FansclubType = iota
	FansclubType_Upgrade               // 升级
	FansclubType_JoinFans              // 加入粉丝团
)

type ShareType uint32

// 未知 = 0,
// 微信 = 1,
// 朋友圈 = 2,
// 微博 = 3,
// QQ空间 = 4,
// QQ = 5,
// 抖音好友 = 112
const (
	ShareType_Unkonwn      ShareType = iota // 未知
	ShareType_Wechat                        // 微信
	ShareType_WechatQuan                    // 朋友圈
	ShareType_Weibo                         // 微博
	ShareType_QQZone                        // qq空间
	ShareType_QQ                            // qq
	ShareType_DouyinFriend ShareType = 112  // 抖音好友
)

type MsgPack struct {
	MsgType MsgType `json:"type"`
	Data    string  `json:"data"`
}

func CreateMsgPack(data string, msgType MsgType) *MsgPack {
	return &MsgPack{
		MsgType: msgType,
		Data:    data,
	}
}
