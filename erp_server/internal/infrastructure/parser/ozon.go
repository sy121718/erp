package parser

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// OzonParse 解析HTML，返回map供调用方使用
func OzonParse(html []byte) (map[string]interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	// 提取JSON-LD数据
	jsonLD, err := extractJSONLD(doc)
	if err != nil {
		return nil, err
	}

	// 清理description中的换行和多余空格
	desc := jsonLD.Description
	desc = strings.ReplaceAll(desc, "\n", " ")
	desc = strings.ReplaceAll(desc, "\r", " ")
	for strings.Contains(desc, "  ") {
		desc = strings.ReplaceAll(desc, "  ", " ")
	}
	desc = strings.TrimSpace(desc)

	result := map[string]interface{}{
		"source_product_id": jsonLD.SKU,
		"title":             jsonLD.Name,
		"description":       desc,
		"main_image":        extractMainImage(jsonLD.Image),
		"brand":             jsonLD.Brand,
		"price":             jsonLD.Offers.Price,
		"currency":          jsonLD.Offers.PriceCurrency,
		"stock_status":      jsonLD.Offers.Availability,
		"rating":            jsonLD.AggregateRating.RatingValue,
		"review_count":      jsonLD.AggregateRating.ReviewCount,
		"source_url":        jsonLD.Offers.URL,
	}

	// 提取SKU变体
	result["skus"] = extractSKUs(doc)

	// 提取图片
	result["images"] = extractImages(doc, result["main_image"].(string))

	// 提取属性
	result["properties"] = extractProperties(jsonLD.Description)

	return result, nil
}

// JSONLD ozon的JSON-LD结构
type JSONLD struct {
	SKU             string          `json:"sku"`
	Name            string          `json:"name"`
	Image           any             `json:"image"`
	Brand           string          `json:"brand"`
	Offers          Offers          `json:"offers"`
	AggregateRating AggregateRating `json:"aggregateRating"`
	Description     string          `json:"description"`
}

type Offers struct {
	URL           string `json:"url"`
	Availability  string `json:"availability"`
	Price         string `json:"price"`
	PriceCurrency string `json:"priceCurrency"`
}

type AggregateRating struct {
	RatingValue string `json:"ratingValue"`
	ReviewCount string `json:"reviewCount"`
}

// extractJSONLD 提取JSON-LD数据
func extractJSONLD(doc *goquery.Document) (*JSONLD, error) {
	var jsonLD JSONLD

	doc.Find("script[type=\"application/ld+json\"]").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		// 处理HTML实体
		content = bytesToString([]byte(content))
		json.Unmarshal([]byte(content), &jsonLD)
	})

	return &jsonLD, nil
}

func bytesToString(data []byte) string {
	return string(data)
}

// extractMainImage 提取主图
func extractMainImage(image any) string {
	if image == nil {
		return ""
	}
	switch img := image.(type) {
	case string:
		return img
	case []interface{}:
		if len(img) > 0 {
			if s, ok := img[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// extractSKUs 提取SKU变体
func extractSKUs(doc *goquery.Document) []map[string]interface{} {
	skuMap := make(map[string]map[string]interface{})

	// 从页面中提取SKU数据
	doc.Find("div[data-widget=\"webProductVariants\"]").Each(func(i int, s *goquery.Selection) {
		s.Find("div[data-state]").Each(func(j int, div *goquery.Selection) {
			state, _ := div.Attr("data-state")
			if state == "" {
				return
			}

			var data map[string]interface{}
			json.Unmarshal([]byte(state), &data)

			if items, ok := data["items"].([]interface{}); ok {
				for _, item := range items {
					if m, ok := item.(map[string]interface{}); ok {
						skuID := ""
						if id, ok := m["sku"].(string); ok {
							skuID = id
						}

						price := ""
						if p, ok := m["price"].(map[string]interface{}); ok {
							if basic, ok := p["basic"].(string); ok {
								price = basic
							}
						}

						color := ""
						if c, ok := m["searchableText"].(string); ok {
							color = c
						}

						img := ""
						if im, ok := m["coverImage"].(string); ok {
							img = cleanURL(im)
						}

						stock := 1
						if a, ok := m["availability"].(string); ok {
							if a != "inStock" {
								stock = 0
							}
						}

						if skuID != "" {
							skuMap[skuID] = map[string]interface{}{
								"sku_id":       skuID,
								"price":        price,
								"stock_status": stock,
								"image_url":    img,
								"color":        color,
							}
						}
					}
				}
			}
		})
	})

	skus := make([]map[string]interface{}, 0, len(skuMap))
	for _, sku := range skuMap {
		skus = append(skus, sku)
	}
	return skus
}

func cleanURL(url string) string {
	return replaceAll(url, `\\u002F`, "/")
}

func replaceAll(s, old, new string) string {
	result := s
	for {
		i := findSubstring(result, old)
		if i == -1 {
			break
		}
		result = result[:i] + new + result[i+len(old):]
	}
	return result
}

// extractImages 提取图片列表
func extractImages(doc *goquery.Document, mainImage string) []string {
	imgSet := make(map[string]bool)

	if mainImage != "" {
		imgSet[mainImage] = true
	}

	// 从SKU中提取图片
	skus := extractSKUs(doc)
	for _, sku := range skus {
		if img, ok := sku["image_url"].(string); ok && img != "" {
			imgSet[img] = true
		}
	}

	images := make([]string, 0, len(imgSet))
	for img := range imgSet {
		images = append(images, img)
	}
	return images
}

// extractProperties 从描述中提取属性
func extractProperties(desc string) []map[string]string {
	if desc == "" {
		return nil
	}

	var properties []map[string]string

	// 修复正则，直接匹配
	patterns := []struct {
		regex string
		name  string
	}{
		{`Размеры\s*\(см\):\s*примерно\s*([^\.]+)`, "尺寸"},
		{`Размеры:\s*([^\n]+)`, "尺寸"},
		{`Размер\s*\(см\):\s*([^\.]+)`, "尺寸"},
		{`Цвет:\s*([^\n]+)`, "颜色"},
		{`Количество:\s*([^\n]+)`, "数量"},
		{`Материал:\s*([^\n]+)`, "材料"},
		{`В\s*упаковке:\s*([^\.]+)`, "包装"},
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.regex)
		matches := re.FindAllStringSubmatch(desc, -1)
		for _, m := range matches {
			if len(m) > 1 {
				value := strings.TrimSpace(m[1])
				if len(value) > 0 && len(value) < 100 {
					// 避免重复
					exists := false
					for _, prop := range properties {
						if prop["name"] == p.name {
							exists = true
							break
						}
					}
					if !exists {
						properties = append(properties, map[string]string{
							"name":  p.name,
							"value": value,
						})
					}
				}
			}
		}
	}

	if len(properties) == 0 {
		return nil
	}
	return properties
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
