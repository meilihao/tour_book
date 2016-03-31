# alias

## 描述

alias是一个系统自建的shell命令，允许你为比较长或者经常使用的命令指定别名(新别名将取代同名的旧别名).alias只在当前shell有用,通常将其放入"~/.bashrc".

别名可能会导致安全问题(命令被替换),可用字符`\`对命令进行转义,而使我们可以执行原本的命令.

## 格式

    alias new_command='command seq'
    unalias new_command

## 例

    # alias ll='ls -alF'
    # ll
