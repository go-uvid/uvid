import {get} from './request';

type TimeRangeDTO = {
	start: string;
	end: string;
	unit: 'hour' | 'day' | 'week' | 'month';
};

export async function getPageview(timeRange: TimeRangeDTO) {
	return get('/dash/pageview', timeRange);
}
