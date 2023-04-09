import {request} from '../session';
import { ErrorDTO } from '../types/span';

export async function error(error: Error) {
	const {name, message, stack} = error;
	const span: ErrorDTO = {
		name,
		message,
		stack: stack ?? '',
	};
	return request('/span/error', span);
}

export function listenError() {
	// TODO
}
