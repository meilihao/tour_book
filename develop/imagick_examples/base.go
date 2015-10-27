package main

import (
	"fmt"

	"github.com/gographics/imagick/imagick"
)

func main() {
	//PrintExifInfo()
	//Base()
	//createUserCard()
	BaseOther()
}

//获取图片属性
func PrintExifInfo() {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	fmt.Println(mw.ReadImage("https://ss1.bdstatic.com/kvoZeXSm1A5BphGlnYG/skin_zoom/53.jpg?2"))

	//图片的所有属性
	allProperties := mw.GetImageProperties("*")
	fmt.Println(allProperties)
	//图片的exif信息
	exifProperties := mw.GetImageProperties("exif:*")
	fmt.Println(exifProperties)
	//图片信息概述
	profile := mw.GetImageProfiles("*")
	fmt.Println(profile)
	//fmt.Println(mw.GetImageProfile("exif"))

	exifMap := make(map[string]string)
	//相机相机产商
	exifMap["exif:Make"] = mw.GetImageProperty("exif:Make")
	//相机型号
	exifMap["exif:Model"] = mw.GetImageProperty("exif:Model")
	//拍摄时间
	exifMap["exif:DateTimeOriginal"] = mw.GetImageProperty("exif:DateTimeOriginal")
	//曝光方式
	exifMap["exif:ExposureProgram"] = mw.GetImageProperty("exif:ExposureProgram")
	switch mw.GetImageProperty("exif:ExposureProgram") {
	case "1":
		fmt.Println("手动(M)")
	case "2":
		fmt.Println("程序自动(P)")
	case "3":
		fmt.Println("光圈优先(A,Av)")
	case "4":
		fmt.Println("快门优先(S,Tv)")
	case "5":
		fmt.Println("艺术程序(景深优先)")
	case "6":
		fmt.Println("运动模式")
	case "7":
		fmt.Println("肖像模式")
	case "8":
		fmt.Println("风景模式")
	default:
		fmt.Println("Not defined")
	}
	//光圈系数(F值)
	exifMap["exif:FNumber"] = mw.GetImageProperty("exif:FNumber")
	//曝光时间	即快门速度
	exifMap["exif:ExposureTime"] = mw.GetImageProperty("exif:ExposureTime")
	//曝光补偿
	exifMap["exif:ExposureBiasValue"] = mw.GetImageProperty("exif:ExposureBiasValue")
	//ISO感光度
	exifMap["exif:ISOSpeedRatings"] = mw.GetImageProperty("exif:ISOSpeedRatings")
	//测光模式
	exifMap["exif:MeteringMode"] = mw.GetImageProperty("exif:MeteringMode")
	switch mw.GetImageProperty("exif:MeteringMode") {
	case "0":
		fmt.Println("Unknown")
	case "1":
		fmt.Println("平均测光")
	case "2":
		fmt.Println("中央重点测光")
	case "3":
		fmt.Println("点测光")
	case "4":
		fmt.Println("多点测光")
	case "5":
		fmt.Println("评价测光")
	case "6":
		fmt.Println("局部测光")
	case "255":
		fmt.Println("other")
	}
	//闪光灯状态
	exifMap["exif:Flash"] = mw.GetImageProperty("exif:Flash")
	//镜头焦距
	exifMap["exif:FocalLength"] = mw.GetImageProperty("exif:FocalLength")
	//图像宽度
	exifMap["exif:ExifImageWidth"] = mw.GetImageProperty("exif:ExifImageWidth")
	//软件信息(固件Firmware版本或编辑软件)
	exifMap["exif:Software"] = mw.GetImageProperty("exif:Software")
	fmt.Println(exifMap)
}

func Base() {
	img := imagick.NewMagickWand()
	defer img.Destroy()
	fmt.Println(img.ReadImage("https://ss1.bdstatic.com/kvoZeXSm1A5BphGlnYG/skin_zoom/53.jpg?2"))
	//图片格式转换
	//img.SetImageFormat("png")
	//压缩质量
	//img.SetCompressionQuality(100)
	//去噪点,增益图片质量
	//img.EnhanceImage()
	//获取图片宽度
	fmt.Println(img.GetImageWidth())
	//缩放图片
	img.ScaleImage(800, 800*img.GetImageHeight()/img.GetImageWidth())
	//模拟3D按钮般的效果,和FrameImage类似
	img.RaiseImage(8, 8, 100, 50, true)
	//---绘制图形
	pwBlack := imagick.NewPixelWand()
	defer pwBlack.Destroy()
	pwBlack.SetColor("black")

	drawBack := imagick.NewDrawingWand()
	defer drawBack.Destroy()
	//填充时颜色
	drawBack.SetFillColor(pwBlack)
	//填充时的不透明度
	drawBack.SetFillOpacity(0.6)
	//绘制矩形
	drawBack.Rectangle(100, 100, 400, 400)
	img.DrawImage(drawBack)

	//---绘制文字
	drawWord := imagick.NewDrawingWand()
	//文本字体
	drawWord.SetFont("WenQuanYi Micro Hei")
	//字体大小
	drawWord.SetFontSize(16)
	//文字对齐方式."LEFT", 1;"CENTER", 2;"RIGHT", 3.
	drawWord.SetTextAlignment(1)
	//文字颜色
	pwBlack.SetColor("white")
	drawWord.SetFillColor(pwBlack)
	//AnnotateImage(dw,x,y,文本旋转的角度,文本内容)
	img.AnnotateImage(drawWord, 100, 100, 0, "hello")
	img.AnnotateImage(drawWord, 100, 150, 30, "world")
	img.WriteImage("/home/chen/test.png")
}

//复合图片
func createUserCard() {
	// 新建一个空白图片用来做画布
	canvas := imagick.NewMagickWand()
	defer canvas.Destroy()

	pwWhite := imagick.NewPixelWand()
	defer pwWhite.Destroy()
	pwWhite.SetColor("white")

	canvas.NewImage(588, 684, pwWhite)
	canvas.SetImageFormat("jpg")
	//--- 头像
	face := imagick.NewMagickWand()
	defer face.Destroy()
	face.ReadImage("http://c.hiphotos.baidu.com/image/pic/item/d01373f082025aaf56f1b11affedab64024f1a24.jpg")
	face.CropImage(450, 450, 100, 100)
	face.ThumbnailImage(200, 200)

	//--- 二维码
	qr := imagick.NewMagickWand()
	defer qr.Destroy()
	qr.ReadImage("http://a.hiphotos.baidu.com/baike/w%3D268/sign=b0ccc8a2239759ee4a5067cd8afa434e/2934349b033b5bb571dc8c5133d3d539b600bc12.jpg")
	qr.ThumbnailImage(256, 256)

	//--- 背景
	backgroud := imagick.NewMagickWand()
	defer backgroud.Destroy()
	backgroud.ReadImage("https://camo.githubusercontent.com/e12576de22dd40723f7cdbec3f6012af536fec4c/687474703a2f2f6d656e676b616e672e6e65742f75706c6f61642f696d6167652f323031352f303631372f32303135303631373132313231305f37393032312e706e67")

	//将图片合并到画布
	canvas.CompositeImage(face, imagick.COMPOSITE_OP_OVER, 194, 0)
	canvas.CompositeImage(qr, imagick.COMPOSITE_OP_OVER, (588-256)/2, 684-256-31)
	canvas.CompositeImage(backgroud, imagick.COMPOSITE_OP_OVER, 0, 0)

	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("black")

	draw := imagick.NewDrawingWand()
	defer draw.Destroy()
	//选择正确的字体(即支持中文的字体)防止中文乱码
	draw.SetFont("/home/chen/.local/share/fonts/msyh.ttf")
	draw.SetFontSize(20)
	draw.SetFillColor(pw)
	draw.SetTextAlignment(2)
	canvas.AnnotateImage(draw, 588/2, 230, 0, "hao")

	draw2 := draw.Clone()
	defer draw2.Destroy()
	//含换行符
	canvas.AnnotateImage(draw, 588/2, 310, 0, "二维条码/二维码（2-dimensional bar code）\n是用某种特定的几何图形按一定规律在平面（二维方向上）分布的黑白相间的图形记录数据符号信息的")

	canvas.WriteImage("/home/chen/canvas")
}

//
func BaseOther() {
	p := imagick.NewPixelWand()
	defer p.Destroy()
	p.SetColor("white")

	pwo := imagick.NewPixelWand()
	defer pwo.Destroy()
	pwo.SetColor("red")

	im := imagick.NewMagickWand()
	defer im.Destroy()
	im.ReadImage("https://ss1.bdstatic.com/kvoZeXSm1A5BphGlnYG/skin_zoom/53.jpg?2")
	//偏移图像
	//im.RollImage(20, 39)
	//生成缩略图
	//im.ThumbnailImage(100, 100)/*-4阿
	//给图像添加噪声
	//im.AddNoiseImageChannel(imagick.CHANNEL_OPACITY, imagick.NOISE_POISSON)
	//图像模糊处理
	//im.BlurImage(5, 3)
	//添加边框,(左右宽度,上下宽度)
	//im.BorderImage(p, 200, 300)
	//模拟炭笔绘图,radius ：越小越薄;sigma： 越大 墨越深 反之。
	//im.CharcoalImage(0.0001, 0.001)
	//反色(胶卷底片效果)
	//im.NegateImage(false)
	//?
	//im.ColorizeImage(pwo, p)
	//返回带三维效果的灰度图片
	//im.EmbossImage(1, 1)
	//垂直翻转
	//im.FlipImage()
	//水平翻转
	//im.FlopImage()
	//给图片模拟一个三维边框(颜色,水平方向边框宽度,垂直方向边框宽度,内斜面宽度,外斜面宽度)
	//im.FrameImage(pwo, 50, 50, 10, 20)
	//伽玛校正图像
	//im.GammaImage(2)
	//高斯模糊(高斯模糊的半径,像素,不包括中心象素;高斯的标准偏差,以像素为单位)
	//im.GaussianBlurImage(30, 3)
	//运动模糊.float angle:模糊角度
	//im.MotionBlurImage(61, 10, 10)
	//径向模糊
	//im.RadialBlurImage(5)
	//调整图像的色阶
	//im.LevelImageChannel(imagick.CHANNEL_GREEN, 200, 200, 200)
	//简便的图像等比例放大2倍
	//im.MagnifyImage()
	//简便的图像等比例缩小到一半
	//im.MinifyImage()
	//控制调整图像的 亮度、饱和度、色调
	//im.ModulateImage(100, 1, 100)
	//模拟油画滤镜
	//im.OilPaintImage(10)
	im.WriteImage("/home/chen/t.jpg")
}

