import process from 'node:process';
import path from 'node:path';
import {test, expect} from '@playwright/test';
import {type BaseSessionDTO, type PerformanceDTO} from '../lib/types/span';

const host = 'http://localhost:3000';

const htmlPath = path.join(process.cwd(), 'tests/basic.html');
const pageUrl = `file://${htmlPath}`;

type Data = {
	url: string;
	body: any;
};
const actualData: Data[] = [];

test.beforeEach(async ({page}) => {
	await page.goto(pageUrl);
	await page.route(`${host}/**/*`, async (route) => {
		const postData = route.request().postData();
		const url = route.request().url();
		expect(postData).toBeTruthy();
		if (!postData) return;
		// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
		const body = JSON.parse(postData);
		actualData.push({
			url,
			// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
			body,
		});

		await route.fulfill({
			status: 204,
		});
	});
});

test('basic', async ({page}) => {
	const referrer = await page.evaluate(() => document.referrer);
	const language = await page.evaluate(() => navigator.language);
	const screen = await page.evaluate(
		() => `${window.screen.width}*${window.screen.height}`,
	);

	const button = await page.$('#fid');
	await button?.click();
	await page.waitForRequest(`${host}/span/performance`);

	const session: BaseSessionDTO = {
		url: atob(pageUrl),
		screen,
		referrer,
		language,
		meta: '{}',
		appVersion: '',
	};
	const lcpValue = await page.evaluate(
		() => window.uvid.performanceMetrics.find((m) => m.name === 'LCP')!.value,
	);
	const fidValue = await page.evaluate(
		() => window.uvid.performanceMetrics.find((m) => m.name === 'FID')!.value,
	);
	const url = atob(await page.evaluate(() => window.location.href));

	const lcp: PerformanceDTO = {
		name: 'LCP',
		value: lcpValue,
		url,
	};
	const fid: PerformanceDTO = {
		name: 'FID',
		value: fidValue,
		url,
	};
	const expectData: Data[] = [
		{
			url: `${host}/span/session`,
			body: session,
		},
		{
			url: `${host}/span/performance`,
			body: lcp,
		},
		{
			url: `${host}/span/performance`,
			body: fid,
		},
	];
	expect(actualData).toEqual(expectData);
});
