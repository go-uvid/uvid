import useSWR, {type Fetcher, type Key} from 'swr';
import {useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import {message} from 'antd';
import {endTime, startTime} from '../store';
import {type RequestError} from './request';
import {getPageview} from './api';

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

export function usePageview() {
	return useRequest(
		() => `/dash/pageview/${startTime.value}/${endTime.value}`,
		async () =>
			getPageview({start: startTime.value, end: endTime.value, unit: 'hour'}),
	);
}
