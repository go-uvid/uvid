import dayjs from 'dayjs';
import {atom, useAtom, useAtomValue, useSetAtom} from 'jotai';
import {type TimeUnit} from './lib/api';
import {setupDayjs} from './lib/dayjs';

const [defaultStartTime, defaultEndTime, defaultTimeRange] = thisWeek();
export const startTimeAtom = atom(defaultStartTime);
export const endTimeAtom = atom(defaultEndTime);
export const timeUnitAtom = atom<TimeUnit>('day');

const timeRangeAtom = atom<TimeRange>(defaultTimeRange);
export const intervalTypeAtom = atom<IntervalType>('uv');

export function useSpanFilterPayload() {
	const startTime = useAtomValue(startTimeAtom);
	const endTime = useAtomValue(endTimeAtom);
	return {
		startTime,
		endTime,
	};
}

export function useTimeIntervalPayload() {
	const {startTime, endTime} = useSpanFilterPayload();
	const timeUnit = useAtomValue(timeUnitAtom);
	return {
		startTime,
		endTime,
		timeUnit,
	};
}

export function useTimeRange() {
	const setStartTime = useSetAtom(startTimeAtom);
	const setEndTime = useSetAtom(endTimeAtom);
	const [timeRange, _setTimeRange] = useAtom(timeRangeAtom);
	const setTimeUnit = useSetAtom(timeUnitAtom);

	function setTimeRange(range: TimeRange) {
		_setTimeRange(range);
		const now = dayjs();
		switch (range) {
			case 'today': {
				setStartTime(now.startOf('day').toISOString());
				setEndTime(now.endOf('day').toISOString());
				setTimeUnit('hour');
				break;
			}

			case 'yesterday': {
				setStartTime(now.subtract(1, 'day').startOf('day').toISOString());
				setEndTime(now.subtract(1, 'day').endOf('day').toISOString());
				setTimeUnit('hour');
				break;
			}

			case 'thisWeek': {
				setStartTime(now.startOf('week').toISOString());
				setEndTime(now.endOf('week').toISOString());
				setTimeUnit('day');
				break;
			}

			case 'thisMonth': {
				setStartTime(now.startOf('month').toISOString());
				setEndTime(now.endOf('month').toISOString());
				setTimeUnit('day');
				break;
			}

			case 'thisYear': {
				setStartTime(now.startOf('year').toISOString());
				setEndTime(now.endOf('year').toISOString());
				setTimeUnit('day');
				break;
			}

			case 'allTime': {
				setStartTime(dayjs(0).toISOString());
				setEndTime(now.toISOString());
				setTimeUnit('day');
				break;
			}

			default: {
				throw new Error('Invalid time range');
			}
		}
	}

	return {
		timeRange,
		setTimeRange,
	};
}

export type TimeRange =
	| 'today'
	| 'yesterday'
	| 'thisWeek'
	| 'thisMonth'
	| 'thisYear'
	| 'allTime';

export type IntervalType = 'uv' | 'pv' | 'jsError' | 'httpError';

function thisWeek(): [string, string, TimeRange] {
	setupDayjs();
	const now = dayjs();
	return [
		now.startOf('week').toISOString(),
		now.endOf('week').toISOString(),
		'thisWeek',
	];
}
