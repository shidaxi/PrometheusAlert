package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis/v8"
)

type AppConfig struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Ou          string            `json:"ou"`
	Authz       map[string]string `json:"authz"`
	WorkflowId  string            `json:"workflowId"`
	Notify      struct {
		Alert map[string]map[string]string `json:"alert"`
		Event map[string]map[string]string `json:"event"`
	} `json:"notify"`

	Owners map[string]struct {
		Mobihex string `json:"mobihex"`
	} `json:"owners"`

	Git struct {
		Repo          string `json:"repo"`
		DefaultBranch string `json:"defaultBranch"`
	}
}

var cicdServiceRedisKey = "appworks:cicd:services"
var ctx = context.Background()

func getAppConfig(app string, env string) AppConfig {
	var appConfig AppConfig
	redisHost := beego.AppConfig.String("APPWORKS_REDIS_HOST")
	var rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       0,
	})

	res, err := rdb.HGet(ctx, cicdServiceRedisKey, fmt.Sprintf("%s:%s", app, env)).Result()
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		return AppConfig{}
	}
	if err := json.Unmarshal([]byte(res), &appConfig); err != nil {
		panic(err)
	}
	return appConfig
}
func getAppNotifyLarkIds(app string, env string) []string {
	var larkBotIds []string
	for k, _ := range getAppConfig(app, env).Notify.Alert["larkBot"] {
		larkBotIds = append(larkBotIds, k)
	}
	return larkBotIds
}

func removeDuplicateString(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
