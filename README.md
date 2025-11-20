# puny-io

一个简单的oss系统，包含前后端

## 技术栈

### 前端

* react + vite + TS

### 后端

* golang + **gin**

## 目录结构

```
puny-io/                    # 项目根目录
├── main.go                 # 后端入口文件
├── go.mod                  # Go 模块依赖管理
├── Dockerfile              # Docker 构建文件
├── build.sh                # 完整构建脚本
├── buildui.sh              # 前端构建脚本
├── builddocker.sh          # Docker 构建脚本
├── README.md               # 项目说明文档
├── conf/                   # 配置文件目录
│   └── config.yaml         # 主配置文件
├── data/                   # 数据存储目录（运行时创建）
├── internal/               # 后端核心代码（Go 内部包）
│   ├── config/             # 配置管理
│   ├── handler/            # HTTP 处理器
│   │   ├── auth.go         # 认证相关接口
│   │   └── oss.go          # 对象存储接口
│   ├── middleware/         # 中间件
│   │   └── auth.go         # JWT 认证中间件
│   ├── model/              # 数据模型
│   ├── service/            # 业务逻辑服务
│   │   └── jwt.go          # JWT 服务
│   └── storage/            # 存储抽象层
│       ├── storage.go      # 存储接口定义
│       └── file.go         # 文件系统存储实现
├── ui/                     # 前端项目目录
│   ├── package.json        # Node.js 依赖管理
│   ├── vite.config.ts      # Vite 构建配置
│   ├── tsconfig.json       # TypeScript 配置
│   ├── eslint.config.js    # ESLint 代码检查配置
│   ├── index.html          # HTML 模板
│   ├── public/             # 静态资源
│   └── src/                # React 源码
│       ├── main.tsx        # React 应用入口
│       ├── App.tsx         # 主应用组件
│       ├── components/     # 可复用组件
│       ├── pages/          # 页面组件
│       │   ├── Login.tsx   # 登录页面
│       │   └── Home.tsx    # 主页面
│       └── services/       # 前端服务层
│           ├── api.ts      # API 请求封装
│           └── auth.ts     # 认证服务
└── build/                  # 构建输出目录
    └── puny-io/            # 完整构建包
        ├── bin/            # 后端可执行文件
        ├── conf/           # 配置文件
        └── ui/             # 前端构建产物
```

**架构说明：**
- **前后端分离设计**: 前端在 `ui/` 目录独立开发，后端在根目录
- **静态文件集成**: 生产环境下，前端构建产物由 Go 后端直接服务
- **分层架构**: 后端采用 handler -> service -> storage 的清晰分层
- **接口抽象**: 存储层定义接口，便于扩展不同存储后端

## 本地开发

### 环境要求

- **Go**: 1.23 或更高版本
- **Node.js**: 20 或更高版本  
- **npm**: 用于前端依赖管理

### 开发流程

#### 1. 克隆项目
```bash
git clone https://github.com/liur118/puny-io.git
cd puny-io
```

#### 2. 后端开发

```bash
# 安装 Go 依赖
go mod tidy

# 启动后端服务（开发模式）
go run main.go
```

后端服务默认在 `http://localhost:8080` 启动

**后端开发说明：**
- 配置文件位于 `conf/config.yaml`
- 数据存储在 `./data` 目录（自动创建）
- 支持热重载：使用 `air` 工具可实现代码变更自动重启

#### 3. 前端开发

```bash
# 进入前端目录
cd ui

# 安装前端依赖
npm install

# 启动前端开发服务器（支持热重载）
npm run dev
```

前端开发服务器默认在 `http://localhost:5173` 启动

**前端开发说明：**
- Vite 提供热重载和快速构建
- 开发时前端独立运行，通过代理访问后端 API
- TypeScript 提供类型检查和更好的开发体验

#### 4. 全栈开发

同时运行前后端服务进行全栈开发：

```bash
# 终端1：启动后端
go run main.go

# 终端2：启动前端
cd ui && npm run dev
```

**API 代理配置：** 前端开发时，Vite 会将 `/api` 请求代理到后端服务

#### 5. 构建部署

```bash
# 构建前端项目
./buildui.sh

# 完整构建（包含前后端）
./build.sh

# 构建 Docker 镜像
./builddocker.sh
```

**构建流程：**
1. `buildui.sh` - 构建 React 项目生成静态文件
2. `build.sh` - 编译 Go 后端，整合前端静态文件，打包完整应用
3. 生产环境下，Go 服务器直接服务前端静态文件

#### 6. Docker 部署

```bash
# 构建镜像
docker build -t puny-io .

# 运行容器
docker run -p 8080:8080 -v ./data:/app/data puny-io
```

**多阶段构建：**
- 第一阶段：Node.js 环境构建前端
- 第二阶段：Go 环境编译后端，整合前端构建产物
- 第三阶段：Alpine 精简镜像运行应用

### 开发调试

- **后端日志**: 使用 `gin.Default()` 提供请求日志
- **前端调试**: 浏览器开发者工具，支持 React DevTools
- **API 测试**: 推荐使用 Postman 或 curl 测试后端接口
- **热重载**: 前端自动热重载，后端可配置 air 实现热重载