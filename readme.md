#Node happy script


## 这是一个nodeping的操作脚本,这样你就可以解放双手,让它帮你完成

- [ ] 优化文件目录
- [ ] 优化代码结构
- [ ] 缩小每个操作时间,能做到更短时间完成任务

###   First step
    修改账号
    const user = "jiangyu.huang@xxx.io"
    const pwd = "xxxxxxx"
###  Second step
    设置好对应的文件excel 生成,和要替换的路径

### Third  steo
    运行脚本
#### 获取 nodeping alert信息,这个剩下的excel格式你自己定义即可
```shell

go run main.go all  //

```
#### 根据你的excel去修改对应的信息
 例如:
```
Target	|Change to Name
 xx	    |bb
```
就会把名字 从 xx - > bb
那么问题来了,如果有其他条件筛选怎么办
搜索代码,搜索这个代码区域,去扩充你想定义的规则
```TODO 这是排除 同一个label(title) 有不同的条件 可以通过以下的方法去排除
```

```shell

go run main.go  update //

```

 