import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AudioService {
  private audioContext: AudioContext | null = null;
  private ws: WebSocket | null = null;
  private audioProcessor: AudioWorkletNode | null = null;
  private audioBufferQueue: Float32Array[] = [];
  private isPlaying: boolean = false;
  private playbackTime: number = 0;

  constructor() { }

  async startSession() {
    console.log("Initializing audio...");
    this.audioContext = new AudioContext({ sampleRate: 24000 });
    this.playbackTime = this.audioContext.currentTime;
    await this.audioContext.audioWorklet.addModule('assets/js/processor.js');

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    const source = this.audioContext.createMediaStreamSource(stream);

    this.audioProcessor = new AudioWorkletNode(this.audioContext, 'audio-processor');
    source.connect(this.audioProcessor);

    this.ws = new WebSocket('ws://localhost:8080/stream');
    this.ws.binaryType = 'arraybuffer';

    this.ws.onopen = () => {
      console.log("WebSocket connected");

      this.audioProcessor!.port.onmessage = (event) => {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
          this.ws.send(event.data);
        }
      };
    };

    this.ws.onclose = () => {
      if (this.audioProcessor) {
        source.disconnect(this.audioProcessor);
      }
      console.log("WebSocket closed");
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      alert("WebSocket connection error. Check the server and try again.");
    };

    this.ws.onmessage = (event) => {
      console.log('Received audio data from server');

      if (event.data instanceof ArrayBuffer) {
        const pcm16Data = new Int16Array(event.data);
        const float32Data = this.pcm16ToFloat32(pcm16Data);

        this.audioBufferQueue.push(float32Data);

        this.playAudioQueue();
      } else {
        console.error('Unexpected data format:', event.data);
      }
    };
  }

  stopSession() {
    console.log("Stopping session...");

    if (this.ws) {
      if (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING) {
        this.ws.close();
      }
      this.ws = null;
    }

    if (this.audioProcessor) {
      this.audioProcessor.port.onmessage = null;
      this.audioProcessor.disconnect();
      this.audioProcessor = null;
    }

    if (this.audioContext) {
      this.audioContext.close().then(() => {
        this.audioContext = null;
      });
    }

    this.audioBufferQueue = [];
    this.isPlaying = false;
    this.playbackTime = 0;

    console.log("Session stopped.");
  }

  private pcm16ToFloat32(pcm16Data: Int16Array): Float32Array {
    const float32Data = new Float32Array(pcm16Data.length);
    for (let i = 0; i < pcm16Data.length; i++) {
      float32Data[i] = pcm16Data[i] / 32768;
    }
    return float32Data;
  }

  private playAudioQueue() {
    if (this.audioBufferQueue.length === 0 || !this.audioContext) return;

    if (this.playbackTime < this.audioContext.currentTime) {
      this.playbackTime = this.audioContext.currentTime;
    }

    while (this.audioBufferQueue.length > 0) {
      const audioData = this.audioBufferQueue.shift()!;

      const buffer = this.audioContext.createBuffer(1, audioData.length, this.audioContext.sampleRate);
      buffer.copyToChannel(audioData, 0);

      const source = this.audioContext.createBufferSource();
      source.buffer = buffer;
      source.connect(this.audioContext.destination);

      source.start(this.playbackTime);

      this.playbackTime += buffer.duration;
    }
  }
}
