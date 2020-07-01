# python
## vsocde
参考:
- [Using Python environments in VS Code](https://code.visualstudio.com/docs/python/environments#_environment-variable-definitions-file)

`${project}/.vscode/settings.json`:
```json
{
	"python.envFile": "${workspaceFolder}/.vscode/.env",
	"python.pythonPath":"/usr/bin/python3"
}
```

`${workspaceFolder}/.vscode/.env`:
```conf
PYTHONPATH=/xxx1:/xxx2
```
