import {Card, Layout, Select, Statistic} from 'antd';
import {useAtom} from 'jotai';
import {type PropsWithChildren} from 'react';
import {
	useAvgPerformance,
	useErrorCount,
	useHttpErrorCount,
	useIntervalData,
	usePageviewCount,
	useUniqueVisitorCount,
} from '../../lib/useApi';
import {
	type TimeRange,
	intervalTypeAtom,
	useTimeRange,
	type IntervalType,
} from '../../store';
import {IntervalLineChart} from '../../components/lineChart';

const {Content, Header} = Layout;
const gridStyle: React.CSSProperties = {
	width: '14.285%',
	textAlign: 'center',
};
const options: Array<{value: TimeRange; label: string}> = [
	{value: 'today', label: 'Today'},
	{value: 'yesterday', label: 'Yesterday'},
	{value: 'thisWeek', label: 'This Week'},
	{value: 'thisMonth', label: 'This Month'},
	{value: 'thisYear', label: 'This Year'},
	{value: 'allTime', label: 'All Time'},
];

export function Home() {
	const {data: intervalData} = useIntervalData();
	const {timeRange, setTimeRange} = useTimeRange();
	const [intervalType, setIntervalType] = useAtom(intervalTypeAtom);
	const {data: uv} = useUniqueVisitorCount();
	const {data: pv} = usePageviewCount();
	const {data: errorCount} = useErrorCount();
	const {data: httpErrorCount} = useHttpErrorCount();
	const {data: performance} = useAvgPerformance();

	function handleChange(value: TimeRange) {
		setTimeRange(value);
	}

	function changeIntervalType(type: IntervalType) {
		setIntervalType(type);
	}

	console.log('intervalData :>>', intervalData);

	return (
		<Layout>
			<header className="sticky top-0 w-screen">
				<div className="w-main m-auto flex items-center justify-between">
					<h1 className="text-xl">{location.hostname}</h1>
					<Select
						style={{width: 160}}
						value={timeRange}
						onChange={handleChange}
						options={options}
					/>
				</div>
			</header>
			<Content className="w-main m-auto mt-2">
				<Card>
					<Card.Grid
						style={gridStyle}
						className="cursor-pointer"
						onClick={() => {
							changeIntervalType('uv');
						}}
					>
						<Statistic
							title={
								<StatisticTitle active={intervalType === 'uv'}>
									Unique visitors
								</StatisticTitle>
							}
							value={uv}
						/>
					</Card.Grid>
					<Card.Grid
						style={gridStyle}
						className="cursor-pointer"
						onClick={() => {
							changeIntervalType('pv');
						}}
					>
						<Statistic
							title={
								<StatisticTitle active={intervalType === 'pv'}>
									Total pageviews
								</StatisticTitle>
							}
							value={pv}
						/>
					</Card.Grid>
					<Card.Grid
						style={gridStyle}
						className="cursor-pointer"
						onClick={() => {
							changeIntervalType('jsError');
						}}
					>
						<Statistic
							title={
								<StatisticTitle active={intervalType === 'jsError'}>
									JS errors
								</StatisticTitle>
							}
							value={errorCount}
						/>
					</Card.Grid>
					<Card.Grid
						style={gridStyle}
						className="cursor-pointer"
						onClick={() => {
							changeIntervalType('httpError');
						}}
					>
						<Statistic
							title={
								<StatisticTitle active={intervalType === 'httpError'}>
									HTTP errors
								</StatisticTitle>
							}
							value={httpErrorCount}
						/>
					</Card.Grid>
					<Card.Grid style={gridStyle} hoverable={false}>
						<Statistic title="LCP" value={performance?.LCP} />
					</Card.Grid>
					<Card.Grid style={gridStyle} hoverable={false}>
						<Statistic title="CLS" value={performance?.CLS} />
					</Card.Grid>
					<Card.Grid style={gridStyle} hoverable={false}>
						<Statistic title="FID" value={performance?.FID} />
					</Card.Grid>
					<IntervalLineChart data={intervalData} />
				</Card>
			</Content>
		</Layout>
	);
}

function StatisticTitle(
	props: PropsWithChildren & {
		active: boolean;
	},
) {
	const {active, children} = props;
	return (
		<div
			className={
				active
					? 'underline text-primary underline-offset-4 font-semibold transition-all'
					: ''
			}
		>
			{children}
		</div>
	);
}
