import {RouterProvider, createBrowserRouter} from 'react-router-dom';
import {ConfigProvider, theme} from 'antd';
import './App.css';
import {Home} from './pages/home';
import {Login} from './pages/login';
import {Setting} from './pages/setting';
import {RootLayout} from './components/layout';

const router = createBrowserRouter([
	{
		path: '/',
		element: (
			<RootLayout>
				<Home />
			</RootLayout>
		),
	},
	{
		path: '/setting',
		element: (
			<RootLayout>
				<Setting />
			</RootLayout>
		),
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
