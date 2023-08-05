import {Card, Col, Divider, Layout, Row, Select, Statistic} from 'antd';
import {useAtom} from 'jotai';
import {type ReactNode, type PropsWithChildren} from 'react';
import {
	useAvgPerformance,
	useBrowsers,
	useDevice,
	useErrorCount,
	useEventGroup,
	useHttpErrorCount,
	useIntervalData,
	useOSs,
	usePageViewCount,
	usePageViews,
	useReferrers,
	useUniqueVisitorCount,
} from '../../lib/useApi';
import {
	type TimeRange,
	intervalTypeAtom,
	useTimeRange,
	type IntervalType,
} from '../../store';
import {IntervalLineChart} from '../../components/lineChart';
import {GroupBarChart} from '../../components/barChart';

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
	const {data: uvCount} = useUniqueVisitorCount();
	const {data: pv} = usePageViews();
	const {data: pvCount} = usePageViewCount();
	const {data: errorCount} = useErrorCount();
	const {data: httpErrorCount} = useHttpErrorCount();
	const {data: performance} = useAvgPerformance();
	const {data: events} = useEventGroup();
	const referrers = useReferrers();
	const oss = useOSs();
	const devices = useDevice();
	const browsers = useBrowsers();

	function handleChange(value: TimeRange) {
		setTimeRange(value);
	}

	function changeIntervalType(type: IntervalType) {
		setIntervalType(type);
	}

	return (
		<Layout>
			<header>
				<div className="w-main m-auto flex items-center justify-between">
					<h1 className="text-xl"></h1>
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
							value={uvCount}
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
									Total page views
								</StatisticTitle>
							}
							value={pvCount}
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
						<Statistic title="LCP" value={performance?.LCP} suffix="s" />
					</Card.Grid>
					<Card.Grid style={gridStyle} hoverable={false}>
						<Statistic title="CLS" value={performance?.CLS} suffix="s" />
					</Card.Grid>
					<Card.Grid style={gridStyle} hoverable={false}>
						<Statistic title="FID" value={performance?.FID} suffix="s" />
					</Card.Grid>
					<div className="w-main px-3">
						<IntervalLineChart data={intervalData} />
					</div>
					<Divider />
					<ChartGroup
						left={
							<>
								<GroupChartTitle left="Events" />
								<GroupBarChart data={events} />
							</>
						}
						right={
							<>
								<GroupChartTitle left="Pages" />
								<GroupBarChart data={pv} />
							</>
						}
					/>
					<Divider />
					<ChartGroup
						left={
							<>
								<GroupChartTitle left="Referrers" />
								<GroupBarChart data={referrers} />
							</>
						}
						right={
							<>
								<GroupChartTitle left="OS" />
								<GroupBarChart data={oss} />
							</>
						}
					/>
					<ChartGroup
						left={
							<>
								<GroupChartTitle left="Browsers" />
								<GroupBarChart data={browsers} />
							</>
						}
						right={
							<>
								<GroupChartTitle left="Devices" />
								<GroupBarChart data={devices} />
							</>
						}
					/>
				</Card>
			</Content>
		</Layout>
	);
}

function ChartGroup({left, right}: {left: ReactNode; right: ReactNode}) {
	return (
		<Row className="w-main px-6">
			<Col span={11}>{left}</Col>
			<Col span={2} className="flex items-center justify-center">
				<Divider type="vertical" className="h-full" />
			</Col>
			<Col span={11}>{right}</Col>
		</Row>
	);
}

function GroupChartTitle({left}: {left: string}) {
	return (
		<h4 className="flex justify-between">
			<span className="text-base text-primary">{left}</span>
		</h4>
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
