// Node cuộn: hướng, số nấc / quãng px, cuộn tại chuột / điểm, bản đồ cuộn tới đâu.
import NumInput from '../NumInput';
import Segmented from '../Segmented';
import AnchorField from '../steps/AnchorField';
import ScrollMap from '../steps/ScrollMap';
import { boundsOf, endpoints, setAnchor, setScrollAt, toScreen, updateStep, useScript } from '../../store';
import { pinnedInput, setPinnedValue } from '../../scriptInputs';
import { insideBounds } from '../../anchor';
import { SCROLL_NOTCH_PX, ScrollDirection, Step } from '../../types';
import { DIRECTION_OPTIONS, errText, Field, NumField, PickButton, SetPoint } from './shared';
import PointSource from './PointSource';

export function ScrollFields({ step: s, setPoint }: { step: Step; setPoint: SetPoint }) {
  const script = useScript();
  const { region } = script;
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
          <>
        <Field label="Hướng" wide group>
          <Segmented value={s.direction ?? 'down'} onChange={(direction: ScrollDirection) => set({ direction })} full label="Hướng cuộn"
            options={DIRECTION_OPTIONS} />
        </Field>
        <NumField label="Số nấc" step={s} field="amount" unit="nấc" min={1} stepper />
        <Field label="Quãng cuộn ≈" help={`Gõ px, tự đổi ra nấc (${SCROLL_NOTCH_PX} px / nấc). Cuộn dài hơn màn hình cũng được.`}>
          <NumInput value={(s.amount ?? 0) * SCROLL_NOTCH_PX} unit="px" min={SCROLL_NOTCH_PX} pinned={!!pinnedInput(script, s.id, 'amount')}
            onChange={(px) => {
              const amount = Math.max(1, Math.round(px / SCROLL_NOTCH_PX));
              const inp = pinnedInput(script, s.id, 'amount');
              if (inp) setPinnedValue(s.id, 'amount', amount);
              else set({ amount });
            }} />
        </Field>
        <Field label="Cuộn tại" group help={s.points.length ? 'Đặt chuột tại điểm rồi lăn (không phải chỗ cuộn tới).' : 'Cuộn ở chỗ chuột đang đứng – thường sau một bước click vào trang / danh sách cần cuộn.'}>
          <Segmented value={s.points.length ? 'point' : 'mouse'} onChange={(v) => setScrollAt(s.id, v === 'point')} full label="Cuộn tại"
            options={[{ id: 'mouse', label: 'Chuột' }, { id: 'point', label: 'Điểm' }]} />
        </Field>
        {s.points[0] ? (
          <PointSource step={s} index={0} setPoint={setPoint}
            label="Đặt chuột tại">
            {(menu) => <>
            <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
              pick={<PickButton step={s} />} menu={menu} />
            {!insideBounds(toScreen(s.points[0], region), boundsOf(region)) && (
              <span className={errText}>
                Điểm nằm ngoài {region ? 'vùng' : 'màn hình'} ({boundsOf(region).w} × {boundsOf(region).h}) – chuột không tới được, chạy sẽ báo lỗi.
                Muốn cuộn xa hơn thì tăng số nấc / quãng cuộn.
              </span>
            )}
            </>}
          </PointSource>
        ) : null}
        <Field label="Cuộn tới đâu" wide group>
          <ScrollMap bounds={boundsOf(region)} label={region ? 'Vùng' : 'Màn hình'} direction={s.direction ?? 'down'} amount={s.amount ?? 0}
            at={endpoints(script, region).get(s.id)?.start ?? null} />
        </Field>
      </>
  );
}
