package main

import (
	"context"
	"fmt"
	"github.com/h2non/bimg"
	"golang.org/x/sync/semaphore"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Request struct {
	Epoch          int
	ReqId          int
	ImageBytes     []byte
	WatermarkBytes []byte
	Wg             *sync.WaitGroup
}

func main() {
	fmt.Println("process is starting")

	imageBytes, err := os.ReadFile("./image.jpg")
	if err != nil {
		panic(err)
	}
	watermarkBytes, err := os.ReadFile("./watermark.png")
	if err != nil {
		panic(err)
	}

	var (
		concurrentReq = runtime.NumCPU() * 4
		reqQueue      = make(chan *Request, 1000)
		wg            = &sync.WaitGroup{}
	)

	concurrentReqString := os.Getenv("CONCURRENT_REQ")
	if concurrentReqString != "" {
		if temp, err := strconv.Atoi(concurrentReqString); err == nil {
			concurrentReq = temp
		}
	}

	go func() {
		worker(reqQueue)
	}()

	for epoch := 0; epoch < 10; epoch++ {
		fmt.Printf("queueing #%v\n", epoch)
		for req := 0; req < concurrentReq; req++ {
			wg.Add(1)
			go func(epoch, req int, imageBytes, watermarkBytes []byte) {
				var (
					cloneImage     = make([]byte, len(imageBytes))
					cloneWatermark = make([]byte, len(watermarkBytes))
				)

				copy(cloneImage, imageBytes)
				copy(cloneWatermark, watermarkBytes)

				reqQueue <- &Request{
					Epoch:          epoch,
					ReqId:          req,
					ImageBytes:     cloneImage,
					WatermarkBytes: cloneWatermark,
					Wg:             wg,
				}
			}(epoch, req, imageBytes, watermarkBytes)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("waiting for all requests")

	wg.Wait()
	fmt.Println("process is done")
	close(reqQueue)
}

func worker(queue chan *Request) {
	workerSize := int64(runtime.NumCPU())

	workerSizeString := os.Getenv("WORKER")
	if workerSizeString != "" {
		if temp, err := strconv.Atoi(workerSizeString); err == nil {
			workerSize = int64(temp)
		}
	}

	sem := semaphore.NewWeighted(workerSize)

	for {
		req, ok := <-queue
		if !ok {
			break
		}

		err := sem.Acquire(context.Background(), 1)
		if err != nil {
			panic(err)
		}

		go func(request *Request) {
			defer fmt.Printf("done epoch #%1d req #%2d\n", request.Epoch, request.ReqId)
			defer sem.Release(1)
			defer request.Wg.Done()

			img := bimg.NewImage(request.ImageBytes)
			_, err := img.Process(bimg.Options{
				WatermarkImage: bimg.WatermarkImage{
					Left:    1000,
					Top:     1000,
					Buf:     request.WatermarkBytes,
					Opacity: 0.5,
				},
			})
			if err != nil {
				panic(err)
			}
		}(req)
	}
}
