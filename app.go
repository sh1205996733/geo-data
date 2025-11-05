package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed data/jp.csv
var jpRaw string

//go:embed data/cn.csv
var cnRaw string

var geoData = map[string]map[int]GeoData{
	"jp": makeGeoData(jpRaw),
	"cn": makeGeoData(cnRaw),
}

var postData = map[string]map[string][]GeoData{
	"jp": makePostData(geoData["jp"]),
	"cn": makePostData(geoData["cn"]),
}

type GeoData struct {
	Id       int       `json:"id"`                 // 地区id
	Level    int       `json:"level"`              // 地区级别
	Parent   int       `json:"parent"`             // 上级地区id
	Post     string    `json:"post"`               // 邮政编码
	Name     string    `json:"name"`               // 地区名称
	Address  string    `json:"address"`            // 地区全称
	Spell    string    `json:"spell"`              // 地区拼音或读音
	Children []GeoData `json:"children,omitempty"` // 子级地区
}

func makePostData(data map[int]GeoData) map[string][]GeoData {
	postDataMap := make(map[string][]GeoData)

	for _, data := range data {
		if data.Post != "" {
			postDataMap[data.Post] = append(postDataMap[data.Post], data)
		}
	}

	return postDataMap
}

func makeGeoData(rawData string) map[int]GeoData {
	lines := strings.Split(strings.TrimSpace(rawData), "\n")

	geoDataMapId := make(map[int]GeoData)

	for _, line := range lines {
		fields := strings.Split(line, ",")
		id, _ := strconv.Atoi(strings.TrimSpace(fields[0]))
		parent, _ := strconv.Atoi(strings.TrimSpace(fields[1]))
		name := strings.TrimSpace(fields[2])
		spell := strings.TrimSpace(fields[3])
		address := strings.TrimSpace(fields[4])
		post := strings.TrimSpace(fields[5])

		parentData, parentExists := geoDataMapId[parent]
		oldData, oldDataExists := geoDataMapId[id]

		level := 1
		if parentExists {
			level = parentData.Level + 1
		}

		children := []GeoData{}
		if oldDataExists {
			children = oldData.Children
		}

		cur := GeoData{
			Id:       id,
			Level:    level,
			Parent:   parent,
			Post:     post,
			Address:  address,
			Name:     name,
			Spell:    spell,
			Children: children,
		}

		geoDataMapId[id] = cur

		if parentExists {
			parentData.Children = append(parentData.Children, cur)
			geoDataMapId[parent] = parentData
		} else {
			geoDataMapId[parent] = GeoData{
				Level:    level - 1,
				Children: []GeoData{cur},
			}
		}
	}

	return geoDataMapId
}

func ginGetGeoData(c *gin.Context) {
	country := c.Param("country")
	idStr := c.Param("id")

	var id int
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
	} else {
		// 默认使用 0（根节点），以兼容之前的行为
		id = 0
	}

	// 检查国家是否存在
	countryMap, countryExists := geoData[country]
	if !countryExists {
		c.JSON(404, gin.H{"error": "Country not found"})
		return
	}

	// 返回地区或 404
	if info, exists := countryMap[id]; exists {
		c.JSON(200, info)
	} else {
		c.JSON(404, gin.H{"error": "ID not found"})
	}
}

func ginSearchByPost(c *gin.Context) {
	country := c.Query("country")
	postcode := c.Query("postcode")

	// 检查国家是否存在
	countryMap, countryExists := postData[country]
	if !countryExists {
		c.JSON(404, gin.H{"error": "Country not found"})
		return
	}

	// 搜索匹配的邮政编码
	var results = countryMap[postcode]

	if len(results) == 0 {
		c.JSON(404, gin.H{"error": "No matching post code found"})
	} else {
		c.JSON(200, results)
	}
}

func main() {
	r := gin.Default()

	// 路由：支持 /:country 和 /:country/:id
	r.GET("/:country", ginGetGeoData)
	r.GET("/:country/:id", ginGetGeoData)
	r.GET("/search_postcode/:country/:postcode", ginSearchByPost)

	// 输出加载的数据统计
	totalCodes := 0
	for countryCode, countryData := range geoData {
		fmt.Printf("Loaded %d %s codes\n", len(countryData), countryCode)
		totalCodes += len(countryData)
	}
	fmt.Printf("Total loaded %d codes\n", totalCodes)

	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
