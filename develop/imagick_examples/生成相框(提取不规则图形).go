package main

import (
	"log"

	"github.com/gographics/imagick/imagick"
)
//根据线条素材生成相框
func main() {
	//素材
	texturePath := "top.png"

	line := GenerateLineBox(texturePath, 900, 600, 80)
	defer line.Destroy()

	line.WriteImage("line.png")
}

func GenerateLineBox(texturePath string, width, height, border int) *imagick.MagickWand {
	mw := imagick.NewMagickWand()
	defer func() {
		if err := recover(); err != nil {
			mw.Destroy()
			mw = nil
			panic(err)
		}
	}()

	//生成底图
	transparent := imagick.NewPixelWand()
	defer transparent.Destroy()
	transparent.SetColor("transparent")
	log.Println(mw.NewImage(uint(width), uint(height), transparent))

	{
		texture := imagick.NewMagickWand()
		defer texture.Destroy()
		log.Println(texture.ReadImage(texturePath))
		//将素材缩放到实际所需尺寸大小
		log.Println(texture.ScaleImage(texture.GetImageWidth()*uint(border)/texture.GetImageHeight(), uint(border)))

		paint := imagick.NewMagickWand()
		defer paint.Destroy()
		max := width
		if width < height {
			max = height
		}
		//生成所需最长的线条
		log.Println(paint.NewImage(uint(max), uint(border), transparent))
		//反复平铺素材
		paint = paint.TextureImage(texture)
		defer paint.Destroy()

		//将素材放在底图的顶部,其他多余部分会被隐藏
		log.Println(mw.CompositeImage(paint, imagick.COMPOSITE_OP_SRC_OVER, 0, 0))
		//旋转素材180度
		log.Println(paint.RotateImage(transparent, 180))
		//将素材放在底图的底部
		log.Println(mw.CompositeImage(paint, imagick.COMPOSITE_OP_SRC_OVER, 0, height-border))
		{
			//生成left线框
			log.Println(paint.RotateImage(transparent, 90))
			if height < max {
				//裁剪到合适尺寸,否则下面旋转后生成的右边线条会出现位置错落
				log.Println(paint.CropImage(uint(border), uint(height), 0, 0))
				//?
				//log.Println(paint.ResetImagePage(""))
			}

			mask := imagick.NewMagickWand()
			defer mask.Destroy()
			log.Println(mask.NewImage(uint(border), uint(height), transparent))

			dw := imagick.NewDrawingWand()
			defer dw.Destroy()

			black := imagick.NewPixelWand()
			defer black.Destroy()
			black.SetColor("black")
			dw.SetFillColor(black)

			//绘制一个梯形,left_top->left_bottom->right_bottom->right_top
			dw.PathStart()
			dw.PathMoveToAbsolute(0, 0)
			dw.PathMoveToAbsolute(0, float64(height))
			dw.PathMoveToAbsolute(float64(border), float64(height-border))
			dw.PathMoveToAbsolute(float64(border), float64(border))
			dw.PathClose()
			dw.PathFinish()

			//构建蒙版
			log.Println(mask.DrawImage(dw))
			//用蒙版提取出所需的图样(去右边上下角的两个三角形)
			log.Println(paint.CompositeImage(mask, imagick.COMPOSITE_OP_COPY_OPACITY, 0, 0))
		}

		log.Println(mw.CompositeImage(paint, imagick.COMPOSITE_OP_SRC_OVER, 0, 0))
		log.Println(paint.RotateImage(transparent, 180))
		log.Println(mw.CompositeImage(paint, imagick.COMPOSITE_OP_SRC_OVER, width-border, 0))
	}

	return mw
}
/*提取多边形的其他方法
		// --- matte
		pwBlack := imagick.NewPixelWand()
		defer pwBlack.Destroy()
		pwBlack.SetColor("black")
		pwWhite := pwBlack.Clone()
		defer pwWhite.Destroy()
		pwWhite.SetColor("white")

		matte := imagick.NewMagickWand()
		defer matte.Destroy()
		log.ErrPanic(matte.NewImage(uint(width), uint(border), pwBlack))

		dw := imagick.NewDrawingWand()
		defer dw.Destroy()
                //反锯齿
		dw.SetStrokeAntialias(true)
		dw.SetFillColor(pwWhite)

		//needed polygon
		points := []imagick.PointInfo{{0, 0}, {float64(width), 0}, {float64(width - border), float64(border)}, {float64(border), float64(border)}}
		dw.Polygon(points)

		log.ErrPanic(matte.DrawImage(dw))
		log.ErrPanic(matte.SetImageMatte(false))
                //topLine为原图
		log.ErrPanic(topLine.CompositeImage(matte, imagick.COMPOSITE_OP_COPY_OPACITY, 0, 0))
*/
