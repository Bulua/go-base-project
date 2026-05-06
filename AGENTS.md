# AGENTS.md

本文件是给各类自动化助手和 Agent 的入口说明。项目长期规则写在 [CLAUDE.md](./CLAUDE.md)，请先阅读并遵守它。

英文版本见 [AGENTS-EN.md](./AGENTS-EN.md)。如中英文内容不一致，以本文件为准。

## 快速上下文

- **项目类型**：Go + Vue 3 后台管理基础项目。
- **后端**：Go 1.22，标准库 `net/http`，MySQL，JWT，bcrypt。
- **前端**：Vue 3，TypeScript，Vite，Element Plus，Pinia，Vue Router，Axios。
- **默认远程依赖启动**：`docker compose up -d --build`
- **完整本地启动**：`docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build`
- **后端测试**：`cd server && go test ./...`
- **前端构建检查**：`cd web && npm run build`
- **前端测试**：`cd web && npm run test`
- **默认演示账号**：`admin` / `Admin@123456`

## Agent 工作规则

1. 修改任何文件前，先读 [CLAUDE.md](./CLAUDE.md)。
2. 不要把长期规则复制到本文件，避免和 `CLAUDE.md` 分叉。
3. 不要提交或写入真实密码、Token、密钥。
4. 不要让正式代码、脚本、Docker 配置或可提交文档依赖 `development-plan/`。
5. 工作区可能已有他人改动，修改前先确认状态，不要回退自己没有创建的改动。
