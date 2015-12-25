## css
- css规则=选择器(selector)+声明块(declaration block)
- css注释为`/*...*/`

### css分类:
- 内联样式(Inline Style),**不推荐**
- 内部样式表(Internal Style Sheet,`<style></style>,html5写法,可省略type`)
- 外部样式表(External Style Sheet,`<link rel="stylesheet" href="style.css" />,html5写法,可省略type`),**推荐**

### css层叠(覆盖)的原则
- **特殊性(specificity)**,指定选择器的具体程度。**选择器越特殊,规则就越强**。遇到冲突时,优先应用特殊性强的规则.
```css
#id > .class1.class2 > class1 > p
```
- **顺序(order)**,在特殊性还不足以判断相互冲突的规则中应该优先应用哪一个的情况下,顺序规则就可以起到决定作用:**晚出现的优先级高**.
- **重要性(importance)**,如果以上两条规则还不够,可以声明一条特殊的规则覆盖整个系统中的规则,这条规则的重要程度要比其他所有规则高,即在某条声明的末尾加上`!important`(**除非是在特殊情况下,否则不推荐使用这种方法**).

ps: 如果某元素没有直接指定某条规则,则使用继承的值(如果有的话)。

### css3的hsl
构想 HSL:选择一个 0 到 360 之间的色相,并将饱和度设为 100,亮度设为 50,就会得到这种颜色最纯的形式。降低饱和度,颜色就会向灰色变化。增加亮度,颜色就会向白色变化;减少亮度,颜色就会向黑色变化。

### 选择器
选择器可以定义五个不同的标准来选择要进行格式化的元素:
1. 元素的类型或名称
2. 元素所在的上下文
3. 元素的类或ID
4. 元素的伪类或伪元素
5. 元素是否有某些属性和值

为了指出目标元素,选择器可以使用这五个标准的任意组合。
ps:
- 如果要定位的元素有多个类名,可以在选择器中将它们连在一起,就像`.classA.classB { color: blue;}`.
- 在 class 选择器和 id 选择器之间作选择时,建议尽可能地使用 class,因为class易于复用.
- E:first-child/E:last-child选取的是E的父元素(最高为body)的第一个/最后的子元素,如果该元素是E类型则应用规则.
- 由于链接可能同时处于多种状态(如同时处于激活和鼠标停留状态),且晚出现的规则会覆盖前面出现的规则,所以,一定要按照下面的顺序定义规则:**link 、 visited 、focus 、 hover 、 active ( 缩写为LVFHA)**.
- `a[href][title~="howdy"]{...}`通过联合使用多种方法,这个选择器选择所有既有任意 href 属性,又有任意属性值包含单词 howdy 的 title 属性的 a 元素.
- 指定元素组,可以将将相同的样式规则应用于多个元素,比如`h1,h2{...}`.

>在 CSS3 中, :first-line的语法为 ::first-line , :first-letter的语法为 ::first-letter 。注意,它们用两个冒号代替了单个冒号。这样修改的目的是将伪元素(有四个,包括 ::first-line 、 ::first-letter 、 ::before 和::after )与伪类( 如 :first-child 、:link 、 :hover 等)区分开。推荐使用::first-line 和 ::first-letter 这样的双冒号语法.
>
>**伪元素(pseudo-element)**是 HTML 中并不存在的元素。例如,定义第一个字母或第一行文字时,并未在 HTML 中作相应的标记。它们是另一个元素(比如p元素)的部分内容。
相反,**伪类(pseudo-class)**则应用于一组 HTML 元素,而你无需在 HTML 代码中用类标记它们。例如,使用 :first-child 可以选择某元素的第一个子元素,你就不用写成class="first-child".

### 为文本添加css
background-attachment:控制背景图像是否随页面滚动
使用 background 简记法(推荐)：`background: #004 url(../img/ufo.png) no-repeat 170px 20px;`．
text-align适用于块元素．
white-space，设置如何处理元素内的空白，规定段落中的文本不进行换行．

### 用css进行布局
#### display
- 设置为inline的元素会忽略任何width、height、margin-top 和 margin-bottom,但 padding-top 和padding-bottom 会越界进入相邻元素的区域,而不是局限于该元素本身的空间.
- float并不会影响父元素或其他祖先元素的高度,因此从这一点来说,它不属于文档流的一部分了.
- clear 属性:对某个元素使用该属性,该元素和它后面的元素就会显示在浮动元素的下面.
- clearfix 和 overflow 方法是应用于浮动元素的父元素或祖先元素的.clearfix 或overflow 应用到浮动元素的任何一个非父元素的祖先元素,这样做不会让父元素变高,但祖先元素的高度会包含浮动元素.推荐clearfix,因为overflow: hidden会将内容截断,对此要多加注意。有时使用overflow: auto; 也有效,但这样做可能会出现一个滚动条.
- relative是相对于原始(文档流)的位置进行定位.其他元素不会受到偏移的影响,仍然按照这个元素原来的盒子进行排列。设置了相对定位的内容可能与其他内容重叠,这取决于 top 、 right 、bottom 和 left 的值.
- 使用相对定位、绝对定位或固定定位时,对于相互重叠的元素,可以用 z-index属性指定它们的叠放次序.
- static 没有定位，元素出现在正常的流中（忽略 top, bottom, left, right 或者 z-index 声明）

## ps
>### 会被继承的属性
#### 文本
- color (颜色, a 元素除外)，不常见的格式：rgba(r,g,b,a),hsla(h,s,l,a).
- direction (方向)
- font (字体),格式为`font-style font-variant font-weight font-size/line-height font-family`,font简写属性至少包括字体系列和字体大小属性
- font-family (字体系列),但表单的select,textarea,input除外.
- font-size (字体大小),px/百分数/相对大小rem等/xx-small 、 x-small 、 small 、 medium 、 large 、x-large 或 xx-large/smaller,larger
- font-style (用于设置斜体)
- font-variant (用于设置小型大写字母)
- font-weight (用于设置粗体),normal/lighter/bold/bolder/100~900之间的100的倍数,其中400代表正常粗细,700代表粗体。
- letter-spacing (字母间距)
- line-height (行高),段落的行距,即段落内每行之间的距离.使用数字设定行高时,其所有的子元素都会继承这个因子n,lh=当前元素的font-size*n;而使用百分数,相对大小时,子元素直接继承父元素的行高.
- text-align (用于设置对齐方式)
- text-indent (用于设置首行缩进)
- text-transform (用于修改大小写)
- visibility (可见性)
- white-space (用于指定如何处理空格)
- word-spacing (字间距)
#### 列表
- list-style (列表样式)
- list-style-image (用于为列表指定定制的标记)
- list-style-position (用于确定列表标记的位置)
- list-style-type (用于设置列表的标记)
#### 表格
- border-collapse (用于控制表格相邻单元格的边框是否合并为单一边框)
- border-spacing (用于指定表格边框之间的空隙大小)
- caption-side (用于设置表格标题的位置)
- empty-cells (用于设置是否显示表格中的空单元格)
#### 页面设置(对于印刷物)
- orphans (用于设置当元素内部发生分页时在页面底部需要保留的最少行数)
- page-break-inside (用于设置元素内部的分页方式)
- widows (用于设置当元素内部发生分页时在页面顶部需要保留的最少行数)
#### 其他
- cursor (鼠标指针)
- quotes (用于指定引号样式)
