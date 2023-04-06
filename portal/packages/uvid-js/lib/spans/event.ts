import {request} from '../session';

export async function event(name: string, value: string) {
	const span = {
		name,
		value,
	};
	return request('/span/event', span);
}
