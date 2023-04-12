import { buildUrl } from './url';

const authTokenKey = 'authToken';

export function apiRequest(
  method: string,
  url: string,
  body?: string,
  headers?: object
): Promise<{ ok: boolean; status: number; data?: any; error?: any }> {
  return fetch(url, {
    method,
    cache: 'no-cache',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem(authTokenKey)}`,
      ...headers,
    },
    body,
  }).then((res) => {
    if (res.ok) {
      return res
        .json()
        .then((data) => ({ ok: res.ok, status: res.status, data }));
    }

    return res
      .text()
      .then((text) =>
        Promise.reject({ ok: res.ok, status: res.status, error: text })
      );
  });
}

export function get(url: string, params?: object, headers?: object) {
  return apiRequest('get', buildUrl(url, params), undefined, headers);
}

export function del(url: string, params?: object, headers?: object) {
  return apiRequest('delete', buildUrl(url, params), undefined, headers);
}

export function post(url: string, params?: object, headers?: object) {
  return apiRequest('post', url, JSON.stringify(params), headers);
}

export function put(url: string, params?: object, headers?: object) {
  return apiRequest('put', url, JSON.stringify(params), headers);
}

export function login(name: string, password: string) {
  return post('/dash/user/login', { name, password }).then((response) => {
    if (response.ok) {
      const { token } = response.data;
      localStorage.setItem(authTokenKey, token);
    }
  });
}
