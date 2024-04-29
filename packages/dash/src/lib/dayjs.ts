import dayjs from 'dayjs';
import localeData from 'dayjs/plugin/localeData';

export function setupDayjs() {
	dayjs.extend(localeData);
	dayjs.locale(navigator.language, {
		weekStart: 1,
	});
}
