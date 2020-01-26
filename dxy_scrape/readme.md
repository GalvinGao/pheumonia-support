# 抓取丁香园数据

## Usage

```
git clone https://github.com/GalvinGao/pneumoniaSupport
cd pneumoniaSupport/dxy_scrape
go run dxy_scrape.go
```

## Description

- 项目将自动在当前目录下创建 `data` 文件夹
- 项目将每分钟获取一次丁香园api，并与本地数据对比，发现不一致则向 `data` 文件夹内写入 `jsonlines` 格式的数据文件
    - 其中，数据文件命名格式与丁香园html内的 `id` 字段一致
    - 写入的文件中包含四个字段：
    
    | 字段名 | Go 数据类型 | 描述 |
    |-------------------|-------------|----------------------------------------------------------------------------------|
    | `time` | `time.Time` | 获取到此记录的时间 |
    | `session` | `string` | 项目在开始运行前会生成一段16字符长度的随机字符串，以便在多线程情景下区分数据源 |
    | `session_started` | `time.Time` | 项目开始运行的时间 |
    | `content` | `string` | 所获取到的记录内容；请注意本字段并非是json对象，而是序列化为`string`后的记录内容 |  
