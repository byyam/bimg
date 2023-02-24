package main

import "C"
import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/h2non/bimg"
)

var (
	framesDir = "/Users/zhangyan/Documents/images/frames2/"
	delayMs   = flag.Int("d", 1000, "delay ms")
	loop      = flag.Int("l", 1, "loop times")
	frames    = flag.Int("f", 0, "max frames")
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// daemon server checking for leaking memory issues
func main() {
	flag.Parse()
	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	// go func(c chan os.Signal) {
	//      sig := <-c
	//      log.Printf("Caught signal %s: shutting down.", sig)
	//      // release resources
	//      mmsvips.Shutdown()
	//      os.Exit(0)
	// }(signals)

	defer func() {
		fmt.Printf("shutdown and exit\n")
		bimg.Shutdown()
	}()

	startTime := time.Now()

	// for i := 0; i < 100; i++ {
	// 	time.Sleep(time.Second * 10)
	// 	fmt.Printf("[%d]start read frames\n", i)
	// 	readFrames()
	// }

	readFrames()

	// resizer()
	fmt.Printf("time cost:%+v, loop:%d\n", time.Since(startTime), *loop)

	// select {}
}

func readFrames() {
	buffers := make([][]byte, 0)
	files, err := os.ReadDir(framesDir)
	check(err)
	i := 0
	for _, file := range files {
		switch file.Name() {
		case ".DS_Store", ".gitignore":
			continue
		}
		i++
		filePath := framesDir + file.Name()
		buf, err := os.ReadFile(filePath)
		check(err)
		buffers = append(buffers, buf)
		// fmt.Printf("read %s buf len:%d\n", filePath, len(buf))
		if *frames != 0 && i+1 >= *frames {
			fmt.Printf("frames limited:%d\n", *frames)
			break
		}
	}
	fmt.Printf("frames:%d\n", len(buffers))

	if *loop == 1 {
		a := &bimg.AnimatedGif{}
		buf, err := a.Join(buffers, *delayMs, 0)
		check(err)
		err = os.WriteFile("/Users/zhangyan/Documents/images/out.gif", buf, 0644)
		fmt.Printf("write out gif\n")
		check(err)
	} else {
		var wg sync.WaitGroup
		for i := 0; i < *loop; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				createGif(buffers)
				fmt.Printf("task %d done\n", n)
			}(i)
		}
		wg.Wait()
	}
}

func createGif(buffers [][]byte) {
	a := &bimg.AnimatedGif{}
	_, err := a.Join(buffers, *delayMs, 0)
	check(err)
}
