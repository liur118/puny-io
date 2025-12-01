import { useState, useEffect } from 'react';
import { Layout, Menu, Button, Table, Upload, Modal, Input, message, notification } from 'antd';
import {
  FolderOutlined,
  UploadOutlined,
  DeleteOutlined,
  DownloadOutlined,
  PlusOutlined,
  LinkOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons';
import { listBuckets, listObjects, createBucket, uploadObject, downloadObject, deleteObject, getObjectURL } from '../services/api';
import { AxiosError } from 'axios';

const { Header, Sider, Content } = Layout;

export default function Home() {
  const [buckets, setBuckets] = useState<string[]>([]);
  const [selectedBucket, setSelectedBucket] = useState<string>('');
  const [objects, setObjects] = useState<string[]>([]);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [newBucketName, setNewBucketName] = useState('');

  useEffect(() => {
    loadBuckets();
  }, []);

  useEffect(() => {
    if (selectedBucket) {
      loadObjects(selectedBucket);
    }
  }, [selectedBucket]);

  const loadBuckets = async () => {
    try {
      const { data } = await listBuckets();
      setBuckets(Array.isArray(data.buckets) ? data.buckets : []);
    } catch {
      message.error('加载桶列表失败');
    }
  };

  const loadObjects = async (bucket: string) => {
    try {
      const { data } = await listObjects(bucket);
      setObjects(data.objects || []);
    } catch (error) {
      if (error instanceof AxiosError && error.response?.data?.error) {
        message.error(error.response.data.error);
      } else {
        message.error('加载对象列表失败');
      }
      setObjects([]);
    }
  };

  const handleCreateBucket = async () => {
    try {
      await createBucket(newBucketName);
      await loadBuckets();
      setIsModalVisible(false);
      setNewBucketName('');
      message.success('创建桶成功');
    } catch {
      message.error('创建桶失败');
    }
  };

  const handleUpload = async (file: File) => {
    if (!selectedBucket) {
      message.error('请先选择桶');
      return;
    }
    try {
      await uploadObject(selectedBucket, file.name, file);
      loadObjects(selectedBucket);
      message.success('上传成功');
    } catch {
      message.error('上传失败');
    }
  };

  const handleDownload = async (key: string) => {
    try {
      const response = await downloadObject(selectedBucket, key);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', key);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch {
      message.error('下载失败');
    }
  };

  const handleDelete = async (key: string) => {
    try {
      await deleteObject(selectedBucket, key);
      loadObjects(selectedBucket);
      message.success('删除成功');
    } catch {
      message.error('删除失败');
    }
  };

  const handleGetURL = async (key: string) => {
    try {
      const { data } = await getObjectURL(selectedBucket, key);
      
      // 尝试使用 Clipboard API
      try {
        await navigator.clipboard.writeText(data.url);
        notification.success({
          message: '链接已复制',
          description: (
            <div>
              <p>访问链接已成功复制到剪贴板：</p>
              <code style={{ 
                display: 'block', 
                padding: '8px', 
                background: '#f5f5f5', 
                borderRadius: '4px',
                wordBreak: 'break-all',
                fontSize: '12px'
              }}>
                {data.url}
              </code>
            </div>
          ),
          icon: <CheckCircleOutlined style={{ color: '#52c41a' }} />,
          duration: 4,
          placement: 'topRight',
        });
      } catch {
        // 如果 Clipboard API 失败，使用传统方法
        const textArea = document.createElement('textarea');
        textArea.value = data.url;
        document.body.appendChild(textArea);
        textArea.select();
        try {
          document.execCommand('copy');
          notification.success({
            message: '链接已复制',
            description: (
              <div>
                <p>访问链接已成功复制到剪贴板：</p>
                <code style={{ 
                  display: 'block', 
                  padding: '8px', 
                  background: '#f5f5f5', 
                  borderRadius: '4px',
                  wordBreak: 'break-all',
                  fontSize: '12px'
                }}>
                  {data.url}
                </code>
              </div>
            ),
            icon: <CheckCircleOutlined style={{ color: '#52c41a' }} />,
            duration: 4,
            placement: 'topRight',
          });
        } catch {
          // 所有复制方法都失败，显示链接供手动复制
          Modal.info({
            title: '请手动复制链接',
            content: (
              <div>
                <p>自动复制失败，请手动复制以下链接：</p>
                <Input.TextArea
                  value={data.url}
                  readOnly
                  autoSize={{ minRows: 2, maxRows: 4 }}
                  style={{ marginTop: '8px' }}
                  onFocus={(e) => e.target.select()}
                />
              </div>
            ),
            okText: '知道了',
          });
        }
        document.body.removeChild(textArea);
      }
    } catch (error) {
      if (error instanceof AxiosError && error.response?.data?.error) {
        message.error(error.response.data.error);
      } else {
        message.error('获取访问链接失败');
      }
    }
  };

  const columns = [
    {
      title: '名称',
      dataIndex: 'key',
      key: 'key',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: { key: string }) => (
        <>
          <Button
            type="link"
            icon={<DownloadOutlined />}
            onClick={() => handleDownload(record.key)}
          >
            下载
          </Button>
          <Button
            type="link"
            icon={<LinkOutlined />}
            onClick={() => handleGetURL(record.key)}
          >
            获取链接
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.key)}
          >
            删除
          </Button>
        </>
      ),
    },
  ];

  return (
    <Layout style={{ minHeight: '100vh', width: '100vw' }}>
      <Header style={{ 
        background: '#fff', 
        padding: '0 16px',
        position: 'fixed',
        width: '100%',
        zIndex: 1,
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <h1 style={{ margin: 0, lineHeight: '64px' }}>OSS 管理系统</h1>
      </Header>
      <Layout style={{ marginTop: 64 }}>
        <Sider width={200} style={{ 
          background: '#fff',
          position: 'fixed',
          height: 'calc(100vh - 64px)',
          left: 0,
          overflow: 'auto'
        }}>
          <div style={{ padding: '16px' }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setIsModalVisible(true)}
              block
            >
              新建桶
            </Button>
          </div>
          <Menu
            mode="inline"
            selectedKeys={[selectedBucket]}
            style={{ height: '100%', borderRight: 0 }}
            items={(buckets || []).map(bucket => ({
              key: bucket,
              icon: <FolderOutlined />,
              label: bucket,
              onClick: () => setSelectedBucket(bucket),
            }))}
          />
        </Sider>
        <Layout style={{ marginLeft: 200 }}>
          <Content style={{ 
            background: '#fff', 
            padding: 24, 
            margin: 24,
            minHeight: 280,
            borderRadius: 4,
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}>
            {selectedBucket ? (
              <>
                <div style={{ marginBottom: 16 }}>
                  <Upload
                    customRequest={({ file }) => handleUpload(file as File)}
                    showUploadList={false}
                  >
                    <Button icon={<UploadOutlined />}>上传文件</Button>
                  </Upload>
                </div>
                <Table
                  columns={columns}
                  dataSource={objects.map(key => ({ key }))}
                  rowKey="key"
                />
              </>
            ) : (
              <div style={{ textAlign: 'center', color: '#999' }}>
                请选择一个桶
              </div>
            )}
          </Content>
        </Layout>
      </Layout>

      <Modal
        title="新建桶"
        open={isModalVisible}
        onOk={handleCreateBucket}
        onCancel={() => setIsModalVisible(false)}
      >
        <Input
          placeholder="请输入桶名称"
          value={newBucketName}
          onChange={e => setNewBucketName(e.target.value)}
        />
      </Modal>
    </Layout>
  );
} 