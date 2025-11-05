# 全球行政区划数据

这是一个轻量的 HTTP 服务，提供行政区划数据查询（当前包含日本 `jp` 与中国 `cn` 两份示例数据）。

## 路由与用法

- GET /:country
	- 返回该国家的顶级行政地区。
	- 例如：`GET /jp`。

- GET /:country/:id
	- 返回指定 id 的地区信息（包含 children 字段，若有子项）。
	- 例如：`GET /cn/110000`。

- GET /search_postcode/:country/:postcode
	- 通过邮政编码查找匹配的地区。
	- 例如：`GET /search_postcode/jp/1000001`。

## 示例（使用 curl）

```bash
# 获取日本根节点
curl http(s)://geo-data.easit.cc/jp

# 获取中国 id=110000
curl http(s)://geo-data.easit.cc/cn/110000

# 按邮编查询
curl http(s)://geo-data.easit.cc/search_postcode/jp/1000001
```

## 数据格式说明（CSV）

当前代码期望数据为逗号分隔的 CSV，字段顺序如下：

1. id (字符串，唯一标识)
2. parent (字符串，父节点 id，根节点可为空)
3. name (名称)
4. spell (拼音或读音)
5. address (完整地址)
6. post (邮政编码 / 可选)

## 构建与发布
```bash
export KO_DOCKER_REPO=pch18/geo-data
ko build --platform=linux/amd64,linux/arm64 --sbom=none --bare .
```