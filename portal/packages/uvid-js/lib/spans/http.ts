import {request} from '../session';
import {type HTTPDTO} from '../types/span';

export async function http(model: HTTPDTO) {
	await request('/span/http', model);
}

export function listenHTTP() {
	// TODO
}
