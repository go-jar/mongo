# 思路

- query.go：提供 mongo 查询语句格式
- client.go：将查询语句转换为 mongo 格式并执行；执行其它 crud 语句
- dao.go：借助 client，提供方便操作 mongo 的方法
- config.go：配置 client

# 参考

- https://github.com/goinbox/mongo
