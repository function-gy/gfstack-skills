:::warning
该功能是实验性特性。使用时应以 `logic` 下的模块划分为主，先梳理模块边界与依赖关系，避免循环依赖。
:::

:::tip
该功能特性从 `v2.1` 版本开始提供。
:::

## 基本介绍

`gf gen service` 用于分析 `internal/logic` 中的业务逻辑模块代码，自动生成 `internal/service` 目录下的接口定义和服务注册代码。

推荐的工程职责：

- `internal/logic/<module>`：具体业务实现，按业务变化原因组织 package。
- `internal/service`：业务能力接口目录，供 controller 或其他业务模块依赖。
- `controller`：处理外部 API 入参、校验、响应组装，不直接承载可复用业务逻辑。

这种模式让业务实现和接口定义分离，模块间通过接口解耦。也支持先手写 `service` 接口，再编码 `logic` 具体实现；手写接口文件不要保留工具生成文件顶部的可覆盖注释。

## 注意事项

- `gf gen service` 根据 `logic` 实现生成 `service` 接口，这不是唯一标准做法，而是官方提供的便捷管理方式。
- 命令默认只解析二级目录下的 Go 文件，例如 `internal/logic/user/*.go`，不会无限递归扫描。
- 不同业务模块中的结构体名称不要冲突，否则生成的 service 接口名可能互相覆盖。
- 结构体嵌套、继承式组合等复杂场景可能无法完整自动生成接口；这类接口应手动维护。
- 业务模块之间不要直接导入其他 `logic` package，优先依赖 `service` 接口，避免形成循环依赖。

## 命令使用

在项目根目录执行：

```text
gf gen service
```

命令帮助：

```text
$ gf gen service -h
USAGE
    gf gen service [OPTION]

OPTION
    -s, --srcFolder         source folder path to be parsed. default: internal/logic
    -d, --dstFolder         destination folder path storing automatically generated go files. default: internal/service
    -f, --dstFileNameCase   destination file name storing automatically generated go files. default: Snake
    -w, --watchFile         used in file watcher, it re-generates all service go files only if given file is under srcFolder
    -a, --stPattern         regular expression matching struct name for generating service. default: ^s([A-Z]\w+)$
    -p, --packages          produce go files only for given source packages
    -i, --importPrefix      custom import prefix to calculate import path for generated importing go file of logic
    -l, --clear             delete all generated go files that are not used any further
    -h, --help              more information about this command

EXAMPLE
    gf gen service
    gf gen service -f Snake
```

如果使用官方工程脚手架并安装了 `make`，也可以执行：

```text
make service
```

## 生成规则

默认 `stPattern` 为 `^s([A-Z]\w+)$`，也就是 `logic` 中小写 `s` 开头、后跟大写字母的结构体会被识别为业务模块实现，并生成对应 service 接口。

示例：

| logic 结构体名称 | service 接口名称 |
| --- | --- |
| `sUser` | `User` |
| `sMetaData` | `MetaData` |

## 开发流程

1. 在 `internal/logic/<module>` 编写业务实现。
2. 确保需要暴露的方法定义在符合规则的业务结构体上。
3. 执行 `gf gen service` 生成或更新 `internal/service` 接口。
4. 在对应 `logic` 模块中完成接口实现注册。
5. 在启动入口按项目约定引入生成的 service 注册代码。

## 自动模式

IDE 可以配置文件监听，在 `logic` 代码变化时自动执行：

```text
gf gen service
```

VS Code 可使用 RunOnSave 类插件，对 `logic` 下的 Go 文件变更触发该命令。

## 官方文档

- https://goframe.org/docs/cli/gen-service
