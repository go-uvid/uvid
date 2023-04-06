/* eslint-disable @typescript-eslint/no-floating-promises */
import {onLCP, onFID, onCLS} from 'web-vitals';
import {request} from '../session';
import {type PerformanceDTO} from '../types/span';
import {pageview} from './pageview';

export async function performance(name: PerformanceDTO['name'], value: number) {
	const span = {
		name,
		value,
	};
	return request('/span/performance', span as PerformanceDTO);
}

// A metric callback may be called more than once:
const metricIds = new Set<string>();

export function listenPerformance() {
	onLCP((metric) => {
		if (metricIds.has(metric.id)) return;
		performance('LCP', metric.value);
		// https://github.com/GoogleChrome/web-vitals#basic-usage
		// do not call any of the Web Vitals functions (e.g. onCLS(), onFID(), onLCP()) more than once per page load.
		pageview();
		metricIds.add(metric.id);
	});
	onCLS((metric) => {
		if (metricIds.has(metric.id)) return;
		performance('CLS', metric.value);
		metricIds.add(metric.id);
	});
	onFID((metric) => {
		if (metricIds.has(metric.id)) return;
		performance('FID', metric.value);
		metricIds.add(metric.id);
	});
}
