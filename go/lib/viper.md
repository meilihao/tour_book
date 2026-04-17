# viper
## FAQ
### viper 多YAML优先级加载（后加载覆盖前加载）
```go
type Configs struct {
	path      string
	name      string
	err       error
	vp        *viper.Viper
	configMap map[string]any
}

func (c *Configs) ConfigRead() error {
	c.vp.SetConfigName(config.name)
	c.vp.AddConfigPath(config.path)
	c.vp.SetConfigType("yaml")

	c.configMap = make(map[string]any)

	if err := c.vp.ReadInConfig(); err != nil {
		return err
	}

	c.vp.SetConfigFile(filepath.Join(config.path, "configs.local.yaml"))
	return c.vp.MergeInConfig()
}
```