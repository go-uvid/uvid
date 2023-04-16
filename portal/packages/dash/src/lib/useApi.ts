import useSWR, {type Fetcher, type Key} from 'swr';
import {useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import {message} from 'antd';
import {useAtom} from 'jotai';
import {
	type IntervalType,
	intervalTypeAtom,
	useSpanFilterPayload,
	useTimeIntervalPayload,
	useTimeRange,
} from '../store';
import {type RequestError} from './request';
import {
	ApiPath,
	type IntervalData,
	type TimeRangeDTO,
	getAvgPerformance,
	getError,
	getErrorCount,
	getEvent,
	getHttpError,
	getHttpErrorCount,
	getPageview,
	getPageviewCount,
	getUniqueVisitor,
	getUniqueVisitorCount,
	PerformanceName,
} from './api';

export function useRequest<Data = any>(key: Key, fetcher: Fetcher<Data>) {
	const navigate = useNavigate();

	const {data, error, isLoading, isValidating} = useSWR<Data, RequestError>(
		key,
		fetcher,
	);
	useEffect(() => {
		if (error?.status === 401) {
			navigate('/login');
			void message.warning('Please login first');
		}
	}, [error]);

	return {data, error, isLoading, isValidating};
}

const intervalFetcher: Record<
	IntervalType,
	(timeRange: TimeRangeDTO) => Promise<IntervalData[]>
> = {
	uv: getUniqueVisitor,
	pv: getPageview,
	jsError: getError,
	httpError: getHttpError,
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

export function useEvent() {
	const {startTime, endTime} = useSpanFilterPayload();
	return useRequest([ApiPath.getEvent, startTime, endTime], async () =>
		getEvent({start: startTime, end: endTime}),
	);
}
