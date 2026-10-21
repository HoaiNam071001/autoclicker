// Hoạt hoạ nhóm web qua Chrome của app (bước 47): luôn có khung trình duyệt + nhãn selector.
import { Browser, Btn, Chip, Clip, Cur, Lines, Move, Ring, Ripple, Scenes, Show, Spinner, Stage, T, Tick, Tween, Typed } from './parts';

const sel = (label: string, from = 0.2) => <Show from={from}><Chip x={20} y={112} w={Math.max(label.length * 5 + 14, 50)} label={label} cls="fill-surface stroke-accent" /></Show>;

export const web: Scenes = {
  webWait: () => (
    <Stage><Browser url="shop.vn">
      <Show to={0.5}><Spinner x={120} y={78} r={11} /></Show>
      <Show from={0.5}><Btn x={92} y={62} w={56} h={22} label="Đặt hàng" /><Ring x={88} y={58} w={64} h={30} from={0} /><Tick x={158} y={58} from={0} /></Show>
    </Browser>{sel('button#order', 0.5)}</Stage>
  ),
  webClick: () => (
    <Stage><Browser url="shop.vn"><Lines x={34} y={46} ws={[110, 70]} />
      <Btn x={92} y={68} w={56} h={22} label="Đặt hàng" />
      <Show from={0.45} to={0.7}><rect x={92} y={68} width="56" height="22" rx="5" className="fill-accent-soft stroke-accent" /><T x={120} y={82} anchor="middle" s={8.5}>Đặt hàng</T></Show>
      <Cur pts={[[0, 40, 116], [0.45, 118, 82]]} /><Ripple x={118} y={82} t={0.45} />
    </Browser>{sel('text=Đặt hàng', 0.1)}</Stage>
  ),
  webType: () => (
    <Stage><Browser url="shop.vn/dang-nhap"><T x={34} y={52} s={8} cls="fill-muted">Email</T>
      <rect x={34} y={58} width="170" height="22" rx="5" className="fill-field stroke-accent" strokeWidth="1.5" />
      <Typed x={42} y={73} w={120} text="ban@mail.com" from={0.2} to={0.7} />
      <Btn x={34} y={88} w={60} h={20} label="Đăng nhập" />
    </Browser>{sel('input[type=email]')}</Stage>
  ),
  webRead: () => (
    <Stage><Browser url="shop.vn"><T x={34} y={50} s={8} cls="fill-muted">Tổng tiền</T><T x={34} y={72} s={16}>1.250.000 đ</T>
      <Ring x={28} y={56} w={104} h={22} from={0.25} />
      <Show from={0.5}><Chip x={146} y={58} w={64} label="{{tổng}}" /></Show>
    </Browser>{sel('.total')}</Stage>
  ),
  webKey: () => (
    <Stage><Browser url="shop.vn"><rect x={34} y={48} width="170" height="22" rx="5" className="fill-field stroke-accent" strokeWidth="1.5" /><T x={42} y={63} s={9}>iphone 15</T>
      <rect x={84} y={82} width="72" height="26" rx="6" className="fill-field stroke-line-strong" /><T x={120} y={99} anchor="middle" s={10}>Enter ↵</T>
      <Show from={0.4}><rect x={84} y={82} width="72" height="26" rx="6" className="fill-accent" /><T x={120} y={99} anchor="middle" s={10} cls="fill-on-accent">Enter ↵</T></Show>
    </Browser><Show from={0.55}><Chip x={20} y={112} w={110} label="gửi, không cần Chrome ở trước" cls="fill-surface stroke-ok" /></Show></Stage>
  ),
  webScroll: () => (
    <Stage><Browser url="shop.vn/feed">
      <Clip x={28} y={38} w={184} h={78}>
        <Move pts={[[0, 0, 0], [0.2, 0, 0], [0.7, 0, -56]]}>
          {[0, 1, 2, 3, 4].map((i) => <g key={i}><rect x={28} y={40 + i * 28} width="184" height="22" rx="4" className="fill-field" /><Lines x={36} y={47 + i * 28} ws={[90 - i * 6]} /></g>)}
          <rect x={28} y={152} width="184" height="22" rx="4" className="fill-accent-soft stroke-accent" strokeWidth="1.6" />
        </Move>
      </Clip>
    </Browser>{sel('tới khi thấy: bài cần tìm', 0.2)}</Stage>
  ),
};
