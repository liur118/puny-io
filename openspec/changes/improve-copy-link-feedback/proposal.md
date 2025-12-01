# Change: Improve Copy Link Feedback

## Why
当前在点击"获取链接"按钮复制对象访问链接时，虽然有 message.success 提示，但用户反馈体验不够明显，容易被忽略。需要提供更清晰、更直观的反馈机制，让用户明确知道链接已成功复制到剪贴板。

## What Changes
- 保留现有的 message.success 提示作为基础反馈
- 添加更明显的视觉反馈，可能包括：
  - Modal 对话框显示复制成功及完整链接
  - 或使用 Ant Design 的 notification 组件提供更持久、更显眼的通知
  - 或在按钮上显示短暂的"已复制"状态变化
- 确保在复制失败时也有清晰的错误提示
- 保持现有的降级方案（clipboard API → execCommand → 显示链接）

## Impact
- Affected specs: ui-feedback (新增)
- Affected code:
  - `ui/src/pages/Home.tsx` - handleGetURL 函数的反馈逻辑
- 不涉及后端代码变更
- 不影响现有功能逻辑，仅增强用户体验
