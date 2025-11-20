# Project Context

## Purpose
puny-io 是一个简单的对象存储系统（OSS），提供类似 AWS S3 的文件存储服务。项目包含完整的前后端实现，支持桶管理、对象上传下载、认证等核心功能。

## Tech Stack

### Backend
- **Go 1.23** - 主要编程语言
- **Gin** - Web 框架
- **JWT** - 用户认证 (golang-jwt/jwt/v5)
- **Viper** - 配置管理
- **File System** - 本地文件存储

### Frontend  
- **React 19** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Ant Design** - UI 组件库
- **React Router** - 路由管理
- **Axios** - HTTP 客户端

## Project Conventions

### Code Style
- **Go**: 遵循 Go 官方代码规范，使用 gofmt 格式化
- **TypeScript**: 使用 ESLint 进行代码检查
- **命名**: 使用 kebab-case 命名文件和目录，camelCase 命名变量和函数
- **包结构**: 按功能模块组织代码 (`internal/handler`, `internal/service`, `internal/storage`)

### Architecture Patterns
- **Clean Architecture**: 分层架构，handler -> service -> storage
- **依赖注入**: 通过构造函数注入依赖
- **接口抽象**: 存储层使用接口，便于扩展不同存储后端
- **中间件模式**: 使用 Gin 中间件处理认证等横切关注点
- **RESTful API**: 遵循 REST 设计原则

### Testing Strategy
- 单元测试覆盖核心业务逻辑
- 集成测试验证 API 接口
- 前端组件测试使用 React Testing Library

### Git Workflow
- **main** 分支为稳定版本
- feature 分支进行功能开发
- 使用有意义的提交信息，遵循 Conventional Commits

## Domain Context
- **Bucket**: 存储桶，用于组织文件的容器
- **Object**: 对象，存储在桶中的文件
- **Key**: 对象的唯一标识符（路径）
- **Metadata**: 对象的元数据信息（文件类型、大小、修改时间等）
- **ETag**: 对象的哈希值，用于完整性校验
- **ACL**: 访问控制，目前通过 JWT 实现简单的用户认证

## Important Constraints
- **单用户系统**: 目前只支持管理员账户，配置在 config.yaml 中
- **本地存储**: 仅支持本地文件系统存储，数据存放在 `./data` 目录
- **无数据库**: 元数据通过文件系统管理，无独立数据库
- **简单认证**: 使用硬编码用户账号，JWT token 认证

## External Dependencies
- **无外部服务依赖**: 系统完全自包含
- **静态文件服务**: 前端构建产物通过 Gin 静态文件服务提供
- **配置文件**: 依赖 `conf/config.yaml` 进行系统配置
- **构建脚本**: 使用 `build.sh`, `builddocker.sh`, `buildui.sh` 进行构建
