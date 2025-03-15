const loadingElement = document.getElementById("loading");
function doneLoading() {
  loadingElement.classList.add("hidden");

  window.addEventListener("message", function (event) {
    if (event.data && event.data.type === "loadROM") {
      const uint8Array = new Uint8Array(event.data.data);
      loadROM(uint8Array);
    }
  });

  if (window.parent) {
    window.parent.postMessage({ type: "wasmReady" }, window.location.origin);
  }
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
  .then((result) => {
    go.run(result.instance);
    doneLoading();
  })
  .catch((err) => {
    console.error("Error loading WASM:", err);
    loadingElement.innerHTML = `<div style="color: white; text-align: center;">
      <p>Error loading emulator</p>
      <p>${err.message}</p>
    </div>`;
  });
