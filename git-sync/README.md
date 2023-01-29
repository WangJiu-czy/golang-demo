### 为什么使用git-sync

---

- 0个外部依赖(git除外)
- 每隔一段时间,自动备份
- 基于JSON配置
- 合理的默认值

---

#### 安装git-sync
> 要求有go环境

```shell
git clone https://github.com/WangJiu-czy/golang-demo/tree/main/git-sync
cd git-sync
go build
```
```shell
./git-sync #unix
git-sync.exe #windows
```

#### 运行git-sync
先决条件:
- 需要安装git,否则git-sync会panic
    1. 项目需要是具有远程设置的git存储库
    2. git用户需要通过远程身份验证
    3. 在使用git-sync之前,您应该能够在项目中毫无问题地运行以下命令:
        - `git add -A`
        - `git commit -m "test"`
        - `git push`
    4. 你现在可以在你的项目中使用git-sync
1. 切换到要备份的git项目
2. 在终端运行 git-sync
> 如果没有`gitsync.json`,git-sync将使用其默认配置

#### 配置路径

- 在Unix系统上,`$HOME/.config/gitsync.json`
- 在Windows上, `%AppData%/gitsync.json`

#### 配置选项和默认值
如果git-sync找不到其配置文件(`gitsync.json`),它将回退到其默认配置
```json5
{
  //将在提交title前面添加前缀, 默认:"backup: %date%"
  "commit_format":     "backup:%date%",
  
  //指定日期格式，默认为："2023-01-01 15:15:15"
  "commit_date":       "2023-01-01 15:15:15",
  
  //列出提交主体中受提交影响的文件名，默认值：true
  "add_affected_files": true,
  
  //备份之间的时间间隔（以秒为单位），默认值：300
  "backup_interval":   300,
  
  //提交命令，默认："git commit -m"
  "commit_command":    "git commit -m",
  
  //启用调试模式（详细日志记录、额外信息等），默认值：false
  "debug":        false,
  
  //启动时启用从远程拉取，默认值：true
  "pull_on_start":      true
	}

```
