# PFlow 平台

PFlow 是一个基于 Go 与 React 构建的流程引擎 & 工单管理平台样板，结合 Camunda BPMN2.0、PostgreSQL 与 RabbitMQ，提供拖拽式流程建模、工单生命周期管理以及 OpenAPI 能力，便于在业务系统中快速集成高并发、可编排的工单流转能力。

## 项目结构

```
├── backend        # Go 服务：配置、流程/工单领域服务、Camunda 与队列集成
│   ├── cmd        # 入口程序
│   ├── internal   # 领域逻辑、HTTP Handler、存储与外部依赖
│   └── go.mod
└── frontend       # React + Vite 前端，提供可视化建模与工单面板
```

## 后端特性

- **BPMN2.0 对接**：通过 `internal/camunda` 与 Camunda 引擎交互，完成流程部署、实例启动与重试。
- **持久化层**：使用 PostgreSQL 存储流程定义与工单实例，提供迁移脚本 `internal/persistence/migrations/0001_init.sql`。
- **消息队列**：基于 RabbitMQ 推送流程/工单事件，便于与外部系统集成或构建审计流水。
- **分层架构**：`service` + `repository` + `handler` 分离，接口驱动，有利于替换 Camunda、存储或队列实现。

### 本地运行

1. 准备依赖服务（可使用 docker-compose）：
   - PostgreSQL
   - RabbitMQ
   - Camunda Platform 8 (或 7) REST API
2. 复制配置模板并根据环境调整：

   ```bash
   cp backend/config.example.yaml backend/config.yaml
   ```

3. 启动后端服务：

   ```bash
   cd backend
   go run ./cmd/server
   ```

   > 若网络受限无法下载依赖，可在本地预先配置私有 Go Proxy。

## 前端特性

- 采用 React、MUI 与 React Flow 实现拖拽式流程建模体验。
- 使用 React Query 管理 API 数据，实时刷新工单状态。
- 内置工单创建、失败重试入口，可作为集成在主业务系统的 iframe 或独立页面。

### 本地运行

```bash
cd frontend
npm install
npm run dev
```

前端默认通过 Vite 代理将 `/api` 请求转发到 `http://localhost:8080`，可按需在 `vite.config.ts` 中调整。

## OpenAPI & 扩展

后端默认提供如下 REST 接口（可通过网关或自定义鉴权扩展）：

- `GET /api/flows`：获取流程列表
- `POST /api/flows`：创建/部署流程
- `GET /api/flows/:id` / `PUT /api/flows/:id`
- `GET /api/workorders`：获取工单列表
- `POST /api/workorders`：创建工单实例
- `POST /api/workorders/:id/retry`：重试失败工单

结合 `internal/mq` 可将事件推送给其他系统，或通过 npm 包方式封装前端能力嵌入自有平台。

## 下一步建议

- 接入统一身份认证体系，扩展权限模型。
- 引入 OpenTelemetry 采集端到端链路信息。
- 编写 docker-compose 与 Helm Chart，完善一键部署能力。
