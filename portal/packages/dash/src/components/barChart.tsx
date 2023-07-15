import {
	ResponsiveContainer,
	Bar,
	YAxis,
	XAxis,
	ComposedChart,
	LabelList,
} from 'recharts';
import {Empty} from 'antd';
import {type IntervalData} from '../lib/api';
import {Theme} from '../lib/theme';

export function GroupBarChart(props: {data?: IntervalData[]}) {
	const {data} = props;
	return (
		<ResponsiveContainer
			width="100%"
			height={400}
			className="flex items-center justify-center"
		>
			{data ? (
				<ComposedChart
					layout="vertical"
					data={data}
					margin={{
						right: 20,
					}}
				>
					<XAxis type="number" hide />
					<YAxis dataKey="x" type="category" hide />
					<Bar dataKey="y" barSize={20} fill={Theme.color.primary}>
						<LabelList dataKey="x" position="insideLeft" fill="#ffffff" />
						<LabelList dataKey="y" position="right" />
					</Bar>
				</ComposedChart>
			) : (
				<Empty />
			)}
		</ResponsiveContainer>
	);
}
