import {
	LineChart,
	Line,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip,
	ResponsiveContainer,
} from 'recharts';
import {type IntervalData} from '../lib/api';

export function IntervalLineChart(props: {data?: IntervalData[]}) {
	const {data} = props;
	if (!data) return null;
	return (
		<ResponsiveContainer width="100%" height={368} className="mt-4">
			<LineChart
				data={data}
				margin={{
					top: 10,
					right: 60,
					left: 0,
				}}
			>
				<CartesianGrid strokeDasharray="3 3" />
				<XAxis dataKey="x" />
				<YAxis />
				<Tooltip />
				<Line type="monotone" dataKey="y" stroke="#1890ff" activeDot={{r: 8}} />
			</LineChart>
		</ResponsiveContainer>
	);
}
