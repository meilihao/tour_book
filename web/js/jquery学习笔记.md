# jquery入门

## jQuery 能做什么
jQuery库为Web脚本编程提供了**通用的抽象层**,使得它几乎适用于任何脚本编程的情形。由于它**容易扩展而且不断有新插件面世增强它的功能**,所以根本无法描述涵盖它所有可能的用途和功能。抛开这些不谈,仅就其核心特性而言,jQuery能够满足下列需求：
- **取得文档中的元素**。如果不使用JavaScript库,遍历DOM(Document Object Model,文档对象模型)树,以及查找HTML文档结构中某个特殊的部分,必须编写很多行代码。 jQuery为准确地获取需要检查或操纵的文档元素,提供了可靠而富有效率的选择符机制。
      $('div.content').find('p');
- **修改页面的外观**。CSS虽然为影响文档呈现的方式提供了一种强大的手段,但当所有浏览器不完全支持相同的标准时,单纯使用CSS就会显得力不从心。 jQuery可以弥补这一不足,**它提供了跨浏览器的标准解决方案**。而且,即使在页面已经呈现之后,jQuery仍然能够改变文档中某个部分的类或者个别的样式属性。
      $('ul > li:first').addClass('active');
- **改变文档的内容**。jQuery能够影响的范围并不局限于简单的外观变化,使用少量的代码,jQuery就能改变文档的内容。可以改变文本、插入或翻转图像、列表重新排序,甚至对HTML文档的整个结构都能重写和扩充——所有这些只需一个简单易用的API。
      $('#container').append('<a href="more.html">more</a>');
- **响应用户的交互操作**。即使是最强大和最精心设计的行为,如果我们无法控制它何时发生,那它也毫无用处。jQuery提供了截获形形色色的页面事件(比如用户单击某个链接)的适当方式,而不需要使用事件处理程序拆散HTML代码。此外,它的事件处理API也消除了经常困扰Web开发人员浏览器的不一致性。
      $('button.show-details').click(function() {
        $('div.details').show();
      });
- **为页面添加动态效果**。为了实现某种交互式行为,设计者也必须向用户提供视觉上的反馈。jQuery中内置的一批淡入、擦除之类的效果,以及制作新效果的工具包,为此提供了便利。
      $('div.details').slideDown();
- **无需刷新页面从服务器获取信息**。这种编程模式就是众所周知的Ajax(AsynchronousJavaScript and XML,异步JavaScript和XML),它是一系列在客户端和服务器之间传输数据的强大技术。jQuery通过消除这一过程中的浏览器特定的复杂性,使开发人员得以专注于服务器端的功能设计,从而得以创建出反应灵敏、功能丰富的网站。
      $('div.details').load('more.html #content');
- **简化常见的JavaScript任务**。除了这些完全针对文档的特性之外,jQuery也改进了对基本的JavaScript数据结构的操作(例如迭代和数组操作等)。
      $.each(obj, function(key, value) {
        total += value;
      });
## jQuery 为什么如此出色
为了在维持上述各种特性的同时仍然保持紧凑的代码,jQuery采用了如下策略:

- **利用CSS的优势**。通过将查找页面元素的机制构建于CSS选择符之上, jQuery继承了简明清晰地表达文档结构的方式。由于进行专业Web开发的一个必要条件是掌握CSS语法,因而jQuery成为希望向页面中添加行为的设计者们的切入点。

- **支持扩展**。为了避免特性蠕变(feature creep,有人译为特性蔓延,指软件应用开发中过分强调新的功能以至于损害了其他的设计目标,例如简洁性、轻巧性、稳定性及错误出现率等),jQuery将特殊情况下使用的工具归入插件当中。创建新插件的方法很简单,而且拥有完备的文档说明,这促进了大量有创意且有实用价值的模块的开发。甚至在下载的基本jQuery库文件当中,多数特性在内部都是通过插件架构实现的。而且,如有必要,可以移除这些内部插件,从而生成更小的库文件。

- **抽象浏览器不一致性**。Web开发领域中一个令人遗憾的事实是,每种浏览器对颁布的标准都有自己的一套不太一致的实现方案。任何Web应用程序中都会包含一个用于处理这些平台间特性差异的重要组成部分。虽然不断发展的浏览器前景,使得为某些高级特性提供浏览器中立的完美的基础代码(code base)变得不大可能,但jQuery添加一个抽象层来标准化常见的任务,从而有效地减少了代码量,同时,也极大地简化了这些任务。

- **总是面向集合**。当我们指示jQuery“找到带有 collapsible 类的全部元素,然后隐藏它们”时,不需要循环遍历每一个返回的元素。相反, .hide() 之类的方法被设计成自动操作对象集合,而不是单独的对象。利用这种称作隐式迭代(implicit iteration)的技术,就可以抛弃那些臃肿的循环结构,从而大幅地减少代码量。

- **将多重操作集于一行**。为了避免过度使用临时变量或不必要的代码重复,jQuery在其多数方法中采用了一种称作连缀(chaining)的编程模式。这种模式意味着基于一个对象进行的多数操作的结果,都会返回这个对象自身,以便为该对象应用下一次操作。


这些策略不仅保证了jQuery包的小型化,也为我们使用这个库创建简洁的自定义代码提供了技术保障。
## 其他
通常,JavaScript代码在浏览器初次遇到它们时就会执行,而在浏览器处理头部时,HTML还不会呈现样式。因此,我们需要将代码延迟到DOM可用时再执行。**通过使用 `$(document).ready()` 方法, jQuery支持我们预定在DOM加载完毕后调用某个函数,而不必等待页面中的图像加载**。尽管不使用jQuery,也可以做到这种预定,但`$(document).ready()`为我们提供了很好的跨浏览器解决方案,涉及如下功能:
- 尽可能使用浏览器原生的DOM就绪实现,并以 window.onload 事件处理程序作为后备;
- 可以多次调用 $(document).ready() 并按照调用它们的顺序执行;
- 即便是在浏览器事件发生之后把函数传给 $(document).ready() ,这些函数也会执行;
- 异步处理事件的预定,必要时脚本可以延迟执行;
- 通过重复检查一个几乎与DOM同时可用的方法,在较早版本的浏览器中模拟DOM就绪事件。

`.ready()`方法的参数可以是一个已经定义好的函数的引用,不过通常使用的是匿名函数，特别适合传递那些不会被重用的函数。

## 让jQuery与其他JavaScript库协同工作
在jQuery中,美元符号 `$` 只不过标识符 jQuery 的“别名”。由于 $() 在JavaScript库中很常见,所以,如果在一个页面中使用了几个这样的库,那么就会导致冲突。在这种情况下,可以在我们自定义的jQuery代码中,通过将每个 $ 的实例替换成
jQuery 来避免这种冲突。

# 选择元素
DOM的的树形结构：
```html
<html>
  <head>
    <title>the title</title>
  </head>
  <body>
    <div>
      <p>This is a paragraph.</p>
      <p>This is another paragraph.</p>
      <p>This is yet another paragraph.</p>
    </div>
  </body>
</html>
```
`<html>` 是其他所有元素的**祖先元素**,换句话说,其他所有元素都是 `<html>` 的**后代元素**。 `<head>` 和 `<body>` 元素是 `<html>` 的**子元素**(但并不是它唯一的子元素)。因此除了作为`<head>` 和 `<body>` 的祖先元素之外, `<html>` 也是它们的**父元素**。 而 `<p>` 元素则是 `<div>` 的子元素(也是后代元素),是 `<body>` 和 `<html>` 的后代元素,是其他 `<p>` 元素的**同辈元素**。

我们通过`$()`函数来获取页面的任何元素，其会返回一个新的**jQuery对象实例**,它是我们从现在开始就要打交道的基本的构建块。jQuery对象中会封装零个或多个DOM元素,并允许我们以多种不同的方式与这些DOM元素进行交互。这个函数通常接受CSS选择符作为参数.

```js
//使用否定式伪类选择符来识别没有 horizontal 类的所有列表项
$('#selected-plays li:not(.horizontal)').addClass('sub-level');
//选择带有 alt 属性的所有图像元素
$('img[alt]')
//寻找所有带 href 属性( [href] )且以 mailto 开头( ^="mailto:"] )的锚元素( a )
$('a[href^="mailto:"]').addClass('mailto');
//为所有指向PDF文件的链接添加类
$('a[href$=".pdf"]').addClass('pdflink');
//为 href 属性即以 http 开头且任意位置包含 henry的所有链接添加一个 henrylink 类
$('a[href^="http"][href*="henry"]').addClass('henrylink');
//从带有horizontal 类的< div> 集合中选择第2项（使用自定义选择符）
$('div.horizontal:eq(1)')
//表格中的奇数行添加alt类
$('tr:even').addClass('alt');<=>$('tr').filter(':even').addClass('alt');
//当网页至少有两个table时，上面的代码就有问题了（第二个表格是偶数行添加了alt类）。改进如下
$('tr:nth-child(odd)').addClass('alt');
//提到"Henry"的所有表格单元添加highlight类
$('td:contains(Henry)').addClass('highlight');
//选择所有选中的单选按钮(而不是复选框)
$('input[type="radio"]:checked')
//选择所有密码输入字段和禁用的文本输入字段
$('input[type="password"],input[type="text"]:disabled')
```
>因为JavaScript数组采用从0开始的编号方式,所以 eq(1) 取得的是集合中的第2个元素。而CSS则是从1开始的,因此CSS选择符 $('div:nth-child(1)') 取得的是作为每个div父元素的第1个子元素是div的所有该div子元素。
>
>E:nth-child(n) { sRules },匹配父元素的第n个子元素E，假设该子元素不是E，则选择符无效。**该选择符相对于元素的父元素而非当前选择的所有元素来计算位置**,它可以接受数值、 odd 或 even 作为参数。

## CSS选择符
jQuery支持CSS规范1到规范3中的几乎所有选择符,具体内容可以参考[W3C (World Wide Web Consortium,万维网联盟)网站](http://www.w3.org/Style/CSS/specs),[CSS参考手册](http://css.doyoe.com/),[CSS 选择器](http://www.runoob.com/cssref/css-selectors.html)或者[jQuery API文档#Selectors](http://api.jquery.com/category/selectors/。)。常见的有3种基本的选择符:**标签名、ID和类**。这些选择符可以单独使用,也可以与其他选择符组合使用。
## 属性选择符
属性选择符是CSS选择符中特别有用的一类选择符，是通过HTML元素的属性选择元素.它使用一种从正则表达式中借鉴来的通配符语法,以`^`表示值在字符串的开始,以`$`表示值在字符串的结尾。而且,也是用星号`*`表示要匹配的值可以出现在字符串中的任意位置,用叹号`!`表示对值取反。
## 自定义选择符
除了各种CSS选择符之外,jQuery还添加了独有的完全不同的自定义选择符。
>只要可能,jQuery就会使用浏览器原生的DOM选择符引擎去查找元素。但在使用自定义选择符的时候,就无法使用速度最快的原生方法了。因此,建议在能够使用原生方法的情况下,就不要频繁地使用自定义选择符,以确保性能。

自定义选择符通常跟在一个CSS选择符后面,基于已经选择的元素集的位置来查找元素。自定义选择符的语法与CSS中的伪类选择符语法相同,即选择符以冒号(`:`)开头。比如`:eq()`选择符，`:odd`，`:even`选择符和`:contains() 选择符(区分大小写)`，它们都使用JavaScript内置从0开始的编号方式。

自定义选择符并不局限于基于元素的位置选择元素,还有些适用于表单：

| 选择符     | 匹配    |
|-----------|--------|
|`:input`   |输入字段、文本区、选择列表和按钮元素|
|`:button`  |按钮元素或type属性值为button的输入元素|
|`:enabled` |启用的表单元素|
|`:disabled`|禁用的表单元素|
|`:checked` |勾选的单选按钮或复选框|
|`:selected`|选择的选项元素|
## DOM遍历
参考：[jQuery 遍历](http://www.runoob.com/jquery/jquery-traversing.html)和[jQuery 遍历方法](http://www.runoob.com/jquery/jquery-ref-traversing.html)

### `.filter()`
`.filter()`的功能十分强大,因为它可以接受函数参数。通过传入的函数,可以执行复杂的测试,以决定相应元素是否应该保留在匹配的数组集合中(如果函数返回 false ,则从匹配集合中删除相应元素;如果返回 true ,则保留相应元素)，其通过迭代测试所有匹配的元素来实现。
```js
$('a').filter(function() {
return this.hostname && this.hostname != location.hostname;
}).addClass('external');
```
上面代码可以筛选出符合下面两个条件的 `<a>` 元素：
- 必须包含一个带有域名( this.hostname )的 href 属性。这个测试可以排除 mailto 及类似链接。
- 链接指向的域名(还是 this.hostname )必须不等于( != )页面当前所在域的名称( location.hostname )。

### jquery其他遍历方法
`.next()` 方法只选择下一个最接近的同辈元素；`.nextAll()`方法是返回被选元素之后的所有同辈元素。
`.next()`和`.nextAll()`方法分别有一个对应方法,即`.prev()`和`.prevAll()` 。
`.siblings() 方法`是返回被选元素的所有同胞元素。
`.addBack()`是将先前的元素集合重新加入到当前集合。
`.parent()`，返回被选元素的直接父元素。
`.children()`,返回被选元素的所有直接子元素
`end()`,结束当前链中最近的一次筛选操作，并把匹配元素集合返回到前一次的状态
```js
//给每个包含Henry的单元格的下一个单元格添加样式,可以将已经编写好的选择符作为起点,然后连缀一个 .next() 方法即可
$('td:contains(Henry)').next().addClass('highlight');
//突出显示Henry所在单元格后面的全部单元格,可以使用 .nextAll() 方法
$('td:contains(Henry)').nextAll().addClass('highlight');
//$('td:contains(Henry)').nextAll()是$('td:contains(Henry)')遍历之后的所有同辈元素，addBack()是将$('td:contains(Henry)')加回到当前集合(即nextAll()获取到的集合)中。
$('td:contains(Henry)').nextAll().addBack().addClass('highlight');
//通过 .parent() 方法在DOM中上溯一层到达 <tr> ,然后再通过 .children() 选择该行的所有单元格
$('td:contains(Henry)').parent().children().addClass('highlight');
```
### 连缀
```js
//这里只用于演示,不建议使用如此复杂的连缀方式
$('td:contains(Henry)') //取得包含Henry的所有单元格
.parent() //取得它的父元素
.find('td:eq(1)') //在父元素中查找第2个单元格
.addClass('highlight') //为该单元格添加hightlight类
.end() //恢复到包含Henry的单元格的父元素(筛选动作find,回到find()之前)
.find('td:eq(2)') //在父元素中查找第3个单元格
.addClass('highlight'); //为该单元格添加hightlight类
```
jQuery的连缀能力:在jQuery中,可以通过一行代码取得多个元素集合并对这些元素集合执行多次操作。jQuery的这种连缀能力不仅有助于保持代码简洁,而且在替代组合重新指定选择符时,还有助于提升脚本性能。
>方法连缀的原理:
几乎所有jQuery方法都会返回一个jQuery对象,因而可连缀调用多个jQuery方法.

## 访问DOM元素
jQuery提供了`.get()`方法来直接访问DOM.
```js
//获取带有 id="my-element" 属性的元素的标签名
$('#my-element').get(0).tagName <=> $('#my-element')[0].tagName
```