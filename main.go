package main

import (
	"fmt"
	"github.com/h2non/bimg"
	"os"
	"runtime"
	"sync"
	"time"
)

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
		wg            sync.WaitGroup
	)

	for epoch := 0; epoch < 10; epoch++ {
		fmt.Printf("queueing #%v\n", epoch)
		for req := 0; req < concurrentReq; req++ {
			wg.Add(1)
			go func(epoch, req int, imageBytes, watermarkBytes []byte) {
				defer fmt.Printf("done epoch #%1d req #%2d\n", epoch, req)
				defer wg.Done()

				var (
					cloneImage     = make([]byte, len(imageBytes))
					cloneWatermark = make([]byte, len(watermarkBytes))
				)

				copy(cloneImage, imageBytes)
				copy(cloneWatermark, watermarkBytes)

				img := bimg.NewImage(cloneImage)
				output, err := img.Process(bimg.Options{
					WatermarkImage: bimg.WatermarkImage{
						Left:    1000,
						Top:     1000,
						Buf:     cloneWatermark,
						Opacity: 0.5,
					},
				})
				if err != nil {
					panic(err)
				}
				_ = output
			}(epoch, req, imageBytes, watermarkBytes)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("waiting for all requests")

	wg.Wait()
	fmt.Println("process is done")
}
