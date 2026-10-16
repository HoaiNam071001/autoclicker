// Thanh bên thư viện kịch bản (artboard Main). Hiện chỉ có kịch bản đang soạn; tìm, danh sách, mở gần đây,
// Mới / Mở / Lưu dựng ở bước 27b, nối dữ liệu ở bước 21.
// Thu gọn / mở rộng bằng nút ở đầu cột hoặc Ctrl/Cmd+B; trạng thái nhớ trên máy (localStorage).
import { useEffect, useState } from 'react';
import { useScript } from '../scriptStore';
import { SidebarToggle } from './icons';

const KEY = 'library-collapsed';
const isMac = navigator.platform.startsWith('Mac');
const shortcut = isMac ? '⌘B' : 'Ctrl+B';

const loadCollapsed = () => {
  try {
    return localStorage.getItem(KEY) === '1';
  } catch {
    return false;
  }
};

export default function Library() {
  const { name } = useScript();
  const [collapsed, setCollapsed] = useState(loadCollapsed);
  useEffect(() => {
    try {
      localStorage.setItem(KEY, collapsed ? '1' : '0');
    } catch { /* không lưu được thì thôi – chỉ là tiện ích */ }
  }, [collapsed]);
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() !== 'b' || !(isMac ? e.metaKey : e.ctrlKey)) return;
      e.preventDefault();
      setCollapsed((c) => !c);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const toggle = (
    <button onClick={() => setCollapsed(!collapsed)} aria-expanded={!collapsed}
      aria-label={collapsed ? 'Mở rộng thư viện' : 'Thu gọn thư viện'} title={`${collapsed ? 'Mở rộng' : 'Thu gọn'} thư viện (${shortcut})`}
      className="grid h-7 w-7 flex-none place-items-center rounded-md text-muted hover:bg-hover hover:text-ink">
      <SidebarToggle open={!collapsed} />
    </button>
  );

  if (collapsed) {
    return (
      <nav aria-label="Thư viện kịch bản (thu gọn)" className="flex w-11 flex-none flex-col items-center border-r border-line bg-side py-3">
        {toggle}
      </nav>
    );
  }
  return (
    <nav aria-label="Thư viện kịch bản" className="flex w-[236px] flex-none flex-col gap-3 overflow-y-auto border-r border-line bg-side px-2 py-3 max-[1100px]:w-[200px]">
      <div className="flex flex-col gap-0.5">
        <div className="flex items-center justify-between pr-0.5 pl-3">
          <span className="text-xs font-semibold text-muted">Kịch bản</span>
          {toggle}
        </div>
        <span className="relative flex h-9 items-center rounded-md bg-hover px-3 font-semibold">
          <span className="absolute top-2.5 left-0 h-4 w-[3px] rounded-sm bg-accent" />
          <span className="truncate">{name}</span>
        </span>
      </div>
    </nav>
  );
}
