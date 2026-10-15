import {ComponentType, useEffect, useState} from 'react';
import MainScreen from './screens/MainScreen';
import PreviewScreen from './screens/PreviewScreen';
import OverlayScreen from './screens/OverlayScreen';
import MiniBar from './screens/MiniBar';

// Tạm thời: chọn màn bằng hash URL (#preview, #overlay, #minibar) để xem bản thiết kế.
// Sẽ thay bằng state chế độ cửa sổ (main / overlay / minibar) khi làm spike overlay (bước 5).
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
