package bimg

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"
import (
	"fmt"
	"runtime"
)

type AnimatedGif struct {
}

func (i *AnimatedGif) Join(buffers [][]byte, delayMs, loop int) ([]byte, error) {
	defer runtime.KeepAlive(buffers)
	defer C.vips_thread_shutdown()

	var pageHeight int
	images := make([]*C.VipsImage, len(buffers))
	fmt.Println("image length:", len(buffers))
	nPages := len(buffers)
	for i, buf := range buffers {
		image, _, err := loadImage(buf)
		if err != nil {
			return nil, err
		}
		//image, err = vipsSRGB(image)
		//if err != nil {
		//	return nil, err
		//}
		// 只能用下标赋值，append拷贝是nil
		images[i] = image
		// 获取最大高度作为gif画布的高度，所有图片从画布左上角对齐，右下补黑边
		height := int(image.Ysize)
		if height > pageHeight {
			pageHeight = height
		}
	}
	// todo: from image
	image, err := vipsAnimatedGifJoin(images, pageHeight, nPages, delayMs, loop)
	if err != nil {
		return nil, err
	}
	saveOptions := vipsSaveOptions{
		Type: GIF,
	}
	// Finally get the resultant buffer
	return vipsSave(image, saveOptions)
}

func ValidAnimatedGifDelay(delayMs int) bool {
	if delayMs >= 10 && delayMs <= 0xffff*10 {
		return true
	}
	return false
}

func ValidAnimatedGifLoop(loop int) bool {
	if loop >= 0 && loop <= 0xffff {
		return true
	}
	return false
}
