/* eslint-disable @typescript-eslint/no-floating-promises */
import {init} from '../lib/main';

const sdk = init({
	host: 'http://localhost:3000',
	sessionMeta: {
		userId: '123',
	},
	appVersion: '1.0.0',
	// @ts-expect-error __internal__request is not part of the public API
	__internal__request: console.log,
});

const errorButton = document.querySelector('#error-button');
const httpButton = document.querySelector('#http-button');
const eventButton = document.querySelector('#event-button');

errorButton?.addEventListener('click', () => {
	try {
		throw new Error('This is an error!');
	} catch (error) {
		sdk.error(error as Error);
	}
});

httpButton?.addEventListener('click', async () => {
	sdk.http({
		resource: 'http://google.com',
		status: 400,
		method: 'GET',
		headers: '',
		response: '',
	});
});

eventButton?.addEventListener('click', (event) => {
	sdk.event('click', event.timeStamp.toString());
});
