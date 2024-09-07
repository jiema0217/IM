package models

import (
	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model
	Name          string `json:"name" gorm:"type:varchar(100);not null;comment:'名字'"`
	Password      string `json:"password" gorm:"type:varchar(100);not null;comment:'密码'"`
	Phone         string `json:"phone" gorm:"type:char(20);comment:'手机号码'"`
	Email         string `json:"email" gorm:"type:varchar(50);comment:'邮箱'"`
	Identity      string `json:"identity" gorm:"type:varchar(50);comment:'唯一标识'"`
	ClientIP      string `json:"client_IP" gorm:"type:char(15);comment:'设备IP'"`
	ClientPort    string `json:"client_port" gorm:"type:char(15);comment:'设备端口'"`
	OnlineTime    uint   `json:"online_time" gorm:"type:bigint;comment:'登陆时间'"`
	OfflineTime   uint   `json:"offline_time" gorm:"type:bigint;comment:'下线时间'"`
	HeartBeatTime uint   `json:"heart_beat_time" gorm:"type:bigint;comment:'心跳时间'"`
	IsOffline     bool   `json:"is_offline" gorm:"type:tinyint(1);comment:'是否下线'"`
	DeviceInfo    string `json:"deviceInfo" gorm:"type:varchar(50);comment:'设备信息'"`
}

// TableName 返回用户信息表的表名
// 该方法主要用于数据库操作，确保对正确的表进行操作
func (userInfo *UserInfo) TableName() string {
	return "user_info"
}
