import {ComponentType, useEffect, useState} from 'react';
import MainScreen from './screens/MainScreen';
import PreviewScreen from './screens/PreviewScreen';
import OverlayScreen from './screens/OverlayScreen';
import MiniBar from './screens/MiniBar';

// Chọn màn theo hash URL (#preview, #overlay, #minibar). setWindowMode() (windowMode.ts) đổi
// kích thước cửa sổ bên Go rồi đổi hash; trên trình duyệt đổi hash tay để xem từng màn.
const screens: Record<string, ComponentType> = {
    main: MainScreen,
    preview: PreviewScreen,
    overlay: OverlayScreen,
    minibar: MiniBar,
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
