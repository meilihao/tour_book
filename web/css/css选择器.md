# CSS选择器汇总及其CSS版本

[原文](http://www.runoob.com/cssref/css-selectors.html)

# CSS3选择器分类

- 基本选择器,最频繁,最基础的选择器
- 层次选择器,通过DOM元素间的层次关系匹配元素
- 伪元素选择器
- 属性选择器
- 伪类选择器,其语法和其他css选择器有所区别,都以冒号`:`开头
  - 动态伪类选择器
  - 目标伪类选择器
  - 语言伪类选择器
  - UI元素状态伪类选择器
  - 结构伪类选择器
  - 否定伪类选择器

## 一、基本选择器

| 选择器 | 类型 |含义 |
|-------|--------|
|`*`    |通配选择器  |匹配所有元素|
|E      |元素选择器  |匹配所有指定类型的元素|
|#id    |id选择器   |匹配指定id属性的元素(唯一)|
|.class |class选择器|匹配指定class属性的所有元素|
|selector1,...,selectorN|群组选择器|将每个选择器匹配的元素集合并

实例：
```css
* { margin:0; padding:0; }
p { font-size:2em; }
.info { background:#ff0; }
p.info { background:#ff0; }
p.info.error { color:#900; font-weight:bold; }
#info { background:#ff0; }
p#info { background:#ff0; }
```

## 二、层次选择器
| 选择器 |类型| 含义 |
|-------|--------|
|E F	  |后代选择器|匹配所有属于E元素后代的F元素|
|E > F  |子选择器|匹配所有E元素的子元素F|
|E + F  |相邻同胞选择器|匹配所有紧随E元素之后的同级元素F|
|E ~ F  |通用同胞选择器|匹配所有E元素之后的同级元素F|

实例：
```css
div p { color:#f00; }
#nav li { display:inline; }
#nav a { font-weight:bold; }
div > strong { color:#f00; }
.active + div { color:#f00; }
.active ~ div { color:#f00; }
```

## 三、伪类选择器

### 1. 动态伪类选择器

| 选择器 |类型| 含义 |
|E:link	   |链接伪类选择器   |匹配所有未被点击的链接|
|E:visited |链接伪类选择器   |匹配所有已被点击的链接|
|E:active  |用户行为伪类选择器|匹配鼠标已经其上按下、还没有释放的E元素|
|E:hover	 |用户行为伪类选择器|匹配鼠标悬停其上的E元素|
|E:focus	 |用户行为伪类选择器|匹配获得当前焦点的E元素|

### 2. 目标伪类选择器

| 选择器 | 含义 |
|E:target|匹配文档中特定"id"点击后的效果|

```css
h2:target {
	color: white;
	background: #090;
}
```

请参看HTML DOG上关于该选择器的[详细解释](http://htmldog.com/articles/suckerfish/target/)和[实例](http://htmldog.com/articles/suckerfish/target/example/).

### 3. 语言伪类选择器

| 选择器 | 含义 |
|E:lang(xxx)|匹配lang属性为xxx的E元素|

```css
p:lang(sv) { quotes: '"' '"'; }
:lang(sv) { quotes: '"' '"'; }
```

### 4. UI元素状态伪类选择器

| 选择器 | 含义 |
|E:enabled |匹配表单中启用的元素|
|E:disabled|匹配表单中禁用的元素|
|E:checked |匹配表单中被选中的radio（单选框）或checkbox（复选框）元素|

```css
input[type="text"]:disabled { background:#ddd; }
input:disabled {...} // <=> input[disabled] {...}
```

### 5. 结构伪类选择器

| 选择器 | 含义 |
|E:first-child|匹配父元素的第一个子元素E,与`E:nth-child(1)`等同|
|E:last-child|匹配父元素的最后一个子元素E,与`E:nth-last-child(1)`等同|
|E:root               |匹配文档的根元素，对于HTML文档，就是HTML元素|
|E:nth-child(n)       |匹配其父元素的第n个子元素|
|E:nth-last-child(n)  |匹配其父元素的倒数第n个子元素|
|E:nth-of-type(n)     |匹配父元素下具有指定类型的第n个E元素|
|E:nth-last-of-type(n)|匹配父元素下具有指定类型的倒数第n个E元素|
|E:first-of-type      |匹配父元素下使用同种标签的第一个子元素E，与`E:nth-of-type(1)`等同|
|E:last-of-type       |匹配父元素下使用同种标签的最后一个子元素E，与`E:nth-last-of-type(1)`等同|
|E:only-child         |匹配父元素下仅有的一个子元素E|
|E:only-of-type       |匹配父元素下只有一个子元素E|
|E:empty              |匹配一个不包含任何子元素的元素，注意，文本节点也被看作子元素|

> `n`的初始值是`0`,只是计算结果`<1`时,选择器不匹配任何元素.

```css
p:nth-child(3) { color:#f00; }
p:nth-child(odd) { color:#f00; }
p:nth-child(even) { color:#f00; }
p:nth-child(3n+0) { color:#f00; }
p:nth-child(3n) { color:#f00; }
tr:nth-child(2n+11) { background:#ff0; }
tr:nth-last-child(2) { background:#ff0; }
p:last-child { background:#ff0; }
p:only-child { background:#ff0; }
p:empty { background:#ff0; }
```

[测试工具](http://lea.verou.me/demos/nth.html)

### 6. 否定伪类选择器

| 选择器 | 含义 |
|E:not(s)|匹配不符合当前选择器的任何元素|

```css
:not(p) { border:1px solid #ccc; }
```

## 四、伪元素选择器

| 选择器 | 含义 |
|E::first-line  |匹配E元素的第一行|
|E::first-letter|匹配E元素的第一个字母|
|E::before      |在E元素之前插入生成的内容|
|E::after	      |在E元素之后插入生成的内容|
|E::selection   |匹配元素中被用户选中或处于高亮状态的部分|

```css
p::first-line { font-weight:bold; color;#600; }
.preamble::first-letter { font-size:1.5em; font-weight:bold; }
.cbb::before { content:""; display:block; height:17px; width:18px; background:url(top.png) no-repeat 0 0; margin:0 0 0 -18px; }
a:link::after { content: " (" attr(href) ") "; }
::selection { background:red; color:#fff;}
```

> 双冒号和单冒号在CSS3中主要用来区别伪类和伪元素.

## 五、属性选择器

| 选择器 | 含义 |
|E[att] 	    |匹配所有具有att属性的E元素，不考虑它的值.（注意：E在此处可以省略，比如"[cheacked]",以下相同）|
|E[att=val]   |匹配所有att属性等于"val"的E元素|
|`E[att|=val]`|匹配所有att属性为"val"或"val-"开头的E元素，主要用于lang属性，比如"en"、"en-us"、"en-gb"等等,**不推荐,用`E[att^="val"]`代替**|
|E[att~=val]  |匹配所有att属性具有多个空格分隔的值,且其中一个值等于"val"的E元素|
|E[att*="val"]|属性att的值包含"val"字符串的元素|
|E[att^="val"]|属性att的值以"val"开头的元素,可用于替换E[att|="val"]|
|E[att$="val"]|属性att的值以"val"结尾的元素|

```css
p[title] { color:#f00; }
div[class=error] { color:#f00; }
td[headers~=col1] { color:#f00; }
p[lang|=en] { color:#f00; }
blockquote[class=quote][cite] { color:#f00; }
div[id^="nav"] { background:#ff0; }
```
