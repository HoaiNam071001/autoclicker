import {useState} from 'react';
import {Greet} from "../wailsjs/go/main/App";

function App() {
    const [result, setResult] = useState("Auto Clicker – project skeleton");
    const [name, setName] = useState('');

    function greet() {
        Greet(name).then(setResult);
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center gap-4 p-6">
            <h1 className="text-2xl font-bold">Auto Clicker</h1>
            <p className="text-slate-300">{result}</p>
            <div className="flex gap-2">
                <input
                    className="rounded bg-slate-700 px-3 py-1.5 outline-none focus:ring-2 focus:ring-sky-500"
                    onChange={e => setName(e.target.value)}
                    autoComplete="off"
                    type="text"
                />
                <button className="rounded bg-sky-600 px-4 py-1.5 hover:bg-sky-500" onClick={greet}>
                    Greet
                </button>
            </div>
        </div>
    )
}

export default App
