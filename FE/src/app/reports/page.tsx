import { useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { Table, Tabs, Tag, Typography } from 'antd'
import { reportsApi, type TReportStockIn, type TReportStockOut } from '@/modules/reports'
import { productsApi, type TProduct } from '@/modules/products'
import { toStockInDetail, toStockOutDetail } from '@/commons/route'

const { Title } = Typography

function useProductsMap() {
  const { data } = useQuery({
    queryKey: ['products'],
    queryFn: () => productsApi.list({ perPage: 100 }),
  })
  return useMemo(() => {
    const map = new Map<string, TProduct>()
    for (const p of data?.items ?? []) map.set(p.id, p)
    return map
  }, [data])
}

function StockInReportsTab() {
  const [page, setPage] = useState(1)
  const productsMap = useProductsMap()
  const navigate = useNavigate()

  const { data, isFetching } = useQuery({
    queryKey: ['reports-stock-in', page],
    queryFn: () => reportsApi.stockIn({ page, perPage: 10 }),
  })

  const columns = [
    {
      title: 'Product',
      key: 'product',
      render: (_: unknown, r: TReportStockIn) =>
        productsMap.get(r.product_id)?.name ?? r.product_id,
    },
    { title: 'Quantity', dataIndex: 'quantity', key: 'qty' },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: () => <Tag color="green">received</Tag>,
    },
    { title: 'Notes', dataIndex: 'notes', key: 'notes', ellipsis: true },
    {
      title: 'Completed Date',
      dataIndex: 'updated_at',
      key: 'updated_at',
      render: (v: string) => new Date(v).toLocaleString(),
    },
  ]

  return (
    <Table
      rowKey="id"
      columns={columns}
      dataSource={data?.items ?? []}
      loading={isFetching}
      onRow={(record) => ({
        onClick: () => navigate(toStockInDetail(record.id)),
        style: { cursor: 'pointer' },
      })}
      pagination={{
        current: page,
        pageSize: 10,
        total: data?.pagination?.total,
        onChange: setPage,
        showSizeChanger: false,
      }}
    />
  )
}

function StockOutReportsTab() {
  const [page, setPage] = useState(1)
  const productsMap = useProductsMap()
  const navigate = useNavigate()

  const { data, isFetching } = useQuery({
    queryKey: ['reports-stock-out', page],
    queryFn: () => reportsApi.stockOut({ page, perPage: 10 }),
  })

  const columns = [
    {
      title: 'Product',
      key: 'product',
      render: (_: unknown, r: TReportStockOut) =>
        productsMap.get(r.product_id)?.name ?? r.product_id,
    },
    { title: 'Quantity', dataIndex: 'quantity', key: 'qty' },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: () => <Tag color="green">shipped</Tag>,
    },
    { title: 'Notes', dataIndex: 'notes', key: 'notes', ellipsis: true },
    {
      title: 'Completed Date',
      dataIndex: 'updated_at',
      key: 'updated_at',
      render: (v: string) => new Date(v).toLocaleString(),
    },
  ]

  return (
    <Table
      rowKey="id"
      columns={columns}
      dataSource={data?.items ?? []}
      loading={isFetching}
      onRow={(record) => ({
        onClick: () => navigate(toStockOutDetail(record.id)),
        style: { cursor: 'pointer' },
      })}
      pagination={{
        current: page,
        pageSize: 10,
        total: data?.pagination?.total,
        onChange: setPage,
        showSizeChanger: false,
      }}
    />
  )
}

export default function ReportsPage() {
  const tabItems = [
    {
      key: 'stock-in',
      label: 'Stock In Reports',
      children: <StockInReportsTab />,
    },
    {
      key: 'stock-out',
      label: 'Stock Out Reports',
      children: <StockOutReportsTab />,
    },
  ]

  return (
    <>
      <Title level={4} style={{ margin: '0 0 24px' }}>
        Reports
      </Title>
      <div style={{ background: '#fff', borderRadius: 8, padding: 24 }}>
        <Tabs items={tabItems} />
      </div>
    </>
  )
}
