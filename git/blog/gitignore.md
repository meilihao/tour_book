# gitignore
参考:
- [官方doc](https://git-scm.com/docs/gitignore)

`.gitignore`用于忽略git repo中的某些文件, 比如:
- 含有敏感信息的文件
- 编译出的代码，如 .dll 或 .class
- 系统文件，如 .DS_Store 或 Thumbs.db
- 含有临时信息的文件，如日志、缓存等
- 生成的文件，如 dist 文件夹

`.gitignore`的优势:
- 确保特定的文件不被 Git 追踪
- 通过忽略不需要的文件，它可以帮助你保持代码库的干净
- 它可以控制代码库的大小，这在正在做一个大项目的时候特别有用
- 每一次提交、推送和拉取请求都将是干净的


`.gitignore`中的每一行都指定了一个模式. 这里的`模式`可以指一个特定的文件名，或者指文件名的某些部分结合上通配符.

`.gitignore`的基本规则：
- 任何以哈希（#）开头的行都是注释
- \ 字符可以转义特殊字符
- / 字符表示该规则只适用于位于同一文件夹中的文件和文件夹
- 星号（*）表示任意数量的字符（零个或更多）
- 两个星号（**）表示任意数量的子目录
- 一个问号（?）代替零个或一个字符
- 一个感叹号（!）会反转特定的规则（即包括了任何被前一个模式排除的文件）
- 空行会被忽略，所以你可以用它们来增加空间，使你的文件更容易阅读
- 在末尾添加 / 会忽略整个目录路径

有两种类型的 .gitignore 文件：
- 本地：放在 Git 仓库的根目录下，只在该仓库中工作，并且必须提交到该仓库中
- 全局：放在`$HOME`下，影响`$HOME`中使用的每个仓库，不需要提交
