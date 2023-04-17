import {ResponsiveContainer, Bar, YAxis, XAxis, ComposedChart} from 'recharts';
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
				<ComposedChart layout="vertical" data={data}>
					<XAxis type="number" hide />
					<YAxis dataKey="x" type="category" />
					<Bar
						dataKey="y"
						barSize={20}
						fill={Theme.color.primary}
						label={Label}
					/>
				</ComposedChart>
			) : (
				<Empty />
			)}
		</ResponsiveContainer>
	);
}

function Label(props: LabelProps) {
	const {x, y, width, height, value} = props;
	return (
		<text
			x={x + width}
			y={y + height / 2}
			dx={5}
			fill={Theme.text.primary}
			fontSize={12}
		>
			{value}
		</text>
	);
}

export type LabelProps = {
	x: number;
	y: number;
	width: number;
	height: number;
	index: number;
	value: number;
	viewBox: ViewBox;
	offset: number;
	stroke?: string;
};

export type ViewBox = {
	x: number;
	y: number;
	width: number;
	height: number;
};
