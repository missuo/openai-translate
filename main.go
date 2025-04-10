/*
 * @Author: Vincent Young
 * @Date: 2023-11-13 11:16:26
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-04-10 08:37:58
 * @FilePath: /openai-translate/main.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2023 by Vincent, All Rights Reserved.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
)

type ResData struct {
	TransText  string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

func tokenCount(text string) (int, error) {
	tkm, err := tiktoken.EncodingForModel("gpt-4o")
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return 0, err
	}
	token := len(tkm.Encode(text, nil, nil))
	return token, nil
}

func translator(apiKey string, targetLang string, transText string, baseUrl string) (string, error) {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseUrl
	c := openai.NewClientWithConfig(config)
	// c := openai.NewClient(apiKey)
	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// Model: openai.GPT3Dot5Turbo,
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You're a translator. Translate to " + targetLang + ". Please note not to have any unnecessary explanations, just return the translated text directly.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: transText,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	// Define a command line flag
	apiKeyFlag := flag.String("apiKey", "", "API key for OpenAI")
	baseUrlFlag := flag.String("baseUrl", "", "Base URL, default is https://api.openai.com/v1")
	// modelFlag := flag.String("model", "", "Model to use, default is gpt-3.5-turbo")
	flag.Parse()

	// First try to get the API key from the command line flag
	apiKey := *apiKeyFlag
	baseUrl := *baseUrlFlag
	// model := *modelFlag

	// If it's not provided, try to get it from the environment variable
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_KEY")
	}

	if baseUrl == "" {
		if os.Getenv("BASE_URL") != "" {
			baseUrl = os.Getenv("BASE_URL")
		} else {
			baseUrl = "https://api.openai.com/v1"
		}
	}

	// if model == "" {
	// 	if os.Getenv("MODEL") != "" {
	// 		model = os.Getenv("MODEL")
	// 	} else {
	// 		model = "gpt-3.5-turbo"
	// 	}
	// }

	// fmt.Printf(model)

	// If the API key is still empty, return an error and exit
	if apiKey == "" {
		fmt.Println("Error: No API key provided. Set the apiKey flag or the OPENAI_KEY environment variable.")
		os.Exit(1)
	}

	fmt.Println("Starting server on port 23333. Made by Vincent.")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Hello",
		})
	})

	r.POST("/translate", func(c *gin.Context) {
		req := ResData{}
		c.BindJSON(&req)
		sourceLang := req.SourceLang
		targetLang := req.TargetLang
		translateText := req.TransText
		targetText, _ := translator(apiKey, targetLang, translateText, baseUrl)

		if targetText == "" {
			c.JSON(http.StatusTooManyRequests, gin.H{ // 429 Too Many Requests
				"code":    http.StatusTooManyRequests,
				"message": "Translation limit exceeded or service unavailable",
			})
			return
		}

		// if sourceLang is empty, detect the language
		if sourceLang == "" {
			sourceLang = strings.ToUpper(whatlanggo.DetectLang(req.TransText).Iso6391())
		}

		importToken, _ := tokenCount(translateText)
		importToken += 9 // 9 token for the prompt
		exportToken, _ := tokenCount(targetText)
		tokenConsumed := importToken + exportToken

		c.JSON(http.StatusOK, gin.H{
			"code":           http.StatusOK,
			"data":           targetText,
			"source_lang":    sourceLang,
			"target_lang":    targetLang,
			"token_consumed": tokenConsumed,
		})

	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})

	r.Run(":23333")
}
