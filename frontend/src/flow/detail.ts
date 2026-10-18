// Dòng chi tiết thẻ node hành động (vị trí / nội dung / cách tìm) – flow/view.ts ghép thêm thời gian chờ sau node.
import { fmtAxis, fmtPoint, usesVar } from '../anchor';
import { fmtKey } from '../keys';
import { fmtShot } from './shot';
import { fmtClip } from './clipboard';
import { fmtHTTP } from './http';
import { fmtFile } from './file';
import { fmtSystem } from './system';
import { fmtWaitUntil } from './waitUntil';
import { fmtPath } from './shape';
import { fmtWebStep, isWebStep } from './web';
import { fmtElement } from './element';
import { fmtMs, numVarOf, SCROLL_ARROW, SCROLL_NOTCH_PX, type Step, type StepPoint } from '../types';

const fmtSigned = (v: number) => (v > 0 ? `+${v}` : `${v}`);
const fmtPair = (p: StepPoint) => (usesVar(p) ? `${p.xVar}, ${p.yVar}` : `${fmtAxis(p.x, p.ax)}, ${fmtAxis(p.y, p.ay)}`);
// refOf: id node → mã (n2) để ghi nguồn chỗ thấy của click theo ảnh.
export const fmtPos = (s: Step, refOf: (id: string) => string = () => '?') =>
  s.type === 'drag'
    ? `${s.dragFromPrev ? 'Từ điểm trước' : `A ${fmtPair(s.points[0])}`} → ${
      s.dragRelative ? `Δ ${usesVar(s.points[1]) ? fmtPair(s.points[1]) : `${fmtSigned(s.points[1].x)},${fmtSigned(s.points[1].y)}`}` : `B ${fmtPair(s.points[1])}`}`
    : s.type === 'type' ? `"${s.text}"`
    : s.type === 'key' ? fmtKey(s)
    : s.type === 'scroll' ? `${SCROLL_ARROW[s.direction ?? 'down']} ${numVarOf(s, 'amount') !== undefined ? `${numVarOf(s, 'amount')} nấc` // nấc từ biến: chỉ biết lúc chạy
      : `${s.amount ?? 0} nấc ≈ ${(s.amount ?? 0) * SCROLL_NOTCH_PX} px`} ${s.points[0] ? `tại ${fmtPoint(s.points[0])}` : 'tại chuột'}`
    : s.type === 'open' ? `${s.target || '—'}${s.webBrowser ? ' · Chrome của app' : s.app ? ` · ${s.app}` : ''}`
    : isWebStep(s.type) ? fmtWebStep(s)
    : s.type === 'screenshot' ? fmtShot(s, refOf)
    : s.type === 'clipboard' ? fmtClip(s)
    : s.type === 'http' ? fmtHTTP(s)
    : s.type === 'file' ? fmtFile(s)
    : s.type === 'system' ? fmtSystem(s)
    : s.type === 'waitUntil' ? fmtWaitUntil(s)
    : s.type === 'path' ? fmtPath(s)
    : s.type === 'elementRead' ? fmtElement(s)
    : s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search')
      ? `${s.type === 'clickImage' ? 'tự tìm · ' : ''}${s.findBy === 'element' ? fmtElement(s) : s.findBy === 'text'
        ? (s.findText?.trim() ? `chữ “${s.findText}”${(s.occurrence ?? 0) > 1 ? ` (chỗ thứ ${s.occurrence})` : ''} · tối đa ${fmtMs(s.timeoutMs ?? 0)}` : 'chưa nhập chữ cần tìm')
        : s.template ? `khớp ≥ ${s.threshold}% · tối đa ${fmtMs(s.timeoutMs ?? 0)}` : 'chưa có ảnh mẫu'}`
    : s.type === 'clickImage' ? `chỗ ${s.imageFrom ? refOf(s.imageFrom) : '(chưa chọn node)'} thấy${s.dx || s.dy ? ` ${fmtSigned(s.dx ?? 0)},${fmtSigned(s.dy ?? 0)}` : ''}`
    : s.type === 'waitPixel' && s.points[0] ? `${fmtPoint(s.points[0])} = ${s.color} ±${s.tolerance} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`
    : s.points[0] ? fmtPoint(s.points[0]) : '—';
