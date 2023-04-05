export type ErrorDTO = {
	name: string;
	message: string;
	stack: string;
	cause?: string;
};

export type HTTPDTO = {
	resource: string;
	method: string;
	headers: string;
	status: number;
	data?: string;
	response?: string;
};

export type EventDTO = {
	name: string;
	value?: string;
};

export type PerformanceDTO = {
	name: 'LCP' | 'CLS' | 'FID';
	value: number;
	url: string;
};

export type PageViewDTO = {
	url: string;
};

export type BaseSessionDTO = {
	appVersion?: string;
	url: string;
	screen: string;
	referrer: string;
	language: string;
	meta?: string;
};

export type SessionDTO = {
	ua: string;
	ip: string;
} & BaseSessionDTO;

export type HTTPHeaders = Record<string, string>;
