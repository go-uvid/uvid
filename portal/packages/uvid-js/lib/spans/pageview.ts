import {request} from '../session';
import {domReady} from '../utils';

export async function pageview() {
	const span = {
		url: location.href,
	};
	return request('/span/pageview', span);
}

export function listenPageview() {
	domReady(pageview);
}
