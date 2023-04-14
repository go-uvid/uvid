export function getQueryString(parameters: Record<string, any> = {}): string {
	return Object.keys(parameters)
		.reduce((array: string[], key: string) => {
			if (parameters[key] !== undefined) {
				// eslint-disable-next-line @typescript-eslint/no-unsafe-argument
				return array.concat(`${key}=${encodeURIComponent(parameters[key])}`);
			}

			return array;
		}, [])
		.join('&');
}

export function buildUrl(
	url: string,
	parameters: Record<string, unknown> = {},
): string {
	const queryString = getQueryString(parameters);
	return `${url}${queryString && '?' + queryString}`;
}
