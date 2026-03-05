import { useState, useMemo } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
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
  Input
} from "antd";
import { PlusOutlined } from "@ant-design/icons";
import {
  stockOutApi,
  type TStockOut,
  type TStockOutStatus
} from "@/modules/stock-out";
import { productsApi, type TProduct } from "@/modules/products";
import { useMutation } from "@/app/_hooks/request/use-mutation";

const { Title } = Typography;

const STATUS_COLOR: Record<TStockOutStatus, string> = {
  DRAFT: "blue",
  IN_PROGRESS: "orange",
  DONE: "green",
  CANCELLED: "default"
};

function AllocateModal({
  productsMap,
  onClose,
  onSuccess
}: {
  productsMap: Map<string, TProduct>;
  onClose: () => void;
  onSuccess: () => void;
}) {
  const [form] = Form.useForm();

  const { mutate, isPending } = useMutation({
    mutationFn: stockOutApi.allocate,
    onSuccess: () => {
      message.success("Stock allocated");
      onClose();
      onSuccess();
    },
    onError: (err) => message.error(err.error.message)
  });

  const productOptions = Array.from(productsMap.values()).map((p) => ({
    value: p.id,
    label: `${p.name} (${p.sku})`
  }));

  return (
    <Modal
      open
      title="Allocate Stock Out"
      onCancel={onClose}
      onOk={() => form.submit()}
      confirmLoading={isPending}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={(values) =>
          mutate(values as Parameters<typeof stockOutApi.allocate>[0])
        }
      >
        <Form.Item
          label="Product"
          name="product_id"
          rules={[{ required: true }]}
        >
          <Select
            options={productOptions}
            placeholder="Select a product"
            showSearch
            filterOption={(input, option) =>
              (option?.label ?? "").toLowerCase().includes(input.toLowerCase())
            }
          />
        </Form.Item>
        <Form.Item
          label="Quantity"
          name="quantity"
          rules={[{ required: true }]}
        >
          <InputNumber min={1} style={{ width: "100%" }} />
        </Form.Item>
        <Form.Item label="Notes" name="notes" rules={[{ required: true }]}>
          <Input.TextArea rows={3} />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default function StockOutPage() {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<TStockOutStatus | undefined>(
    undefined
  );
  const [showAllocate, setShowAllocate] = useState(false);
  const queryClient = useQueryClient();

  const { data, isFetching } = useQuery({
    queryKey: ["stock-out", page, statusFilter],
    queryFn: () =>
      stockOutApi.list({
        page,
        perPage: 10,
        ...(statusFilter ? { status: statusFilter } : {})
      })
  });
  const { data: productsData } = useQuery({
    queryKey: ["products"],
    queryFn: () => productsApi.list({ perPage: 100 })
  });

  const productsMap = useMemo(() => {
    const map = new Map<string, TProduct>();
    for (const p of productsData?.items ?? []) map.set(p.id, p);
    return map;
  }, [productsData]);

  const { mutate: executePicked } = useMutation({
    mutationFn: (id: string) => stockOutApi.execute(id, { status: "IN_PROGRESS" }),
    onSuccess: () => {
      message.success("Moved to in progress");
      queryClient.invalidateQueries({ queryKey: ["stock-out"] });
    },
    onError: (err) => message.error(err.error.message)
  });

  const { mutate: executeShipped } = useMutation({
    mutationFn: (id: string) => stockOutApi.execute(id, { status: "DONE" }),
    onSuccess: () => {
      message.success("Marked as done");
      queryClient.invalidateQueries({ queryKey: ["stock-out"] });
    },
    onError: (err) => message.error(err.error.message)
  });

  const { mutate: cancel } = useMutation({
    mutationFn: (id: string) => stockOutApi.remove(id),
    onSuccess: () => {
      message.success("Cancelled — available stock restored");
      queryClient.invalidateQueries({ queryKey: ["stock-out"] });
      queryClient.invalidateQueries({ queryKey: ["inventory"] });
    },
    onError: (err) => message.error(err.error.message)
  });

  const columns = [
    {
      title: "Product",
      key: "product",
      render: (_: unknown, r: TStockOut) =>
        productsMap.get(r.product_id)?.name ?? r.product_id
    },
    { title: "Quantity", dataIndex: "quantity", key: "qty" },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      render: (s: TStockOutStatus) => <Tag color={STATUS_COLOR[s]}>{s.replace("_", " ")}</Tag>
    },
    { title: "Notes", dataIndex: "notes", key: "notes", ellipsis: true },
    {
      title: "Created",
      dataIndex: "created_at",
      key: "created_at",
      render: (v: string) => new Date(v).toLocaleString()
    },
    {
      title: "Actions",
      key: "actions",
      render: (_: unknown, r: TStockOut) => {
        if (r.status === "DRAFT") {
          return (
            <Space>
              <Popconfirm
                title="Move to in progress?"
                onConfirm={() => executePicked(r.id)}
              >
                <Button type="primary" size="small">
                  Start
                </Button>
              </Popconfirm>
              <Popconfirm
                title="Cancel allocation? Available stock will be restored."
                onConfirm={() => cancel(r.id)}
              >
                <Button danger size="small">
                  Cancel
                </Button>
              </Popconfirm>
            </Space>
          );
        }
        if (r.status === "IN_PROGRESS") {
          return (
            <Space>
              <Popconfirm
                title="Mark as done?"
                onConfirm={() => executeShipped(r.id)}
              >
                <Button type="primary" size="small">
                  Complete
                </Button>
              </Popconfirm>
              <Popconfirm
                title="Cancel? Available stock will be restored."
                onConfirm={() => cancel(r.id)}
              >
                <Button danger size="small">
                  Cancel
                </Button>
              </Popconfirm>
            </Space>
          );
        }
        return null;
      }
    }
  ];

  return (
    <>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          marginBottom: 24
        }}
      >
        <Title level={4} style={{ margin: 0 }}>
          Stock Out
        </Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setShowAllocate(true)}
        >
          Allocate Stock Out
        </Button>
      </div>

      <div style={{ background: "#fff", borderRadius: 8, padding: 24 }}>
        <Space style={{ marginBottom: 16 }}>
          <Select
            placeholder="All Status"
            allowClear
            style={{ width: 160 }}
            value={statusFilter}
            onChange={(v) => {
              setStatusFilter(v);
              setPage(1);
            }}
            options={[
              { value: "DRAFT", label: "Draft" },
              { value: "IN_PROGRESS", label: "In Progress" },
              { value: "DONE", label: "Done" },
              { value: "CANCELLED", label: "Cancelled" }
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
            showSizeChanger: false
          }}
        />
      </div>

      {showAllocate && (
        <AllocateModal
          productsMap={productsMap}
          onClose={() => setShowAllocate(false)}
          onSuccess={() =>
            queryClient.invalidateQueries({ queryKey: ["stock-out"] })
          }
        />
      )}
    </>
  );
}
