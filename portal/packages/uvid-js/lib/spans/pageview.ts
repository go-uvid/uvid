import {request} from '../session';

export async function pageview() {
	const span = {
		url: location.href,
	};
	await request('/span/pageview', span);
}
