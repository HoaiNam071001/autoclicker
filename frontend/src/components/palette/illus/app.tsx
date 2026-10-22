// Hoạt hoạ nhóm ứng dụng / cửa sổ / trình duyệt thường (bước 47).
import { Arr, Browser, Chip, Cur, Lines, Move, Ripple, Scenes, Show, Spinner, Stage, T, Tick, Tween, Typed, Win } from './parts';

const tab = (x: number, w: number, label: string, on = false) => (
  <g><path d={`M${x} 40v-6a4 4 0 0 1 4-4h${w - 8}a4 4 0 0 1 4 4v6z`} className={on ? 'fill-accent-soft stroke-accent' : 'fill-field stroke-line-strong'} /><T x={x + w / 2} y={38} anchor="middle" s={7.5}>{label}</T></g>
);

export const app: Scenes = {
  open: () => (
    <Stage>
      <Win x={16} y={12} w={208} h={111} />
      <rect x={28} y={22} width="184" height="11" rx="5.5" className="fill-field" />
      <Typed x={36} y={30.2} w={110} text="youtube.com/watch?v=…" s={7} from={0.1} to={0.45} />
      <Show from={0.5}><rect x={28} y={42} width="184" height="70" rx="4" className="fill-field" /><path d="M108 58v38l34-19z" className="fill-danger" /></Show>
      <Show to={0.5}><T x={120} y={80} anchor="middle" cls="fill-muted" s={8}>đang mở…</T></Show>
    </Stage>
  ),
  openWeb: () => (
    <Stage>
      <Browser url="https://trang-dich.vn" />
      <Show from={0.4}><Lines x={34} y={46} ws={[120, 80, 150]} /><rect x={34} y={78} width="60" height="22" rx="5" className="fill-accent-soft stroke-accent" /></Show>
      <Show from={0.4}><Chip x={150} y={96} w={66} label="Chrome của app" cls="fill-surface stroke-ok" /></Show>
      <Show to={0.4}><Spinner x={120} y={74} r={11} /></Show>
    </Stage>
  ),
  window: () => (
    <Stage>
      <rect x={120} y={30} width="104" height="86" rx="6" className="fill-none stroke-accent" strokeWidth="2" strokeDasharray="5 3" />
      <T x={172} y={126} anchor="middle" s={7.5} cls="fill-muted">vị trí + cỡ đặt sẵn</T>
      <Move pts={[[0, 20, 40], [0.25, 20, 40], [0.7, 120, 30]]}>
        <rect width="64" height="50" rx="5" className="fill-surface stroke-ink"><Tween attr="width" pts={[[0, 64], [0.25, 64], [0.7, 104]]} /><Tween attr="height" pts={[[0, 50], [0.25, 50], [0.7, 86]]} /></rect>
        <rect width="64" height="9" rx="4" className="fill-chrome"><Tween attr="width" pts={[[0, 64], [0.25, 64], [0.7, 104]]} /></rect>
      </Move>
      <Show from={0.7}><Tick x={214} y={22} from={0} /></Show>
    </Stage>
  ),
  app: () => (
    <Stage><Win />
      <rect x={90} y={36} width="60" height="60" rx="14" className="fill-accent-soft stroke-accent" strokeWidth="1.6" /><T x={120} y={74} anchor="middle" s={24} cls="fill-accent-text">A</T>
      <circle cx="146" cy="40" r="7" className="fill-ok" /><Show from={0.55}><circle cx="146" cy="40" r="7" className="fill-line-strong" /></Show>
      <Show to={0.55}><T x={120} y={112} anchor="middle" s={8.5} cls="fill-ok">đang chạy</T></Show>
      <Show from={0.55}><T x={120} y={112} anchor="middle" s={8.5} cls="fill-muted">đã thoát</T></Show>
    </Stage>
  ),
  browserTab: () => (
    <Stage>
      <Win /><rect x={16} y={12} width="208" height="30" rx="6" className="fill-chrome" />
      {tab(24, 60, 'Mua sắm')}{tab(88, 60, 'Giỏ hàng', true)}{tab(152, 60, 'Tin tức')}
      <Show from={0.5}><rect x={150} y={26} width="64" height="16" className="fill-bg" /><T x={182} y={38} anchor="middle" cls="fill-muted" s={8}>đã đóng</T></Show>
      <Lines x={30} y={56} ws={[120, 90]} />
      <Cur pts={[[0, 40, 90], [0.2, 118, 38], [0.55, 205, 34]]} /><Ripple x={118} y={36} t={0.2} r={8} />
    </Stage>
  ),
  browserGo: () => (
    <Stage>
      <Browser url="" />
      <Show to={0.4}><Typed x={36} y={30.2} w={60} text="shop.vn/a" s={7} from={0} to={0.01} /><Lines x={34} y={46} ws={[100, 70]} /></Show>
      <Show from={0.4} to={0.6}><Spinner x={120} y={76} r={11} /></Show>
      <Show from={0.4}><text x={36} y={30.2} fontSize="7" className="fill-ink font-mono">shop.vn/b</text></Show>
      <Show from={0.6}><Lines x={34} y={46} ws={[130, 60, 100]} /><Tick x={200} y={50} from={0} /></Show>
    </Stage>
  ),
  browserNorm: () => (
    <Stage>
      <rect x={0} y={0} width="240" height="135" rx="8" className="fill-surface stroke-line-strong"><Tween attr="width" pts={[[0, 240], [0.45, 240], [0.8, 190]]} /><Tween attr="height" pts={[[0, 135], [0.45, 135], [0.8, 105]]} /></rect>
      <Show to={0.45}><T x={120} y={70} anchor="middle" cls="fill-muted" s={9}>toàn màn hình</T></Show>
      <Show from={0.8}><Chip x={70} y={50} w={52} label="zoom 100%" cls="fill-surface stroke-ok" /></Show>
    </Stage>
  ),
  readTable: () => (
    <Stage><Win w={130} />
      {[0, 1, 2].map((r) => [0, 1, 2].map((c) => <rect key={`${r}${c}`} x={26 + c * 38} y={30 + r * 24} width="36" height="22" rx="3" className={r === 0 ? 'fill-soft stroke-line-strong' : 'fill-field stroke-line'} />))}
      {[1, 2].map((r) => <Show key={r} from={0.2 + r * 0.25} to={0.95}><rect x={26} y={30 + r * 24} width="112" height="22" rx="3" className="fill-accent-soft stroke-accent" /></Show>)}
      <Arr x1={150} y1={70} x2={168} y2={70} cls="stroke-accent" />
      <Show from={0.5}><Chip x={172} y={56} w={56} label="rows" /></Show>
    </Stage>
  ),
  pageText: () => (
    <Stage><Win w={130} />
      {[0, 1, 2, 3].map((i) => <g key={i}><rect x={26} y={28 + i * 22} width="112" height="16" rx="3" className="fill-field" /><rect x={32} y={34 + i * 22} width={70 + (i % 2) * 24} height="4" rx="2" className="fill-line-strong" />
        <Show from={0.15 + i * 0.17} to={0.95}><rect x={26} y={28 + i * 22} width="112" height="16" rx="3" className="fill-accent-soft stroke-accent" /></Show></g>)}
      <Arr x1={150} y1={70} x2={168} y2={70} cls="stroke-accent" />
      <rect x={172} y={26} width="56" height="86" rx="6" className="fill-surface stroke-line-strong" />
      {[0, 1, 2, 3].map((i) => <Show key={i} from={0.2 + i * 0.17}><rect x={178} y={32 + i * 20} width="44" height="14" rx="3" className="fill-accent-soft" /></Show>)}
    </Stage>
  ),
  elementRead: () => (
    <Stage><Win />
      <T x={30} y={44} s={8} cls="fill-muted">Số dư</T>
      <rect x={30} y={50} width="110" height="24" rx="5" className="fill-field stroke-line-strong" /><T x={40} y={67} s={12}>1.250.000 đ</T>
      <Show from={0.3}><rect x={26} y={46} width="118" height="32" rx="4" className="fill-none stroke-accent" strokeWidth="1.6" strokeDasharray="4 2" /></Show>
      <Show from={0.55}><Arr x1={146} y1={62} x2={164} y2={62} cls="stroke-accent" /><Chip x={166} y={55} w={56} label="{{số_dư}}" /></Show>
    </Stage>
  ),
};
