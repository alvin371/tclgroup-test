import { useState, useMemo } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Table,
  Input,
  Select,
  Button,
  Modal,
  Form,
  InputNumber,
  Space,
  Typography,
  Tag,
  message,
} from 'antd'
import { EditOutlined, ExportOutlined, SyncOutlined, SearchOutlined } from '@ant-design/icons'
import { inventoryApi, type TInventoryItem } from '@/modules/inventory'
import { useMutation } from '@/app/_hooks/request/use-mutation'

const { Title, Text } = Typography

function getStockTag(item: TInventoryItem) {
  const { available_stock, physical_stock } = item
  if (available_stock === 0) return <Tag color="red">{available_stock}</Tag>
  if (physical_stock > 0 && available_stock / physical_stock < 0.3)
    return <Tag color="orange">{available_stock}</Tag>
  return <Tag color="green">{available_stock}</Tag>
}

function AdjustStockModal({
  item,
  onClose,
  onSuccess,
}: {
  item: TInventoryItem
  onClose: () => void
  onSuccess: () => void
}) {
  const [form] = Form.useForm()

  const { mutate, isPending } = useMutation({
    mutationFn: (values: { new_qty: number; notes: string }) =>
      inventoryApi.adjust(item.product_id, values),
    onSuccess: () => {
      message.success('Stock adjusted successfully')
      onClose()
      onSuccess()
    },
    onError: (err) => {
      message.error(err.error.message)
    },
  })

  return (
    <Modal
      open
      title="Adjust Stock"
      onCancel={onClose}
      onOk={() => form.submit()}
      confirmLoading={isPending}
      destroyOnClose
    >
      <Space direction="vertical" style={{ width: '100%', marginBottom: 16 }}>
        <Text type="secondary">
          Physical: <strong>{item.physical_stock}</strong> · Reserved:{' '}
          <strong>{item.reserved}</strong> · Available:{' '}
          <strong>{item.available_stock}</strong>
        </Text>
      </Space>
      <Form
        form={form}
        layout="vertical"
        initialValues={{ new_qty: item.physical_stock }}
        onFinish={(values) => mutate(values as { new_qty: number; notes: string })}
      >
        <Form.Item label="New Quantity" name="new_qty" rules={[{ required: true }]}>
          <InputNumber min={0} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item
          label="Notes"
          name="notes"
          rules={[{ required: true, message: 'Notes are required' }]}
        >
          <Input.TextArea rows={3} />
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default function InventoryPage() {
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [stockStatus, setStockStatus] = useState<string | undefined>(undefined)
  const [adjustItem, setAdjustItem] = useState<TInventoryItem | null>(null)
  const queryClient = useQueryClient()

  const { data, isFetching } = useQuery({
    queryKey: ['inventory', page],
    queryFn: () => inventoryApi.list({ page, perPage: 10 }),
  })

  const pagination = data?.pagination

  const filtered = useMemo(() => {
    const items = data?.items ?? []
    return items.filter((item) => {
      const matchSearch =
        !search ||
        item.product_name.toLowerCase().includes(search.toLowerCase()) ||
        item.product_sku.toLowerCase().includes(search.toLowerCase())

      const matchStatus = (() => {
        if (!stockStatus) return true
        if (stockStatus === 'out') return item.available_stock === 0
        if (stockStatus === 'low')
          return (
            item.available_stock > 0 &&
            item.physical_stock > 0 &&
            item.available_stock / item.physical_stock < 0.3
          )
        if (stockStatus === 'in')
          return (
            item.available_stock > 0 &&
            (item.physical_stock === 0 || item.available_stock / item.physical_stock >= 0.3)
          )
        return true
      })()

      return matchSearch && matchStatus
    })
  }, [data?.items, search, stockStatus])

  const columns = [
    { title: 'SKU', dataIndex: 'product_sku', key: 'sku' },
    { title: 'ITEM NAME', dataIndex: 'product_name', key: 'name' },
    {
      title: 'PHYSICAL STOCK',
      dataIndex: 'physical_stock',
      key: 'physical',
      align: 'right' as const,
    },
    {
      title: 'AVAILABLE STOCK',
      key: 'available',
      render: (_: unknown, record: TInventoryItem) => getStockTag(record),
    },
    {
      title: 'ACTIONS',
      key: 'actions',
      render: (_: unknown, record: TInventoryItem) => (
        <Button type="link" icon={<EditOutlined />} onClick={() => setAdjustItem(record)}>
          Adjust Stock
        </Button>
      ),
    },
  ]

  return (
    <>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          marginBottom: 24,
        }}
      >
        <div>
          <Title level={4} style={{ margin: 0 }}>
            Inventory and Stock Levels
          </Title>
          <Text type="secondary">Track and manage your inventory levels</Text>
        </div>
        <Space>
          <Button icon={<ExportOutlined />}>Export</Button>
          <Button icon={<SyncOutlined />}>Sync</Button>
        </Space>
      </div>

      <div style={{ background: '#fff', borderRadius: 8, padding: 24 }}>
        <Space style={{ marginBottom: 16 }}>
          <Input
            placeholder="Search by name or SKU..."
            prefix={<SearchOutlined />}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ width: 280 }}
          />
          <Select
            placeholder="Stock Status"
            allowClear
            style={{ width: 160 }}
            value={stockStatus}
            onChange={(v) => setStockStatus(v)}
            options={[
              { value: 'in', label: 'In Stock' },
              { value: 'low', label: 'Low Stock' },
              { value: 'out', label: 'Out of Stock' },
            ]}
          />
        </Space>

        <Table
          rowKey="id"
          columns={columns}
          dataSource={filtered}
          loading={isFetching}
          pagination={{
            current: page,
            pageSize: 10,
            total: pagination?.total,
            onChange: setPage,
            showSizeChanger: false,
          }}
        />
      </div>

      {adjustItem && (
        <AdjustStockModal
          item={adjustItem}
          onClose={() => setAdjustItem(null)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ['inventory'] })}
        />
      )}
    </>
  )
}
