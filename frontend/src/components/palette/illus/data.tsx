// Hoạt hoạ nhóm dữ liệu: clipboard, HTTP, tệp, âm lượng, tính toán, xử lý chữ, đặt biến (bước 47).
import { Arr, Chip, Move, Scenes, Show, Stage, T, Tween, Typed, Win } from './parts';

const box = (x: number, y: number, w: number, h: number, label?: string) => (
  <g><rect x={x} y={y} width={w} height={h} rx="6" className="fill-surface stroke-line-strong" />
    {label && <T x={x + w / 2} y={y + h / 2 + 3} anchor="middle" s={8} cls="fill-muted">{label}</T>}</g>
);
const packet = (pts: [number, number, number][], label: string) => (
  <Move pts={pts}><rect x={-14} y={-7} width="28" height="14" rx="7" className="fill-accent" /><T x={0} y={3} anchor="middle" s={7.5} cls="fill-on-accent">{label}</T></Move>
);

export const data: Scenes = {
  clipboard: () => (
    <Stage>
      {box(14, 40, 62, 56)}<T x={45} y={66} anchor="middle" s={10}>abc 123</T>
      <rect x={96} y={32} width="48" height="70" rx="6" className="fill-field stroke-line-strong" /><rect x={110} y={27} width="20" height="9" rx="3" className="fill-accent" />
      <Show from={0.4}><T x={120} y={72} anchor="middle" s={9}>abc 123</T></Show>
      {box(164, 40, 62, 56)}<Show from={0.8}><T x={195} y={66} anchor="middle" s={10}>abc 123</T></Show>
      <Arr x1={78} y1={68} x2={94} y2={68} /><Arr x1={146} y1={68} x2={162} y2={68} />
      <T x={45} y={112} anchor="middle" s={8} cls="fill-muted">Copy</T><T x={195} y={112} anchor="middle" s={8} cls="fill-muted">Paste</T>
    </Stage>
  ),
  http: () => (
    <Stage>
      {box(14, 38, 56, 58, 'Ghosthand')}{box(170, 38, 56, 58, 'Server')}
      <Arr x1={74} y1={56} x2={164} y2={56} /><Arr x1={166} y1={80} x2={76} y2={80} />
      {packet([[0.1, 88, 56], [0.4, 150, 56]], 'GET')}{packet([[0.5, 150, 80], [0.85, 90, 80]], '{ }')}
      <Show from={0.85}><Chip x={14} y={104} w={56} label="200 OK ✓" cls="fill-surface stroke-ok" /></Show>
    </Stage>
  ),
  file: () => (
    <Stage><Win />
      {[0, 1, 2, 3].map((i) => <g key={i}><rect x={28} y={30 + i * 18} width="184" height="16" className={i === 0 ? 'fill-soft' : 'fill-field'} />
        {[0, 1, 2].map((j) => <rect key={j} x={34 + j * 60} y={36 + i * 18} width={i === 0 ? 40 : 30 + ((i + j) % 3) * 8} height="4" rx="2" className="fill-line-strong" />)}</g>)}
      <Show from={0.4}><rect x={28} y={102} width="184" height="16" className="fill-accent-soft stroke-accent" /><T x={34} y={113} s={8}>ĐH-4821 · 1.250.000 · xong</T></Show>
      <Show from={0.4}><T x={206} y={113} anchor="end" cls="fill-accent-text">+ dòng mới</T></Show>
    </Stage>
  ),
  shell: () => (
    <Stage>
      <rect x={30} y={22} width="180" height="94" rx="7" className="fill-ink" opacity="0.92" />
      <Typed x={40} y={46} w={120} text="$ ls | wc -l" from={0.1} to={0.45} s={10} />
      <Show from={0.55}><T x={40} y={68} s={10} cls="fill-ok font-mono">42</T></Show>
      <Show from={0.7}><Chip x={150} y={92} w={52} label="exit 0 ✓" cls="fill-surface stroke-ok" /></Show>
    </Stage>
  ),
  system: () => (
    <Stage>
      <path d="M70 58h14l18-14v46L84 76H70z" className="fill-ink" />
      {[0, 1, 2].map((i) => <path key={i} d={`M${112 + i * 10} ${62 - i * 6}a${10 + i * 6} ${10 + i * 6} 0 0 1 0 ${12 + i * 12}`} fill="none" className="stroke-accent" strokeWidth="3" strokeLinecap="round">
        <Tween attr="opacity" pts={[[0, 0.2], [0.2 + i * 0.15, 0.2], [0.3 + i * 0.15, 1], [0.9, 1], [0.95, 0.2]]} /></path>)}
      <rect x={60} y={108} width="120" height="6" rx="3" className="fill-field" />
      <rect x={60} y={108} width="20" height="6" rx="3" className="fill-accent"><Tween attr="width" pts={[[0, 20], [0.25, 20], [0.75, 100]]} /></rect>
    </Stage>
  ),
  calc: () => (
    <Stage><Win />
      <Show from={0.05}><T x={44} y={62} s={20}>12</T></Show>
      <Show from={0.25}><T x={84} y={62} s={20} cls="fill-accent-text">×</T></Show>
      <Show from={0.4}><T x={110} y={62} s={20}>3</T></Show>
      <Show from={0.55}><T x={136} y={62} s={20} cls="fill-accent-text">=</T></Show>
      <Show from={0.7}><rect x={156} y={40} width="52" height="30" rx="6" className="fill-accent-soft stroke-accent" /><T x={182} y={61} s={20} anchor="middle">36</T></Show>
      <Show from={0.85}><Chip x={150} y={84} w={64} label="{{vars.tong}}" /></Show>
    </Stage>
  ),
  textNumber: () => (
    <Stage><Win />
      <T x={120} y={52} anchor="middle" s={15}>1.234.000 đ</T>
      <Show from={0.4}><Arr x1={120} y1={60} x2={120} y2={80} cls="stroke-accent" /></Show>
      <Show from={0.5}><rect x={78} y={84} width="84" height="24" rx="6" className="fill-accent-soft stroke-accent" /><T x={120} y={101} anchor="middle" s={14} cls="fill-ink">1234000</T></Show>
    </Stage>
  ),
  textMatch: () => (
    <Stage><Win />
      <T x={30} y={50} s={10} cls="fill-muted">Đơn của bạn mã</T>
      <Show from={0.25}><rect x={122} y={39} width="64" height="16" rx="3" className="fill-accent-soft stroke-accent" /></Show>
      <T x={126} y={51} s={10}>ĐH-4821</T><T x={190} y={51} s={10} cls="fill-muted">đã</T>
      <T x={30} y={68} s={10} cls="fill-muted">giao thành công.</T>
      <Show from={0.5}><Arr x1={154} y1={58} x2={154} y2={84} cls="stroke-accent" /><Chip x={118} y={88} w={72} label="ĐH-4821" /></Show>
    </Stage>
  ),
  textSplit: () => (
    <Stage><Win />
      <Show to={0.4}><T x={120} y={70} anchor="middle" s={14}>táo, cam, nho</T></Show>
      {['táo', 'cam', 'nho'].map((w, i) => (
        <Show key={w} from={0.4}><Move pts={[[0.4, 120, 70], [0.65, 60 + i * 60, 84]]}><rect x={-20} y={-12} width="40" height="22" rx="6" className="fill-accent-soft stroke-accent" /><T x={0} y={3} anchor="middle">{w}</T></Move></Show>
      ))}
    </Stage>
  ),
  textClean: () => (
    <Stage><Win />
      <Show to={0.45}><T x={120} y={56} anchor="middle" s={15}>  Đà  Nẵng </T><T x={120} y={86} anchor="middle" s={8} cls="fill-muted">thừa khoảng trắng · có dấu · Hoa</T></Show>
      <Show from={0.45}><T x={120} y={56} anchor="middle" s={15} cls="fill-accent-text">da-nang</T><T x={120} y={86} anchor="middle" s={8} cls="fill-ok">sạch · không dấu · thường</T></Show>
    </Stage>
  ),
  setVar: () => (
    <Stage><Win />
      <rect x={76} y={40} width="88" height="40" rx="8" className="fill-field stroke-accent" strokeWidth="1.6" />
      <T x={88} y={56} s={8} cls="fill-muted">đã_xử_lý</T>
      {['0', '1', '2'].map((v, i) => <Show key={v} from={i === 0 ? 0 : 0.33 * i} to={i === 2 ? 1 : 0.33 * (i + 1)}><T x={140} y={72} anchor="middle" s={20}>{v}</T></Show>)}
      <Chip x={86} y={92} w={68} label="{{vars.đã_xử_lý}}" />
    </Stage>
  ),
};
