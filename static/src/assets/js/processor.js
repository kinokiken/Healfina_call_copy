class AudioProcessor extends AudioWorkletProcessor {
    constructor() {
        super();
        this.sampleRate = 24000; // Устанавливаем частоту дискретизации
        this.bufferSize = 1024;  // Размер буфера
    }

    process(inputs, outputs, parameters) {
        const input = inputs[0];
        const output = outputs[0];

        if (input && input[0]) {
            const inputData = input[0];
            const outputData = output[0];

            // Преобразование входного сигнала Float32 в PCM16
            const audioData = new Int16Array(inputData.length);
            for (let i = 0; i < inputData.length; i++) {
                audioData[i] = Math.max(-32768, Math.min(32767, inputData[i] * 32767));
                outputData[i] = inputData[i]; // Копируем данные для мониторинга
            }
            console.log(inputData.length)

            // Отправляем данные в основной поток
            this.port.postMessage(audioData);
        }
        return true; // Продолжить обработку
    }
}

registerProcessor('audio-processor', AudioProcessor);
