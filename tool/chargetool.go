package tool

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/w3liu/go-common/constant/timeformat"
	"github.com/zeromicro/go-zero/core/logx"
)

// 设备id十六进制转设备id10进制
func ConverDeviceID(hexDeviceID string) int {
	hexDeviceID = hexDeviceID[0:6]
	swappedSubstring := hexDeviceID[4:6] + hexDeviceID[2:4] + hexDeviceID[0:2]
	decimalID, err := strconv.ParseInt(swappedSubstring, 16, 64)
	if err != nil {
		return 0
	}
	return int(decimalID)
}

// 生成24位订单号
// 前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
var num int64

func Generate(t time.Time) string {
	s := t.Format(timeformat.Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

// 对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

// 友电订单号生成规则
func GenerateChargeOrderID(deviceID string, port int64) string {
	timestamp := time.Now().Format("20060102150405") // 年月日时分秒
	randomNum := generateRandomNum(4)                // 生成8位随机数
	devicePort := convertDevicePort(deviceID, port)  // 将设备号和端口号转换为10进制

	return timestamp + randomNum + devicePort
}

func generateRandomNum(length int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func convertDevicePort(deviceID string, port int64) string {
	deviceNum, _ := strconv.ParseInt(deviceID, 16, 64)
	// portNum, _ := strconv.ParseInt(port, 10, 64)

	devicePort := (deviceNum << 8) + int64(port)

	return fmt.Sprintf("%010d", devicePort)
}

// 设备端口号处理，设备默认是从0开始，且为string，如0为“00”，以此类推
func ConvertDevicePort(port int64) string {
	if port < 10 {
		return fmt.Sprintf("0%v", port-1)
	}
	return fmt.Sprintf("%v", port-1)
}

// 充电时间转化，设备接收的是秒,且为string
func ConvertChargeTime(chargeTime int64) string {
	return fmt.Sprintf("%v", chargeTime*60)
}

// 生成设备messageId：随机数 4位
func GenerateMessageID() string {
	return generateRandomNum(4)
}

// 根据当前功率得出当前功率所属档次的金额--特定格式！！！
func GetPowerPrice(power int, combo string) *int {
	var money *int
	type PowerLevel struct {
		Grade    string `json:"grade"`
		MinPower int    `json:"minPower"`
		MaxPower int    `json:"maxPower"`
		Money    int    `json:"money"`
	}
	type PowerLevelsData struct {
		PowerLevels []PowerLevel `json:"powerLevels"`
	}
	var data PowerLevelsData
	err := json.Unmarshal([]byte(combo), &data)
	if err != nil {
		logx.Error(err)
		return nil
	}
	var selectedLevel PowerLevel
	for _, level := range data.PowerLevels {
		if power >= level.MinPower && power <= level.MaxPower {
			selectedLevel = level
			break
		}
	}
	if selectedLevel.Grade != "" {
		money = &selectedLevel.Money
		fmt.Printf("当前功率：%d,所需消费金额档次：%s,金额：%d\n", power, selectedLevel.Grade, selectedLevel.Money)
	} else {
		money = nil
		fmt.Println("未找到匹配的功率档次")
	}
	return money
}
