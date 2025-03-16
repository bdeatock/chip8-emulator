let wasmReady = false;
window.addEventListener("message", function (event) {
  if (event.data && event.data.type === "wasmReady") {
    wasmReady = true;
  }
});

function refocusEmulator() {
  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.focus();

  iframe.contentWindow.postMessage(
    {
      type: "focus",
    },
    window.location.origin
  );
}

function handleFileSelect(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = function (e) {
    const arrayBuffer = e.target.result;
    iframe.contentWindow.postMessage(
      {
        type: "loadROM",
        data: arrayBuffer,
      },
      window.location.origin
    );
  };
  reader.readAsArrayBuffer(file);
}

function handleRomSelect(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  const rom = event.target.value;
  if (!rom) return;

  fetch(`/roms/${rom}`)
    .then((response) => {
      if (!response.ok) {
        throw new Error(
          `Failed to load ROM: ${response.status} ${response.statusText}`
        );
      }
      return response.arrayBuffer();
    })
    .then((arrayBuffer) => {
      iframe.contentWindow.postMessage(
        {
          type: "loadROM",
          data: arrayBuffer,
        },
        window.location.origin
      );
      refocusEmulator();
    })
    .catch((error) => {
      console.error("Error loading ROM:", error);
    });
}

function handleToggleLegacyShift(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  event.target.checked = iframe.contentWindow.toggleLegacyShift();

  refocusEmulator();
}

function handleToggleLegacyJump(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  event.target.checked = iframe.contentWindow.toggleLegacyJump();

  refocusEmulator();
}

function handleToggleLegacyStoreLoad(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  event.target.checked = iframe.contentWindow.toggleLegacyStoreLoad();

  refocusEmulator();
}

function handleSwitchMode(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.contentWindow.switchMode();
  refocusEmulator();
}

function handleSetCycleRate(event) {
  if (!wasmReady) return;

  const iframe = document.querySelector("iframe");
  if (!iframe) return;

  iframe.contentWindow.updateCycleRate(parseInt(event.target.value));
  refocusEmulator();
}
