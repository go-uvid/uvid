import process from 'node:process';
import path from 'node:path';
import {test, expect} from '@playwright/test';
import {
	type ErrorDTO,
	type BaseSessionDTO,
	type PerformanceDTO,
	type EventDTO,
} from '../lib/types/span';

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

	// Using globalThis.atob will cause error: TypeError: Cannot set property message of  which has only a getter
	const session: BaseSessionDTO = {
		url: pageUrl,
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
	const url = await page.evaluate(() => window.location.href);
	const errorMessage = 'test error';
	const errorStack = await page.evaluate(async (errorMessage_: string) => {
		const error = new Error(errorMessage_);
		await window.uvid.error(error);
		return error.stack;
	}, errorMessage);
	const testEvent: EventDTO = {
		action: 'event-name',
		value: 'event-value',
	};
	await page.evaluate(async (testEvent_) => {
		await window.uvid.event(testEvent_.action, testEvent_.value);
	}, testEvent);
	const registerEvent: EventDTO = {
		action: 'register',
		value: 'test',
	};
	const registerButton = await page.$('#register');
	await registerButton?.click();

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
	const error: ErrorDTO = {
		name: 'Error',
		message: errorMessage,
		stack: errorStack!,
	};
	const expectData: Data[] = [
		{
			url: `${host}/span/session`,
			body: session,
		},
		{
			url: `${host}/span/pageview`,
			body: {
				url,
			},
		},
		{
			url: `${host}/span/performance`,
			body: lcp,
		},
		{
			url: `${host}/span/performance`,
			body: fid,
		},
		{
			url: `${host}/span/error`,
			body: error,
		},
		{
			url: `${host}/span/event`,
			body: testEvent,
		},
		{
			url: `${host}/span/event`,
			body: registerEvent,
		},
	];
	expect(actualData).toEqual(expectData);
});
