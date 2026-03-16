package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"erp-server/internal/infrastructure/parser"
)

func main() {
	// 读取HTML目录
	htmlDir := "storage/html"
	files, err := os.ReadDir(htmlDir)
	if err != nil {
		fmt.Println("读取目录失败:", err)
		return
	}

	var results []map[string]interface{}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".html" {
			continue
		}

		fmt.Println("处理:", f.Name())

		html, err := os.ReadFile(filepath.Join(htmlDir, f.Name()))
		if err != nil {
			fmt.Printf("读取文件失败 %s: %v\n", f.Name(), err)
			continue
		}

		data, err := parser.OzonParse(html)
		if err != nil {
			fmt.Printf("解析失败 %s: %v\n", f.Name(), err)
			continue
		}

		results = append(results, data)
	}

	// 保存JSON
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	os.WriteFile("storage/all-products.json", jsonData, 0644)

	fmt.Printf("\n成功提取 %d 个商品\n", len(results))
	fmt.Println("已保存到: storage/all-products.json")
}
