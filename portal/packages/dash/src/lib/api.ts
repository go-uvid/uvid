import {AUTH_TOKEN_KEY, baseRequest, get, goLogin, post} from './request';

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
type BaseSessionDTO = {
	appVersion: string;
	url: string;
	screen: string;
	referrer: string;
	language: string;
	meta?: string;
};

type SessionDTO = {
	ua: string;
	ip: string;
} & BaseSessionDTO;

export enum ApiPath {
	changeUserPassword = '/dash/user/password',
	getPageViews = '/dash/pvs',
	getPageViewInterval = '/dash/pvs/interval',
	getPageViewCount = '/dash/pvs/count',
	getUniqueVisitorInterval = '/dash/uvs/interval',
	getUniqueVisitorCount = '/dash/uvs/count',
	getErrorInterval = '/dash/errors/interval',
	getErrorCount = '/dash/errors/count',
	getHttpErrorInterval = '/dash/https/errors/interval',
	getHttpErrorCount = '/dash/https/errors/count',
	getAvgPerformance = '/dash/performances',
	getEventGroup = '/dash/events/group',
	getSessions = '/dash/sessions',
}

export type ChangePasswordPayload = {
	currentPassword: string;
	newPassword: string;
};

export async function login(name: string, password: string) {
	return baseRequest<{token: string}>(
		'post',
		'/dash/user/login',
		JSON.stringify({name, password}),
	).then((response) => {
		const {token} = response.data;
		localStorage.setItem(AUTH_TOKEN_KEY, token);
	});
}

export function logout() {
	localStorage.removeItem(AUTH_TOKEN_KEY);
	void goLogin();
}

export async function changeUserPassword(data: ChangePasswordPayload) {
	return post<void>(ApiPath.changeUserPassword, data);
}

type PageViewItem = {
	url: string;
};
export async function getPageViews(timeRange: SpanFilter) {
	const pvs = await get<PageViewItem[]>(ApiPath.getPageViews, timeRange);
	return pvs;
}

export async function getPageViewInterval(timeRange: TimeRangeDTO) {
	return get<IntervalData[]>(ApiPath.getPageViewInterval, timeRange);
}

export async function getPageViewCount(timeRange: SpanFilter) {
	return get<number>(ApiPath.getPageViewCount, timeRange);
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

export async function getSessions(timeRange: SpanFilter) {
	return get<SessionDTO[]>(ApiPath.getSessions, timeRange);
}
