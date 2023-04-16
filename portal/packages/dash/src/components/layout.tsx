import {Layout} from 'antd';
import {HomeOutlined, SettingOutlined} from '@ant-design/icons';
import {type PropsWithChildren} from 'react';
import {Link} from 'react-router-dom';

const {Header, Content} = Layout;

export function RootLayout(props: PropsWithChildren) {
	return (
		<Layout>
			<Header className="flex justify-center">
				<Content className="w-main flex justify-between flex-grow-0">
					<Link to="/">
						<HomeOutlined className="text-base" />
					</Link>
					<Link to="/setting">
						<SettingOutlined className="text-base" />
					</Link>
				</Content>
			</Header>
			<Content className="m-auto">{props.children}</Content>
		</Layout>
	);
}
