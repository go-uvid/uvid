import useSWR, {type Fetcher, type Key} from 'swr';
import {useAtom} from 'jotai';
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
	getErrorCount,
	getEventGroup,
	getHttpErrorInterval,
	getHttpErrorCount,
	getPageviewInterval,
	getPageviewCount,
	getUniqueVisitorInterval,
	getUniqueVisitorCount,
	PerformanceName,
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
	pv: getPageviewInterval,
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

export function usePageviewCount() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getPageviewCount, startTime, endTime], async () =>
		getPageviewCount({start: startTime, end: endTime}),
	);
}

export function useUniqueVisitorCount() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest(
		[ApiPath.getUniqueVisitorCount, startTime, endTime],
		async () => getUniqueVisitorCount({start: startTime, end: endTime}),
	);
}

export function useErrorCount() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getErrorCount, startTime, endTime], async () =>
		getErrorCount({start: startTime, end: endTime}),
	);
}

export function useHttpErrorCount() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getHttpErrorCount, startTime, endTime], async () =>
		getHttpErrorCount({start: startTime, end: endTime}),
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
