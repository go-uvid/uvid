import {randomUUID} from 'node:crypto';
import {test, expect, type Page} from '@playwright/test';
import {
	type ErrorDTO,
	type BaseSessionDTO,
	type PerformanceDTO,
	type EventDTO,
} from '../lib/types/span';

const apiHost = 'http://localhost:3000';

const pageUrl = `http://localhost:4000/basic.test`;

type Data = {
	url: string;
	body: any;
};
const actualData: Data[] = [];

test.beforeEach(async ({page}) => {
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
			const uuid = randomUUID();
			await route.fulfill({body: uuid, status: 204});
		} else {
			await route.fulfill({
				status: 204,
			});
		}
	});
	// WARN route before goto page, otherwise we can't catch session request
	await page.goto(pageUrl);
});

test('basic', async ({page}) => {
	const {referrer, language, screen} = await page.evaluate(() => ({
		referrer: document.referrer,
		language: navigator.language,
		screen: `${window.screen.width}*${window.screen.height}`,
	}));
	const url = page.url();

	await page.locator('#fid').click();
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
	await page.locator('#register').click();
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
