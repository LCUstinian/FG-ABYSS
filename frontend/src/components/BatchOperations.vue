<template>
  <div class="batch-operations">
    <n-card title="批量操作工具" :bordered="false">
      <n-tabs type="line" animated>
        <!-- 批量导入 -->
        <n-tab-pane name="import" tab="批量导入">
          <n-space vertical>
            <n-alert type="info" title="支持的格式">
              支持 JSON、CSV、XML 格式的 WebShell 数据导入
            </n-alert>

            <n-form label-placement="top">
              <n-form-item label="导入格式">
                <n-select
                  v-model:value="importFormat"
                  :options="formatOptions"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="数据内容">
                <n-input
                  v-model:value="importData"
                  type="textarea"
                  :rows="10"
                  placeholder="粘贴 JSON/CSV/XML 数据，或点击下方按钮下载模板"
                  style="font-family: monospace"
                />
              </n-form-item>

              <n-space>
                <n-button type="primary" @click="handleImport" :loading="importing">
                  <template #icon>
                    <n-icon :component="CloudUploadOutline" />
                  </template>
                  导入
                </n-button>
                <n-button @click="handleDownloadTemplate">
                  <template #icon>
                    <n-icon :component="DocumentTextOutline" />
                  </template>
                  下载模板
                </n-button>
                <n-button @click="handleValidateImport">
                  <template #icon>
                    <n-icon :component="CheckmarkCircleOutline" />
                  </template>
                  验证数据
                </n-button>
              </n-space>
            </n-form>

            <n-divider />

            <n-space v-if="importResult" vertical>
              <n-result
                :status="importResult.success > 0 ? 'success' : 'error'"
                :title="`导入完成 - 成功 ${importResult.success}/${importResult.total}`"
                :description="importResult.errors.join(', ')"
              />
            </n-space>
          </n-space>
        </n-tab-pane>

        <!-- 批量导出 -->
        <n-tab-pane name="export" tab="批量导出">
          <n-space vertical>
            <n-alert type="info" title="导出说明">
              选择要导出的 WebShell，支持 JSON、CSV、XML 格式
            </n-alert>

            <n-form label-placement="top">
              <n-form-item label="导出格式">
                <n-select
                  v-model:value="exportFormat"
                  :options="formatOptions"
                  style="width: 200px"
                />
              </n-form-item>

              <n-form-item label="选择 WebShell">
                <n-select
                  v-model:value="selectedIds"
                  :options="webshellOptions"
                  multiple
                  placeholder="选择要导出的 WebShell"
                  style="width: 400px"
                />
              </n-form-item>

              <n-form-item label="输出文件名">
                <n-input
                  v-model:value="exportFilename"
                  placeholder="webshells_export"
                  style="width: 300px"
                />
              </n-form-item>

              <n-space>
                <n-button type="primary" @click="handleExport" :loading="exporting">
                  <template #icon>
                    <n-icon :component="CloudDownloadOutline" />
                  </template>
                  导出
                </n-button>
              </n-space>
            </n-form>
          </n-space>
        </n-tab-pane>

        <!-- 批量删除 -->
        <n-tab-pane name="delete" tab="批量删除">
          <n-space vertical>
            <n-alert type="warning" title="危险操作">
              批量删除操作不可恢复，请谨慎操作
            </n-alert>

            <n-form label-placement="top">
              <n-form-item label="选择 WebShell">
                <n-select
                  v-model:value="deleteIds"
                  :options="webshellOptions"
                  multiple
                  placeholder="选择要删除的 WebShell"
                  style="width: 400px"
                />
              </n-form-item>

              <n-space>
                <n-button
                  type="error"
                  @click="handleBatchDelete"
                  :loading="deleting"
                >
                  <template #icon>
                    <n-icon :component="TrashOutline" />
                  </template>
                  批量删除
                </n-button>
              </n-space>
            </n-form>

            <n-space v-if="deleteResult" vertical>
              <n-result
                status="success"
                :title="`删除完成 - 成功 ${deleteResult.success}/${deleteResult.total}`"
              />
            </n-space>
          </n-space>
        </n-tab-pane>

        <!-- 批量测试 -->
        <n-tab-pane name="test" tab="批量测试">
          <n-space vertical>
            <n-alert type="info" title="批量测试">
              批量测试选中的 WebShell 连接状态
            </n-alert>

            <n-form label-placement="top">
              <n-form-item label="选择 WebShell">
                <n-select
                  v-model:value="testIds"
                  :options="webshellOptions"
                  multiple
                  placeholder="选择要测试的 WebShell"
                  style="width: 400px"
                />
              </n-form-item>

              <n-space>
                <n-button
                  type="primary"
                  @click="handleBatchTest"
                  :loading="testing"
                >
                  <template #icon>
                    <n-icon :component="PulseOutline" />
                  </template>
                  批量测试
                </n-button>
              </n-space>
            </n-form>

            <n-space v-if="testResult" vertical>
              <n-result
                :status="testResult.success > 0 ? 'success' : 'error'"
                :title="`测试完成 - 成功 ${testResult.success}/${testResult.total}`"
                :description="testResult.errors.join(', ')"
              />
            </n-space>
          </n-space>
        </n-tab-pane>
      </n-tabs>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useMessage } from 'naive-ui'
import {
  CloudUploadOutline,
  CloudDownloadOutline,
  DocumentTextOutline,
  CheckmarkCircleOutline,
  TrashOutline,
  PulseOutline,
} from '@vicons/ionicons5'
import {
  Import,
  Export,
  BatchDelete,
  BatchTest,
  GetImportTemplate,
  ValidateImport,
} from '@bindings/BatchHandler'

const message = useMessage()

const importFormat = ref('json')
const importData = ref('')
const importing = ref(false)
const importResult = ref<any>(null)

const exportFormat = ref('json')
const selectedIds = ref<number[]>([])
const exportFilename = ref('webshells_export')
const exporting = ref(false)

const deleteIds = ref<number[]>([])
const deleting = ref(false)
const deleteResult = ref<any>(null)

const testIds = ref<number[]>([])
const testing = ref(false)
const testResult = ref<any>(null)

const formatOptions = [
  { label: 'JSON', value: 'json' },
  { label: 'CSV', value: 'csv' },
  { label: 'XML', value: 'xml' },
]

const webshellOptions = [
  { label: '示例 Shell 1', value: 1 },
  { label: '示例 Shell 2', value: 2 },
  { label: '示例 Shell 3', value: 3 },
]

const handleImport = async () => {
  if (!importData.value) {
    message.warning('请输入导入数据')
    return
  }

  try {
    importing.value = true
    const response = await Import({
      data: importData.value,
      format: importFormat.value,
      file_type: 'webshell',
    })

    importResult.value = response
    message.success(`导入完成：成功 ${response.success}/${response.total}`)
  } catch (error: any) {
    message.error('导入失败：' + (error.message || '未知错误'))
  } finally {
    importing.value = false
  }
}

const handleValidateImport = async () => {
  if (!importData.value) {
    message.warning('请输入导入数据')
    return
  }

  try {
    const response = await ValidateImport({
      data: importData.value,
      format: importFormat.value,
    })

    if (response.success > 0) {
      message.success(`验证通过：${response.success} 条有效数据`)
    } else {
      message.error(`验证失败：${response.errors.join(', ')}`)
    }
  } catch (error: any) {
    message.error('验证失败')
  }
}

const handleDownloadTemplate = async () => {
  try {
    const response = await GetImportTemplate({
      format: importFormat.value,
    })

    const blob = new Blob([response.template], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `import_template.${importFormat.value}`
    a.click()
    URL.revokeObjectURL(url)

    message.success('模板下载成功')
  } catch (error: any) {
    message.error('下载模板失败')
  }
}

const handleExport = async () => {
  if (selectedIds.value.length === 0) {
    message.warning('请选择要导出的 WebShell')
    return
  }

  try {
    exporting.value = true
    const response = await Export({
      ids: selectedIds.value,
      format: exportFormat.value,
      filename: `${exportFilename.value}.${exportFormat.value}`,
      file_type: 'webshell',
    })

    if (response.success) {
      message.success('导出成功')
    } else {
      message.error('导出失败：' + response.message)
    }
  } catch (error: any) {
    message.error('导出失败')
  } finally {
    exporting.value = false
  }
}

const handleBatchDelete = async () => {
  if (deleteIds.value.length === 0) {
    message.warning('请选择要删除的 WebShell')
    return
  }

  try {
    deleting.value = true
    const response = await BatchDelete({
      ids: deleteIds.value,
    })

    deleteResult.value = response
    message.success(`删除完成：成功 ${response.success}/${response.total}`)
  } catch (error: any) {
    message.error('删除失败')
  } finally {
    deleting.value = false
  }
}

const handleBatchTest = async () => {
  if (testIds.value.length === 0) {
    message.warning('请选择要测试的 WebShell')
    return
  }

  try {
    testing.value = true
    const response = await BatchTest({
      ids: testIds.value,
    })

    testResult.value = response
    if (response.success > 0) {
      message.success(`测试完成：${response.success} 个连接正常`)
    } else {
      message.error(`测试失败：${response.errors.join(', ')}`)
    }
  } catch (error: any) {
    message.error('测试失败')
  } finally {
    testing.value = false
  }
}
</script>

<style scoped lang="scss">
.batch-operations {
  padding: 20px;
}
</style>
