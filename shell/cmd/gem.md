# gem

```bash
# 查看gem使用的网络源
gem sources -l
# 删除某个网络源
gem sources --remove https://rubygems.org/
# 添加国内的网络源
gem sources --add https://gems.ruby-china.org/
# 查看源上所有的软件包
gem query --remote
# 根据名字匹配包查看，可以匹配正则
gem query --remote --name-matches '^redis$'
# 查看详细信息
gem query --remote --name-matches '^redis$' -d
````

## FAQ
### gem install 报错`unable to convert xxx to US-ASCII for xxx, skipping`
```
# gem install ffi -v 1.9.10
# unable to convert "\xE2" to UTF-8 in conversion from ASCII-8BIT to UTF-8 to US-ASCII for lib/ffi/library.rb, skipping
# gem install fpm -v 1.3.3
# unable to convert U+91CE from UTF-8 to US-ASCII for lib/backports/tools/normalize.rb, skipping
# gem install ffi -v 1.9.10 --no-rdoc --no-ri # gem install生成文档时报错, 不生成文档即可.
```

### 构建离线部署
```bash
# gem install bundle
# mkdir gem_cache && cd gem_cache
# cat << EOF > Gemfile
source 'https://rubygems.org'
gem 'sassc'
gem 'ffi', '1.9.10'
# bundle cache # 构建出gem cache, 在当前目录的`vendor/cache`下
# bundle install --local # 会安装Gemfile中全部组件
# cd vendor/cache && gem install --local ffi* # 逐个安装组件
```

Gemfile中必须有source.

### `bundle cache`报`Could not find gem 'childprocess (= 0.5.6)' in any of the gem sources listed in your Gemfile`
不知咋回事, 执行`bundle cache`报错, 可用`bundle install && bundle package`

执行`bundle install && bundle package`后, `bundle cache`又可用了.
