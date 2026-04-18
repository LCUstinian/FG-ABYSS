import { createI18n } from 'vue-i18n'

const messages = {
  'zh-CN': {
    nav: {
      home:    '首页',
      project: '项目',
      payload: '载荷',
      plugin:  '插件',
      settings: '设置',
    },
    app: {
      name: '渊渟',
      sub:  'ABYSS',
    },
    common: {
      loading:  '加载中...',
      empty:    '暂无数据',
      confirm:  '确认',
      cancel:   '取消',
      save:     '保存',
      delete:   '删除',
      edit:     '编辑',
      create:   '新建',
      search:   '搜索',
      refresh:  '刷新',
      copy:     '复制',
      download: '下载',
    },
  },
  'en-US': {
    nav: {
      home:    'Home',
      project: 'Projects',
      payload: 'Payload',
      plugin:  'Plugins',
      settings: 'Settings',
    },
    app: {
      name: 'FG',
      sub:  'ABYSS',
    },
    common: {
      loading:  'Loading...',
      empty:    'No data',
      confirm:  'Confirm',
      cancel:   'Cancel',
      save:     'Save',
      delete:   'Delete',
      edit:     'Edit',
      create:   'New',
      search:   'Search',
      refresh:  'Refresh',
      copy:     'Copy',
      download: 'Download',
    },
  },
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en-US',
  messages,
})
