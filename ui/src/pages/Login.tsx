import { Form, Input, Button, Card, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { handleLogin } from '../services/auth';
import { AxiosError } from 'axios';

interface LoginForm {
  username: string;
  password: string;
}

export default function Login() {
  const onFinish = async (values: LoginForm) => {
    try {
      const success = await handleLogin(values.username, values.password);
      if (success) {
        message.success('登录成功');
        // 修改重定向路径，确保包含 /ui 前缀
        window.location.href = '/ui';
      }
    } catch (error) {
      if (error instanceof AxiosError && error.response?.data?.error) {
        message.error(error.response.data.error);
      } else {
        message.error('登录失败');
      }
    }
  };

  return (
    <div style={{
      height: '100vh',
      width: '100vw',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: '#f0f2f5'
    }}>
      <Card style={{ width: 400, boxShadow: '0 4px 8px rgba(0,0,0,0.1)' }}>
        <h2 style={{ textAlign: 'center', marginBottom: 24 }}>OSS 管理系统</h2>
        <Form
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
              size="large"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
              size="large"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large">
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
} 