import {request} from '../session';

export async function event(name: string, value: string) {
	const span = {
		name,
		value,
	};
	await request('/span/event', span);
}
