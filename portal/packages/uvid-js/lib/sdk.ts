import {type Metric} from 'web-vitals';
import {type SDKConfig, sdkConfig} from './config';
import {error} from './spans/error';
import {event} from './spans/event';
import {http} from './spans/http';
import {pageview} from './spans/pageview';
import {performance} from './spans/performance';
import {type HTTPDTO, type PerformanceDTO} from './types/span';

export type SDK = {
	error: (error: Error) => Promise<Response>;
	http: (model: HTTPDTO) => Promise<Response>;
	event: (name: string, value: string) => Promise<Response>;
	performance: (
		name: PerformanceDTO['name'],
		value: number,
	) => Promise<Response>;
	pageview: () => Promise<Response>;
	setSessionMeta: (meta: SDKConfig['sessionMeta']) => void;
	performanceMetrics: Metric[];
};

export const sdk: SDK = {
	error,
	http,
	event,
	performance,
	pageview,
	setSessionMeta,
	performanceMetrics: [],
};

/**
 * Set meta data for the current session
 * @param params
 */
function setSessionMeta(meta: SDKConfig['sessionMeta']) {
	Object.assign(sdkConfig.sessionMeta, meta);
}
