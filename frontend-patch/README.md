# 前端 3D 功能补丁

将以下文件应用到前端仓库 [UncleBloom/chatroom](https://github.com/UncleBloom/chatroom)

## 新增文件（直接复制到 src/ 对应目录）

```
src/components/model/modelViewer.tsx    ← Three.js 3D模型查看器
src/components/model/modelViewer.scss   ← 样式
src/components/model/generate3D.tsx     ← 生成按钮+弹窗+进度条
src/components/model/generate3D.scss    ← 样式
src/api/threeD.ts                       ← 3D API调用+轮询
```

## 修改文件（替换原文件）

```
src/components/input/input.tsx          ← 加了🎨3D按钮
src/components/message/message.tsx      ← 支持3D模型消息渲染
```

## 安装依赖

```bash
npm install three@0.160.0 @types/three@0.160.0 --save
```

> 注意：必须用 0.160.0 版本，新版 three.js 需要 TypeScript 5.0+

## 已验证

- [x] npm install 成功
- [x] react-scripts build 编译通过
