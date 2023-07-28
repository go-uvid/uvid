import {type SDK} from './lib/sdk';

declare global {
	// eslint-disable-next-line @typescript-eslint/consistent-type-definitions
	interface Window {
		uvid: SDK;
	}
}
