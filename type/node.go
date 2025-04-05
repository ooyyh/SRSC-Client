// ./nodes/node.go
package nodes

import "encoding/json"

// Node 结构体定义JSON中的节点信息
type Node struct {
	NodeName  string `json:"NodeName"`
	EndPoint  string `json:"EndPoint"`
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
	Region    string `json:"Region"`
}

// GetNodes 从JSON内容解析所有节点
func GetNodes(fileContent []byte) ([]Node, error) {
	var nodes []Node
	err := json.Unmarshal(fileContent, &nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
func AddNode(nodes []Node, node Node) []Node {
	nodes = append(nodes, node)
	return nodes
}
