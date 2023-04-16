import {Link, RouterProvider, createBrowserRouter} from 'react-router-dom';
import {ConfigProvider, Layout, theme} from 'antd';
import {HomeOutlined, SettingOutlined} from '@ant-design/icons';
import {type CSSProperties} from 'react';
import './App.css';
import {Home} from './pages/home';
import {Login} from './pages/login';

const {Header, Content} = Layout;

function HomeWithHeader() {
	return (
		<Layout>
			<Header className="flex justify-center">
				<Content className="w-main flex justify-between flex-grow-0">
					<Link to="/">
						<HomeOutlined className="text-base" />
					</Link>
					<Link to="/setting">
						<SettingOutlined className="text-base" />
					</Link>
				</Content>
			</Header>
			<Content className="m-auto">
				<Home />
			</Content>
		</Layout>
	);
}

const router = createBrowserRouter([
	{
		path: '/',
		element: <HomeWithHeader />,
	},
	{
		path: '/login',
		element: <Login />,
	},
]);

function App() {
	return (
		<ConfigProvider theme={{algorithm: theme.defaultAlgorithm}}>
			<RouterProvider router={router} />
		</ConfigProvider>
	);
}

export default App;
