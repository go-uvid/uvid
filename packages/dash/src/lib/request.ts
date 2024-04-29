import {buildUrl} from './url';

export class RequestError extends Error {
	constructor(public message: string, public status: number) {
		super(message);
	}
}

const host = import.meta.env.DEV ? 'http://localhost:3000' : '';

export async function baseRequest<T = any>(
	method: string,
	url: string,
	body?: string,
	headers?: Record<string, unknown>,
): Promise<{ok: boolean; status: number; data: T}> {
	const response = await fetch(host + url, {
		method,
		cache: 'no-cache',
		headers: {
			Accept: 'application/json',
			'Content-Type': 'application/json',
			...headers,
		},
		body,
	});
	if (response.ok) {
		// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
		const data = await response.json();
		return {ok: response.ok, status: response.status, data: data as T};
	}

	const text = await response.text();
	throw new RequestError(text, response.status);
}

export async function apiRequest<T = any>(
	method: string,
	url: string,
	body?: string,
	headers?: Record<string, unknown>,
): Promise<{ok: boolean; status: number; data: T}> {
	const response = await baseRequest<T>(method, url, body, {
		...headers,
	});
	return response;
}

export async function get<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	const response = await apiRequest<T>(
		'get',
		buildUrl(url, parameters),
		undefined,
		headers,
	);
	return response.data;
}

export async function del<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	const response = await apiRequest<T>(
		'delete',
		buildUrl(url, parameters),
		undefined,
		headers,
	);
	return response.data;
}

export async function post<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	const response = await apiRequest<T>(
		'post',
		url,
		JSON.stringify(parameters),
		headers,
	);
	return response.data;
}

export async function put<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	const response = await apiRequest<T>(
		'put',
		url,
		JSON.stringify(parameters),
		headers,
	);
	return response.data;
}
