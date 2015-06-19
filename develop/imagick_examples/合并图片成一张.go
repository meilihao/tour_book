package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/gographics/imagick/imagick"
)

func main() {
	oldBasePath := "/home/chen/Wallpaper"
	mode := []string{"1.jpg", "2.jpg"}
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	p := make([]string, 0)
	p = append(p, "Name")

	y := 0

	for _, v := range mode {
		mw.ReadImage(path.Join(oldBasePath, v))
		tW, tH := int(mw.GetImageWidth()), int(mw.GetImageHeight())
		p = append(p, v, fmt.Sprintf("%d,%d,%d,%d", 0, y, tW, y+tH))
		y += tH
	}

	mw.ResetIterator()
	n := mw.AppendImages(true)
	n.SetImageFormat("jpg")

	//也可将分割信息写入图片的元数据,n.SetImageProperty()
	n.WriteImage(oldBasePath + "/" + strings.Join(p, "-") + ".jpg")
}
