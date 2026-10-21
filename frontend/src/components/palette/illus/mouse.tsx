// Hoạt hoạ nhóm chuột (bước 47): click, double, giữ, kéo thả, vẽ đường, cuộn, click theo ảnh / chữ / phần tử.
import { Arr, Btn, Clip, Cur, Draw, Gem, Move, Ln, Lines, Ring, Ripple, Scenes, Show, Stage, T, Tween, Win } from './parts';

const page = <><Win /><Lines x={28} y={26} ws={[90, 60]} /></>;

export const mouse: Scenes = {
  click: () => (
    <Stage>{page}
      <Btn x={90} y={62} w={60} h={22} label="Mua ngay" />
      <Show from={0.42} to={0.7}><rect x={90} y={62} width="60" height="22" rx="5" className="fill-accent-soft stroke-accent" /><T x={120} y={76} anchor="middle">Mua ngay</T></Show>
      <Cur pts={[[0, 40, 112], [0.42, 118, 74]]} /><Ripple x={118} y={74} t={0.42} />
    </Stage>
  ),
  double: () => (
    <Stage>{page}
      <rect x={96} y={52} width="48" height="40" rx="6" className="fill-field stroke-line-strong" />
      <path d="M108 60h14l6 6v18h-20z" className="fill-surface stroke-muted" />
      <Show from={0.6}><rect x={56} y={26} width="128" height="82" rx="6" className="fill-surface stroke-accent" strokeWidth="2" /><Lines x={68} y={40} ws={[90, 70, 80]} /></Show>
      <Cur pts={[[0, 40, 112], [0.3, 124, 76]]} />
      <Ripple x={124} y={76} t={0.34} /><Ripple x={124} y={76} t={0.46} />
    </Stage>
  ),
  hold: () => (
    <Stage>{page}
      <circle cx="120" cy="64" r="22" className="fill-field stroke-line-strong" strokeWidth="2" />
      <circle cx="120" cy="64" r="22" fill="none" className="stroke-accent" strokeWidth="4" strokeLinecap="round" strokeDasharray="138" strokeDashoffset="138" transform="rotate(-90 120 64)">
        <Tween attr="stroke-dashoffset" pts={[[0, 138], [0.3, 138], [0.75, 0], [1, 0]]} />
      </circle>
      <Show from={0.75}><T x={120} y={67} anchor="middle" cls="fill-ok">✓</T></Show>
      <Cur pts={[[0, 50, 112], [0.28, 120, 66]]} /><Ripple x={120} y={64} t={0.28} r={26} />
      <T x={120} y={104} anchor="middle" cls="fill-muted" s={8}>giữ chuột… rồi nhả</T>
    </Stage>
  ),
  drag: () => (
    <Stage>{page}
      <rect x={156} y={52} width="50" height="42" rx="6" className="fill-accent-soft stroke-accent" strokeDasharray="4 2" />
      <T x={181} y={77} anchor="middle" cls="fill-muted" s={8}>Thư mục</T>
      <Move pts={[[0, 50, 84], [0.28, 50, 84], [0.72, 170, 84]]}>
        <rect x={-14} y={-14} width="28" height="24" rx="4" className="fill-surface stroke-ink" /><Ln x={-9} y={-6} w={18} />
      </Move>
      <Cur pts={[[0, 40, 112], [0.25, 50, 84], [0.72, 170, 84]]} />
      <Arr x1={74} y1={106} x2={150} y2={106} />
    </Stage>
  ),
  path: () => (
    <Stage>{page}
      <Draw d="M50 96 L88 44 L120 88 L152 44 L190 96" from={0.12} to={0.82} />
      <Cur pts={[[0, 40, 112], [0.12, 50, 96], [0.32, 88, 44], [0.47, 120, 88], [0.62, 152, 44], [0.82, 190, 96]]} />
    </Stage>
  ),
  scroll: () => (
    <Stage>{page}
      <Clip x={28} y={26} w={150} h={88}>
        <Move pts={[[0, 0, 0], [0.2, 0, 0], [0.7, 0, -44]]}>
          {[0, 1, 2, 3, 4, 5].map((i) => <g key={i}><rect x={28} y={28 + i * 22} width="150" height="18" rx="4" className="fill-field" /><Ln x={36} y={35 + i * 22} w={60 + (i % 3) * 25} /></g>)}
        </Move>
      </Clip>
      <rect x={196} y={30} width="8" height="80" rx="4" className="fill-field" />
      <rect x={196} y={32} width="8" height="26" rx="4" className="fill-accent"><Tween attr="y" pts={[[0, 32], [0.2, 32], [0.7, 78]]} /></rect>
      <Arr x1={120} y1={118} x2={120} y2={128} cls="stroke-accent" />
    </Stage>
  ),
  clickImage: () => (
    <Stage>{page}
      <Gem x={70} y={70} s={1.3} /><circle cx="120" cy="70" r="9" className="fill-ok" /><rect x={152} y={61} width="18" height="18" rx="3" className="fill-danger" />
      <Ring x={102} y={52} w={36} h={36} from={0.25} />
      <Cur pts={[[0, 40, 116], [0.45, 120, 70]]} /><Ripple x={120} y={70} t={0.45} />
    </Stage>
  ),
  clickText: () => (
    <Stage>{page}
      <Lines x={28} y={70} ws={[70]} /><T x={120} y={78} anchor="middle" s={11}>Đăng nhập</T>
      <Ring x={82} y={62} w={76} h={22} from={0.25} />
      <Cur pts={[[0, 40, 116], [0.45, 120, 76]]} /><Ripple x={120} y={76} t={0.45} />
    </Stage>
  ),
  clickElement: () => (
    <Stage>{page}
      <Btn x={92} y={62} w={56} h={22} label="Gửi" />
      <Ring x={88} y={58} w={64} h={30} from={0.25} />
      <Show from={0.25}><rect x={84} y={40} width="72" height="12" rx="3" className="fill-ink" /><text x={120} y={49} fontSize="7" textAnchor="middle" className="fill-surface font-mono">button · Gửi</text></Show>
      <Cur pts={[[0, 40, 116], [0.5, 120, 74]]} /><Ripple x={120} y={74} t={0.5} />
    </Stage>
  ),
};
