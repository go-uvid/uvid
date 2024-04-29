import { get } from "./request";

type SpanFilter = {
  start: string;
  end: string;
};
export type TimeUnit = "hour" | "day" | "month" | "year";
export type TimeRangeDTO = SpanFilter & {
  unit: TimeUnit;
};

export enum PerformanceName {
  LCP = "LCP",
  CLS = "CLS",
  FID = "FID",
}

export type IntervalData = {
  x: string;
  y: number;
};
type BaseSessionDTO = {
  appVersion: string;
  url: string;
  screen: string;
  referrer: string;
  language: string;
  meta?: string;
};

type SessionDTO = {
  ua: string;
  ip: string;
} & BaseSessionDTO;

export enum ApiPath {
  getPageViews = "/dash/pvs",
  getPageViewInterval = "/dash/pvs/interval",
  getUniqueVisitorInterval = "/dash/uvs/interval",
  getErrorInterval = "/dash/errors/interval",
  getHttpErrorInterval = "/dash/https/errors/interval",
  getAvgPerformance = "/dash/performances",
  getEventGroup = "/dash/events/group",
  getSessions = "/dash/sessions",
  metricCount = "/dash/metric/count",
}

type PageViewItem = {
  url: string;
};
export async function getPageViews(timeRange: SpanFilter) {
  const pvs = await get<PageViewItem[]>(ApiPath.getPageViews, timeRange);
  return pvs;
}

export async function getPageViewInterval(timeRange: TimeRangeDTO) {
  return get<IntervalData[]>(ApiPath.getPageViewInterval, timeRange);
}

export async function getUniqueVisitorInterval(timeRange: TimeRangeDTO) {
  return get<IntervalData[]>(ApiPath.getUniqueVisitorInterval, timeRange);
}

export async function getErrorInterval(timeRange: TimeRangeDTO) {
  return get<IntervalData[]>(ApiPath.getErrorInterval, timeRange);
}

export async function getHttpErrorInterval(timeRange: TimeRangeDTO) {
  return get<IntervalData[]>(ApiPath.getHttpErrorInterval, timeRange);
}

export async function getMetricCount(timeRange: SpanFilter) {
  return get<{
    pv: number;
    uv: number;
    jsError: number;
    httpError: number;
  }>(ApiPath.metricCount, timeRange);
}

export async function getAvgPerformance(timeRange: SpanFilter) {
  return get<IntervalData[]>(ApiPath.getAvgPerformance, timeRange);
}

export async function getEventGroup(timeRange: SpanFilter) {
  return get<IntervalData[]>(ApiPath.getEventGroup, timeRange);
}

export async function getSessions(timeRange: SpanFilter) {
  return get<SessionDTO[]>(ApiPath.getSessions, timeRange);
}
