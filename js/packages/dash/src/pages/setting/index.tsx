import {Form, Input, Button, message, Layout, Tabs, type TabsProps} from 'antd';
import {LockOutlined} from '@ant-design/icons';
import {
	changeUserPassword,
	type ChangePasswordPayload,
	logout,
} from '../../lib/api';

const operations = <Button onClick={logout}>Log out</Button>;

export function Setting() {
	const [form] = Form.useForm();

	const onFinish = async (values: ChangePasswordPayload) => {
		try {
			await changeUserPassword(values);
			form.resetFields();
			await message.success('Change password success');
		} catch (error) {
			console.error(error);
			await message.error('Change password failed');
		}
	};

	const tabItems: TabsProps['items'] = [
		{
			label: `Change password`,
			key: '0',
			children: (
				<Form form={form} onFinish={onFinish} className="w-80">
					<Form.Item
						name="currentPassword"
						rules={[
							{
								required: true,
							},
						]}
					>
						<Input.Password
							prefix={<LockOutlined />}
							type="password"
							placeholder="Current password"
						/>
					</Form.Item>
					<Form.Item
						name="newPassword"
						rules={[
							{
								required: true,
							},
						]}
					>
						<Input.Password
							prefix={<LockOutlined />}
							type="password"
							placeholder="New password"
						/>
					</Form.Item>
					<Form.Item
						name="newPassword2"
						rules={[
							{
								required: true,
							},
							({getFieldValue}) => ({
								async validator(_, value) {
									if (!value || getFieldValue('newPassword') === value) {
										return;
									}

									throw new Error(
										'The two passwords that you entered do not match!',
									);
								},
							}),
						]}
					>
						<Input.Password
							prefix={<LockOutlined />}
							type="password"
							placeholder="Confirm password"
						/>
					</Form.Item>

					<Button type="primary" htmlType="submit">
						Submit
					</Button>
				</Form>
			),
		},
	];
	return (
		<div className="w-main">
			<Tabs
				tabPosition="left"
				items={tabItems}
				tabBarExtraContent={operations}
			/>
		</div>
	);
}
