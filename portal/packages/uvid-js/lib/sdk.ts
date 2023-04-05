import {type SDKConfig, sdkConfig} from './config';
import {error} from './spans/error';
import {event} from './spans/event';
import {http} from './spans/http';
import {pageview} from './spans/pageview';
import {performance} from './spans/performance';

export const sdk = {
	error,
	http,
	event,
	performance,
	pageview,
	setSessionMeta,
};

/**
 * Set meta data for the current session
 * @param params
 */
function setSessionMeta(meta: SDKConfig['sessionMeta']) {
	Object.assign(sdkConfig.sessionMeta, meta);
}
