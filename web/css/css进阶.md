# css

## [border-image](https://www.w3.org/TR/css3-background/#border-image-slice)

### border-image-slice

```css
border-image-slice: [<number> | <percentage>]{1,4} && fill?
```

指定顶部，右，底部，左边缘的图像向内偏移的距离，将图片分为九个区域(像九宫格),除非追加`fill`,否则图像中间部分将被丢弃（完全透明的处理).

- number : 数字表示图像的像素（位图图像）或向量的坐标（如果图像是矢量图像）
- % : 百分比图像的大小是相对的：水平偏移图像的宽度，垂直偏移图像的高度
- fill :	保留图像的中间部分

> number|% 用法与padding,margin或border-width类似.

### border-image-width

等同border-width,width和图片大小不一致时会缩放图片.

### boder-image-repeat

```css
border-image-repeat: [stretch|repeat|round]{1,2}
```

接受两个(或一个)值,第一个值表示水平方向的排列方式,第二个值是垂直方向的排列方式.

- stretch : 默认值,拉伸图像来填充区域
- repeat : 平铺,图像来填充区域
- round	: 类似repeat值,如果**无法完整平铺所有图像，则对图像进行缩放以适应区域**

## [border-radius](https://www.w3.org/TR/css3-background/#border-radius)

border-radius其实生成的是椭圆.

```css
border-radius	[ <length> | <percentage> ]{1,4} [ / [ <length> | <percentage> ]{1,4} ]? // 如果`/`存在,前值是圆角的水平方向半径,后值是垂直方向的半径,如果没有即水平半径=垂直半径
border-top-left-radius, border-top-right-radius, border-bottom-right-radius, border-bottom-left-radius [ <length> | <percentage> ]{1,2}
```

取值:
- 四个值: 第一个值为左上角，第二个值为右上角，第三个值为右下角，第四个值为左下角。
- 三个值: 第一个值为左上角, 第二个值为右上角和左下角，第三个值为右下角
- 两个值: 第一个值为左上角与右下角，第二个值为右上角与左下角
- 一个值： 四个圆角值相同

注意点:
1. 当border-radius半径值小于或等于border的厚度时,元素边框内部就不出现圆角效果,即内边框的圆角半径=border-radius的半径-对应边框的宽度
2. 相邻边有不同宽度时,元素内角会从宽的边平滑多度到窄的一边,其中一边甚至可以是0,相邻转角是由大向小转
3. 当元素相邻俩边的颜色和线条样式不同时,其转变的中心点在一个和两边宽度成正比的角上

## border-shadow

```css
border-shadow : none | <shadow> [ , <shadow> ]* //   <shadow> = inset? && <length>{2,4} && <color>? ,越靠前优先级越高即允许多重阴影
```

值:
- inset : 阴影类型,inset表示内阴影(即阴影在元素内)
- 1st <length> : 水平阴影的偏移
- 2nd <length> : 垂直阴影的偏移
- 3rd <length> : 阴影模糊半径,只能是正值，如果为0表示阴影不具有模糊效果，其值越大阴影的边缘就越模糊,[原理](http://www.ruanyifeng.com/blog/2012/11/gaussian_blur.html)
- 4th <length> : 阴影的扩展半径,在阴影的基础上向外扩展,类似于content-box的padding

> 层次关系: 边框>内阴影>背景图片>背景色>外阴影

## background

### background-origin

指定background-position属性的参考原点

```css
background-origin: padding-box|border-box|content-box
```

取值:
- border-box	从border的外边缘开始显示背景图
- padding-box	默认,从padding的外边缘(即border的内边缘)开始显示图片
- content-box	从content的外边缘(即padding的内边缘)开始显示

### background-clip

```css
background-clip: border-box|padding-box|content-box
```

裁剪元素的背景（背景图片或颜色)

取值:
- border-box	默认,元素背景图从元素的border区域向外裁剪,即border之外的背景会被裁掉
- padding-box	padding之外的背景会被裁掉
- content-box	content之外的背景会被裁掉

### background-size

```css
background-size: auto||<length> || <%> ||cover ||contain
```

取值:
- auto : 保持背景图的原始高度和宽度
- length : 指定背景图的大小
- % : 根据元素长宽的百分比来指定背景图的大小
- cover : 将图片按比例缩放来铺满整个容器,会失真
- contain : 将图片按比例缩放到宽度或高度正好适应所定义背景容器的区域

### background

```css
// 支持多背景,图片叠放优先级:从左往右.
background : [ <bg-layer> , ]* <final-bg-layer>
// <bg-layer> = <bg-image> || <position> [ / <bg-size> ]? || <repeat-style> || <attachment> || <box> || <box>
// <final-bg-layer> = <bg-image> || <position> [ / <bg-size> ]? || <repeat-style> || <attachment> || <box> || <box> || <'background-color'>
```
