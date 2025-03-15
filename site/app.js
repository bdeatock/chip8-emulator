let goWasm;
const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then(
  (result) => {
    goWasm = result.instance;
    go.run(goWasm);
  }
);

function handleFileSelect(event) {
  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = function (e) {
    const arrayBuffer = e.target.result;
    const uint8Array = new Uint8Array(arrayBuffer);
    const result = loadROM(uint8Array);
    if (result && result.error) {
      console.error("Error loading ROM:", result.error);
    }
  };
  reader.readAsArrayBuffer(file);
}

function handleSwitchMode(event) {
  switchMode();
}

function handleSetCycleRate(event) {
  updateCycleRate(parseInt(event.target.value));
}
