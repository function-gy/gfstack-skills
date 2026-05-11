// ================================================================
// 示例: SMS 短信发送 (internal/library/sms/)
// 驱动模式：支持阿里云和腾讯云短信
// ================================================================

package sms

import (
	"context"
	"fmt"

	"hotgo/internal/consts"
	"hotgo/internal/model"
	"hotgo/internal/model/input/sysin"
)

// Drive 短信驱动接口
type Drive interface {
	SendCode(ctx context.Context, in *sysin.SendCodeInp, conf *model.SmsConfig) (err error)
}

// New 创建短信驱动实例
func New(name ...string) Drive {
	var (
		instanceName = consts.SmsDriveAliYun
		drive        Drive
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	switch instanceName {
	case consts.SmsDriveAliYun:
		drive = &AliYunDrive{}
	case consts.SmsDriveTencent:
		drive = &TencentDrive{}
	default:
		panic(fmt.Sprintf("暂不支持短信驱动:%v", instanceName))
	}
	return drive
}
