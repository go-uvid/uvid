'use client';

import { Layout } from 'antd';
import { LoginForm } from './form';

const { Content } = Layout;

export default function Home() {
  return (
    <Layout
      className="flex justify-center items-center"
      style={{
        height: '100vh',
      }}
    >
      <Content className="bg-white rounded shadow p-8 w-96 max-h-72">
        <h1 className="mb-4 text-lg text-black">Login</h1>
        <LoginForm />
      </Content>
    </Layout>
  );
}
