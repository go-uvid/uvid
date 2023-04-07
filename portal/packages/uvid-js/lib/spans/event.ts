import {request} from '../session';

export async function event(name: string, value?: string) {
	const span = {
		name,
		value,
	};
	return request('/span/event', span);
}

const eventAttribute = 'data-uvid-name';
const eventAttributeSelector = `[${eventAttribute}]`;

/**
 * @example `<button data-uvid-name="register" data-uvid-value="">Register</button>`
 * @description When user click the button, uvid-js will track it and call `uvid.event('register', '')`, you can provide a optional event value by `data-uvid-value`
 */
export function listenEvent() {
	document.body.addEventListener('click', (evt) => {
		if (evt.target instanceof HTMLElement) {
			const source = evt.target.closest(eventAttributeSelector);
			if (source instanceof HTMLElement) {
				const name = source.dataset.uvidName!;
				const value = source.dataset.uvidValue!;
				void event(name, value);
			}
		}
	});
}
