package nodes

import (
	"time"
)

// BucketInfo 存储桶信息结构体
type BucketInfo struct {
	Name                string    `json:"name"`                // 桶名称
	CreationDate        time.Time `json:"creationDate"`        // 创建时间
	UsedSpace           int64     `json:"usedSpace"`           // 已用空间(字节)
	TotalObjects        int64     `json:"totalObjects"`        // 对象总数
	VersioningEnabled   bool      `json:"versioningEnabled"`   // 版本控制是否启用
	PublicAccessBlocked bool      `json:"publicAccessBlocked"` // 是否阻止公共访问
	HasPolicy           bool      `json:"hasPolicy"`           // 是否有桶策略
	EncryptionEnabled   bool      `json:"encryptionEnabled"`   // 是否启用加密
	EncryptionType      string    `json:"encryptionType"`      // 加密类型
	HasLifecycleRules   bool      `json:"hasLifecycleRules"`   // 是否有生命周期规则
	LifecycleRulesCount int       `json:"lifecycleRulesCount"` // 生命周期规则数量
	Region              string    `json:"region"`              // 桶所在区域
	WebsiteEnabled      bool      `json:"websiteEnabled"`      // 是否启用静态网站
}

// NodeBucketInfo 节点桶信息结构体
type NodeBucketInfo struct {
	NodeName string       `json:"nodeName"` // 节点名称
	EndPoint string       `json:"endPoint"` // 节点端点
	Buckets  []BucketInfo `json:"buckets"`  // 桶信息列表
}
