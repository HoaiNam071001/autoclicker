import {ComponentType, useEffect, useState} from 'react';
import MainScreen from './screens/MainScreen';
import PreviewScreen from './screens/PreviewScreen';
import OverlayScreen from './screens/OverlayScreen';
import MiniBar from './screens/MiniBar';
import DemoScreen from './screens/DemoScreen';
import DemoOverlay from './screens/DemoOverlay';
import DemoMiniBar from './screens/DemoMiniBar';

// Chọn màn theo hash URL (#preview, #overlay, #minibar, #demo, #demo-overlay, #demo-minibar). setWindowMode() (windowMode.ts) đổi
// kích thước cửa sổ bên Go rồi đổi hash; trên trình duyệt đổi hash tay để xem từng màn.
const screens: Record<string, ComponentType> = {
    main: MainScreen,
    preview: PreviewScreen,
    overlay: OverlayScreen,
    minibar: MiniBar,
    // giao diện đầy đủ với dữ liệu mẫu, kể cả tính năng chưa có logic (bước 27b)
    demo: DemoScreen,
    'demo-overlay': DemoOverlay,
    'demo-minibar': DemoMiniBar,
};

function App() {
    const [hash, setHash] = useState(location.hash.slice(1));

    useEffect(() => {
        const onHash = () => setHash(location.hash.slice(1));
        window.addEventListener('hashchange', onHash);
        return () => window.removeEventListener('hashchange', onHash);
    }, []);

    const Screen = screens[hash] ?? MainScreen;
    return <Screen/>;
}

export default App
