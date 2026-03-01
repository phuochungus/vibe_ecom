import type { ThemeConfig } from 'antd'

export const shopeeTheme: ThemeConfig = {
  token: {
    colorPrimary: '#EE4D2D', // Shopee orange/red
    colorSuccess: '#26AA99', // Green for success states
    colorWarning: '#FFBF00', // Yellow for warnings
    colorError: '#FF424F', // Red for errors
    colorInfo: '#1890FF', // Blue for info
    borderRadius: 4,
    fontSize: 14,
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
  },
  components: {
    Button: {
      controlHeight: 40,
      fontWeight: 500,
      primaryShadow: '0 2px 0 rgba(238, 77, 45, 0.1)',
    },
    Card: {
      borderRadiusLG: 8,
      boxShadowTertiary: '0 1px 4px 0 rgba(0,0,0,.08)',
    },
    Input: {
      controlHeight: 40,
      borderRadius: 2,
    },
    InputNumber: {
      controlHeight: 40,
    },
    Select: {
      controlHeight: 40,
    },
    Layout: {
      headerBg: '#fff',
      headerHeight: 64,
      headerPadding: '0 24px',
    },
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#fff1ec',
      itemSelectedColor: '#EE4D2D',
      itemHoverBg: '#fff1ec',
      itemHoverColor: '#EE4D2D',
    },
    Tag: {
      borderRadiusSM: 2,
    },
    Badge: {
      dotSize: 8,
    },
  },
}
