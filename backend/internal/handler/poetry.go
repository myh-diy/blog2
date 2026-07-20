package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type poem struct {
	Title   string   `json:"title"`
	Content []string `json:"content"`
	Author  string   `json:"author"`
	Dynasty string   `json:"dynasty"`
}

var fallbackPoems = []poem{
	{Title: "山居秋暝", Author: "王维", Dynasty: "唐", Content: []string{"空山新雨后，天气晚来秋。", "明月松间照，清泉石上流。", "竹喧归浣女，莲动下渔舟。", "随意春芳歇，王孙自可留。"}},
	{Title: "静夜思", Author: "李白", Dynasty: "唐", Content: []string{"床前明月光，疑是地上霜。", "举头望明月，低头思故乡。"}},
	{Title: "饮湖上初晴后雨", Author: "苏轼", Dynasty: "宋", Content: []string{"水光潋滟晴方好，山色空蒙雨亦奇。", "欲把西湖比西子，淡妆浓抹总相宜。"}},
	{Title: "登鹳雀楼", Author: "王之涣", Dynasty: "唐", Content: []string{"白日依山尽，黄河入海流。", "欲穷千里目，更上一层楼。"}},
}

func GetRandomPoem(apiURL string) gin.HandlerFunc {
	client := &http.Client{Timeout: 4 * time.Second}

	return func(c *gin.Context) {
		result, err := fetchPoem(client, apiURL)
		if err != nil {
			result = fallbackPoems[rand.Intn(len(fallbackPoems))]
		}
		c.Header("Cache-Control", "no-store")
		c.JSON(http.StatusOK, gin.H{"poem": result})
	}
}

func fetchPoem(client *http.Client, apiURL string) (poem, error) {
	response, err := client.Get(apiURL)
	if err != nil {
		return poem{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return poem{}, fmt.Errorf("poetry API returned %s", response.Status)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return poem{}, err
	}
	var payload struct {
		Data struct {
			Title   string   `json:"title"`
			Content []string `json:"content"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
			Dynasty struct {
				Name string `json:"name"`
			} `json:"dynasty"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return poem{}, err
	}
	if payload.Data.Title == "" || len(payload.Data.Content) == 0 {
		return poem{}, fmt.Errorf("poetry API returned an empty poem")
	}
	return poem{
		Title:   payload.Data.Title,
		Content: payload.Data.Content,
		Author:  payload.Data.Author.Name,
		Dynasty: payload.Data.Dynasty.Name,
	}, nil
}
