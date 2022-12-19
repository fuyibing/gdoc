# GDOC

文档生成工具 GDOC（Generate Document Client）原理是通过扫描文件中的注解、标签
获取关键信息. 本工具依赖注解(Annotation),标签(Tag), 反射(Reflection).

### 关于依赖

1. 注解
    1. 控制器
        1. `@RoutePrefix(path string)` - 路由前缀
    2. 接口方法
        1. `@Error(code int, message, description string)` - 错误码
        2. `@Header(key, value, description message)` - 请求头
        3. `@Ignore(on bool)` - 忽略文档
        4. `@Request(struct string)` - 请求入参结构体
        5. `@Response(struct string)` - 请求出参结构体/单条数据
        6. `@ResponseData(struct string)` - 请求出参结构体/单条数据
        7. `@ResponseList(struct string)` - 请求出参结构体/数据列列
        8. `@ResponsePage(struct string)` - 请求出参结构体/带分页
        9. `@Version(str string)` - 版本号
2. 标签
    1. `desc` - 字段描述
    2. `form` - 从POST表单中取值
    3. `json` - 从JSON表单中取值
    4. `label` - 字段标题/名称
    5. `mock` - Mock数据值
    6. `url` - 从URL参数中取值
    7. `validate` - 字段格式校验

### 如何使用?

```go
// Example code.

package main

import (
	"github.com/fuyibing/gdoc/adapters/markdown"
	"github.com/fuyibing/gdoc/base"
	"github.com/fuyibing/gdoc/conf"
	"github.com/fuyibing/gdoc/reflectors"
	"github.com/fuyibing/gdoc/scanners"
)

func main() {
	// 设置项目根目录.
	conf.Path.SetBasePath("/home/applications/sketch")

	// 设置控制器目录.
	conf.Path.SetControllerPath("/app/controllers")

	// 加载配置.
	// 
	// 工具的执行时, 将读取此项目的的 gdoc.json 和 config/app.yaml
	// 文件.
	conf.Config.Load()

	// 扫描文件.
	//
	// 扫描控制器目录下的GO源码, 然后从源码中匹配控制器
	// 与接口方法列表.
	scanners.Scanner.Scan()

	// 反射结构.
	reflectors.New(base.Mapper).
		Configure().
		Make()

	// 调用适配器.
	//
	// 下列调用Markdown适配器, 触发后将把扫描结果与反射结果
	// 格式化并生成MD文件.
	markdown.New(base.Mapper).Run()
}
```

