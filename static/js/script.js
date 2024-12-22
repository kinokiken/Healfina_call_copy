let audioContext;
let ws;
let audioProcessor;
let audioBufferQueue = [];
let isPlaying = false;
let playbackTime = 0;
let isSessionActive = false;

document.getElementById("startButton").addEventListener("click", async function () {
    const button = document.getElementById("startButton");

    if (isSessionActive) {
        stopSession();
        button.innerHTML = '<i class="fa fa-phone"></i>    Начать звонок';
        button.classList.remove("active-call");
    } else {
        try {
            await startSession();
            button.innerHTML = '<i class="fa fa-bell-slash"></i>    Завершить звонок';
            button.classList.add("active-call");
        } catch (err) {
            console.error("Error starting session:", err);
            alert("Failed to start session. Please check permissions or server connection.");
        }
    }

    isSessionActive = !isSessionActive;
});

async function startSession() {
    console.log("Initializing audio...");
    audioContext = new AudioContext({ sampleRate: 24000 });
    playbackTime = audioContext.currentTime;
    await audioContext.audioWorklet.addModule('/static/js/processor.js');

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    const source = audioContext.createMediaStreamSource(stream);

    audioProcessor = new AudioWorkletNode(audioContext, 'audio-processor');
    source.connect(audioProcessor);

    ws = new WebSocket('ws://localhost:8080/stream');
    ws.binaryType = 'arraybuffer';

    ws.onopen = () => {
        console.log("WebSocket connected");

        audioProcessor.port.onmessage = (event) => {
            if (ws.readyState === WebSocket.OPEN) {
                ws.send(event.data);
            }
        };
    };

    ws.onclose = () => {
        source.disconnect(audioProcessor);
        console.log("WebSocket closed");
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        alert("WebSocket connection error. Check the server and try again.");
    };

    ws.onmessage = (event) => {
        console.log('Received audio data from server');

        if (event.data instanceof ArrayBuffer) {
            const pcm16Data = new Int16Array(event.data);
            const float32Data = pcm16ToFloat32(pcm16Data);

            audioBufferQueue.push(float32Data);

            playAudioQueue();
        } else {
            console.error('Unexpected data format:', event.data);
        }
    };
}

function stopSession() {
    console.log("Stopping session...");

    if (ws) {
        if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
            ws.close();
        }
        ws = null;
    }

    if (audioProcessor) {
        audioProcessor.port.onmessage = null;
        audioProcessor.disconnect();
        audioProcessor = null;
    }

    if (audioContext) {
        audioContext.close().then(() => {
            audioContext = null;
        });
    }

    audioBufferQueue = [];
    isPlaying = false;
    playbackTime = 0;

    console.log("Session stopped.");
}

function pcm16ToFloat32(pcm16Data) {
    const float32Data = new Float32Array(pcm16Data.length);
    for (let i = 0; i < pcm16Data.length; i++) {
        float32Data[i] = pcm16Data[i] / 32768;
    }
    return float32Data;
}

function playAudioQueue() {
    if (audioBufferQueue.length === 0 || !audioContext) return;

    if (playbackTime < audioContext.currentTime) {
        playbackTime = audioContext.currentTime;
    }

    while (audioBufferQueue.length > 0) {
        const audioData = audioBufferQueue.shift();

        const buffer = audioContext.createBuffer(1, audioData.length, audioContext.sampleRate);
        buffer.copyToChannel(audioData, 0);

        const source = audioContext.createBufferSource();
        source.buffer = buffer;
        source.connect(audioContext.destination);

        source.start(playbackTime);

        playbackTime += buffer.duration;
    }
}
