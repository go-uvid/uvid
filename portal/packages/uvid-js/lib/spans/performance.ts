/* eslint-disable @typescript-eslint/no-floating-promises */
import {onLCP, onFID, onCLS, type Metric} from 'web-vitals';
import {request} from '../session';
import {type PerformanceDTO} from '../types/span';

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
		if (validateMetric(metric)) {
			performance('LCP', metric.value);
		}
	});
	onCLS((metric) => {
		if (validateMetric(metric)) {
			performance('CLS', metric.value);
		}
	});
	onFID((metric) => {
		if (validateMetric(metric)) {
			performance('FID', metric.value);
		}
	});
}

function validateMetric(metric: Metric): boolean {
	if (metricIds.has(metric.id)) return false;
	metricIds.add(metric.id);
	// eslint-disable-next-line @typescript-eslint/no-unsafe-call
	window.uvid.performanceMetrics.push(metric);
	return true;
}
