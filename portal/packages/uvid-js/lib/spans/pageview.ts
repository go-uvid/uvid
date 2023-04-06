import {request} from '../session';

export async function pageview() {
	const span = {
		url: location.href,
	};
	return request('/span/pageview', span);
}
