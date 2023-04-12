'use client';

import { Layout } from 'antd';
import { LoginForm } from './form';

const { Content } = Layout;

export default function Home() {
  return (
    <Layout className="" style={{
		height: '100vh',
	}}>
      <Content style={{ padding: '50px' }}>
        <div
          style={{
            backgroundColor: '#fff',
            padding: '50px',
            maxWidth: '500px',
            margin: '0 auto',
          }}
        >
          <h2 className="mr-4">Login</h2>
          <LoginForm />
        </div>
      </Content>
    </Layout>
  );
}
