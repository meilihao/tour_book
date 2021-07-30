# base64
base64编解码

## 选项
- -d : 解码

## example

    ```bash
    $ echo "Hello World" | base64
    SGVsbG8gV29ybGQK

    $ echo SGVsbG8gV29ybGQK | base64 -d
    Hello World
    ```