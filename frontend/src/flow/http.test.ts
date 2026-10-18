import { describe, expect, it } from 'vitest';
import { fmtHTTP, isSecretKey, jsonBodyValid } from './http';
// Cùng dữ liệu với internal/model/http_test.go
import data from '../../../testdata/secret_keys.json';

describe('http', () => {
  it.each(data.cases)('isSecretKey($name) = $secret', (c) => expect(isSecretKey(c.name)).toBe(c.secret));
  it('jsonBodyValid', () => {
    expect(jsonBodyValid('')).toBeNull();
    expect(jsonBodyValid('{"a": "{{ten}}"}')).toBeNull();
    expect(jsonBodyValid('{"a": 1}')).toBe(true);
    expect(jsonBodyValid('{a: 1}')).toBe(false);
  });
  it('fmtHTTP', () => {
    expect(fmtHTTP({ url: 'https://' })).toBe('GET (chưa nhập URL)');
    expect(fmtHTTP({ method: 'POST', url: 'https://api.x.com/orders', bodyType: 'json', query: [{ key: 'a', value: '1' }, { key: ' ', value: '' }] }))
      .toBe('POST api.x.com/orders · 1 tham số · JSON');
  });
});
