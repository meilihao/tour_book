# /home/chen/.config/fish/conf.d
set -gx GOROOT /opt/go
set -gx GOPATH /home/chen/git/go
set -gx PATH $GOROOT/bin $GOPATH/bin $PATH

## /etc/profile
## golang
#export GOROOT=/opt/go
#export GOPATH=/home/chen/git/go
#export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

alias psqltest "psql -h rm-xxx.pg.rds.aliyuncs.com -p 5555 -U postgres -d postgres -W"
