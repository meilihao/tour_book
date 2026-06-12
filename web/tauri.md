# tauri

## 调试
- `ctrl + shift + i`/右键->`检查元素`: 打开调试工具

## FAQ
### pnpm tauri dev窗口白屏, 但chrome访问http://localhost:1420正常
执行一次`pnpm tauri dev --verbose`后, 后续执行`pnpm tauri dev`等正常了. 但之后`pnpm tauri dev`又出现白屏, 且`pnpm tauri dev --verbose`也是白屏.
右键->`重新加载`后恢复正常