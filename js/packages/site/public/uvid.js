/* eslint-disable no-undef */
import {init} from 'https://www.unpkg.com/uvid-js?module';

/** @type {import('../../uvid-js/lib/main').init} */
const _init = init;

const uvid = _init({
	host: 'https://uvid-demo.applet.ink',
	sessionMeta: {
		from: 'uvid-site',
	},
});

window.addEventListener('click', (event) => {
	const {target} = event;
	if (
		target instanceof HTMLElement &&
		target.textContent &&
		target.closest('.action > .VPButton')
	) {
		uvid.event(target.textContent);
	}
});

// eslint-disable-next-line unicorn/prefer-add-event-listener, max-params
window.onerror = function (a, b, c, d, error) {
	uvid.error(error);
};
