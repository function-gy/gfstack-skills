// ================================================================
// 示例: Storager 文件存储 (internal/library/storager/)
// 驱动模式：支持本地/UCloud/COS/OSS/七牛/Minio
// ================================================================

package storager

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model/entity"
	"hotgo/utility/format"
	"hotgo/utility/url"
	"hotgo/utility/validate"
)

// UploadDrive 存储驱动接口
type UploadDrive interface {
	Upload(ctx context.Context, file *ghttp.UploadFile) (fullPath string, err error)
	CreateMultipart(ctx context.Context, in *CheckMultipartParams) (res *MultipartProgress, err error)
	UploadPart(ctx context.Context, in *UploadPartParams) (res *UploadPartModel, err error)
}

// New 初始化存储驱动
func New(name ...string) UploadDrive {
	driveType := consts.UploadDriveLocal
	if len(name) > 0 && name[0] != "" {
		driveType = name[0]
	}
	switch driveType {
	case consts.UploadDriveLocal:
		return &LocalDrive{}
	case consts.UploadDriveUCloud:
		return &UCloudDrive{}
	case consts.UploadDriveCos:
		return &CosDrive{}
	case consts.UploadDriveOss:
		return &OssDrive{}
	case consts.UploadDriveQiNiu:
		return &QiNiuDrive{}
	case consts.UploadDriveMinio:
		return &MinioDrive{}
	default:
		panic(fmt.Sprintf("暂不支持的存储驱动:%v", driveType))
	}
}

// DoUpload 上传入口：验证 → 存储 → 写DB
func DoUpload(ctx context.Context, typ string, file *ghttp.UploadFile) (result *entity.SysAttachment, err error) {
	if file == nil {
		return nil, gerror.New("文件必须!")
	}

	meta, err := GetFileMeta(file)
	if err != nil {
		return nil, err
	}

	// 生成文件名: {date}/{hash}.{ext}
	name, err := GenerateFileName(ctx, file)
	if err != nil {
		return nil, err
	}

	// 执行存储
	fullPath, err := New().Upload(ctx, file)
	if err != nil {
		return nil, err
	}

	// 写入数据库
	result = &entity.SysAttachment{
		FileUrl:  fullPath,
		FileName: file.Filename,
		FileSize: file.Size,
		FileType: meta.Ext,
		FileMd5:  meta.Md5,
		Status:   consts.StatusEnabled,
		MemberId: contexts.GetUserId(ctx),
	}
	_, err = dao.SysAttachment.Ctx(ctx).Data(result).Insert()
	return
}

// GenerateFileName 生成随机文件名: {date}/{hash}.{ext}
func GenerateFileName(ctx context.Context, file *ghttp.UploadFile) (string, error) {
	meta, err := GetFileMeta(file)
	if err != nil {
		return "", err
	}
	dateDir := gtime.Now().Format("Ymd")
	hashName := strings.ToLower(format.Md5(grand.Letters(16) + gtime.Now().String())) + meta.Ext
	return url.GenFullPath("file", dateDir+"/"+hashName), nil
}

// GetFileMeta 获取文件元数据
func GetFileMeta(file *ghttp.UploadFile) (*FileMeta, error) {
	meta := &FileMeta{}
	meta.Filename = file.Filename
	meta.Size = file.Size

	var pos = strings.LastIndex(file.Filename, ".")
	if pos < 0 {
		return nil, gerror.New("仅支持带后缀的文件")
	}
	meta.Ext = file.Filename[pos:]

	if !validate.IsSliceExistStr(GetAllowExt(), meta.Ext) {
		return nil, gerror.Newf("不支持的文件类型:%v", meta.Ext)
	}

	meta.Md5, _ = format.Md5File(file.FilePath)
	return meta, nil
}
