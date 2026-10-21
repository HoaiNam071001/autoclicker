// Bảng đăng ký hoạt hoạ minh hoạ: loại node → cảnh (bước 47). Thiếu loại nào thì illus.test.ts đỏ.
import type { ReactNode } from 'react';
import type { InfoKey } from '../../../flow/nodeInfo';
import { app } from './app';
import { data } from './data';
import { flow } from './flow';
import { mouse } from './mouse';
import { wait } from './wait';
import { web } from './web';

export const ILLUS: Partial<Record<InfoKey, () => ReactNode>> = { ...mouse, ...wait, ...data, ...app, ...flow, ...web };
