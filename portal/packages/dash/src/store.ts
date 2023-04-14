import {signal} from '@preact/signals-core';
import dayjs from 'dayjs';

export const startTime = signal(dayjs().subtract(1, 'day').toISOString());
export const endTime = signal(dayjs().toISOString());
