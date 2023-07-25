package uniqueid

import (
	"fmt"
	"github.com/nononai/mycommon/tool"
	"time"
)

// 生成单号 第三方支付：THD 系统支付SYS
func GenSn(snPrefix string) string {
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().Format("20060102150405"), tool.Krand(8, tool.KC_RAND_KIND_NUM))
}
