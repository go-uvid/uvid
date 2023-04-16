import {get, post} from './request';

type SpanFilter = {
	start: string;
	end: string;
};
export type TimeUnit = 'hour' | 'day' | 'month' | 'year';
export type TimeRangeDTO = SpanFilter & {
	unit: TimeUnit;
};

export enum PerformanceName {
	LCP = 'LCP',
	CLS = 'CLS',
	FID = 'FID',
}

export type IntervalData = {
	x: string;
	y: number;
};

export enum ApiPath {
	updateUserPassword = '/dash/user/password',
	getPageview = '/dash/pv',
	getPageviewCount = '/dash/pv/count',
	getUniqueVisitor = '/dash/uv',
	getUniqueVisitorCount = '/dash/uv/count',
	getError = '/dash/error',
	getErrorCount = '/dash/error/count',
	getHttpError = '/dash/http/error',
	getHttpErrorCount = '/dash/http/error/count',
	getAvgPerformance = '/dash/performance',
	getEvent = '/dash/event',
}

export async function updateUserPassword(data: {password: string}) {
	return post<void>(ApiPath.updateUserPassword, data);
}

export async function getPageview(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getPageview, timeRange);
}

export async function getPageviewCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getPageviewCount, timeRange);
}

export async function getUniqueVisitor(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getUniqueVisitor, timeRange);
}

export async function getUniqueVisitorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getUniqueVisitorCount, timeRange);
}

export async function getError(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getError, timeRange);
}

export async function getErrorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getErrorCount, timeRange);
}

export async function getHttpError(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getHttpError, timeRange);
}

export async function getHttpErrorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getHttpErrorCount, timeRange);
}

export async function getAvgPerformance(timeRange: SpanFilter) {
	return get<IntervalData[]>(ApiPath.getAvgPerformance, timeRange);
}

export async function getEvent(timeRange: SpanFilter) {
	return get<IntervalData[]>(ApiPath.getEvent, timeRange);
}
