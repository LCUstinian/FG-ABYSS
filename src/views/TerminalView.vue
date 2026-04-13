<template>
  <div class="terminal-container">
    <n-card :title="t('terminal')">
      <div ref="terminalRef" class="terminal"></div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { NCard } from 'naive-ui'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'

const { t } = useI18n()
const terminalRef = ref<HTMLElement | null>(null)

onMounted(() => {
  if (terminalRef.value) {
    const terminal = new Terminal({
      theme: {
        background: '#1e1e1e',
        foreground: '#ffffff',
      },
      fontSize: 14,
      fontFamily: 'Consolas, "Courier New", monospace',
    })

    const fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)
    terminal.open(terminalRef.value)
    fitAddon.fit()

    terminal.writeln('Welcome to FG-ABYSS Terminal')
    terminal.writeln('Type your commands...')
  }
})
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  padding: 16px;
  box-sizing: border-box;
}

.terminal {
  width: 100%;
  height: calc(100% - 100px);
}
</style>
