import {request} from '../session';

export async function error(error: Error) {
	const {name, message, stack} = error;
	const span = {
		name,
		message,
		stack,
	};
	await request('/span/error', span);
}

export function listenError() {
	// TODO
}
