// Màn #demo-illus (bước 47): lưới mọi hoạt hoạ minh hoạ node để duyệt hình trong trình duyệt, không cần mở bảng Thêm node.
import { ILLUS } from '../components/palette/illus';
import { NODE_INFO } from '../flow/nodeInfo';

export default function DemoIllus() {
  return (
    <div className="min-h-screen bg-bg p-4 text-ink">
      <div className="grid grid-cols-[repeat(auto-fill,minmax(220px,1fr))] gap-3">
        {Object.entries(ILLUS).map(([k, scene]) => (
          <figure key={k} className="overflow-hidden rounded-card border border-line bg-surface">
            <div className="aspect-[240/135]">{scene()}</div>
            <figcaption className="px-2.5 py-1.5 text-xs"><span className="font-mono font-semibold">{k}</span> · <span className="text-muted">{NODE_INFO[k as keyof typeof NODE_INFO].desc.slice(0, 50)}</span></figcaption>
          </figure>
        ))}
      </div>
    </div>
  );
}
