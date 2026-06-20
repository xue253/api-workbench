---
AIGC:
    Label: "1"
    ContentProducer: 001191440300708461136T1XGW3
    ProduceID: d1cc4a2bedeea9339be462491361c702_03170b4b6cae11f1aa625254006c9bbf
    ReservedCode1: nWrSXTs5vxfXeGL+3kvGmdJOlUdFYydP/05g5P8JF4jZv6Wsch1VSitIn8LqPqSpXIiU/4meLXfzSfLHEm3DvT4hma8wbhmzgxWdxDic60MNqkggLzHadZUFYT4VDp04PDe38xJf9slPoJwuLoWVHCaNE+RYQ9vtqBWDSnruuso8iHDsorh1uovd2HY=
    ContentPropagator: 001191440300708461136T1XGW3
    PropagateID: d1cc4a2bedeea9339be462491361c702_03170b4b6cae11f1aa625254006c9bbf
    ReservedCode2: nWrSXTs5vxfXeGL+3kvGmdJOlUdFYydP/05g5P8JF4jZv6Wsch1VSitIn8LqPqSpXIiU/4meLXfzSfLHEm3DvT4hma8wbhmzgxWdxDic60MNqkggLzHadZUFYT4VDp04PDe38xJf9slPoJwuLoWVHCaNE+RYQ9vtqBWDSnruuso8iHDsorh1uovd2HY=
---

# API Manager —— 接口管理与自动化测试平台

> 版本：v0.1  
> 类型：Go 全栈单体应用  
> 部署：单二进制 + MySQL

---

## 1. 项目概述

API Manager 是一个面向个人/小团队的接口管理与自动化测试平台。支持 HTTP/REST/gRPC/WebSocket 多协议接口的统一管理，提供手动触发与定时调度的自动化测试能力。

---

## 2. 技术栈

| 层 | 技术选型 | 说明 |
|---|---|---|
| 后端框架 | Gin v1.9+ | 路由、中间件、RESTful API |
| ORM | GORM v2 | MySQL 数据持久化 |
| 数据库 | MySQL 8.0+ | 接口定义、用例、测试结果、调度任务等 |
| 配置管理 | Viper | YAML 配置文件 |
| 定时调度 | robfig/cron v3 | 内置 cron 调度，同进程执行测试 |
| HTTP 客户端 | go-resty | 发送 HTTP 测试请求 |
| gRPC 客户端 | google.golang.org/grpc | 反射调用 + JSON 输入 |
| WebSocket 客户端 | gorilla/websocket | WS 连接管理与消息收发 |
| 断言引擎 | expr-lang/expr | 响应断言表达式求值 |
| 前后端通信 | RESTful API + JSON | 开发期 Vite proxy，生产 embed |
| 前端框架 | React 18 + TypeScript | SPA |
| UI 组件 | Ant Design 5 | 表单/表格/布局/图表 |
| 状态管理 | Zustand | 轻量全局状态 |
| HTTP 客户端(前端) | Axios | API 请求 |
| 构建 | Vite | 前端构建成静态资源 |
| 打包 | Go embed | 单二进制发布 |

---

## 3. 核心功能模块

### 3.1 项目管理

- 创建/编辑/删除项目
- 项目列表含名称、描述、环境变量、创建时间
- 每个项目下管理接口集合（类似 Postman Collection）
- 环境变量：定义键值对，请求中 `{{key}}` 占位符自动替换
- 多环境支持：开发/测试/预发布/生产，切换环境时变量值联动

### 3.2 接口库

- 接口按 Collection 分组管理，Collection 支持树形嵌套
- 接口定义字段：

| 字段 | HTTP | gRPC | WebSocket |
|------|------|------|-----------|
| 名称/描述 | 通用 | 通用 | 通用 |
| 协议类型 | HTTP/HTTPS | gRPC | WS/WSS |
| 请求方法 | GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS | - | - |
| URL / 服务地址 | 完整 URL | host:port | ws:// URL |
| 请求头 | 通用 KV | 通用 KV | 通用 KV |
| 路径参数 | 通用 | - | - |
| Query 参数 | 通用 | - | - |
| 请求体 | JSON/Form/XML/Raw/二进制 | JSON（需 proto 定义） | JSON/Text |
| Proto 服务/方法 | - | 服务名.方法名 | - |
| 断言规则 | 多组 通用 | 多组 通用 | 多组 通用 |
| 前置脚本 | JS(goja) | JS(goja) | JS(goja) |
| 后置脚本 | JS(goja) | JS(goja) | JS(goja) |

- 断言规则支持：
  - 状态码断言（等于/不等于）
  - 响应体 JSONPath 断言（等于/包含/正则/存在/不存在/大于/小于）
  - 响应头断言
  - 响应时间断言（毫秒级）
  - 每条规则独立启用/禁用

### 3.3 在线调试

- 选中接口 → 编辑参数 → 点击「发送」
- 实时展示：状态码、耗时、响应头、响应体（JSON格式化高亮）
- 请求历史：保留最近 20 条，点击可回填参数
- 支持的请求体编辑方式：JSON 编辑器、表单 KV 编辑、Raw 文本

### 3.4 测试用例

- 基于接口定义创建测试用例
- 数据驱动：单个用例可绑定多组测试数据集（场景表）
- 数据集以表格方式编辑，每行一组输入+预期输出
- 用例执行顺序可编排（测试套件内顺序）
- 支持套件管理：将多个用例组合为测试套件，套件可按序/并行执行

### 3.5 测试执行

- 手动触发：选择用例/套件 → 点击「运行」
- 定时调度：CRON 表达式配置，到点自动执行
- 执行模式：单个用例、套件内顺序执行、套件内并行执行
- 重试机制：失败用例最多重试 3 次
- 超时控制：单接口超时可配置（默认 30s）
- 环境隔离：每次执行可选择不同环境
- 并发控制：套件并行执行时可配置最大并发数

### 3.6 测试报告

- 实时报告：执行过程中 WebSocket 推送进度
- 汇总统计：通过/失败/跳过、通过率、总耗时
- 用例级明细：每个接口的请求/响应完整记录
- 断言详情：每个断言通过/失败，预期值 vs 实际值
- 历史报告：按时间查看，支持对比两次执行的差异
- 报告导出：Markdown / HTML 格式

### 3.7 变量管理

- 环境变量：项目级别，不同环境不同值
- 全局变量：跨项目共享
- 动态变量：内置函数生成随机值（时间戳、UUID、随机字符串等）
- 接口间变量传递：后置脚本提取响应字段，写入临时变量供后续接口使用
- 变量优先级：临时变量 > 用例变量 > 环境变量 > 全局变量

---

## 4. 数据库设计概要

### 表结构

```
projects
  id, name, description, created_at, updated_at

environments
  id, project_id, name, description, sort_order

environment_variables
  id, environment_id, key, value

collections
  id, project_id, parent_id, name, description, sort_order

apis
  id, collection_id, name, description, protocol(http/grpc/ws)
  method, url, headers(JSON), path_params(JSON), query_params(JSON)
  body_type, body(JSON/text), proto_service, proto_method
  expected_status, timeout_ms, created_at, updated_at

assertions
  id, api_id, target_type(status_code/response_body/response_header/response_time)
  operator, path, expected, enabled

test_cases
  id, project_id, name, description, created_at

test_case_apis
  id, test_case_id, api_id, sort_order

test_data_sets
  id, test_case_api_id, data(JSON), sort_order

test_suites
  id, project_id, name, description, run_mode(sequential/parallel), max_concurrency

test_suite_cases
  id, test_suite_id, test_case_id, sort_order

scheduled_tasks
  id, project_id, target_type(suite), target_id
  cron_expr, enabled, environment_id

test_runs
  id, target_type, target_id, environment_id
  status(running/done/failed), trigger_type(manual/scheduled)
  total, passed, failed, skipped, duration_ms, started_at, finished_at

test_run_details
  id, test_run_id, api_id, test_case_id, data_index
  status, status_code, response_headers(JSON), response_body(text)
  duration_ms, error_message, retry_count, executed_at
```

---

## 5. API 设计概要

### RESTful 路由规范

```
前缀: /api/v1

项目管理:
  GET    /projects              列表
  POST   /projects              创建
  PUT    /projects/:id          更新
  DELETE /projects/:id          删除

环境管理:
  GET    /projects/:pid/environments
  POST   /projects/:pid/environments
  PUT    /environments/:id
  DELETE /environments/:id

环境变量:
  GET    /environments/:eid/variables
  PUT    /environments/:eid/variables   批量更新

Collection:
  GET    /projects/:pid/collections
  POST   /projects/:pid/collections
  PUT    /collections/:id
  DELETE /collections/:id
  POST   /collections/:id/move          移动到其他父节点

接口:
  GET    /collections/:cid/apis
  POST   /collections/:cid/apis
  GET    /apis/:id
  PUT    /apis/:id
  DELETE /apis/:id
  PUT    /apis/:id/assertions           更新断言

在线调试:
  POST   /apis/:id/debug                发送调试请求

测试用例:
  GET    /projects/:pid/test-cases
  POST   /projects/:pid/test-cases
  PUT    /test-cases/:id
  DELETE /test-cases/:id
  PUT    /test-cases/:id/apis           管理用例关联的接口

测试数据集:
  GET    /test-case-apis/:id/datasets
  PUT    /test-case-apis/:id/datasets

测试套件:
  GET    /projects/:pid/test-suites
  POST   /projects/:pid/test-suites
  PUT    /test-suites/:id
  DELETE /test-suites/:id
  PUT    /test-suites/:id/cases

测试执行:
  POST   /test-cases/:id/run
  POST   /test-suites/:id/run
  GET    /test-runs/:id                执行详情
  GET    /test-runs/:id/report         获取报告
  WS     /ws/test-runs/:id             实时进度推送

调度管理:
  GET    /projects/:pid/schedules
  POST   /projects/:pid/schedules
  PUT    /schedules/:id
  DELETE /schedules/:id

变量管理:
  GET    /projects/:pid/global-variables
  PUT    /projects/:pid/global-variables
  GET    /variables/dynamic             获取可用动态函数列表

报告导出:
  GET    /test-runs/:id/export?format=md
  GET    /test-runs/:id/export?format=html
```

---

## 6. 项目目录结构

```
api-manager/
├── cmd/
│   └── server/
│       └── main.go              # 入口
├── internal/
│   ├── config/                  # 配置加载
│   ├── db/                      # 数据库初始化与迁移
│   ├── model/                   # GORM 模型
│   ├── handler/                 # HTTP 处理器
│   ├── service/                 # 业务逻辑层
│   ├── repository/              # 数据访问层
│   ├── engine/                  # 测试执行引擎
│   │   ├── runner.go           # 执行主控
│   │   ├── http_executor.go    # HTTP 执行器
│   │   ├── grpc_executor.go    # gRPC 执行器
│   │   └── ws_executor.go      # WebSocket 执行器
│   ├── scheduler/              # 定时调度
│   ├── assertion/              # 断言引擎
│   ├── script/                 # 前后置脚本执行（goja）
│   ├── variable/               # 变量解析与替换
│   ├── middleware/              # Gin 中间件
│   └── router/                  # 路由注册
├── web/                         # React 前端源码
│   ├── src/
│   │   ├── pages/              # 页面组件
│   │   ├── components/         # 通用组件
│   │   ├── stores/             # Zustand store
│   │   ├── services/           # API 调用封装
│   │   └── utils/
│   ├── index.html
│   ├── vite.config.ts
│   └── package.json
├── config.yaml                  # 配置示例
├── go.mod
├── go.sum
└── Makefile
```

---

## 7. 非功能需求

- 接口响应：API 列表查询 < 200ms，复杂查询 < 500ms
- 并发测试：单次执行至少支持 20 个接口并行
- 数据安全：敏感环境变量入库前 AES 加密存储
- 日志：测试执行详细日志写入文件，支持按日期滚动
- 配置：数据库连接、端口、日志路径等通过 YAML 配置
- 部署：`go build -o api-manager` 生成单二进制，同目录放置 config.yaml 和 web/ 静态资源或 embed 打包
*（内容由AI生成，仅供参考）*
