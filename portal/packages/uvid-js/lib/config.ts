import {type RequestData} from './session';

export type SDKConfig = {
	host: string;
	sessionMeta: Record<string, any>;
	appVersion?: string;
	/**
	 * DO NOT use this option
	 * @warning
	 */
	__internal__request?: (path: string, data: RequestData) => Promise<Response>;
};

export const sdkConfig: SDKConfig = {
	host: '',
	sessionMeta: {},
	appVersion: '',
};
