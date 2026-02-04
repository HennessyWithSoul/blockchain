---
name: Issue Branch Creator
type: skill
version: 1.0.0
agent: CodeActAgent
triggers: ["新建issue", "创建issue", "new issue", "create issue", "issue创建", "issue新建"]
---

# Issue Branch Creator Microagent

## 目的
当用户在对话中提到"新建issue"时，自动从main分支创建一个新的开发分支，分支名为 `dev/{当前时间}`。

## 功能
当用户提到新建issue时，本microagent会：
1. 从main分支创建一个新的开发分支
2. 使用格式 `dev/{当前时间}` 命名分支
3. 确保分支名称唯一，避免冲突
4. 自动执行分支创建命令

## 执行流程
当触发器被激活时（用户提到"新建issue"等关键词）：

### 1. 检查当前git状态
```bash
git status
git branch -a
```

### 2. 生成时间戳和分支名
```bash
TIMESTAMP=$(date +"%Y-%m-%d-%H-%M-%S")
BRANCH_NAME="dev/$TIMESTAMP"
echo "将创建分支: $BRANCH_NAME"
```

### 3. 从main分支创建新分支
```bash
# 确保在main分支上
git checkout main
# 拉取最新代码
git pull origin main
# 创建并切换到新分支
git checkout -b $BRANCH_NAME
```

### 4. 推送到远程仓库
```bash
git push -u origin $BRANCH_NAME
echo "分支 $BRANCH_NAME 已创建并推送到远程仓库"
```

### 5. 验证分支创建
```bash
git branch -a | grep $BRANCH_NAME
```

## 完整执行命令
```bash
# 完整的一次性命令
TIMESTAMP=$(date +"%Y-%m-%d-%H-%M-%S") && BRANCH_NAME="dev/$TIMESTAMP" && git checkout main && git pull origin main && git checkout -b $BRANCH_NAME && git push -u origin $BRANCH_NAME && echo "✅ 分支 $BRANCH_NAME 创建成功！"
```

## 错误处理
- 如果main分支不存在，提示用户
- 如果网络连接失败，重试或提示用户手动操作
- 如果分支已存在，使用不同的时间戳

## 使用示例
**用户说**: "我新建了一个issue"
**AI助手响应**: "检测到您新建了issue，正在为您创建开发分支..."
然后执行上述命令创建分支。

**用户说**: "创建issue后需要新分支"
**AI助手响应**: "好的，我将从main分支为您创建一个新的开发分支。"
然后执行上述命令创建分支。

## 注意事项
1. 需要git已配置且有权访问远程仓库
2. 使用系统当前时间（UTC时区）
3. 分支格式：dev/2026-02-04-08-15-30
4. 确保有main分支的访问权限