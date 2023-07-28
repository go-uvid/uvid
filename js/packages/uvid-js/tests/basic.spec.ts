import process from 'node:process';
import path from 'node:path';
import {randomUUID} from 'node:crypto';
import {test, expect, type Page} from '@playwright/test';
import {
	type ErrorDTO,
	type BaseSessionDTO,
	type PerformanceDTO,
	type EventDTO,
} from '../lib/types/span';
import {serve} from './serve.js';

const apiHost = 'http://localhost:3000';

const publicPath = path.join(process.cwd());
const pageUrl = `http://localhost:4000/tests/basic`;

type Data = {
	url: string;
	body: any;
};
const actualData: Data[] = [];

test.beforeEach(async ({page}) => {
	// Serve html page, since we cannot add cookie when using file protocol
	await serve(publicPath);
	await page.goto(pageUrl);
	await page.route(`${apiHost}/**/*`, async (route) => {
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

		if (route.request().url().includes('/span/session')) {
			await page.context().addCookies([
				{
					name: 'uvid-session',
					value: randomUUID(),
					url: pageUrl,
				},
			]);
		}

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
	await page.waitForRequest(`${apiHost}/span/performance`);

	// Using globalThis.atob will cause error: TypeError: Cannot set property message of  which has only a getter
	const session: BaseSessionDTO = {
		url: pageUrl,
		screen,
		referrer,
		language,
		meta: '{}',
		appVersion: '',
	};

	const url = await page.evaluate(() => window.location.href);
	const firstLCP = await getPerformanceSpan(page, 'LCP');
	const fid = await getPerformanceSpan(page, 'FID');
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
	// Test if session exist after reload
	await page.reload();
	const registerButton = await page.$('#register');
	await registerButton?.click();
	const secondLCP = await getPerformanceSpan(page, 'LCP');

	const error: ErrorDTO = {
		name: 'Error',
		message: errorMessage,
		stack: errorStack!,
	};
	const pageViewData = {
		url: `${apiHost}/span/pageview`,
		body: {
			url,
		},
	};
	const expectData: Data[] = [
		{
			url: `${apiHost}/span/session`,
			body: session,
		},
		pageViewData,
		{
			url: `${apiHost}/span/performance`,
			body: firstLCP,
		},
		{
			url: `${apiHost}/span/performance`,
			body: fid,
		},
		{
			url: `${apiHost}/span/error`,
			body: error,
		},
		{
			url: `${apiHost}/span/event`,
			body: testEvent,
		},
		// After page reload
		{
			url: `${apiHost}/span/performance`,
			body: {
				name: 'CLS',
				url: pageUrl,
				value: 0,
			},
		},
		pageViewData,
		{
			url: `${apiHost}/span/performance`,
			body: secondLCP,
		},
		{
			url: `${apiHost}/span/event`,
			body: registerEvent,
		},
	];
	expect(actualData).toEqual(expectData);
});

async function getPerformanceSpan(page: Page, name: PerformanceDTO['name']) {
	const url = await page.evaluate(() => window.location.href);
	const value = await page.evaluate(
		([_name]) => {
			return window.uvid.performanceMetrics.find((m) => m.name === _name)!
				.value;
		},
		[name],
	);
	const span: PerformanceDTO = {
		name,
		value,
		url,
	};
	return span;
}
