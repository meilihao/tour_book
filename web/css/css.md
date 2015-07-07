### 声明顺序

相关的属性声明应该以下面的顺序分组处理：

Positioning
Box model 盒模型
Typographic 排版
Visual 外观

Positioning 处在第一位，因为他可以使一个元素脱离正常文本流，并且覆盖盒模型相关的样式。盒模型紧跟其后，因为他决定了一个组件的大小和位置。

其他属性只在组件 内部 起作用或者不会对前面两种情况的结果产生影响，所以他们排在后面。

例:

```css
.declaration-order {
  /* Positioning */
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 100;

  /* Box-model */
  display: block;
  float: right;
  width: 100px;
  height: 100px;

  /* Typography */
  font: normal 13px "Helvetica Neue", sans-serif;
  line-height: 1.5;
  color: #333;
  text-align: center;

  /* Visual */
  background-color: #f5f5f5;
  border: 1px solid #e5e5e5;
  border-radius: 3px;

  /* Misc */
  opacity: 1;
}
```

### 如何减少 CSS 选择器性能损耗？

CSS选择器对性能的影响源于**浏览器匹配选择器和文档元素时所消耗的时间**.

Google 资深web开发工程师 Steve Souders 对 CSS 选择器的执行效率从高到低做了一个排序：

1. id选择器（#myid）
1. 类选择器（.myclassname）
1. 标签选择器（div,h1,p）
1. 相邻选择器（h1+p）
1. 子选择器（ul > li）
1. 后代选择器（li a）
1. 通配符选择器（`*`）
1. 属性选择器（a[rel="external"]）
1. 伪类选择器（a:hover, li:nth-child）

根据以上「选择器匹配」与「选择器执行效率」原则，我们可以通过避免不恰当的使用，提升 CSS 选择器性能。

1、避免使用通用选择器

    .content * {color: red;}

浏览器匹配文档中所有的元素后分别向上逐级匹配 class 为 content 的元素，直到文档的根节点。因此其匹配开销是非常大的，所以应避免使用关键选择器是通配选择器的情况。

2、避免使用标签或 class 选择器限制 id 选择器

    BAD
    button#backButton {…}
    BAD
    .menu-left#newMenuIcon {…}
    GOOD
    #backButton {…}
    GOOD
    #newMenuIcon {…}

3、避免使用标签限制 class 选择器

    BAD
    treecell.indented {…}
    GOOD
    .treecell-indented {…}
    BEST
    .hierarchy-deep {…}

4、避免使用多层标签选择器。使用 class 选择器替换，减少css查找

     BAD
     treeitem[mailfolder="true"] > treerow > treecell {…}
     GOOD
     .treecell-mailfolder {…}

5、避免使用子选择器

    BAD
    treehead treerow treecell {…}
    BETTER, BUT STILL BAD
    treehead > treerow > treecell {…}
    GOOD
    .treecell-header {…}

6、使用继承

    BAD
    #bookmarkMenuItem > .menu-left { list-style-image: url(blah) }
    GOOD
    #bookmarkMenuItem { list-style-image: url(blah) }
