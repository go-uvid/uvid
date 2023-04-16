import React from 'react';
import {createRoot} from 'react-dom/client';
import App from './App';
import {setupDayjs} from './lib/dayjs';

createRoot(document.querySelector('#root')!).render(
	<React.StrictMode>
		<App />
	</React.StrictMode>,
);

setupDayjs();
