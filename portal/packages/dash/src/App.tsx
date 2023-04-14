import {RouterProvider, createBrowserRouter} from 'react-router-dom';
import './App.css';
import {Home} from './pages/home';
import {Login} from './pages/login';

const router = createBrowserRouter([
	{
		path: '/',
		element: <Home />,
	},
	{
		path: '/login',
		element: <Login />,
	},
]);

function App() {
	return (
		<div className="App">
			<RouterProvider router={router} />
		</div>
	);
}

export default App;
