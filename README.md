# API Workbench

React + Go 全栈项目实战

## 技术栈

- **前端**: React 18 + TypeScript + Vite
- **后端**: Go + Gin

## 项目结构

```
api-workbench/
├── frontend/          # React 前端
│   ├── src/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
└── backend/           # Go 后端
    ├── go.mod
    └── main.go
```

## 快速开始

### 后端

```bash
cd backend
go mod tidy
go run main.go
```

后端运行在 http://localhost:8080

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端运行在 http://localhost:3000

## API 接口

- `GET /api/health` - 健康检查
- `GET /api/items` - 获取项目列表
