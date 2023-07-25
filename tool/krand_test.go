package tool

import (
	"fmt"
	"testing"
)

type UserAccount struct {
	BasicAmount   int64 // 基本账户金额
	GiftAmount    int64 // 赠送账户金额
	BasicDeducted int64 // 基本账户扣费金额记录
	GiftDeducted  int64 // 赠送账户扣费金额记录
}

func TestMd5ByString(t *testing.T) {
	user := UserAccount{
		BasicAmount:   0,  // 基本账户初始金额
		GiftAmount:    25, // 赠送账户初始金额
		BasicDeducted: 0,  // 初始化基本账户扣费金额记录
		GiftDeducted:  0,  // 初始化赠送账户扣费金额记录
	}
	var cost int64 = 300 // 消费金额

	// 先尝试从基本账户扣费
	if user.BasicAmount >= cost {
		user.BasicAmount -= cost
		user.BasicDeducted += cost
	} else {
		remaining := cost - user.BasicAmount
		user.BasicDeducted += user.BasicAmount
		user.BasicAmount = 0

		// 检查赠送账户是否足够扣费
		if user.GiftAmount >= remaining {
			user.GiftAmount -= remaining
			user.GiftDeducted += remaining
		} else {
			user.GiftDeducted += user.GiftAmount
			remaining -= user.GiftAmount // 更新剩余金额
			user.GiftAmount = 0

			// 扣除剩余金额从基本账户（可为负数）
			user.BasicAmount -= remaining
			user.BasicDeducted += remaining
		}
	}

	fmt.Println("基本账户剩余金额:", user.BasicAmount)
	fmt.Println("赠送账户剩余金额:", user.GiftAmount)
	fmt.Println("基本账户扣费金额:", user.BasicDeducted)
	fmt.Println("赠送账户扣费金额:", user.GiftDeducted)
}
