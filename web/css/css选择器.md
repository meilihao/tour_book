## CSS 选择器

[原文](http://www.runoob.com/cssref/css-selectors.html)

CSS选择器用于选择你想要的元素的样式的模式。

"CSS"列表示在CSS版本的属性定义（CSS1，CSS2，或对CSS3）

<table class="reference"> <tbody><tr> <th width="22%" align="left">选择器</th> <th width="17%" align="left">示例</th> <th width="56%" align="left">示例说明</th> <th align="left">CSS</th> </tr> <tr> <td><a href="sel-class.html">.<i>class</i></a></td> <td class="notranslate">.intro</td> <td>选择所有class="intro"的元素</td> <td>1</td> </tr> <tr> <td><a href="sel-id.html">#<i>id</i></a></td> <td class="notranslate">#firstname</td> <td>选择所有id="firstname"的元素</td> <td>1</td> </tr> <tr> <td><a href="sel-all.html">*</a></td> <td class="code notranslate">*</td> <td>选择所有元素</td> <td>2</td> </tr> <tr> <td><i><a href="sel-element.html">element</a></i></td> <td class="notranslate">p</td> <td>选择所有&lt;p&gt;元素</td> <td>1</td> </tr> <tr> <td><i><a href="sel-element-comma.html">element,element</a></i></td> <td class="notranslate">div,p</td> <td>选择所有&lt;div&gt;元素和&lt;p&gt;元素</td> <td>1</td> </tr> <tr> <td><a href="sel-element-element.html"><i>element</i> <i>element</i></a></td> <td class="notranslate">div p</td> <td>选择&lt;div&gt;元素内的所有&lt;p&gt;元素</td> <td>1</td> </tr> <tr> <td><a href="sel-element-gt.html"><i>element</i>&gt;<i>element</i></a></td> <td class="notranslate">div&gt;p</td> <td>选择所有父级是 &lt;div&gt; 元素的 &lt;p&gt; 元素</td> <td>2</td> </tr> <tr> <td><a href="sel-element-pluss.html"><i>element</i>+<i>element</i></a></td> <td class="notranslate">div+p</td> <td>选择所有紧接着&lt;div&gt;元素之后的&lt;p&gt;元素</td> <td>2</td> </tr> <tr> <td><a href="sel-attribute.html">[<i>attribute</i>]</a></td> <td class="notranslate">[target]</td> <td>选择所有带有target属性元素</td> <td>2</td> </tr> <tr> <td><a href="sel-attribute-value.html">[<i>attribute</i>=<i>value</i>]</a></td> <td class="notranslate">[target=-blank]</td> <td>选择所有使用target="-blank"的元素</td> <td>2</td> </tr> <tr> <td><a href="sel-attribute-value-contains.html">[<i>attribute</i>~=<i>value</i>]</a></td> <td class="notranslate">[title~=flower]</td> <td>选择标题属性包含单词"flower"的所有元素</td> <td>2</td> </tr> <tr> <td><a href="sel-attribute-value-lang.html">[<i>attribute</i>|=<i>language</i>]</a></td> <td class="notranslate">[lang|=en]</td> <td>选择一个lang属性的起始值="EN"的所有元素</td> <td>2</td> </tr> <tr> <td><a href="sel-link.html">:link</a></td> <td class="notranslate">a:link</td> <td>选择所有未访问链接</td> <td>1</td> </tr> <tr> <td><a href="sel-visited.html">:visited</a></td> <td class="notranslate">a:visited</td> <td>选择所有访问过的链接</td> <td>1</td> </tr> <tr> <td><a href="sel-active.html">:active</a></td> <td class="notranslate">a:active</td> <td>选择活动链接</td> <td>1</td> </tr> <tr> <td><a href="sel-hover.html">:hover</a></td> <td class="notranslate">a:hover</td> <td>选择鼠标在链接上面时</td> <td>1</td> </tr> <tr> <td><a href="sel-focus.html">:focus</a></td> <td class="notranslate">input:focus</td> <td>选择具有焦点的输入元素</td> <td>2</td> </tr> <tr> <td><a href="sel-firstletter.html">:first-letter</a></td> <td class="notranslate">p:first-letter</td> <td>选择每一个&lt;P&gt;元素的第一个字母</td> <td>1</td> </tr> <tr> <td><a href="sel-firstline.html">:first-line</a></td> <td class="notranslate">p:first-line</td> <td>选择每一个&lt;P&gt;元素的第一行</td> <td>1</td> </tr> <tr> <td><a href="sel-firstchild.html">:first-child</a></td> <td class="notranslate">p:first-child</td> <td>指定只有当&lt;p&gt;元素是其父级的第一个子级的样式。</td> <td>2</td> </tr> <tr> <td><a href="sel-before.html">:before</a></td> <td class="notranslate">p:before</td> <td>在每个&lt;p&gt;元素之前插入内容</td> <td>2</td> </tr> <tr> <td><a href="sel-after.html">:after</a></td> <td class="notranslate">p:after</td> <td>在每个&lt;p&gt;元素之后插入内容</td> <td>2</td> </tr> <tr> <td><a href="sel-lang.html">:lang(<i>language</i>)</a></td> <td class="notranslate">p:lang(it)</td> <td>选择一个lang属性的起始值="it"的所有&lt;p&gt;元素</td> <td>2</td> </tr> <tr> <td><a href="sel-gen-sibling.html"><i>element1</i>~<i>element2</i></a></td> <td>p~ul</td> <td>选择p元素之后的每一个ul元素</td> <td>3</td> </tr> <tr> <td><a href="sel-attr-begin.html">[<i>attribute</i>^=<i>value</i>]</a></td> <td>a[src^="https"]</td> <td>选择每一个src属性的值以"https"开头的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-attr-end.html">[<i>attribute</i>$=<i>value</i>]</a></td> <td>a[src$=".pdf"]</td> <td>选择每一个src属性的值以".pdf"结尾的元素 </td> <td>3</td> </tr> <tr> <td><a href="sel-attr-contain.html">[<i>attribute</i>*=<i>value</i>]</a></td> <td>a[src*="44lan"]</td> <td>选择每一个src属性的值包含子字符串"44lan"的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-first-of-type.html">:first-of-type</a></td> <td>p:first-of-type</td> <td>选择每个p元素是其父级的第一个p元素</td> <td>3</td> </tr> <tr> <td><a href="sel-last-of-type.html">:last-of-type</a></td> <td>p:last-of-type</td> <td>选择每个p元素是其父级的最后一个p元素</td> <td>3</td> </tr> <tr> <td><a href="sel-only-of-type.html">:only-of-type</a></td> <td>p:only-of-type</td> <td>选择每个p元素是其父级的唯一p元素</td> <td>3</td> </tr> <tr> <td><a href="sel-only-child.html">:only-child</a></td> <td>p:only-child</td> <td>选择每个p元素是其父级的唯一子元素</td> <td>3</td> </tr> <tr> <td><a href="sel-nth-child.html">:nth-child(<i>n</i>)</a></td> <td>p:nth-child(2)</td> <td>选择每个p元素是其父级的第二个子元素</td> <td>3</td> </tr> <tr> <td><a href="sel-nth-last-child.html">:nth-last-child(<i>n</i>)</a></td> <td>p:nth-last-child(2)</td> <td>选择每个p元素的是其父级的第二个子元素，从最后一个子项计数</td> <td>3</td> </tr> <tr> <td><a href="sel-nth-of-type.html">:nth-of-type(<i>n</i>)</a></td> <td>p:nth-of-type(2)</td> <td>选择每个p元素是其父级的第二个p元素</td> <td>3</td> </tr> <tr> <td><a href="sel-nth-last-of-type.html">:nth-last-of-type(<i>n</i>)</a></td> <td>p:nth-last-of-type(2)</td> <td>选择每个p元素的是其父级的第二个p元素，从最后一个子项计数</td> <td>3</td> </tr> <tr> <td><a href="sel-last-child.html">:last-child</a></td> <td>p:last-child</td> <td>选择每个p元素是其父级的最后一个子级。</td> <td>3</td> </tr> <tr> <td><a href="sel-root.html">:root</a></td> <td>:root</td> <td>选择文档的根元素</td> <td>3</td> </tr> <tr> <td><a href="sel-empty.html">:empty</a></td> <td>p:empty</td> <td>选择每个没有任何子级的p元素（包括文本节点）</td> <td>3</td> </tr> <tr> <td><a href="sel-target.html">:target</a></td> <td>#news:target </td> <td>选择当前活动的#news元素（包含该锚名称的点击的URL）</td> <td>3</td> </tr> <tr> <td><a href="sel-enabled.html">:enabled</a></td> <td>input:enabled</td> <td>选择每一个已启用的输入元素</td> <td>3</td> </tr> <tr> <td><a href="sel-disabled.html">:disabled</a></td> <td>input:disabled</td> <td>选择每一个禁用的输入元素</td> <td>3</td> </tr> <tr> <td><a href="sel-checked.html">:checked</a></td> <td>input:checked</td> <td>选择每个选中的输入元素</td> <td>3</td> </tr> <tr> <td><a href="sel-not.html">:not(<i>selector</i>)</a></td> <td>:not(p)</td> <td>选择每个并非p元素的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-selection.html">::selection</a></td> <td>::selection</td> <td>匹配元素中被用户选中或处于高亮状态的部分</td> <td>3</td> </tr>  <tr> <td><a href="sel-out-of-range.html">:out-of-range</a></td> <td>:out-of-range</td> <td>匹配值在指定区间之外的input元素</td> <td>3</td> </tr>  <tr> <td><a href="sel-in-range.html">:in-range</a></td> <td>:in-range</td> <td>匹配值在指定区间之内的input元素</td> <td>3</td> </tr> <tr> <td><a href="sel-read-write.html">:read-write</a></td> <td>:read-write</td> <td>用于匹配可读及可写的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-read-only.html">:read-only</a></td> <td>:read-only</td> <td>用于匹配设置 "readonly"（只读） 属性的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-optional.html">:optional</a></td> <td>:optional</td> <td>用于匹配可选的输入元素</td> <td>3</td> </tr> <tr> <td><a href="sel-required.html">:required </a></td> <td>:required</td> <td>用于匹配设置了 "required" 属性的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-valid.html">:valid </a></td> <td>:valid</td> <td>用于匹配输入值为合法的元素</td> <td>3</td> </tr> <tr> <td><a href="sel-invalid.html">:invalid </a></td> <td>:invalid</td> <td>用于匹配输入值为非法的元素</td> <td>3</td> </tr> </tbody></table>

## CSS选择器笔记

[原文](http://www.ruanyifeng.com/blog/2009/03/css_selectors.html)

### 一、基本选择器

| 序号 | 选择器 | 含义 |
|--------|-------|--------|
|  1     |`*`    |通用元素选择器，匹配任何元素|
|  2     |E      |标签选择器，匹配所有使用E标签的元素|
|  3     |.info  |class选择器，匹配所有class属性中包含info的元素|
|  4     |#footer|id选择器，匹配所有id属性等于footer的元素|

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

### 二、多元素的组合选择器
| 序号 | 选择器 | 含义 |
|--------|-------|--------|
|  5     |E,F	 |多元素选择器，同时匹配所有E元素或F元素，E和F之间用逗号分隔|
|  6     |E F	 |后代元素选择器，匹配所有属于E元素后代的F元素，E和F之间用空格分隔|
|  7     |E > F  |子元素选择器，匹配所有E元素的子元素F|
|  8     |E + F  |毗邻元素选择器，匹配所有紧随E元素之后的同级元素F|

实例：
```css
div p { color:#f00; }
#nav li { display:inline; }
#nav a { font-weight:bold; }
div > strong { color:#f00; }
p + p { color:#f00; }
```

### 三、CSS 2.1 属性选择器
| 序号 | 选择器 | 含义 |
|--------|-------------|--------|
|    9   |E[att] 	   |匹配所有具有att属性的E元素，不考虑它的值。（注意：E在此处可以省略，比如"[cheacked]"。以下同。）|
|   10   |E[att=val]   |匹配所有att属性等于"val"的E元素|
|   11   |E[att~=val]  |匹配所有att属性具有多个空格分隔的值、其中一个值等于"val"的E元素|
|   12   |`E[att|=val]`|匹配所有att属性具有多个连字号分隔（`-`的值、其中一个值以"val"开头的E元素，主要用于lang属性，比如"en"、"en-us"、"en-gb"等等|

实例：
```css
p[title] { color:#f00; }
div[class=error] { color:#f00; }
td[headers~=col1] { color:#f00; }
p[lang|=en] { color:#f00; }
blockquote[class=quote][cite] { color:#f00; }
```

### 四、CSS 2.1中的伪类
| 序号 | 选择器 | 含义 |
|--------|-------------|--------|
|   13   |E:first-child|匹配父元素的第一个子元素|
|   14   |E:link	   |匹配所有未被点击的链接|
|   15   |E:visited	   |匹配所有已被点击的链接|
|   16   |E:active     |匹配鼠标已经其上按下、还没有释放的E元素|
|   17   |E:hover	   |匹配鼠标悬停其上的E元素|
|   18   |E:focus	   |匹配获得当前焦点的E元素|
|   19   |E:lang(c)	   |匹配lang属性等于c的E元素|

实例：
```css
p:first-child { font-style:italic; }
input[type=text]:focus { color:#000; background:#ffe; }
input[type=text]:focus:hover { background:#fff; }
q:lang(sv) { quotes: "\201D" "\201D" "\2019" "\2019"; }
```

### 五、 CSS 2.1中的伪元素
| 序号 | 选择器 | 含义 |
|--------|--------------|--------|
|   20   |E:first-line  |匹配E元素的第一行|
|   21   |E:first-letter|匹配E元素的第一个字母|
|   22   |E:before      |在E元素之前插入生成的内容|
|   23   |E:after	    |在E元素之后插入生成的内容|

实例：
```css
p:first-line { font-weight:bold; color;#600; }
.preamble:first-letter { font-size:1.5em; font-weight:bold; }
.cbb:before { content:""; display:block; height:17px; width:18px; background:url(top.png) no-repeat 0 0; margin:0 0 0 -18px; }
a:link:after { content: " (" attr(href) ") "; }
```

### 六、CSS 3的同级元素通用选择器
| 序号 | 选择器 | 含义 |
|--------|-------|--------|
|   24   |E ~ F  |匹配任何在E元素之后的同级F元素|

实例：
```
p ~ ul { background:#ff0; }
```

### 七、CSS 3 属性选择器
| 序号 | 选择器 | 含义 |
|--------|-------------|--------|
|   25   |E[att^="val"]|属性att的值以"val"开头的元素,可用于替换E[att|="val"]|
|   26   |E[att$="val"]|属性att的值以"val"结尾的元素|
|   27   |E[att*="val"]|属性att的值包含"val"字符串的元素|

实例：
```css
div[id^="nav"] { background:#ff0; }
```

### 八、CSS 3中与用户界面有关的伪类
| 序号 | 选择器 | 含义 |
|--------|-------------|--------|
|   28   |E:enabled|匹配表单中激活的元素|
|   29   |E:disabled|匹配表单中禁用的元素|
|   30   |E:checked|匹配表单中被选中的radio（单选框）或checkbox（复选框）元素|
|   31   |E::selection|匹配用户当前选中的元素|

实例：
```css
input[type="text"]:disabled { background:#ddd; }
```

### 九、CSS 3中的结构伪类(选择器)
| 序号 | 选择器 | 含义 |
|--------|---------------------|--------|
|   32   |E:root               |匹配文档的根元素，对于HTML文档，就是HTML元素|
|   33   |E:nth-child(n)       |匹配其父元素的第n个子元素，第一个编号为1|
|   34   |E:nth-last-child(n)  |匹配其父元素的倒数第n个子元素，第一个编号为1|
|   35   |E:nth-of-type(n)     |与:nth-child()作用类似，但是仅匹配使用同种标签的元素|
|   36   |E:nth-last-of-type(n)|与:nth-last-child() 作用类似，但是仅匹配使用同种标签的元素|
|   37   |E:last-child         |匹配父元素的最后一个子元素，等同于:nth-last-child(1)|
|   38   |E:first-of-type      |匹配父元素下使用同种标签的第一个子元素，等同于:nth-of-type(1)|
|   39   |E:last-of-type       |匹配父元素下使用同种标签的最后一个子元素，等同于:nth-last-of-type(1)|
|   40   |E:only-child         |匹配父元素下仅有的一个子元素，等同于:first-child:last-child或 :nth-child(1):nth-last-child(1)|
|   41   |E:only-of-type       |匹配父元素下使用同种标签的唯一一个子元素，等同于:first-of-type:last-of-type或 :nth-of-type(1):nth-last-of-type(1)|
|   42   |E:empty              |匹配一个不包含任何子元素的元素，注意，文本节点也被看作子元素|

实例：
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
### 十、CSS 3的反选伪类
| 序号 | 选择器 | 含义 |
|--------|--------|--------|
|   43   |E:not(s)|匹配不符合当前选择器的任何元素|

实例：
```css
:not(p) { border:1px solid #ccc; }
```
### 十一、CSS 3中的 :target 伪类
| 序号 | 选择器 | 含义 |
|--------|--------|--------|
|   44   |E:target|匹配文档中特定"id"点击后的效果|

请参看HTML DOG上关于该选择器的[详细解释](http://htmldog.com/articles/suckerfish/target/)和[实例](http://htmldog.com/articles/suckerfish/target/example/)。
