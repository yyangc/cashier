package model

import (
	"time"
)

// MemberType 會員類型
type MemberType int8

const (
	MemberTypeUnknown MemberType = iota
	MemberTypeVIP                // VIP
	MemberTypePro                // Pro
)

// Member 會員當前等級
type Member struct {
	ID        int32
	UserID    int64      // 用戶ID
	Type      MemberType // 會員類型
	Level     int8       // 會員等級 e.g. 1, 2, 3 ...
	CreatedAt time.Time  // 創建時間
	UpdatedAt time.Time  // 更新時間
}
