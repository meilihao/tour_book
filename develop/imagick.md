## 简介

参考:[ImageMagick简介、GraphicsMagick、命令行使用示例](http://elf8848.iteye.com/blog/382528),[PHP ImageMagick扩展API](http://php.net/manual/zh/class.imagick.php)

### ImageMagick

ImageMagick是一个免费的创建、编辑、合成图片的软件。它可以读取、转换、写入多种格式的图片。图片切割、颜色替换、各种效果的应用，图片的旋转、组合，文本，直线，多边形，椭圆，曲线，附加到图片伸展旋转。其全部源码开放，可以自由使用，复制，修改，发布。它遵守GPL许可协议。它可以运行于大多数的操作系统。
最为重要的是，ImageMagick的大多数功能的使用都来源于命令行工具。

### GraphicsMagick

GraphicsMagick是从 ImageMagick 5.5.2 分支出来的，但是现在他变得更稳定和优秀，GM更小更容易安装、GM更有效率、GM的手册非常丰富GraphicsMagick的命令与ImageMagick基本是一样的。

服务端推荐使用GraphicsMagick.

## imagick API 中文手册

```
imagick 类
imagick ::adaptiveblurimage 向图像中添加 adaptive 模糊滤镜
imagick ::adaptiveresizeimage 自适应调整图像数据依赖关系
imagick ::adaptivesharpenimage自适应锐化图像
imagick ::adaptivethresholdimage 基于范围的选择为每个像素的亮度
imagick ::addimage 图像列表中添加新图像 imagick 对象.
imagick ::addnoiseimage 给图像添加随机噪声
imagick ::affinetransformimage变换图像
imagick ::animateimages 动画图像或图像
imagick ::annotateimage annotates 图像的文本
imagick ::appendimages 追加一组图像
imagick ::averageimages 平均一组图像
imagick ::blackthresholdimage 强制所有的像素低于阈值分为黑色
imagick ::blurimage 向图像中添加模糊滤镜
imagick ::borderimage 四周带有边框的图像
imagick ::charcoalimage模拟一个炭笔绘图
imagick ::chopimage 删除图像和用到的区域
imagick ::clear 清除所有 imagick 对象相关联的资源
imagick ::clipImage 从 8bim profile clipimage 剪辑沿路径
imagick ::clipPathImage 沿 clippathimage 剪辑命名路径从 8bim profile
imagick ::clone 使 imagick 完全复制的对象
imagick ::clutimage 替换图像中的颜色
imagick ::coalesceimages 贡献一组图像
imagick ::colorfloodfillimage 更改任何像素的颜色值相匹配的目标
imagick ::colorizeimage 融合的填充颜色与图像
imagick ::combineimages 将一个或多个图像合并成一个图像
imagick ::commentimage 添加注释以您的图像
imagick ::compareimagechannels 返回 date 对象中一个或多个图像
imagick ::compareimagelayers 返回图像之间的最大亮度区域
imagick ::compareimages 角落比较图像图像
imagick ::compositeimage 复合图像拖到另一个
imagick ::__constructimagick actionscript
imagick ::contrastimage 改变图像的对比度
imagick ::contraststretchimage 增强对比度彩色图像
imagick ::convolveimage 向图像应用一个定制的卷积 kernel
imagick ::cropimage 中提取区域的图像
imagick ::cropthumbnailimage创建一个裁剪高级
imagick ::current 当前返回对当前 imagick 对象
imagick ::cyclecolormapimage鸟图像的色彩表
imagick ::decipherimagedeciphers 图像
imagick ::deconstructimages 返回某些像素的图像之间的差异
imagick ::deleteimageartifact删除图像工件
imagick ::deskewimage 从映像中删除波纹
imagick ::despeckleimage 减少图像中的斑点噪音
imagick ::destroy 销毁销毁此 imagick 对象
imagick ::displayimage显示图像
imagick ::displayimages 显示一个图像或图像序列.
imagick ::distortimage 扭曲图像中的各种变形方法
imagick ::stroke 对象的呈现上的 imagickdraw 对象当前图像
imagick ::edgeimage 增强边缘在 image
imagick ::embossimage 返回一个带有三个尺寸的灰度图像效果
imagick ::encipherimageenciphers 图像
imagick ::enhanceimage 可以提高一个情景图像的质量
imagick ::equalizeimageequalizes 图像直方图
imagick ::evaluateimage 应用图像的表达式
imagick ::exportimagepixels将原始图像像素
imagick ::extentimage设置图像大小
imagick ::flattenimages 合并图像序列的
imagick ::flipimage 创建垂直镜像
imagick ::floodfillpaintimage 更改任何像素的颜色值相匹配的目标
imagick ::flopimage 创建水平镜像
imagick ::frameimage 添加一个模拟特性的边框.
imagick ::functionimage 在图像上应用了一个函数
imagick ::fximage 评估表达式的每个像素在图像
imagick ::gammaimagegamma 校正图像
imagick ::gaussianblurimage调整
imagick ::getcolorspace获取或设置与 colorspace
imagick ::getcompression 压缩获取该对象.
imagick ::getcompressionquality 获取对象的压缩质量
imagick ::getcopyright 返回 imagemagick api 的版权为一个字符串.
imagick ::getfolder 所指定的文件名相关联的图像序列
imagick ::getfont获取字体
imagick ::getFormat 只返回 imagick 对象的格式.
imagick ::getgravity获取或设置 gravity
imagick ::gethomeurl 返回 imagemagick 首页 url
imagick ::getimage 返回一个新 imagick 对象
imagick ::getimagealphachannel 获取图像的 alpha 通道
imagick ::getimageartifact获取图像的工件
imagick ::getimagebackgroundcolor 返回图像的背景颜色.
imagick ::getimageblob 返回序列作为 blob
imagick ::getimageblueprimary 返回 chromaticy 蓝色主点
imagick ::getimagebordercolor 返回图像的边框颜色.
imagick ::getimagechanneldepth 获取为特定图像深度通道
imagick ::getimagechanneldistortion 角落比较图像图像通道的图像
imagick ::getimagechanneldistortions获取通道 distortions
imagick ::getImageChannelExtrema 输入 extrema getimagechannelextrema 获取图像的一个或多个通道
imagick ::getimagechannelkurtosisgetimagechannelkurtosis 目的
imagick ::getimagechannelmean 获取平均值和标准偏差
imagick ::getimagechannelrange获取通道区域
imagick ::getimagechannelstatistics 返回统计每个通道的图像
imagick ::getimageclipmask获取图像剪辑遮罩
imagick ::getimagecolormapcolor 返回指定色彩表索引的颜色.
imagick ::getimagecolors 获取图像中的唯一的颜色数
imagick ::getimagecolorspace获取或设置图像的 colorspace
imagick ::getimagecompose 返回复合运算符的图像关联的
imagick ::getimagecompression 获取当前图像的压缩类型
imagick ::getimagecompressionquality 获取当前图像的压缩质量
imagick ::getimagedelay获取或设置图像延迟
imagick ::getimagedepth获取或设置图像深度
imagick ::getimagedispose 获取图像的处理方法
imagick ::getimagedistortion 角落比较图像图像
imagick ::getimageextrema 获取图像的"extrema
imagick ::getimagefilename 返回序列中的特定图像的文件名
imagick ::getimageformat 返回序列中的特定图像的格式
imagick ::getimagegamma获取或设置图像的 gamma
imagick ::getimagegeometry 获取作为关联数组的宽度和高度
imagick ::getimagegravity获取或设置图像的态
imagick ::getimagegreenprimary 返回 chromaticy 绿色主点
imagick ::getimageheight返回图像的高度
imagick ::getimagehistogram获取或设置图像直方图
imagick ::getimageindex 获取当前活动图像的索引.
imagick ::getimageinterlacescheme 获取非方形像素图像方案
imagick ::getimageinterpolatemethod返回指定插值方法
imagick ::getimageiterations获取或设置图像的迭代
imagick ::getimagelength 返回图像的字节长度
imagick ::getimagemagicklicense 返回一个字符串,该字符串包含 imagemagick 的许可证
imagick ::getimagematte 如果图像有一个 matte channel getimagematte 返回
imagick ::getimagemattecolor 返回 image matte 颜色.
imagick ::getimageorientation获取或设置图像的方向
imagick ::getimagepage返回页面几何
imagick ::getimagepixelcolor 返回指定像素的颜色
imagick ::getimageprofile 返回命名图像
imagick ::getimageprofiles返回图像概要
imagick ::getimageproperties返回指定的图像属性
imagick ::getimageproperty 返回命名图像属性
imagick ::getimageredprimary 返回 chromaticity 红色主点
imagick ::getimageregion 中提取区域的图像
imagick ::getimagerenderingintent 获取图像渲染方法
imagick ::getimageresolution 获取图像的x和y分辨率
imagick ::getimagesblob 返回所有图片序列作为 blob
imagick ::getimagescene获取或设置图像的场景
imagick ::getimagesignature 生成一个 sha 256 message digest
imagick ::getimagesize ()返回的字节长度
imagick ::getimagetickspersecond获取或设置图像每秒 ticks
imagick ::getimagetotalinkdensity 获取图像的总油墨浓度
imagick ::getimagetype 获取潜在的图像类型
imagick ::getimageunits 获取图像的分辨率
imagick ::getimagevirtualpixelmethod 虚拟像素方法返回
imagick ::getimagewhitepoint 返回 chromaticity 白色点
imagick ::getimagewidth返回指定图像宽度
imagick ::getinterlacescheme 获取的对象的非方形像素方案
imagick ::getiteratorindex 获取当前活动图像的索引.
imagick ::getnumberimages 返回的对象中的图像
imagick ::getoption ()返回一个值与指定键相关联的
imagick ::getpackagename 返回 imagemagick 包名称
imagick ::getpage返回页面几何
imagick ::getpixeliterator返回一个 magickpixeliterator
imagick ::getPixelRegionIterator 获取图像的ImagickPixelIterator
imagick ::getpointsize获取的点大小
imagick ::getquantumdepth获取或设置量程深度
imagick ::getquantumrange 返回 imagick 量程范围
imagick ::getreleasedate 返回 imagemagick 的发行日期
imagick ::getResource 获得与指定返回指定的资源的使用状况.
imagick ::getresourcelimit 返回指定的资源限制
imagick ::getsamplingfactors 获取的水平和垂直采样系数
imagick ::getsize 返回 imagick 对象相关联的
imagick ::getsizeoffset返回的偏移
imagick ::getversion 返回 imagemagick api 版本
imagick ::haldclutimage 替换图像中的颜色
imagick ::hasnextimage 检查对象具有有关图像
imagick ::hasPreviousImage 如果对象有一个 haspreviousimage 检查前面的图像
imagick ::identifyImage 一个图像并取得 identifyimage 标识属性
imagick ::implodeImage 作为副本 implodeimage 新建一图像
imagick ::importimagepixels导入图像像素
imagick ::labelimage 向图像添加标签
imagick ::levelimage 调整图像的级别
imagick ::linearstretchimage 可对图像亮度与饱和度
imagick ::liquidrescaleimage 动画图像或图像
imagick ::magnifyImage 调整图像比例2X
imagick ::mapimage 替换图像和最相近的颜色从一个图像的颜色.
imagick ::mattefloodfillimage 更改颜色的透明度值
imagick ::medianfilterimage应用了一个数字滤波器
imagick ::mergeimagelayers合并图像的图层
imagick ::minifyimage 缩放图像以它的大小的一半
imagick ::modulateimage 亮度、饱和度和色调控制
imagick ::montageimage创建复合图像
imagick ::morphimages 方法 morphs 一组图像
imagick ::mosaicimages 股份从图像马赛克
imagick ::motionblurimage模拟运动模糊
imagick ::negateimage 如何处理主键参照图像中的颜色
imagick ::newimage新建一图像
imagick ::newpseudoimage新建一图像
imagick ::nextimage 移至上一图像
imagick ::normalizeimage 增强对比度彩色图像
imagick ::oilpaintimage模拟一个油画
imagick ::opaquepaintimage 更改任何像素的颜色值相匹配的目标
imagick ::optimizeimagelayers 删除重复的图像部分进行优化.
imagick ::orderedposterizeimage执行顺序抖动
imagick ::paintfloodfillimage 更改任何像素的颜色值相匹配的目标
imagick ::paintopaqueimage 更改颜色相匹配的任何像素
imagick ::painttransparentimage 用指定的颜色填充匹配的像素颜色
imagick ::pingImage 取得图像的基本属性
imagick ::pingimageblob快速提取属性
imagick ::pingimagefile 获得基本图像属性中的轻量的方式
imagick ::polaroidimage模拟 polaroid 图片
imagick ::posterizeimage 可以减少图像有限数目的颜色级别
imagick ::previewImages 快速 previewimages 针点相应参数的图像和视频处理
imagick ::previousimage 移到下一图像中的对象
imagick ::profileimage 从图像中添加或删除档案
imagick ::quantizeimage 分析参考一个图像中的颜色
imagick ::quantizeimages 分析一序列的图像中的颜色
imagick ::queryfontmetrics 返回数组表示字体度量
imagick ::queryfonts返回配置字体
imagick ::queryformats imagick 返回格式支持
imagick ::radialblurimage精致式调整
imagick ::raiseimage 创建一个模拟的三维按钮效果
imagick ::randomthresholdimage 创建一个高对比度,两色图像
imagick ::readimage从文件名中读取图像
imagick ::readimageblob 读取图像从一个二进制字符串
imagick ::readimagefile 读取图像从打开 filehandle
imagick ::recolorimagerecolors 图像
imagick ::reducenoiseimage smooths 的图像
imagick ::remapimagerect 图像颜色
imagick ::removeimage 删除从图像列表中的图像
imagick ::removeimageprofile 删除命名图像并将其返回
imagick ::render 渲染图形呈现所有前缀的命令
imagick ::resampleimage 重采样图像所需的分辨率
imagick ::resetimagepage复位图像页
imagick ::resizeimage缩放图像
imagick ::rollimage偏移图像
imagick ::rotateimage旋转图像
imagick ::roundcorners图像将角
imagick ::sampleimage 缩放图像的像素取样
imagick ::scaleimage 缩放图像的大小
imagick ::segmentimage段中的图像
imagick ::separateimagechannel 分开通道图像
imagick ::sepiatoneimage棕褐色色调发出声音图像
imagick ::setbackgroundcolor 设置对象的默认背景色
imagick ::setcolorspace设置 colorspace
imagick ::setcompression 设置对象的默认压缩类型
imagick ::setcompressionquality 设置对象的默认压缩质量
imagick ::setfilename 设置图像的文件名之前读取或写入
imagick ::setfirstiterator 设置 valid imagick 图像
imagick ::setfont设置字体
imagick ::setformat imagick 设置格式的对象.
imagick ::setgravity设置 gravity
imagick ::setimage 替换图像中的对象
imagick ::setimagealphachannel设置图像的 alpha 通道
imagick ::setimageartifact设置图像的工件
imagick ::setimagebackgroundcolor 图像设置背景颜色
imagick ::setimagebias 设置的任何方法 convolves 图像的图像偏置
imagick ::setimageblueprimary 设置图像 chromaticity 蓝色主点
imagick ::setimagebordercolor 设置图像的边框颜色
imagick ::setimagechanneldepth 设置特定图像的通道.
imagick ::setimageclipmask设置图像剪辑遮罩
imagick ::setimagecolormapcolor 设定指定的色彩表索引的颜色
imagick ::setimagecolorspace设置图像的 colorspace
imagick ::setimagecompose 设置图像的复合运算符
imagick ::setimagecompression设置图像压缩
imagick ::setimagecompressionquality 设置图像的压缩质量.
imagick ::setimagedelay设置图像延迟
imagick ::setimagedepth设置图像深度
imagick ::setimagedispose 设置图像的处理方法
imagick ::setimageextent设置图像尺寸
imagick ::setimagefilename 设置特定图像的文件名
imagick ::setimageformat 设置格式的特定的图像
imagick ::setimagegamma设置图像的 gamma
imagick ::setimagegravity设置图像的态
imagick ::setimagegreenprimary 设置图像 chromaticity 绿色主点
imagick ::setimageindex设置 iterator
imagick ::setimageinterlacescheme设置图像压缩
imagick ::setimageinterpolatemethod 设置图像插值像素的方法
imagick ::setimageiterations设置图像的迭代.
imagick ::setimagematte 设置 image matte channel
imagick ::setimagemattecolor 设置 image matte 颜色
imagick ::setimageopacity 设置图像=
imagick ::setimageorientation设置图像的方向
imagick ::setimagepage 设置图像的几何
imagick ::setimageprofile 添加一个名为 profile imagick 对象
imagick ::setimageproperty设置图像属性.
imagick ::setimageredprimary 设置图像 chromaticity 红色主点
imagick ::setimagerenderingintent 设置图像渲染方法.
imagick ::setimageresolution设置图像分辨率
imagick ::setimagescene设置图像的场景
imagick ::setimagetickspersecond设置图像每秒 ticks
imagick ::setimagetype设置图像类型
imagick ::setimageunits 设置图像的分辨率
imagick ::setimagevirtualpixelmethod 设置图像的像素的虚拟方法
imagick ::setimagewhitepoint 设置 chromaticity 的图像
imagick ::setinterlacescheme设置图像压缩
imagick ::setiteratorindex设置 iterator
imagick ::setlastiterator 设置 valid imagick 最后的图像
imagick ::setoption设置一个选项
imagick ::setpage 设置 imagick 的几何对象.
imagick ::setpointsize设置点的大小
imagick ::setresolution设置图像分辨率
imagick ::setresourcelimit 设置以 mb 为单位的限制为特定的资源.
imagick ::setsamplingfactors 设置图像采样的因素.
imagick ::setsize 设置 imagick 对象的大小
imagick ::setsizeoffset 设置相同的大小和偏移 imagick 对象
imagick ::settype 设置属性的图像类型.
imagick ::shadeimage创建三维效果
imagick ::shadowimage模拟图像的阴影
imagick ::sharpenimage增强图像
imagick ::shaveimage shaves 从图像边缘的像素
imagick ::shearimage创建一个平行四边形
imagick ::sigmoidalcontrastimage 调整图像的对比度
imagick ::sketchimage模拟一个素描
imagick ::solarizeimage 给图像应用一个 solarizing 效果
imagick ::sparsecolorimagepermut 颜色
imagick ::spliceimage my 纯色到图像
imagick ::spreadimage 随机鸟块中的每个像素
imagick ::steganoimage 隐藏数字水印在 image
imagick ::stereoimage组合两张图片
imagick ::stripimage 停车图像的所有概要信息和注释
imagick ::swirlimage 漩涡中心的像素的图像
imagick ::textureImage 反复平铺纹理.
imagick ::thresholdimage 改变单个像素的值基于一个阈值
imagick ::thumbnailimage 段更改图像的大小
imagick ::tintimage 颜色矢量图像中的每个像素
imagick ::transformimage 设置裁剪尺寸图像对象的便捷方法
imagick ::transparentpaintimage若要像素透明
imagick ::transposeimage 创建垂直镜像
imagick ::transverseimage 创建水平镜像
imagick ::trimImage 去除图像边缘
imagick ::uniqueimagecolors 丢弃任何像素的颜色
imagick ::unsharpmaskimage 增强图像
imagick ::valid 当前项是否有效
imagick ::vignetteimage 给图像添加效果筛选器
imagick ::waveimage 给图像应用波形的筛选
imagick ::whitethresholdimage 强制所有像素等于阈值为白色
imagick ::writeimage 将图像写入指定的 filename
imagick ::writeimagefile 图像写入到一个 filehandle
imagick ::writeimages 写入一个图像或图像序列.
imagick ::writeimagesfile 写入帧 filehandle
imagickdraw 类
imagickdraw ::affine 调整当前仿射转换矩阵
imagickdraw ::annotation 在图像上绘制文字
imagickdraw ::arc 绘制圆弧
imagickdraw ::bezier 绘制 bezier 曲线
imagickdraw ::circle 圆绘制一个圆
imagickdraw ::clear 清除 imagickdraw
imagickdraw ::clone 克隆使 imagickdraw 对象指定的精确副本.
imagickdraw ::color 绘制图像上的颜色
imagickdraw ::comment 添加注释
imagickdraw ::composite 复合图像的贡献当前图像
imagickdraw ::__constructimagickdraw actionscript
imagickdraw ::destroy 销毁释放关联的所有资源
imagickdraw ::ellipse 在图像上绘制一个椭圆
imagickdraw ::getClipPath 当前剪辑路径 getclippath 获取 id
imagickdraw ::getcliprule 返回当前多边形填充规则
imagickdraw ::getclipunits 返回集合剪辑路径为单位
imagickdraw ::getfillcolor返回指定填充的颜色
imagickdraw ::getfillopacity 返回绘制时同时使用
imagickdraw ::getfillrule返回指定填充规则
imagickdraw ::getfont返回指定的字体
imagickdraw ::getfontfamily返回指定的字体族
imagickdraw ::getfontsize返回指定的字体 pointsize
imagickdraw ::getfontstyle返回指定的字体样式
imagickdraw ::getfontweight返回该字体的 weight
imagickdraw ::getgravity 返回的位置 gravity
imagickdraw ::getstrokeantialias 返回当前笔触对设置.
imagickdraw ::getstrokecolor 返回的勾画轮廓对象所使用的颜色
imagickdraw ::getstrokedasharray 返回数组表示虚线和间隙的图案用于勾画路径
imagickdraw ::getstrokedashoffset 返回偏移到虚线图案以划线开始
imagickdraw ::getstrokelinecap 描画后会返回用于打开网络访问可远程访问的注册表路径末尾的形状
imagickdraw ::getstrokelinejoin 描画后会返回要使用的角的路径的形状
imagickdraw ::getStrokeMiterLimit 笔划 getstrokemiterlimit 返回斜角限制
imagickdraw ::getstrokeopacity 返回在描画对象的轮廓
imagickdraw ::getstrokewidth 返回笔触的宽度用于绘制对象的轮廓
imagickdraw ::gettextalignment返回文本对齐
imagickdraw ::gettextantialias 返回当前文字对设置.
imagickdraw ::gettextdecoration返回文本修饰
imagickdraw ::gettextencoding 返回代码设定文本的注释
imagickdraw ::gettextundercolor 返回的文本颜色.
imagickdraw ::getvectorgraphics 返回一个字符串,该字符串包含矢量图形
imagickdraw ::line 行画一条线.
imagickdraw ::matte 若要在图像的 alpha 通道
imagickdraw ::pathclose 路径添加一个元素到当前路径
imagickdraw ::pathcurvetoabsolute 立方绘制 bezier 曲线
imagickdraw ::pathcurvetoquadraticbezierabsolute 绘制二次 bezier 曲线
imagickdraw ::pathcurvetoquadraticbezierrelative 绘制二次 bezier 曲线
imagickdraw ::pathcurvetoquadraticbeziersmoothabsolute 绘制二次 bezier 曲线
imagickdraw ::pathcurvetoquadraticbeziersmoothrelative 绘制二次 bezier 曲线
imagickdraw ::pathcurvetorelative 立方绘制 bezier 曲线
imagickdraw ::pathcurvetosmoothabsolute 立方绘制 bezier 曲线
imagickdraw ::pathcurvetosmoothrelative 立方绘制 bezier 曲线
imagickdraw ::pathellipticarcabsolute绘制一个椭圆形的弧
imagickdraw ::pathellipticarcrelative绘制一个椭圆形的弧
imagickdraw ::pathfinish终止当前路径
imagickdraw ::pathlinetoabsolute绘制一条直线路径
imagickdraw ::pathlinetohorizontalabsolute 水平直线绘制路径
imagickdraw ::pathlinetohorizontalrelative绘制一条水平线.
imagickdraw ::pathlinetorelative绘制一条直线路径
imagickdraw ::pathlinetoverticalabsolute绘制一个垂直线条
imagickdraw ::pathlinetoverticalrelative 绘制一个垂直线条路径
imagickdraw ::pathmovetoabsolute开始一个新子路径
imagickdraw ::pathmovetorelative 开始一个新子路径
imagickdraw ::pathStart 声明的路径绘制列表
imagickdraw ::point 绘制一个点
imagickdraw ::polygon 多边形画一个多边形
imagickdraw ::polyline 折线绘制折线
imagickdraw ::pop 的代替当前 imagickdraw 堆栈中,并返回到以前推入 imagickdraw
imagickdraw ::popclippath 终止一个剪辑路径定义
imagickdraw ::popdefs终止一个定义列表中的
imagickdraw ::poppattern终止图案定义
imagickdraw ::push 克隆当前ImagickDraw并将它推到堆栈
imagickdraw ::pushclippath 开始一个剪辑路径定义
imagickdraw ::pushdefs 指示以下命令创建名为元素用于后期处理的
imagickdraw ::pushpattern 表示后续命令达 imagickdraw ::oppattern()命令包含一个命名的定义图案
imagickdraw ::rectangle 绘制一个矩形
imagickdraw ::render 渲染的呈现所有前面的绘图命令拖到图像
imagickdraw ::rotate 应用指定的旋转以当前坐标空间
imagickdraw ::roundrectangle绘制一个圆角矩形
imagickdraw ::scale 比例调整比例系数
imagickdraw ::setclippath 将一个名为剪裁路径与映像关联
imagickdraw ::setcliprule 剪裁多边形填充规则使用的路径
imagickdraw ::setclipunits 设定的解释剪辑路径
imagickdraw ::setfillalpha 设置 alpha 使用填充颜色或填充图形时使用的纹理
imagickdraw ::setfillcolor 设置用于绘制填充对象的填充颜色
imagickdraw ::setfillopacity 设置 alpha 使用填充颜色或填充图形时使用的纹理
imagickdraw ::setfillpatternurl 设置 url 用作填充图案填充对象
imagickdraw ::setfillrule 设置填充规则使用在绘制多边形的步骤
imagickdraw ::setFont 设置字体
imagickdraw ::setfontfamily 设置时要使用的字体族使用文字注释
imagickdraw ::setfontsize 设置时要使用的字体 pointsize 使用文字注释
imagickdraw ::setfontstretch 设置时要使用的字体拉伸使用文字注释
imagickdraw ::setfontstyle 设置时要使用的字体样式使用文字注释
imagickdraw ::setfontweight设置该字体 weight
imagickdraw ::setgravity 设置的位置 gravity
imagickdraw ::setstrokealpha 指定要描画对象的轮廓
imagickdraw ::setstrokeantialias 控制是否较高质量描画轮廓
imagickdraw ::setstrokecolor 设置用于勾画对象的颜色了
imagickdraw ::setstrokedasharray 指定虚线和间隙用于勾画路径的图案
imagickdraw ::setstrokedashoffset 划线指定偏移到虚线图案来启动
imagickdraw ::setstrokelinecap 描画后会打开网络访问可远程访问的注册表路径的末尾时,用于指定的形状.
imagickdraw ::setstrokelinejoin 指定角的路径描画时要使用的形状.
imagickdraw ::setstrokemiterlimit指定斜角限制
imagickdraw ::setstrokeopacity 指定要描画对象的轮廓
imagickdraw ::setstrokepatternurl 设置用于勾画对象填充图案
imagickdraw ::setstrokewidth 笔触的宽度不用于绘制对象的轮廓
imagickdraw ::settextalignment指定文字对齐
imagickdraw ::settextantialias 较高质量控制文本是
imagickdraw ::settextdecoration指定修饰
imagickdraw ::settextencoding 指定指定的代码设置
imagickdraw ::settextundercolor 指定矩形的背景颜色
imagickdraw ::setvectorgraphics矢量图形设置
imagickdraw ::setviewbox 整个画布大小设置
imagickdraw ::skewx 该参数设置为"right"水平方向的当前坐标系的
imagickdraw ::skewy 该参数设置为"right"当前在垂直方向坐标系统
imagickdraw ::translate 应用一个翻译翻译到当前坐标系
imagickpixel 类
imagickpixel ::clear 清除与该对象关联的资源
imagickpixel ::__constructimagickpixel actionscript
imagickpixel ::destroy deallocates 资源与该对象关联的
imagickpixel ::getColor 代码中返回指定颜色
imagickpixel ::getcolorasstring 作为字符串返回颜色.
imagickpixel ::getcolorcount 返回此颜色的颜色相关联的计数
imagickpixel ::getcolorvalue 获取颜色通道提供的标准化值.
imagickpixel ::getHSL 标准化 hsl 颜色的 imagickpixel gethsl 返回对象
imagickpixel ::issimilar 检查这个颜色之间的距离和另一
imagickpixel ::setcolor设置颜色.
imagickpixel ::setcolorvalue 标准化的值至少设置通道
imagickpixel ::sethsl 标准化 hsl 颜色设置
imagickpixeliterator 类
imagickpixeliterator ::clear 清除清除 pixeliterator 关联的资源
imagickpixeliterator ::__constructimagickpixeliterator actionscript
imagickpixeliterator ::destroy 销毁 deallocates pixeliterator 关联的资源
imagickpixeliterator ::getcurrentiteratorrow 返回 imagickpixel 对象的当前行.
imagickpixeliterator ::iterator getiteratorrow 返回当前像素的行
imagickpixeliterator ::getnextiteratorrow 返回行的像素的 iterator
imagickpixeliterator ::getpreviousiteratorrow返回请求的行
imagickpixeliterator ::newpixeliterator 返回一个新像素 iterator
imagickpixeliterator ::newpixelregioniterator 返回一个新像素 iterator
imagickpixeliterator ::resetiterator将当前像素 iterator
imagickpixeliterator ::setIteratorFirstRow 设置第一行的像素迭代器
imagickpixeliterator ::setIteratorLastRow 设置最后一行的像素迭代器
imagickpixeliterator ::setiteratorrow 设置像素的 iterator
imagickpixeliterator ::synciteratorsyncs 同步的像素迭代器
```

常用:

- cropimage : 提取/截取区域的图像
- readimage : 载入图像
- writeimage : 保存图像到指定位置

## 开发包

### imagick

[gographics/imagick](https://github.com/gographics/imagick) is a Go bind to ImageMagick's MagickWand C API.

Ubuntu安装见官方README,下面是Fedora的安装方式:

    yum install ImageMagick-devel

32bit OS安装碰到错误:"github.com/gographics/imagick/imagick/draw_info.go:17:19: unexpected: 12-byte float type - long double",待官方解决(issues#19).