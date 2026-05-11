// ================================================================
// 示例: IP 地理位置 (internal/library/location/location.go)
// 支持 Whois 在线查询 + Cz88 本地IP库
// ================================================================

package location

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/kayon/iploc"

	"hotgo/utility/validate"
)

const (
	whoisApi = "https://whois.pconline.com.cn/ipJson.jsp?json=true&ip="
	dyndns   = "http://members.3322.org/dyndns/getip"
)

type IpLocationData struct {
	Ip           string `json:"ip"`
	Country      string `json:"country"`
	Region       string `json:"region"`
	Province     string `json:"province"`
	ProvinceCode int64  `json:"province_code"`
	City         string `json:"city"`
	CityCode     int64  `json:"city_code"`
	Area         string `json:"area"`
	AreaCode     int64  `json:"area_code"`
}

// GetLocation 获取IP归属地（优先Cz88本地库，支持切换whois在线）
func GetLocation(ctx context.Context, ip string) (data *IpLocationData, err error) {
	if !validate.IsIp(ip) {
		return nil, fmt.Errorf("invalid input ip:%v", ip)
	}
	if validate.IsLocalIPAddr(ip) {
		return
	}

	// 双重检查锁缓存
	if cache.Contains(ip) {
		return cache.GetIpCache(ip)
	}
	cache.Lock()
	defer cache.Unlock()
	if cache.Contains(ip) {
		return cache.GetIpCache(ip)
	}

	mode := g.Cfg().MustGet(ctx, "system.ipMethod", "cz88").String()
	switch mode {
	case "whois":
		data, err = WhoisLocation(ctx, ip)
	default:
		data, err = Cz88Find(ctx, ip)
	}

	if err == nil && data != nil {
		cache.SetIpCache(ip, data)
	}
	return
}

// Cz88Find 通过Cz88本地IP库查询（qqwry-utf8.dat）
func Cz88Find(ctx context.Context, ip string) (*IpLocationData, error) {
	loc, err := iploc.OpenWithoutIndexes("./resource/ip/qqwry-utf8.dat")
	if err != nil {
		return nil, err
	}
	detail := loc.Find(ip)
	if detail == nil {
		return nil, fmt.Errorf("no ip data: %v", ip)
	}
	return &IpLocationData{
		Ip: ip, Country: detail.Country, Region: detail.Region,
		Province: detail.Province, City: detail.City, Area: detail.County,
	}, nil
}

// WhoisLocation 通过Whois在线接口查询
func WhoisLocation(ctx context.Context, ip string, retry ...int64) (*IpLocationData, error) {
	response, err := g.Client().Timeout(10 * time.Second).Get(ctx, whoisApi+ip)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	str, err := gcharset.ToUTF8("GBK", response.ReadAllString())
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		retryCount := int64(3)
		if len(retry) > 0 {
			retryCount = retry[0]
		}
		if retryCount > 0 {
			return WhoisLocation(ctx, ip, retryCount-1)
		}
	}

	var who WhoisRegionData
	if err = gconv.Scan([]byte(str), &who); err != nil {
		return nil, gerror.Newf("WhoisLocation Scan err:%v, str:%v", err, str)
	}
	return &IpLocationData{
		Ip: who.Ip, Region: who.Addr, Province: who.Pro,
		ProvinceCode: gconv.Int64(who.ProCode), City: who.City,
		CityCode: gconv.Int64(who.CityCode), Area: who.Region,
		AreaCode: gconv.Int64(who.RegionCode),
	}, nil
}

// GetClientIp 获取客户端真实IP（优先 X-Forwarded-For > 直连IP）
func GetClientIp(r *ghttp.Request) string {
	if r == nil {
		return ""
	}
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.GetClientIp()
	}
	// 多级代理取第一个
	if gstr.Contains(ip, ",") {
		ip = gstr.StrTillEx(ip, ",")
	}
	if gstr.Contains(ip, ", ") {
		ip = gstr.StrTillEx(ip, ", ")
	}
	return ip
}

// GetLocalIP 获取服务器内网IP
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok || ipAddr.IP.IsLoopback() || !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// GetPublicIP 获取公网IP
func GetPublicIP(ctx context.Context) (ip string, err error) {
	var data WhoisRegionData
	err = g.Client().Timeout(10*time.Second).GetVar(ctx, whoisApi).Scan(&data)
	if err != nil || data.Ip == "" {
		return GetPublicIP2()
	}
	return data.Ip, nil
}

func GetPublicIP2() (ip string, err error) {
	response, err := http.Get(dyndns)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	return strings.ReplaceAll(string(body), "\n", ""), nil
}
