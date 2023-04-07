import {type SDKConfig, sdkConfig} from './config';
import {sdk} from './sdk';
import {listenError} from './spans/error';
import {listenEvent} from './spans/event';
import {listenHTTP} from './spans/http';
import {listenPageview} from './spans/pageview';
import {listenPerformance} from './spans/performance';
import {assert} from './utils';

export function init(config: SDKConfig): typeof sdk {
	window.uvid = sdk;
	try {
		assert(config.host, 'No host');
		Object.assign(sdkConfig, config);
		listenError();
		listenHTTP();
		listenPerformance();
		listenPageview();
		listenEvent();
		return sdk;
	} catch (error) {
		console.error(error);
		return sdk;
	}
}
