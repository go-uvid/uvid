import useSWR, {type Fetcher, type Key} from 'swr';
import {useAtom} from 'jotai';
import {groupBy} from 'lodash-es';
import {UAParser} from 'ua-parser-js';
import {
	type IntervalType,
	intervalTypeAtom,
	useSpanFilterPayload,
	useTimeIntervalPayload,
} from '../store';
import {type RequestError} from './request';
import {
	ApiPath,
	type IntervalData,
	type TimeRangeDTO,
	getAvgPerformance,
	getErrorInterval,
	getEventGroup,
	getHttpErrorInterval,
	getPageViewInterval,
	getUniqueVisitorInterval,
	getMetricCount,
	PerformanceName,
	getPageViews,
	getSessions,
} from './api';

export function useRequest<Data = any>(key: Key, fetcher: Fetcher<Data>) {
	const {data, error, isLoading, isValidating} = useSWR<Data, RequestError>(
		key,
		fetcher,
	);

	return {data, error, isLoading, isValidating};
}

const intervalFetcher: Record<
	IntervalType,
	(timeRange: TimeRangeDTO) => Promise<IntervalData[]>
> = {
	uv: getUniqueVisitorInterval,
	pv: getPageViewInterval,
	jsError: getErrorInterval,
	httpError: getHttpErrorInterval,
};

export function useIntervalData() {
	const {startTime, endTime, timeUnit} = useTimeIntervalPayload();
	const [intervalType] = useAtom(intervalTypeAtom);
	return useRequest([intervalType, startTime, endTime], async () =>
		intervalFetcher[intervalType]({
			start: startTime,
			end: endTime,
			unit: timeUnit,
		}),
	);
}

export function usePageViews() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getPageViews, startTime, endTime], async () =>
		getPageViews({start: startTime, end: endTime}).then((pvs) =>
			groupAndCount(pvs, ({url}) => new URL(url).pathname),
		),
	);
}

export function useSessions() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getSessions, startTime, endTime], async () =>
		getSessions({start: startTime, end: endTime}),
	);
}

function groupAndCount<T>(array?: T[], iteratee?: (item: T) => string) {
	if (!array) return undefined;
	const group = groupBy<T>(array, iteratee);
	return Object.keys(group).map((x) => ({
		x,
		y: group[x].length,
	}));
}

export function useReferrers() {
	const {data: sessions} = useSessions();
	return groupAndCount(sessions, (session) => session.referrer);
}

const UNKNOWN_TYPE = 'unknown';

export function useOSs() {
	const {data: sessions} = useSessions();
	return groupAndCount(
		sessions,
		(session) => new UAParser(session.ua).getOS().name ?? UNKNOWN_TYPE,
	);
}

export function useBrowsers() {
	const {data: sessions} = useSessions();
	return groupAndCount(
		sessions,
		(session) => new UAParser(session.ua).getBrowser().name ?? UNKNOWN_TYPE,
	);
}

export function useDevice() {
	const {data: sessions} = useSessions();
	return groupAndCount(
		sessions,
		(session) => new UAParser(session.ua).getDevice().type ?? UNKNOWN_TYPE,
	);
}

export function useMetricCount() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.metricCount, startTime, endTime], async () =>
		getMetricCount({start: startTime, end: endTime}),
	);
}

export function useAvgPerformance() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest(
		[ApiPath.getAvgPerformance, startTime, endTime],
		async () => {
			const performanceMetric = {
				[PerformanceName.LCP]: '0',
				[PerformanceName.CLS]: '0',
				[PerformanceName.FID]: '0',
			};
			const data = await getAvgPerformance({start: startTime, end: endTime});
			for (const {x, y} of data) {
				performanceMetric[x as keyof typeof performanceMetric] = y.toFixed(2);
			}

			return performanceMetric;
		},
	);
}

export function useEventGroup() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getEventGroup, startTime, endTime], async () =>
		getEventGroup({start: startTime, end: endTime}),
	);
}
