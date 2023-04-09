import {request} from '../session';
import { EventDTO } from '../types/span';

export async function event(action: string, value?: string) {
	const span: EventDTO = {
		action,
		value,
	};
	return request('/span/event', span);
}

const eventAttribute = 'data-uvid-action';
const eventAttributeSelector = `[${eventAttribute}]`;

/**
 * @example `<button data-uvid-action="register" data-uvid-value="">Register</button>`
 * @description When user click the button, uvid-js will track it and call `uvid.event('register', '')`, you can provide a optional event value by `data-uvid-value`
 */
export function listenEvent() {
	document.body.addEventListener('click', (evt) => {
		if (evt.target instanceof HTMLElement) {
			const source = evt.target.closest(eventAttributeSelector);
			if (source instanceof HTMLElement) {
				const name = source.dataset.uvidAction!;
				const value = source.dataset.uvidValue!;
				void event(name, value);
			}
		}
	});
}
