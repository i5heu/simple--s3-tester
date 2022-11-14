package main

import (
	"fmt"
	"math/rand"
	"simple-s3-benchmark/m/config"
	"time"

	"github.com/valyala/fasthttp"
)

// build a "image" random image generator

func main() {
	conf := config.GetValues()
	dataset := createDataset(conf)

	go requester(conf, dataset)

	fasthttp.ListenAndServe(":8085", handler)
}

func handler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("image/jpeg; charset=utf8")
	ctx.SetStatusCode(200)
	ctx.SetBody(make([]byte, 50000000))
}

func createDataset(conf config.Config) []string {
	dataset := []string{}
	datasetItems := conf.DatasetSize / conf.FileSize

	for i := 0; i < datasetItems; i++ {
		dataset = append(dataset, RandStringBytes(10))
	}

	return dataset
}

func TookTime(timeStart time.Time) {
	fmt.Println("took time: ", time.Since(timeStart))
}

func requester(conf config.Config, dataset []string) {
	currentRefPosition := uint(0)
	movePerSecond := len(dataset) / conf.TestDuration
	timeTMP := time.Now().Unix()

	fmt.Println("movePerSecond: ", movePerSecond)

	for {
		if time.Now().Unix() > timeTMP {
			timeTMP = time.Now().Unix()
			currentRefPosition += uint(movePerSecond)
		}

		time.Sleep(time.Second / time.Duration(conf.RequestPerSecond))

		go func() {
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			req.SetRequestURI(conf.S3Endpoint + "/" + positionToRequest(dataset, currentRefPosition, conf))
			err := fasthttp.Do(req, resp)
			if err != nil {
				fmt.Println(err)
			}
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
		}()
	}
}

func positionToRequest(dataset []string, currentRefPosition uint, conf config.Config) string {
	movingPopularFileCount := conf.MovingPopularSize / conf.FileSize

	// get a centered distribution of the popular files
	distributionPosition := rand.Intn(movingPopularFileCount/2) + rand.Intn(movingPopularFileCount/2)

	finalPosition := (currentRefPosition - (uint(movingPopularFileCount) / 2)) + uint(distributionPosition)

	if finalPosition < 0 {
		finalPosition = 0
	}
	if finalPosition >= uint(len(dataset)) {
		return dataset[len(dataset)-1]
	}

	fmt.Println("finalPosition: ", finalPosition)
	return dataset[finalPosition]
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
