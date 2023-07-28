export async function retryPromise<T>(
	cb: () => Promise<T>,
	maxAttempts = 3,
	retryIntervalMs = 1000,
): Promise<T> {
	if (maxAttempts <= 0) {
		throw new RetryError('Max attempts reached');
	}

	try {
		const result = await cb();
		return result;
	} catch (error) {
		console.warn(
			`Attempt failed. Retrying in ${retryIntervalMs}ms. Error:`,
			error,
		);

		// eslint-disable-next-line no-promise-executor-return
		await new Promise((resolve) => setTimeout(resolve, retryIntervalMs));
		return retryPromise(cb, --maxAttempts, retryIntervalMs);
	}
}

export class RetryError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'RetryError';
	}
}

export function assert(
	condition: unknown,
	message = 'Assertion failed',
): asserts condition {
	if (!condition) {
		throw new Error(message);
	}
}

export function domReady(callback: () => void) {
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', function fn() {
			document.removeEventListener('DOMContentLoaded', fn);
			callback();
		});
	} else {
		callback();
	}
}
