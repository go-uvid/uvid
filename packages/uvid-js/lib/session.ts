import {sdkConfig} from './config';
import {
	type BaseSessionDTO,
	type ErrorDTO,
	type HTTPDTO,
	type EventDTO,
	type PerformanceDTO,
	type PageViewDTO,
} from './types/span';
import {assert} from './utils';

const sessionHeaderKey = 'X-UVID-Session';

async function initSession() {
	const data: BaseSessionDTO = {
		appVersion: sdkConfig.appVersion,
		url: location.href,
		screen: `${screen.width}*${screen.height}`,
		referrer: document.referrer,
		language: navigator.language,
		meta: JSON.stringify(sdkConfig.sessionMeta),
	};
	const response = await (sdkConfig.__internal__request ?? _request)(
		'/span/session',
		data,
	);
	const sessionValue = await response.text();
	sessionStorage.setItem(sessionHeaderKey, sessionValue);
}

export async function request(path: string, data: RequestData) {
	if (!getSession()) {
		await initSession();
	}

	return (sdkConfig.__internal__request ?? _request)(path, data);
}

export type RequestData =
	| ErrorDTO
	| HTTPDTO
	| EventDTO
	| PerformanceDTO
	| PageViewDTO
	| BaseSessionDTO;

async function _request(path: string, data: RequestData) {
	assert(sdkConfig.host, 'SDK host is not set');
	return fetch(`${sdkConfig.host}${path}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			[sessionHeaderKey]: getSession(),
		},
		keepalive: true,
		body: JSON.stringify(data),
	}).then((response) => {
		if (response.ok) return response;
		throw new Error(response.statusText);
	});
}

function getSession() {
	return sessionStorage.getItem(sessionHeaderKey) ?? '';
}
