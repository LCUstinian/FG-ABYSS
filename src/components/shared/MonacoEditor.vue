<template>
  <div ref="editorContainer" class="monaco-editor-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'php'
  },
  readOnly: {
    type: Boolean,
    default: true
  },
  theme: {
    type: String,
    default: 'vs-dark'
  }
})

const emit = defineEmits(['update:value'])

const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

onMounted(() => {
  if (editorContainer.value) {
    editor = monaco.editor.create(editorContainer.value, {
      value: props.value,
      language: props.language,
      readOnly: props.readOnly,
      theme: props.theme,
      minimap: {
        enabled: false
      },
      scrollBeyondLastLine: false,
      automaticLayout: true,
      fontSize: 13,
      lineNumbers: 'on',
      renderLineHighlight: 'all',
      tabSize: 2,
      wordWrap: 'on'
    })

    editor.onDidChangeModelContent(() => {
      if (editor) {
        emit('update:value', editor.getValue())
      }
    })
  }
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
})

watch(() => props.value, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue)
  }
})

watch(() => props.language, (newLanguage) => {
  if (editor) {
    monaco.editor.setModelLanguage(editor.getModel()!, newLanguage)
  }
})
</script>

<style scoped>
.monaco-editor-container {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border-radius: 12px;
  overflow: hidden;
}
</style>