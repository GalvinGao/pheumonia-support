# 石墨文档 to CSV
很多在做这件事情的小伙伴都在用石墨文档所以来写一个好了；Inspired by [yangkghjh/shimo2csv](https://github.com/yangkghjh/shimo2csv)

## Description
使用了石墨文档的【导出】API 获取 xlsx 文件，解析后转换并保存为 csv 文件

## Usage
- `$ cp config.example.json config.json`
- 在 `config.json` 内填入从石墨文档网站获取的 Cookie
- `$ go run main.go [文档的ID] [保存的文件名称] (可选选项-标签页index值, 默认为0 即第一个标签页)`
- 运行后即可在 `files/` 文件夹下找到生成的csv