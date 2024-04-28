import {RouterProvider, createBrowserRouter} from 'react-router-dom';
import {ConfigProvider, theme} from 'antd';
import './App.css';
import {Home} from './pages/home';
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
]);

function App() {
	return (
		<ConfigProvider theme={{algorithm: theme.defaultAlgorithm}}>
			<RouterProvider router={router} />
		</ConfigProvider>
	);
}

export default App;
