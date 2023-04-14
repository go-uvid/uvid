import {buildUrl} from './url';

const authTokenKey = 'uvid-auth-token';

function getAuthorization() {
	const token = localStorage.getItem(authTokenKey);
	return token ? `Bearer ${token}` : '';
}

export class RequestError extends Error {
	constructor(public message: string, public status: number) {
		super(message);
	}
}

export async function apiRequest<T = any>(
	method: string,
	url: string,
	body?: string,
	headers?: Record<string, unknown>,
): Promise<{ok: boolean; status: number; data: T}> {
	const response = await fetch(url, {
		method,
		cache: 'no-cache',
		headers: {
			Accept: 'application/json',
			'Content-Type': 'application/json',
			Authorization: getAuthorization(),
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

export async function get<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	return apiRequest<T>('get', buildUrl(url, parameters), undefined, headers);
}

export async function del<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	return apiRequest<T>('delete', buildUrl(url, parameters), undefined, headers);
}

export async function post<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	return apiRequest<T>('post', url, JSON.stringify(parameters), headers);
}

export async function put<T>(
	url: string,
	parameters?: Record<string, unknown>,
	headers?: Record<string, unknown>,
) {
	return apiRequest<T>('put', url, JSON.stringify(parameters), headers);
}

export async function login(name: string, password: string) {
	return post<{token: string}>('/dash/user/login', {name, password}).then(
		(response) => {
			if (response.ok) {
				const {token} = response.data;
				localStorage.setItem(authTokenKey, token);
			}
		},
	);
}
