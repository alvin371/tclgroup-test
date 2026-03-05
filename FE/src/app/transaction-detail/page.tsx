import { useParams, useNavigate } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import {
  Button,
  Card,
  Col,
  Row,
  Table,
  Tag,
  Timeline,
  Typography,
} from 'antd'
import {
  ArrowLeftOutlined,
  PrinterOutlined,
  UserOutlined,
  CalendarOutlined,
  EnvironmentOutlined,
  SwapOutlined,
  CheckCircleFilled,
  SyncOutlined,
  CloseCircleOutlined,
  PlusCircleOutlined,
} from '@ant-design/icons'
import { stockInApi, type TStockInDetail } from '@/modules/stock-in'
import { stockOutApi, type TStockOutDetail } from '@/modules/stock-out'
import { ROUTES } from '@/commons/route'

const { Title, Text } = Typography

type TDetailData = TStockInDetail | TStockOutDetail

const STATUS_COLOR: Record<string, string> = {
  CREATED: 'default',
  DRAFT: 'default',
  IN_PROGRESS: 'processing',
  DONE: 'success',
  CANCELLED: 'error',
}

function timelineIcon(status: string) {
  switch (status) {
    case 'DONE':
      return <CheckCircleFilled style={{ color: '#52c41a' }} />
    case 'IN_PROGRESS':
      return <SyncOutlined style={{ color: '#1677ff' }} />
    case 'CANCELLED':
      return <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
    default:
      return <PlusCircleOutlined style={{ color: '#d9d9d9' }} />
  }
}

function InfoItem({
  icon,
  label,
  value,
}: {
  icon: React.ReactNode
  label: string
  value: string
}) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
      <Text type="secondary" style={{ fontSize: 12 }}>
        {icon} {label}
      </Text>
      <Text strong>{value || '—'}</Text>
    </div>
  )
}

type Props = {
  type: 'stock-in' | 'stock-out'
}

export default function TransactionDetailPage({ type }: Props) {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const { data: detail, isLoading } = useQuery<TDetailData>({
    queryKey: [type, 'detail', id],
    queryFn: () =>
      type === 'stock-in'
        ? stockInApi.getDetail(id!).then((r) => r.data)
        : stockOutApi.getDetail(id!).then((r) => r.data),
    enabled: !!id,
  })

  const shortId = id ? id.slice(0, 8).toUpperCase() : ''
  const title = `Transaction TRX-${shortId}`

  const lineItemColumns = [
    { title: 'SKU', dataIndex: 'product_sku', key: 'sku' },
    { title: 'Item Name', dataIndex: 'product_name', key: 'name' },
    {
      title: 'Quantity',
      dataIndex: 'quantity',
      key: 'qty',
      render: (qty: number) => (
        <Text style={{ color: type === 'stock-in' ? '#52c41a' : '#ff4d4f' }}>
          {type === 'stock-in' ? '+' : '-'}
          {qty}
        </Text>
      ),
    },
    {
      title: 'Unit Price',
      dataIndex: 'unit_price',
      key: 'unit_price',
      render: (v: number) =>
        v != null
          ? `$${Number(v).toLocaleString('en-US', { minimumFractionDigits: 2 })}`
          : '—',
    },
  ]

  return (
    <div>
      {/* Header */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          marginBottom: 24,
        }}
      >
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 4 }}>
            <Title level={4} style={{ margin: 0 }}>
              {title}
            </Title>
            {detail && (
              <Tag color={STATUS_COLOR[detail.status] ?? 'default'}>{detail.status}</Tag>
            )}
          </div>
          <Text type="secondary">
            {type === 'stock-in' ? 'Inbound stock transaction' : 'Outbound stock transaction'}
          </Text>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(ROUTES.REPORTS)}>
            Back to Reports
          </Button>
          <Button
            type="primary"
            icon={<PrinterOutlined />}
            onClick={() => window.print()}
          >
            Print Receipt
          </Button>
        </div>
      </div>

      {/* Info Card */}
      <Card style={{ marginBottom: 16 }} loading={isLoading}>
        {detail && (
          <Row gutter={[24, 16]}>
            <Col span={6}>
              <InfoItem
                icon={<SwapOutlined />}
                label="TYPE"
                value={type === 'stock-in' ? 'Stock In' : 'Stock Out'}
              />
            </Col>
            <Col span={6}>
              <InfoItem
                icon={<UserOutlined />}
                label="USER"
                value={detail.performed_by || 'System'}
              />
            </Col>
            <Col span={6}>
              <InfoItem
                icon={<CalendarOutlined />}
                label="DATE / TIME"
                value={new Date(detail.created_at).toLocaleString()}
              />
            </Col>
            <Col span={6}>
              <InfoItem
                icon={<EnvironmentOutlined />}
                label="LOCATION"
                value={detail.location}
              />
            </Col>
          </Row>
        )}
      </Card>

      {/* Line Items */}
      <Card title="Line Items" style={{ marginBottom: 16 }}>
        <Table
          rowKey="id"
          columns={lineItemColumns}
          dataSource={detail ? [detail] : []}
          loading={isLoading}
          pagination={false}
        />
      </Card>

      {/* Transaction Timeline */}
      <Card title="Transaction Timeline">
        {detail && (
          <Timeline
            mode="left"
            items={detail.timeline.map((entry) => ({
              dot: timelineIcon(entry.status),
              label: new Date(entry.occurred_at).toLocaleString(),
              children: (
                <div>
                  <Tag color={STATUS_COLOR[entry.status] ?? 'default'} style={{ marginBottom: 4 }}>
                    {entry.status}
                  </Tag>
                  <div>
                    <Text type="secondary">{entry.description}</Text>
                  </div>
                </div>
              ),
            }))}
          />
        )}
      </Card>
    </div>
  )
}
