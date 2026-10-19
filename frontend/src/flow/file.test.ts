import { expect, it } from 'vitest';
import { fmtFile, folderOf, isDeleteOp } from './file';

it('fmtFile', () => {
  expect(fmtFile({ filePath: '~/Documents/nhat-ky.txt' })).toBe('ghi nhat-ky.txt');
  expect(fmtFile({ filePath: '~/a.txt', fileAppend: true })).toBe('thêm vào cuối a.txt');
  expect(fmtFile({ fileOp: 'row', filePath: 'C:\\Data\\don.xlsx', cells: ['a', 'b'] })).toBe('thêm dòng 2 ô → don.xlsx');
  expect(fmtFile({ fileOp: 'trash', filePath: '~/Downloads/', filePattern: '*.png' })).toBe('xoá *.png trong Downloads → Thùng rác');
  expect(fmtFile({ fileOp: 'mkdir' })).toBe('tạo thư mục (chưa chọn)');
  expect(fmtFile({ fileOp: 'emptyTrash' })).toBe('dọn Thùng rác (xoá hẳn)');
});

it('folderOf', () => {
  expect(folderOf('~/Documents/autoclicker/nhat-ky.txt')).toBe('~/Documents/autoclicker');
  expect(folderOf('C:\\Data\\don.xlsx')).toBe('C:\\Data');
  expect(folderOf('~/Downloads')).toBe('~/Downloads');
  expect(folderOf('{{thu_muc}}')).toBe('{{thu_muc}}');
});

it('isDeleteOp', () => {
  expect(isDeleteOp('trash') && isDeleteOp('emptyTrash') && !isDeleteOp('row') && !isDeleteOp(undefined)).toBe(true);
});
