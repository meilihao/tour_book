## 定位(position)

### static
默认值,即没有定位，元素出现在正常的流中（忽略 top, bottom, left, right, z-index）.

### relative
生成相对定位的元素，相对于其正常位置(即static时的位置)进行定位.

### fixed
生成绝对定位的元素，相对于浏览器窗口进行定位,这意味着即便页面滚动，它还是会停留在相同的位置,且浏览器不会保留它原本在页面应有的空间.
元素的位置通过 "left", "top", "right" 以及 "bottom" 属性进行规定.

### absolute
生成绝对定位的元素，相对于static定位以外的第一个父元素进行定位(如果绝对定位元素没有static定位的祖先元素，那么它是相对于浏览器的窗口区域(html元素)进行定位，并且它会随着页面滚动而移动).
元素的位置通过 "left", "top", "right" 以及 "bottom" 属性进行规定.

### ps
1. 使用absoulte或fixed定位的话，必须指定 left、right、 top、 bottom 属性中的至少一个，否则left/right/top/bottom属性会使用它们的默认值 auto ，这将导致对象遵从正常的HTML布局规则，在前一个对象之后立即被呈递，简单讲就是都变成relative，会占用文档空间，**这点非常重要，很多人使用absolute定位后发现没有脱离文档流就是这个原因，这里要特别注意**.

1. 偏移规则：
 - 如果top和bottom一同存在的话，那么只有**top**生效。
 - 如果left和right一同存在的话，那么只有**left**生效。

1. 占用空间
 - 对象遵循正常文档流 : static,relative
 - 对象脱离正常文档流 : fixed,absolute

## 盒模型(box model)

### box-sizing
- border-box,指定宽度和高度确定元素边框box(即content+padding+border),是web design的趋势,**推荐**.
- content-box,指定元素的宽度和高度适用于box的宽度和高度(即content).
