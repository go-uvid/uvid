import {
	LineChart,
	Line,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip,
	ResponsiveContainer,
} from 'recharts';
import {Empty} from 'antd';
import {type IntervalData} from '../lib/api';
import {Theme} from '../lib/theme';

export function IntervalLineChart(props: {data?: IntervalData[]}) {
	const {data} = props;
	return (
		<ResponsiveContainer
			width="100%"
			height={368}
			className="mt-4 flex items-center justify-center"
		>
			{data ? (
				<LineChart
					data={data}
					margin={{
						top: 10,
						right: 55,
						left: 0,
					}}
				>
					<CartesianGrid strokeDasharray="3 3" />
					<XAxis dataKey="x" />
					<YAxis />
					<Tooltip />
					<Line
						type="monotone"
						dataKey="y"
						stroke={Theme.color.primary}
						activeDot={{r: 8}}
					/>
				</LineChart>
			) : (
				<Empty />
			)}
		</ResponsiveContainer>
	);
}
