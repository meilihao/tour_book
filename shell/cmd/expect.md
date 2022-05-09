# expect
实现交互式执行命令

## example
```bash
# ---- 嵌套调用: 如果需要在shell脚本中嵌套expect代码，就要使用expect -c "expect代码"
expect -c "
　　spawn ssh $user_name@$ip_addr df -P
　　expect {
　　　　\"*(yes/no)?\" {send \"yes\r\" ; exp_continue}
　　　　\"*password:\" {send \"$user_pwd\r\" ; exp_continue}
　　　　#退出
　　}
"
```