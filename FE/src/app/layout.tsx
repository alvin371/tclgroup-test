import { Layout, Menu, Input, Button, Avatar } from 'antd'
import { Outlet } from 'react-router'
import { useRouter } from '@/app/_hooks/router/use-router'
import { ROUTES } from '@/commons/route'
import { AppstoreOutlined, PlusOutlined, UserOutlined, SearchOutlined } from '@ant-design/icons'

const { Header, Content } = Layout

const menuItems = [
  { key: ROUTES.INVENTORY, label: 'Inventory' },
  { key: ROUTES.STOCK_IN, label: 'Stock In' },
  { key: ROUTES.STOCK_OUT, label: 'Stock Out' },
  { key: ROUTES.REPORTS, label: 'Reports' },
]

export default function AppLayout() {
  const router = useRouter()

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header
        style={{
          position: 'sticky',
          top: 0,
          zIndex: 100,
          background: '#fff',
          boxShadow: '0 1px 4px rgba(0,0,0,0.12)',
          display: 'flex',
          alignItems: 'center',
          gap: 16,
          padding: '0 24px',
          height: 64,
        }}
      >
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, flexShrink: 0 }}>
          <AppstoreOutlined style={{ fontSize: 20, color: '#1677ff' }} />
          <span style={{ fontWeight: 700, fontSize: 16, whiteSpace: 'nowrap' }}>StockManager</span>
        </div>

        <Input
          variant="filled"
          prefix={<SearchOutlined />}
          placeholder="Search..."
          style={{ width: 220 }}
        />

        <div style={{ flex: 1 }} />

        <Menu
          mode="horizontal"
          selectedKeys={[router.pathname]}
          items={menuItems}
          onClick={({ key }) => router.push(key as (typeof ROUTES)[keyof typeof ROUTES])}
          style={{ border: 'none', background: 'transparent', minWidth: 360 }}
        />

        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => router.push(ROUTES.STOCK_IN)}
        >
          Add Item
        </Button>

        <Avatar icon={<UserOutlined />} style={{ flexShrink: 0, cursor: 'pointer' }} />
      </Header>

      <Content style={{ background: '#f5f5f5', padding: 24 }}>
        <Outlet />
      </Content>
    </Layout>
  )
}
