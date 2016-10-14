# css特殊性

CSS的特殊性(specificity),即CSS的优先级或权值:对于每个样式表规则,浏览器都会计算选择器的特殊性,从而使元素属性声明在**有冲突**的情况下能够正确显示.

## 特殊性的具体特性

选择器的特殊性分为4个等级，`a.b.c.d`，从左到右，越左边的越优先, 如果一个选择器规则有多个相同类型选择器，则`+1`.

1. 内联样式的特殊性为`1.0.0.0`
2. ID选择器的特殊性为 `0.1.0.0`
3. 类,伪类,属性选择器的特殊性为 `0.0.1.0`
4. 元素和伪元素的特殊性为 `0.0.0.1`

## 例外

CSS中还有一种东西可以无视特殊性，那就是`!important`，使用此标记的CSS属性总是最优先的,**不推荐使用**.

1. 如果两个属性都有 !important 那么由特殊性来决定优先级.
2. 优先级相同的情况下，后边定义的会覆盖前边定义的.

### 使用`!important`时需要注意

- Never 永远不要在全站范围的css上使用`!important`
- Only 只在需要覆盖全站或外部css的特定页面中使用`!important`
- Never 永远不要在你的插件中使用`!important`
- Always 要优化考虑使用样式规则的优先级来解决问题而不是`!important`

### 示例

<table>
    <thead>
        <tr>
            <th>选择器</th>
            <th>特殊性</th>
            <th>以10为基数的特殊性</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>style="color: red"</code></td>
            <td>1, 0, 0, 0</td>
            <td>1000</td>
        </tr>
        <tr>
            <td><code>#id {}</code></td>
            <td>0, 1, 0, 0</td>
            <td>100</td>
        </tr>
        <tr>
            <td><code>#id #aid</code></td>
            <td>0, 2, 0, 0</td>
            <td>200</td>
        </tr>
        <tr>
            <td><code>.sty {}</code></td>
            <td>0, 0, 1, 0</td>
            <td>10</td>
        </tr>
        <tr>
            <td><code>.sty p[title=""] {}</code></td>
            <td>0, 0, 2, 0</td>
            <td>20</td>
        </tr>
        <tr>
            <td><code>p:hover {}</code></td>
            <td>0, 0, 1, 0</td>
            <td>10</td>
        </tr>
        <tr>
            <td><code>p {}</code></td>
            <td>0, 0, 0, 1</td>
            <td>1</td>
        </tr>
        <tr>
            <td><code>ul::after {}</code></td>
            <td>0, 0, 0, 1</td>
            <td>1</td>
        </tr>
        <tr>
            <td><code>div p {}</code></td>
            <td>0, 0, 0, 2</td>
            <td>2</td>
        </tr>
    </tbody>
</table>
