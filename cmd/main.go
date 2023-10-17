package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ScoreSubmission struct {
	LevelID int    `json:"levelID" binding:"required"`
	Score   int    `json:"score" binding:"required"`
	Player  string `json:"player" binding:"required"`
}

type GetLeadersByLevel struct {
	LevelID int `json:"levelID" binding:"required"`
}

func main() {
	rdb, err := MakeRedisClient()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/leaderboards/getleaders", func(c *gin.Context) {
		var submission GetLeadersByLevel
		if err := c.BindJSON(&submission); err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid submission data",
				"error":   err.Error(),
			})
			return
		}
		key := fmt.Sprintf("leader:%d", submission.LevelID)

		res := rdb.ZRangeWithScores(ctx, key, 0, 1)
		if res.Err() != nil {
			c.JSON(500, gin.H{
				"message": "Querying redis",
				"error":   err,
			})
		}

		slice, err := res.Result()
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Querying redis",
				"error":   err,
			})
		}
		c.JSON(200, gin.H{
			"message": "retreived",
			"leaders": slice,
		})
	})

	r.POST("/leaderboards/submitscore", func(c *gin.Context) {
		var submission ScoreSubmission
		if err := c.BindJSON(&submission); err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid submission data",
				"error":   err.Error(),
			})
			return
		}

		key := fmt.Sprintf("leader:%d", submission.LevelID)
		res := rdb.ZAdd(ctx, key, redis.Z{
			Score:  float64(submission.Score),
			Member: submission.Player,
		})
		if res.Err() != nil {
			c.JSON(500, gin.H{
				"message": "Querying redis",
				"error":   err,
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Score received!",
			"levelID": submission.LevelID,
			"score":   submission.Score,
			"player":  submission.Player,
		})
	})

	r.Run(":8080") // Run on port 8080
}

func MakeRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.59:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ping := rdb.Ping(ctx)
	if ping.Err() != nil {
		return nil, ping.Err()
	}

	return rdb, nil
}
