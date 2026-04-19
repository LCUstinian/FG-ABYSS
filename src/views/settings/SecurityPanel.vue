<template>
  <div class="security-panel">
    <div class="settings-card">
      <h3 class="card-title">主密码</h3>
      <div class="field-row">
        <span class="field-desc">主密码用于加密本地敏感数据，遗忘后无法恢复</span>
        <n-button size="small" @click="showPasswordDrawer = true">修改主密码</n-button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">密钥轮换</h3>
      <div class="field-row">
        <div style="flex:1">
          <div class="field-label">轮换周期</div>
          <div class="field-desc">上次轮换：从未</div>
        </div>
        <n-button size="small" type="primary" ghost>立即轮换</n-button>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="card-title">审计日志</h3>
      <div class="field-row">
        <span style="flex:1;font-size:13px;color:var(--text-2)">记录所有操作到本地日志</span>
        <n-switch v-model:value="auditLog" size="small" />
      </div>
    </div>

    <!-- Password change drawer placeholder -->
    <n-drawer v-model:show="showPasswordDrawer" :width="380" placement="right">
      <n-drawer-content title="修改主密码">
        <div class="drawer-form">
          <div class="form-group">
            <label class="field-label">当前密码</label>
            <n-input type="password" placeholder="输入当前密码" size="small" style="font-family:var(--font-mono)" />
          </div>
          <div class="form-group">
            <label class="field-label">新密码</label>
            <n-input type="password" placeholder="输入新密码" size="small" style="font-family:var(--font-mono)" />
          </div>
          <div class="form-group">
            <label class="field-label">确认新密码</label>
            <n-input type="password" placeholder="再次输入新密码" size="small" style="font-family:var(--font-mono)" />
          </div>
          <n-button type="primary" block style="margin-top:8px">确认修改</n-button>
        </div>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const auditLog = ref(true)
const showPasswordDrawer = ref(false)
</script>

<style scoped>
.security-panel { display: flex; flex-direction: column; gap: 16px; }

.settings-card {
  background: var(--bg-elevated);
  border: 1px solid rgba(255,255,255,0.08);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

[data-theme="light"] .settings-card {
  border: none;
  box-shadow: 0 1px 4px rgba(0,0,0,0.07), 0 0 0 1px rgba(0,0,0,0.05);
}

.card-title { font-size: 14px; font-weight: 600; color: var(--text-1); }
.field-label { font-size: 13px; font-weight: 500; color: var(--text-2); margin-bottom: 4px; }
.field-desc { font-size: 12px; color: var(--text-3); }
.field-row { display: flex; align-items: center; gap: 12px; justify-content: space-between; }

.drawer-form { display: flex; flex-direction: column; gap: 16px; padding: 4px 0; }
.form-group { display: flex; flex-direction: column; gap: 6px; }
</style>
