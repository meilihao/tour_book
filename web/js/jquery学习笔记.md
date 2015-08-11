# jquery入门

参考：[jQuery基础教程(第四版)](https://www.packtpub.com/web-development/learning-jquery-fourth-edition)

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
# 事件
>window.onload和jquery.ready区别:
>1. window.onload必须等到**页面内包括图片的所有元素加载完毕后才能执行**，$(document).ready()是**DOM结构绘制完毕后就执行，不必等到加载完毕**.在 jQuery 中也提供与 window.onload 类似的事件触发方法 $(window).load() ,该方法可以绑定到任意元素上。
>2. window.onload不能同时编写多个，如果有多个window.onload方法时，前面的方法会被后面的覆盖，因此只会执行最后一个，$(document).ready()可以同时编写多个，并且都可以得到执行(每次调用这个方法都会向内部的行为队列中添加一个新函数,当页面加载完成后,所有函数都会被执行。而且,这些函数会**按照注册它们的顺序依次执行**).
>3. window.onload没有简化写法，`$(document).ready(function(){})`可以简写成`$(function(){})`,不过**推荐使用较长的形式**,因为较长的形式能够更清楚地表明代码在做什么.

>一般来说,使用 `$(document).ready()` 要优于使用 onload 事件处理程序,但必须要明确的一点是,因为支持文件可能还没有加载完成,所以**类似图像的高度和宽度这样的属性此时则不一定会有效**。如果需要访问这些属性,可能就得选择实现一个onload 事件处理程序(或者是使用jQuery为load事件设置处理程序)。这两种机制能够和平共存。

>加载样式与执行代码:
>为了保证JavaScript代码执行以前页面已经应用了样式,最好是在 `<head>`元素中把 `<link rel="stylesheet">` 标签和`<style>` 标签放在 `<script>` 标签前面.

由于很多库都使用 `$` 标识符(因为它简短方便),因此就需要一种方式来避免名称冲突。jQuery提供了一个 `jQuery.noConflict()` 方法,调用该方法可以把对 $标识符的控制权让渡还给其他库。之后在需要使用jQuery方法时,必须记住要用 `jQuery`而不是 `$` 来调用.
在这种情况下,还有一个在 .ready() 方法中使用 $ 的技巧。我们传递给它的回调函数可以接收一个参数——jQuery对象本身。利用这个参数,可以重新命名 jQuery 为 $ ,而不必担心造成冲突::
```js
jQuery(document).ready(function($) {
//在这里,可以正常使用!
});
```
`on()`在选定的元素上绑定一个或多个事件处理函数.
`.removeClass()`方法的参数是可选的,即当省略参数时,该方法会移除元素中所有的类。
`.toggleClass()`方法检查每个元素中指定的类。如果不存在则添加类，如果已设置则删除之。这就是所谓的切换效果.
`.hover()`方法规定当鼠标指针悬停在被选元素上时要运行的两个函数(mouseenter 和 mouseleave 事件).

当触发任何事件处理程序时,关键字 this 引用的都是携带相应行为的DOM元素(即利用了事件处理程序运行的上下文)。通过在事件处理程序中使用 $(this) ,可以为相应的元素创建jQuery对象,然后就如同使用CSS选择符找到该元素一样对它进行操作.**利用处理程序的上下文将语句通用化,可以使代码更高效**.

## 事件传播
DOM Level3支持DOM标准的事件模型，即捕获与冒泡型.当一个DOM事件被触发的时候，它并不只是在它的起源对象上触发一次，而是会经历三个不同的阶段。简而言之：事件一开始从文档的根节点流向目标对象（捕获阶段），然后在目标对象上被触发（目标阶段），之后再回溯到文档的根节点（冒泡阶段）.
参考:[DOM事件简介#事件阶段(Event Phases)](http://blog.jobbole.com/52430/)
```html
<div class="foo">
  <span class="bar">
    <a href="http://www.example.com/">
      The quick brown fox jumps over the lazy dog.
    </a>
  </span>
  <p>
    How razorback-jumping frogs can level six piqued gymnasts!
  </p>
</div>
```
从逻辑上看,任何事件都可能会有多个元素负责响应。举例来说,如果单击了页面中的链接元素,那么`<div> 、 <span> 和 <a>` 全都应该得到响应这次单击的机会。

允许多个元素响应单击事件的一种策略叫做**事件捕获**。在事件捕获的过程中,事件首先会交给最外层的元素,接着再交给更具体的元素。在这个例子中,意味着单击事件首先会传递给`<div> ,然后是 <span> ,最后是 <a>`.

另一种相反的策略叫做**事件冒泡**.即当事件发生时,会首先发送给最具体的元素,在这个元素获得响应机会之后,事件会向上冒泡到更一般的元素。
>事件捕获和下文中的事件冒泡是“浏览器大战”时期分别由Netscape和微软提出的两种相反的事件传播模型。

因而,最终出台的DOM标准规定应该同时使用这两种策略:首先,事件要从一般元素到具体元素逐层捕获,然后,事件再通过冒泡返回DOM树的顶层。而事件处理程序可以注册到这个过程中的任何一个阶段。为了确保跨浏览器的一致性,而且也为了让人容易理解,**jQuery始终会在模型的冒泡阶段注册事件处理程序**。

### 事件冒泡的副作用
事件冒泡可能会导致始料不及的行为,特别是在错误的元素响应 mouseover 或 mouseout 事件的情况下。

假设在上面例子中,为 `<div>` 添加了一个 mouseout 事件处理程序。当用户的鼠标指针退出这个 `<div>` 时,会按照预期运行 mouseout 处理程序。因为这个过程发生在顶层元素上,所以其他元素不会取得这个事件。但是,当指针从 `<a>` 元素上离开时, `<a>` 元素也会取得一个 mouseout 事件。然后,这个事件会向上冒泡到 `<span>` 和 `<div>` ,从而触发上述的事件处理程序。这种冒泡序列很可能不是我们所希望的。

而jQuery提供了mouseenter和mouseleave事件，使用它们来代替mouseover和mouseout,无论是单独绑定,还是在 `.hover()` 方法中组合绑定,都可以避免这些冒泡问题。在使用它们处理事件的时候,可以不用担心某些非目标元素得到mouseover 或 mouseout 事件导致的问题。

## 通过事件对象改变事件传播

```html
<style>
.hidden {
  display: none;
}
</style>
<div id="switcher" class="switcher">
  <h3>Style Switcher</h3> 
  <button id="switcher-default" class=""> Default </button>
  <button id="switcher-narrow" class=""> Narrow Column </button>
  <button id="switcher-large" class=""> Large Print </button>
</div>
```
```js
$(document).ready(function() {
  $('#switcher').click(function() {
    $('#switcher button').toggleClass('hidden');
  });
});
```
我们期望单击div背景区域才切换,而不包括后代元素.上面代码会使`<div class="switcher">`整个区域都可以通过单击切换其可见性。但同时也造成了一个问题,即单击其后代元素,比如`<button>`也会导致切换。导致这个问题的原因就是事件冒泡,即事件首先被按钮处理,然后又沿着DOM树向上传递,直至到达`<div id="switcher">`激活事件处理程序并隐藏按钮。

要解决这个问题,必须访问**事件对象**。事件对象是一种DOM结构,它会在元素获得处理事件的机会时传递给被调用的事件处理程序。这个对象中包含着与事件有关的信息(例如事件发生时的鼠标指针位置),也提供了可以用来影响事件在DOM中传递进程的一些方法。

为了在处理程序中使用事件对象,需要为函数添加一个参数,通常命名为`event`.

### 事件目标
事件处理程序中的变量 event 保存着事件对象。而 event.target 属性保存着发生事件的目标元素。通过 .target ,可以确定DOM中首先接收到事件的元素(即实际被单击的元素).
```js
$(document).ready(function() {
  $('#switcher').click(function(event) {
    if (event.target == this) {
      $('#switcher button').toggleClass('hidden');
    }
  });
});
```
此时的代码确保了被单击的元素是 `<div class="switcher">`,而不是其他后代元素。现在,单击`<button>`不会再出现切换,而单击div背景区则会触发切换。

### 停止事件传播
事件对象还提供了一个 .stopPropagation() 方法,该方法可以完全阻止事件冒泡。与 .target 类似,这个方法也是一种基本的DOM特性.

### 阻止默认操作
如果我们把单击事件处理程序注册到锚元素(`<a>`),而不是外层的 `<div>` 上,那么就要面对另外一个问题:当用户单击链接时,浏览器会加载一个新页面。这种行为与我们讨论的事件处理程序不是同一个概念,它是单击锚元素的默认操作。类似地,当用户在编辑完表单后按下回车键时,会触发表单的 `submit` 事件,在此事件发生后,表单提交才会真正发生。

即便在事件对象上调用 `.stopPropagation()` 方法也不能禁止这种默认操作,因为默认操作不是在正常的事件传播流中发生的。在这种情况下, `.preventDefault()` 方法则可以在触发默认操作之前终止事件.

事件传播和默认操作是相互独立的两套机制,在二者任何一方发生时,都可以终止另一方。如果想要同时停止事件传播和默认操作,可以在事件处理程序中返回 false ,这是对在事件对象上同时调用 `.stopPropagation()` 和 `.preventDefault()` 的一种简写方式。

### 事件委托
**事件委托就是利用冒泡的一项高级技术**。通俗的讲，事件就是onclick，onmouseover，onmouseout等就是事件，委托，就是让别人来做，这个事件本来是加在某些元素上的，然而你却加到别人身上来做，完成这个事件,即利用冒泡的原理，把事件加到父级上，触发执行效果。其好处:
1. 提高性能
1. 新添加的元素还会有之前指定的事件,如table新添加行的事件.

例如,有一个显示信息的大型表格,每一行都有一项需要注册单击处理程序。虽然不难通过隐式迭代来指定所有单击处理程序,但性能可能会很成问题,因为循环是由jQuery在内部完成的,而且要维护所有处理程序也需要占用很多内存。为解决这个问题,可以只在DOM中的一个祖先元素上指定一个单击处理程序。由于事件会冒泡,未遭拦截的单击事件最终会到达这个祖先元素,而我们可以在此时再作出相应处理。

`.is()`方法用于查看选择的元素是否匹配选择器.
>is() 与 .hasClass()
>要测试元素是否包含某个类,也可以使用另一个简写方法 .hasClass() .不过, .is() 方法则更灵活一些,它可以测试任何选择符表达式,如$(selector).is(".className,.className"),$(selector).is("div").

`$(event.target).is('button')` 测试被单击的元素是否包含 button 标签
### 使用内置的事件委托功能
`.on()` 方法可以接受相应参数实现事件委托.
```js
$('#switcher').on('click', 'button', function() {
  var bodyClass = event.target.id.split('-')[1];
  $('body').removeClass().addClass(bodyClass);
  $('#switcher button').removeClass('selected');
  $(this).addClass('selected');
});
```
如果给 .on() 方法传入的第二个参数是一个选择符表达式,jQuery会把 click 事件处理程序绑定到 #switcher 对象,同时比较 event.target 和选择符表达式(这里的 'button' )。如果匹配,jQuery会把 this 关键字映射到匹配的元素,否则不会执行事件处理程序。

## 移除事件处理程序
`.off()` 方法通常用于移除通过 `.on()` 方法添加的事件处理程序.
```js
$(document).ready(function() {
  $('#switcher').click(function(event) { //事件委托
    if (!$(event.target).is('button')) {
      $('#switcher button').toggleClass('hidden');
    }
  });
  $('#switcher-narrow, #switcher-large').click(function() {
    $('#switcher').off('click');
  });
  $('#switcher').click(function(event) {
    if ($(event.target).is('button')) {
      var bodyClass = event.target.id.split('-')[1];

      $('body').removeClass().addClass(bodyClass);

      $('#switcher button').removeClass('selected');
      $(event.target).addClass('selected');
    }
  });
});
```
`<div id="switcher">` 上的单击处理程序就会被移除。然后,再单击背景区域将不会导致它折叠起来。但是,其里面的按钮本身的作用却失效了(由于为使用事件委托而重写了按钮处理代码的原因).

为了.off() 的调用更有针对性,以避免把注册的两个单击处理程序全都移除。达成目标的一种方式是使用**事件命名空间**,即在绑定事件时引入附加信息,以便将来识别特定的处理程序。
```js
$(document).ready(function() {
  $('#switcher').on('click.collapse', function(event) {
    if (!$(event.target).is('button')) {
      $('#switcher button').toggleClass('hidden');
    }
  });
  $('#switcher-narrow, #switcher-large').click(function() {
    $('#switcher').off('click.collapse');
  });
});
```
对于事件处理系统而言,后缀.collapse是不可见的。换句话说,这里仍然会像编.on('click')一样,让注册的函数响应单击事件。但是,通过附加的命名空间信息,则可以解除对这个特定处理程序的绑定,同时不影响为按钮注册的其他单击处理程序。

### 重新绑定事件
```js
//未完成的代码
$(document).ready(function() {
  var toggleSwitcher = function(event) {
    if (!$(event.target).is('button')) {
      $('#switcher button').toggleClass('hidden');
    }
  };
  $('#switcher').on('click.collapse', toggleSwitcher);
  $('#switcher-narrow, #switcher-large').click(function() {
    $('#switcher').off('click.collapse');
  });
  $('#switcher-default').click(function() {
    $('#switcher').on('click.collapse', toggleSwitcher);
  });
});
```
每次点击 Default(按钮)就会有一个toggleSwitcher的副本被绑定到样式转换器。换句话说,在用户单击 Narrow 或 Large Print 之前(这样就可以一次性地解除对toggleSwitcher的绑定),每多单击Default一次都会多调用一次这个函数。导致在绑定toggleSwitcher偶数次的情况下,单击`<div id="switcher">`(不是按钮),好像一切都没有发生变化。事实上,这是因为切换了hidden类偶数次,结果状态与开始的时候相同。

改进:
```js
$(document).ready(function() {
  var toggleSwitcher = function(event) {  //使用命令函数还有另外一个好处,即不必再使用事件命名空间.因为 .off()可以将这个命名函数作为第二个参数,结果只会解除对特定处理程序的绑定。
    if (!$(event.target).is('button')) {
      $('#switcher button').toggleClass('hidden');
    }
  };
  $('#switcher').on('click', toggleSwitcher);
  $('#switcher button').click(function() {
    $('#switcher').off('click', toggleSwitcher);
    if (this.id == 'switcher-default') {
      $('#switcher').on('click', toggleSwitcher);
    }
  });
});
```

对于只需触发一次,随后要立即解除绑定的情况也有一种简写方法—— `.one()`,这个简写方法的用法如下:

    $('#switcher').one('click', toggleSwitcher);
## 模拟用户操作
通过 `.trigger()` 方法就可以完成模拟事件的操作.
```js
$(document).ready(function() {
  $('#switcher').trigger('click'); //<=>$('#switcher').click()当不带参数时,.trigger()的简写方法
});
```
事件对象的 `.which` 属性包含着被按下的哪个键的标识符.
```js
$(document).ready(function() {
  var triggers = {
    D: 'default',
    N: 'narrow',
    L: 'large'
  };
  $(document).keyup(function(event) {
    var key = String.fromCharCode(event.which); //String.fromCharCode(),接受一个或多个Unicode值，然后返回一个字符串
    if (key in triggers) {
      $('#switcher-' + triggers[key]).click();
    }
  });
});
```
# 样式和动画
## 修改内联CSS
`css()`方法设置或返回被选元素的一个或多个样式属性.
`hide()`方法隐藏被选元素.
`show()`方法显示隐藏的被选元素.

## 隐藏和显示元素
.hide() 方法会将匹配的元素集合的内联 style 属性设置为 display:none 。但它的聪明之处是,它能够在把 display 的值变成 none 之前,记住原先的 display 值,通常是 block 、inline 或 inline-block 。恰好相反, .show() 方法会将匹配的元素集合的 display 属性,恢复为应用 display: none 之前的可见属性。

当在 .show() 或 .hide() 中指定时长(或更准确地说,一个速度)参数时,就会产生动画效果,即效果会在一个特定的时间段内发生。它们都可以指定两种预设的速度参数: 'slow' 和 'fast' 。使用 .show('slow') 会在600毫秒(0.6秒)内完成效果,而 .show('fast') 则是200毫秒(0.2秒)。如果传入的是其他字符串,jQuery就会在默认的400毫秒内完成效果。要指定更精确的速度,可以使用毫秒数值,例如 .show(850) 。

### 淡入和淡出
fadeIn() 方法逐渐改变被选元素的不透明度，从隐藏到可见（褪色效果）;fadeOut() 方法逐渐改变被选元素的不透明度，从可见到隐藏（褪色效果）。

对于本来就处于文档流之外的元素,比较适合使用淡入和淡出动画。例如,对于那些覆盖在页面之上的“亮盒”元素来说,采用淡入和淡出就显得很自然。
### 滑上和滑下
.slideDown() 和 .slideUp() 方法仅改变元素的高度.要让内容以垂直滑入/滑出(折叠)的效果出现.
`.slideToggle()` 方法通过逐渐增加或减少元素高度来显示或隐藏元素,类试".toggle()".

## 自定义动画
`.animate()` 方法执行 CSS 属性集的自定义动画。
`.outWidth()`方法获取元素集合中第一个元素的当前计算宽度值,包括padding，border和选择性的margin,前提是确保该元素在使用.outerWidth()前可见,否则得到的值不能保证准确.

在使用 .animate() 方法时,必须明确CSS对我们要改变的元素所施加的限制。例如,在元素的CSS定位没有设置成 relative 或 absolute 的情况下,调整 left 属性对于匹配的元素毫无作用。所有块级元素默认的CSS定位属性都是 static ,这个值精确地表明:在改变元素的定位属性之前试图移动它们,它们只会保持静止不动。

## 并发与排队效果
.queue() 方法把非效果方法添加到队列中的一种方式.
