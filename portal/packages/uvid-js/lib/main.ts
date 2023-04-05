import {type SDKConfig, sdkConfig} from './config';
import {sdk} from './sdk';
import {listenError} from './spans/error';
import {listenHTTP} from './spans/http';
import {listenPerformance} from './spans/performance';
import {assert} from './utils';

export function init(config: SDKConfig): typeof sdk {
	try {
		assert(config.host, 'No host');
		Object.assign(sdkConfig, config);
		listenError();
		listenHTTP();
		listenPerformance();
		return sdk;
	} catch (error) {
		console.error(error);
		return sdk;
	}
}
