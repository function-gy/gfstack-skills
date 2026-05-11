// ================================================================
// 示例: Dict 字典系统 (internal/library/dict/)
// ================================================================

package dict

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"sync"

	"hotgo/internal/model"
)

const (
	BuiltinId int64 = -1 // 内置字典ID
	EnumsId   int64 = -2 // 枚举字典ID
	FuncId    int64 = -3 // 方法字典ID
)

var NotExistKeyError = errors.New("not exist key")

// ==================== 枚举字典 ====================

type EnumsOption struct {
	Id    int64
	Key   string
	Label string
	Opts  []*model.Option
}

var (
	enumsOptions = make(map[string]*EnumsOption)
	eLock        sync.Mutex
)

// RegisterEnums 注册枚举字典选项（通常在 init() 中调用）
func RegisterEnums(key, label string, opts []*model.Option) {
	eLock.Lock()
	defer eLock.Unlock()
	if len(key) == 0 {
		panic("字典key不能为空")
	}
	if _, ok := enumsOptions[key]; ok {
		panic(fmt.Sprintf("重复注册枚举字典选项:%v", key))
	}
	for _, v := range opts {
		v.Type = key
	}
	enumsOptions[key] = &EnumsOption{
		Id:    GenIdHash(key, EnumsId),
		Key:   key,
		Label: label,
		Opts:  opts,
	}
}

// GetEnumsOptions 获取指定枚举字典的数据选项
func GetEnumsOptions(key string) []*model.Option {
	if enums, ok := enumsOptions[key]; ok {
		return enums.Opts
	}
	return nil
}

func GetAllEnums() map[string]*EnumsOption {
	return enumsOptions
}

// ==================== 字典查找 ====================

// GetOptions 获取内置选项（先查枚举，再查函数字典）
func GetOptions(ctx context.Context, key string) (opts []*model.Option, err error) {
	opts = GetEnumsOptions(key)
	if opts != nil {
		return
	}
	return GetFuncOptions(ctx, key)
}

// GenIdHash 生成字典ID
func GenIdHash(str string, t int64) int64 {
	prefix := 10000 * t
	h := fnv.New32a()
	h.Write([]byte("dict" + str))
	idStr := fmt.Sprintf("%d%d", prefix, int64(h.Sum32()))
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}
