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
	changeUserPassword = '/dash/user/password',
	getPageviewInterval = '/dash/pvs/interval',
	getPageviewCount = '/dash/pvs/count',
	getUniqueVisitorInterval = '/dash/uvs/interval',
	getUniqueVisitorCount = '/dash/uvs/count',
	getErrorInterval = '/dash/errors/interval',
	getErrorCount = '/dash/errors/count',
	getHttpErrorInterval = '/dash/https/errors/interval',
	getHttpErrorCount = '/dash/https/errors/count',
	getAvgPerformance = '/dash/performances',
	getEventGroup = '/dash/events/group',
}

export type ChangePasswordPayload = {
	currentPassword: string;
	newPassword: string;
};
export async function changeUserPassword(data: ChangePasswordPayload) {
	return post<void>(ApiPath.changeUserPassword, data);
}

export async function getPageviewInterval(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getPageviewInterval, timeRange);
}

export async function getPageviewCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getPageviewCount, timeRange);
}

export async function getUniqueVisitorInterval(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getUniqueVisitorInterval, timeRange);
}

export async function getUniqueVisitorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getUniqueVisitorCount, timeRange);
}

export async function getErrorInterval(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getErrorInterval, timeRange);
}

export async function getErrorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getErrorCount, timeRange);
}

export async function getHttpErrorInterval(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getHttpErrorInterval, timeRange);
}

export async function getHttpErrorCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getHttpErrorCount, timeRange);
}

export async function getAvgPerformance(timeRange: SpanFilter) {
	return get<IntervalData[]>(ApiPath.getAvgPerformance, timeRange);
}

export async function getEventGroup(timeRange: SpanFilter) {
	return get<IntervalData[]>(ApiPath.getEventGroup, timeRange);
}
