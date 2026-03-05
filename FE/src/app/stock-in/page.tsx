import { useState, useMemo } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import {
  Table,
  Button,
  Select,
  Modal,
  Form,
  InputNumber,
  Space,
  Typography,
  Tag,
  Popconfirm,
  message,
  Input,
} from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { stockInApi, type TStockIn, type TStockInStatus } from '@/modules/stock-in'
import { productsApi, type TProduct } from '@/modules/products'
import { useMutation } from '@/app/_hooks/request/use-mutation'

const { Title } = Typography

const STATUS_COLOR: Record<TStockInStatus, string> = {
  CREATED: 'blue',
  IN_PROGRESS: 'orange',
  DONE: 'green',
  CANCELLED: 'default',
}

function CreateStockInModal({
  productsMap,
  onClose,
  onSuccess,
}: {
  productsMap: Map<string, TProduct>
  onClose: () => void
  onSuccess: () => void
}) {
  const [form] = Form.useForm()

  const { mutate, isPending } = useMutation({
    mutationFn: stockInApi.create,
    onSuccess: () => {
      message.success('Stock In created')
      onClose()
      onSuccess()
    },
    onError: (err) => message.error(err.error.message),
  })

  const productOptions = Array.from(productsMap.values()).map((p) => ({
    value: p.id,
    label: `${p.name} (${p.sku})`,
  }))

  return (
    <Modal
      open
      title="Create Stock In"
      onCancel={onClose}
      onOk={() => form.submit()}
      confirmLoading={isPending}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={(values) => mutate(values as Parameters<typeof stockInApi.create>[0])}
      >
        <Form.Item label="Product" name="product_id" rules={[{ required: true }]}>
          <Select
            options={productOptions}
            placeholder="Select a product"
            showSearch
            filterOption={(input, option) =>
              (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
            }
          />
        </Form.Item>
        <Form.Item label="Quantity" name="quantity" rules={[{ required: true }]}>
          <InputNumber min={1} style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item label="Notes" name="notes" rules={[{ required: true }]}>
          <Input.TextArea rows={3} />
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default function StockInPage() {
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<TStockInStatus | undefined>(undefined)
  const [showCreate, setShowCreate] = useState(false)
  const queryClient = useQueryClient()

  const { data, isFetching } = useQuery({
    queryKey: ['stock-in', page, statusFilter],
    queryFn: () =>
      stockInApi.list({ page, perPage: 10, ...(statusFilter ? { status: statusFilter } : {}) }),
  })

  const { data: productsData } = useQuery({
    queryKey: ['products'],
    queryFn: () => productsApi.list({ perPage: 100 }),
  })

  const productsMap = useMemo(() => {
    const map = new Map<string, TProduct>()
    for (const p of productsData?.items ?? []) map.set(p.id, p)
    return map
  }, [productsData])

  const { mutate: advance, isPending: isAdvancing } = useMutation({
    mutationFn: ({ id, status }: { id: string; status: 'IN_PROGRESS' | 'DONE' }) =>
      stockInApi.advance(id, { status }),
    onSuccess: () => {
      message.success('Status updated')
      queryClient.invalidateQueries({ queryKey: ['stock-in'] })
    },
    onError: (err) => message.error(err.error.message),
  })

  const { mutate: cancel } = useMutation({
    mutationFn: (id: string) => stockInApi.remove(id),
    onSuccess: () => {
      message.success('Cancelled')
      queryClient.invalidateQueries({ queryKey: ['stock-in'] })
    },
    onError: (err) => message.error(err.error.message),
  })

  const columns = [
    {
      title: 'Product',
      key: 'product',
      render: (_: unknown, r: TStockIn) =>
        productsMap.get(r.product_id)?.name ?? r.product_id,
    },
    { title: 'Quantity', dataIndex: 'quantity', key: 'qty' },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (s: TStockInStatus) => <Tag color={STATUS_COLOR[s]}>{s.replace('_', ' ')}</Tag>,
    },
    { title: 'Notes', dataIndex: 'notes', key: 'notes', ellipsis: true },
    {
      title: 'Created',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (v: string) => new Date(v).toLocaleString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: unknown, r: TStockIn) => {
        if (r.status === 'CREATED') {
          return (
            <Space>
              <Popconfirm
                title="Move to in progress?"
                onConfirm={() => advance({ id: r.id, status: 'IN_PROGRESS' })}
              >
                <Button type="primary" size="small" loading={isAdvancing}>
                  Start
                </Button>
              </Popconfirm>
              <Popconfirm title="Cancel this stock in?" onConfirm={() => cancel(r.id)}>
                <Button danger size="small">
                  Cancel
                </Button>
              </Popconfirm>
            </Space>
          )
        }
        if (r.status === 'IN_PROGRESS') {
          return (
            <Space>
              <Popconfirm
                title="Mark as done (stock received)?"
                onConfirm={() => advance({ id: r.id, status: 'DONE' })}
              >
                <Button type="primary" size="small" loading={isAdvancing}>
                  Receive
                </Button>
              </Popconfirm>
              <Popconfirm title="Cancel this stock in?" onConfirm={() => cancel(r.id)}>
                <Button danger size="small">
                  Cancel
                </Button>
              </Popconfirm>
            </Space>
          )
        }
        return null
      },
    },
  ]

  return (
    <>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 24,
        }}
      >
        <Title level={4} style={{ margin: 0 }}>
          Stock In
        </Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setShowCreate(true)}>
          Create Stock In
        </Button>
      </div>

      <div style={{ background: '#fff', borderRadius: 8, padding: 24 }}>
        <Space style={{ marginBottom: 16 }}>
          <Select
            placeholder="All Status"
            allowClear
            style={{ width: 160 }}
            value={statusFilter}
            onChange={(v) => {
              setStatusFilter(v)
              setPage(1)
            }}
            options={[
              { value: 'CREATED', label: 'Created' },
              { value: 'IN_PROGRESS', label: 'In Progress' },
              { value: 'DONE', label: 'Done' },
              { value: 'CANCELLED', label: 'Cancelled' },
            ]}
          />
        </Space>

        <Table
          rowKey="id"
          columns={columns}
          dataSource={data?.items ?? []}
          loading={isFetching}
          pagination={{
            current: page,
            pageSize: 10,
            total: data?.pagination?.total,
            onChange: setPage,
            showSizeChanger: false,
          }}
        />
      </div>

      {showCreate && (
        <CreateStockInModal
          productsMap={productsMap}
          onClose={() => setShowCreate(false)}
          onSuccess={() => queryClient.invalidateQueries({ queryKey: ['stock-in'] })}
        />
      )}
    </>
  )
}
